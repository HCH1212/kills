package service

import (
	"fmt"
	"log"
	"rpc-client/model"
	"rpc-client/resp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 用户参与活动下订单
func JoinKill(c *gin.Context) {
	var u model.User
	var h model.Hoster
	var err error
	var s string

	//token验证身份
	u.Token = c.GetHeader("usertoken")

	//利用redis分布式锁保证秒杀秩序
	lockKey := "my_lock"            //锁名
	lockValue := u.Token            //锁值
	lockTimeout := 10 * time.Second //锁的超时时间
	// 尝试获取锁
	success, err := model.RDB.SetNX(c, lockKey, lockValue, lockTimeout).Result()
	if err != nil {
		panic(err)
	}
	if success {
		fmt.Println("获取锁成功")
		// 执行业务逻辑

		//获取商家token以实现商品更改
		h.Token, err = model.RDB.Get(c, "hostertoken").Result()
		if err != nil {
			log.Fatalln("get hostertoken error:", err)
		}
		//一个用户只能买一个商品
		s, err = model.RDB.Get(c, h.Token+"keep").Result()
		if err != nil {
			log.Fatalln("get keep error:", err)
		}
		h.Goods.Keep, _ = strconv.Atoi(s)
		if h.Goods.Keep > 0 {
			h.Goods.Keep--
		}
		model.RDB.Set(c, h.Token+"keep", h.Goods.Keep, 0)

		//订单存在30秒
		model.RDB.Set(c, u.Token+"book", true, 30*time.Second)

		// 释放锁
		res, err := model.RDB.Get(c, lockKey).Result()
		if err == nil && res == lockValue {
			model.RDB.Del(c, "my_lock")
		} else {
			log.Fatalln("释放锁错误", err)
		}

	} else {
		log.Fatalln("获取锁失败")
	}
	resp.OKWithData(c, "book success")
}

// 用户付款
func IsBuy(c *gin.Context) {
	var u model.User
	var h model.Hoster
	//token验证身份
	u.Token = c.GetHeader("usertoken")
	_, err := model.RDB.Get(c, u.Token+"book").Result()
	if err != nil {
		//订单已过期
		//获取商家token以实现返还商品库存
		h.Token, err = model.RDB.Get(c, "hostertoken").Result()
		if err != nil {
			log.Fatalln("get hostertoken error:", err)
		}
		//利用redis的incr来增keep避免并发问题
		model.RDB.Incr(c, h.Token+"keep")
		c.JSON(200, gin.H{
			"status": "20001",
			"info":   "kill fail",
		})
	} else {
		//订单未过期，直接付款
		resp.OKWithData(c, "kill success")
	}
}
