package controllers

import "github.com/gin-gonic/gin"

func Query(c *gin.Context) {
	bLat := c.Query("lat")
	bLng := c.Query("lng")
	return
}
