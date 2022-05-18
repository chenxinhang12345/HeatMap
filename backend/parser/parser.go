package parser

import (
	"encoding/csv"
	"fmt"
	"log"
	nc "main/nanocube"
	"math"
	"os"
	"strconv"
)

//ReadCsvFile return records for a csv file with given filename
func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

//ParseObjects parse a csv files into objects with Lat Lon time and type
func ParseObjects(filename string, typeHead string) []nc.Object {
	LngIndex := -1 //column index for Lng
	LatIndex := -1 //column index for Lat
	TypeIndex := -1
	records := ReadCsvFile(filename)
	res := []nc.Object{}
	for i := 0; i < len(records[0]); i++ {
		if records[0][i] == "Longitude" {
			LngIndex = i
		} else if records[0][i] == "Latitude" {
			LatIndex = i
		} else if records[0][i] == typeHead {
			TypeIndex = i
		}
	}
	if LngIndex == -1 || LatIndex == -1 || TypeIndex == -1 {
		return res
	}
	for i := 1; i < len(records); i++ {
		lngstr := records[i][LngIndex]
		if lngstr == "" {
			continue
		}
		latstr := records[i][LatIndex]
		if latstr == "" {
			continue
		}
		lng, err := strconv.ParseFloat(records[i][LngIndex], 64)
		if err != nil {
			log.Fatal("Longitude is not a valid float", err)
		}
		lat, err1 := strconv.ParseFloat(records[i][LatIndex], 64)
		if err1 != nil {
			log.Fatal("Latitude is not a valid float", err)
		}
		ty := records[i][TypeIndex]
		res = append(res, nc.Object{Lng: lng, Lat: lat, Type: ty})
	}
	return res
}

//CreateNanoCubeFromCsvFile return a nanocube pointer
func CreateNanoCubeFromCsvFile(filePath string, typeHead string, maxDepth int, limit int) *nc.Nanocube {
	objects := ParseObjects(filePath, typeHead)
	minLng := 1e9
	maxLng := -1e9
	minLat := 1e9
	maxLat := -1e9
	if limit == -1 {
		limit = len(objects)
	}
	types := make(map[string]bool)
	for i := 0; i < limit; i++ {
		o := objects[i]
		lng := o.Lng
		lat := o.Lat
		types[o.Type] = true
		minLng = math.Min(minLng, lng)
		maxLng = math.Max(maxLng, lng)
		minLat = math.Min(minLat, lat)
		maxLat = math.Max(maxLat, lat)
	}
	//top left corner minlng maxlat
	fmt.Println("Top left corner: ", minLng, maxLat)
	width := math.Abs(minLng - maxLng)
	height := math.Abs(minLat - maxLat)
	b := math.Max(width, height)
	fmt.Println("boarder length: ", b)
	typesArray := make([]string, 0)
	for k := range types {
		typesArray = append(typesArray, k)
	}
	res := nc.SetUpCube(maxDepth, nc.Bounds{Lng: minLng, Lat: maxLat, Width: b, Height: b}, typesArray)
	for i := 0; i < limit; i++ {
		// //d
		// fmt.Println(i, "th", "add")
		// //d
		res.AddObject(objects[i])
		// //debug
		// boxes := nc.QueryAll(res.Root, 2)
		// fmt.Println(boxes)
		// //debug
	}
	fmt.Println("Total number of objects: ", limit)
	return res
}
