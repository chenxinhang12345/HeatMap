package parser

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	nc "main/nanocube"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// func TestReadCSV(t *testing.T) {
// 	fmt.Println(ReadCsvFile("test.csv"))
// }

func TestParseObjects(t *testing.T) {
	fmt.Println(ParseObjects("crime2020.csv", "PrimaryType", "Date")[0])
}

func TestCulSearch(t *testing.T) {
	tc := nc.TemporalCount{TimeStamp: 8, Count: 1}
	// data := []nc.TemporalCount{{TimeStamp: 12, Count: 1}, {TimeStamp: 19, Count: 2}}
	data := []nc.TemporalCount{{9, 1}, {10, 2}}
	data = nc.AddTemporalCount(tc, data)
	for _, t := range data {
		println(t.TimeStamp, t.Count, " ")
	}
	a := data
	a[0].Count++
	for _, t := range data {
		println(t.TimeStamp, t.Count, " ")
	}
	// expected := []nc.TemporalCount{{8, 1}, {9, 2}, {10, 3}}
}

func TestAddTemporalCount(t *testing.T) {
	// data := []nc.TemporalCount{{TimeStamp: 12, Count: 1}, {TimeStamp: 19, Count: 2}}
	adds := []nc.TemporalCount{{9, 1}, {10, 1}, {8, 1}, {6, 1}, {0, 1}, {4, 1}}
	data := []nc.TemporalCount{}
	for _, add := range adds {
		data = nc.AddTemporalCount(add, data)
	}
	nc.PrintTimestampCounts(data)
}

func TestTimeStampRangeQuery(t *testing.T) {
	adds := []nc.TemporalCount{{9, 1}, {10, 1}, {8, 1}, {6, 1}, {3, 1}, {4, 1}}
	data := []nc.TemporalCount{}
	for _, add := range adds {
		data = nc.AddTemporalCount(add, data)
	}
	nc.PrintTimestampCounts(data)
	res := nc.TemporalCountRangeQuery(data, 0, 0)
	println(res)
}

// func TestParseBig(t *testing.T) {
// 	fmt.Println(ParseObjects("crime2019.csv", "Primary Type"))
// }

// func TestNanoCubeFromSmallFile(t *testing.T) {
// 	fmt.Println(CreateNanoCubeFromCsvFile("test.csv", "Primary Type", 16, -1))
// }

