package main

import (
	"main/server/controllers"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 5 {
		println("usage ./main <filepath> <category column head> <date time column head> <max level for nanocubes> <max number of points in nanocube> ")
	}
	filePath := argsWithoutProg[0]
	catColumn := argsWithoutProg[1]
	timeColumn := argsWithoutProg[2]
	level, err := strconv.ParseInt(argsWithoutProg[3], 10, 64)
	if err != nil {
		println("level should be readable int")
		return
	}
	numPoints, err := strconv.ParseInt(argsWithoutProg[4], 10, 64)
	if err != nil {
		println("number of points should be readable int")
	}
	controllers.InitNanoCube(filePath, catColumn, timeColumn, int(level), int(numPoints))
	debug.FreeOSMemory() //grabage collection for non referenced objects
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/cubes", controllers.QueryAll)
	r.GET("/types", controllers.QueryTypes)
	r.GET("/cubes/time", controllers.QueryWithTypeAndTime)
	r.Run()
}
