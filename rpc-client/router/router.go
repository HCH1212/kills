package router

import (
	"log"
	"rpc-client/api"
	"rpc-client/service"

	"github.com/gin-gonic/gin"
)

func InitRouters() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	u := r.Group("/user")
	h := r.Group("/hoster")

	{
		u.POST("/register", api.UserRegister)
		u.POST("/login", api.UserLogin)
		h.POST("/register", api.HosterRegister)
		h.POST("/login", api.HosterLogin)
	}

	r.POST("/creatkill", service.CreatKill)
	r.GET("/showkill", service.ShowKill)
	r.POST("/joinkill", service.JoinKill)
	r.POST("/joinkill/isbuy", service.IsBuy)

	if err := r.Run(":8080"); err != nil {
		log.Fatalln("run server error :", err)
	}
}