func TestNanoCubeFromBigFile(t *testing.T) {
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "PrimaryType", "Date", 20, 10000, true)
	PrintMemUsage()
	fmt.Println("all types:", n.Index)
	num_cats := len(n.Index)
	fmt.Println("we have total of ", len(n.Index), "categories")
	boxes := nc.QueryAll(n.Root, 20)
	// fmt.Println(boxes)
	sum := 0
	for _, box := range boxes {
		sum += int(box.Count)
	}
	fmt.Println("Total sum is ", sum, "for query all")
	sum_query := 0
	boxes = nc.Query(n.Root, nc.Bounds{Lng: -87.9345, Lat: 42.022585817, Width: 0.424, Height: 0.424}, 16)
	for _, box := range boxes {
		sum_query += int(box.Count)
	}
	fmt.Println("Total sum is ", sum_query, "for query")
	if sum_query != sum {
		t.Errorf("these two sum should be equal")
	}
	sum_all_types := 0
	for index := 0; index < num_cats; index++ {
		boxes = nc.QueryType(index, n.Root, nc.Bounds{Lng: -87.9345, Lat: 42.022585817, Width: 0.424, Height: 0.424}, 18)
		sum_type := 0
		for _, box := range boxes {
			sum_type += int(box.Count)
		}
		fmt.Println("Total sum type for", index, sum_type)
		sum_all_types += sum_type
	}
	fmt.Println("all types count:", sum_all_types)
	if sum_all_types != sum {
		t.Errorf("these two sum should be equal")
	}
	sum_all_types = 0
	for index := 0; index < num_cats; index++ {
		boxes = nc.QueryTypeTime(0, 1653016292, index, n.Root, nc.Bounds{Lng: -87.9345, Lat: 42.022585817, Width: 0.424, Height: 0.424}, 18)
		sum_type := 0
		for _, box := range boxes {
			sum_type += int(box.Count)
		}
		fmt.Println("Total sum type time for", index, sum_type)
		sum_all_types += sum_type
	}
	fmt.Println("all types count:", sum_all_types)
	if sum_all_types != sum {
		t.Errorf("these two sum should be equal")
	}

	sum_all_types = 0
	for index := 0; index < num_cats; index++ {
		boxes = nc.QueryTypeTime(1593325815, 1653016292, index, n.Root, nc.Bounds{Lng: -87.9345, Lat: 42.022585817, Width: 0.424, Height: 0.424}, 18)
		sum_type := 0
		for _, box := range boxes {
			sum_type += int(box.Count)
		}
		fmt.Println("Total sum type time for", index, sum_type)
		sum_all_types += sum_type
	}
	fmt.Println("all types count during a Interval:", sum_all_types)

	sum_all_types = 0
	for index := 0; index < num_cats; index++ {
		boxes = nc.QueryTypeTime(0, 1593325814, index, n.Root, nc.Bounds{Lng: -87.9345, Lat: 42.022585817, Width: 0.424, Height: 0.424}, 18)
		sum_type := 0
		for _, box := range boxes {
			sum_type += int(box.Count)
		}
		fmt.Println("Total sum type time for", index, sum_type)
		sum_all_types += sum_type
	}
	fmt.Println("all types count during another Interval:", sum_all_types)

}

//Test Nanocube memory usage
func TestMemUsageAndConsturctionTimeWithSharing(t *testing.T) {
	start := time.Now()
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "PrimaryType", "Date", 10, 10000, true)
	duration := time.Since(start)
	debug.FreeOSMemory()
	PrintMemUsage()
	println(n.Root)
	fmt.Println("time cost:", duration)
}

//Test standard quadtree memory usage
func TestMemUsageAndConstructionTimeWithoutSharing(t *testing.T) {
	start := time.Now()
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "PrimaryType", "Date", 10, 10000, false)
	duration := time.Since(start)
	debug.FreeOSMemory()
	PrintMemUsage()
	println(n.Root)
	fmt.Println("time cost:", duration)
}

//Test query time for Nanocubes
func TestQueryTimeWithSharing(t *testing.T) {
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "PrimaryType", "Date", 10, 20000, true)
	debug.FreeOSMemory()
	PrintMemUsage()
	num_cats := len(n.Index)
	var total_time float64 = 0
	for index := 0; index < num_cats; index++ {
		start := time.Now()
		_ = nc.QueryTypeTime(1593325815, 1653016292, index, n.Root, nc.Bounds{Lng: -87.9345, Lat: 42.022585817, Width: 0.424, Height: 0.424}, 25)
		duration := time.Since(start)
		total_time += duration.Seconds()
	}
	println(n.Root)
	fmt.Println("average query time cost:", total_time/float64(num_cats))
}

//Test query time for standard quadtree memory usage
func TestQueryTimeWithoutSharing(t *testing.T) {
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "PrimaryType", "Date", 10, 20000, false)
	debug.FreeOSMemory()
	PrintMemUsage()
	num_cats := len(n.Index)
	var total_time float64 = 0
	for index := 0; index < num_cats; index++ {
		start := time.Now()
		_ = nc.QueryTypeTime(1593325815, 1653016292, index, n.Root, nc.Bounds{Lng: -87.9345, Lat: 42.022585817, Width: 0.424, Height: 0.424}, 25)
		duration := time.Since(start)
		total_time += duration.Seconds()
	}
	println(n.Root)
	fmt.Println("average query time cost:", total_time/float64(num_cats))
}
