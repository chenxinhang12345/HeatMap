package main

import (
	"main/server/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	controllers.InitNanoCube()
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/cubes", controllers.QueryAll)
	r.GET("/types", controllers.QueryTypes)
	r.GET("/cubes/time", controllers.QueryWithTypeAndTime)
	r.Run()
}
