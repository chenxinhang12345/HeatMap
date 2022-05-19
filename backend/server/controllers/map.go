package controllers

import (
	"math"
	"net/http"
	"strconv"

	nanocube "main/nanocube"
	parser "main/parser"
	"main/server/models"
	"main/utils"

	"github.com/gin-gonic/gin"
)

var Nanocube *nanocube.Nanocube

func InitNanoCube() {
	filePath := "./parser/crime2020.csv"
	catColumn := "Primary Type"
	level := 20
	numPoints := 100000
	Nanocube = parser.CreateNanoCubeFromCsvFile(filePath, catColumn, level, numPoints, true)
}

//computeOpacity helper function for computing the opacity value. Algorithnm is from paper
//https://idl.cs.washington.edu/files/2013-imMens-EuroVis.pdf
func computeOpacity(count int64, max int64, min int64, alpha float64, gamma float64) float64 {
	X := float64(count)
	maxX := float64(max)
	minX := float64(min)
	Y := alpha + math.Pow(((X-minX)/(maxX-minX)), gamma)*(1-alpha)
	return Y
}

func convertGridsToRectangles(grids []nanocube.HeatMapGrid) []models.Rectangle {
	res := []models.Rectangle{}
	var maxCount int64 = 0
	var minCount int64 = 1e9
	for _, grid := range grids {
		b := grid.B
		count := grid.Count
		maxCount = utils.Max(maxCount, count)
		minCount = utils.Min(minCount, count)
		res = append(res, models.Rectangle{N: b.Lat, S: b.Lat - b.Height, E: b.Lng + b.Width, W: b.Lng, Count: count})
	}
	for i := 0; i < len(res); i++ {
		res[i].Opacity = computeOpacity(res[i].Count, maxCount, minCount, 0.15, 0.33333) //use papar's parameters
	}
	return res
}

func QueryAll(c *gin.Context) {
	minLat := c.Query("minLat")
	maxLat := c.Query("maxLat")
	minLng := c.Query("minLng")
	maxLng := c.Query("maxLng")
	zoomStr := c.Query("zoom")
	typeStr := c.Query("type")
	var grids []nanocube.HeatMapGrid
	minlat, err := strconv.ParseFloat(minLat, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": grids})
	}
	maxlat, err := strconv.ParseFloat(maxLat, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": grids})
	}
	minlng, err := strconv.ParseFloat(minLng, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": grids})
	}
	maxlng, err := strconv.ParseFloat(maxLng, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": grids})
	}
	zoom, err := strconv.Atoi(zoomStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": grids})
	}
	typeNum, err := strconv.Atoi(typeStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": grids})
	}
	println(minlat, maxlat, minlng, maxlng, zoom, typeNum)
	lat := maxlat
	lng := minlng
	width := maxlng - minlng
	height := maxlat - minlat

	if typeNum < 0 {
		grids = nanocube.Query(Nanocube.Root, nanocube.Bounds{Lng: lng, Lat: lat, Width: width, Height: height}, zoom-5)
	} else {
		grids = nanocube.QueryType(typeNum, Nanocube.Root, nanocube.Bounds{Lng: lng, Lat: lat, Width: width, Height: height}, zoom-5)
	}
	var rects = convertGridsToRectangles(grids)
	c.JSON(http.StatusOK, gin.H{"data": rects})
}

func QueryTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": Nanocube.Index})
}
