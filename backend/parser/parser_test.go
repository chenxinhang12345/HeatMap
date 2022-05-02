package parser

import (
	"fmt"
	"runtime"
	"testing"

	nc "github.com/chenxinhang12345/backend/nanocube"
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

// func TestParseObjects(t *testing.T) {
// 	fmt.Println(ParseObjects("test.csv", "Primary Type"))
// }

// func TestParseBig(t *testing.T) {
// 	fmt.Println(ParseObjects("crime2019.csv", "Primary Type"))
// }

// func TestNanoCubeFromSmallFile(t *testing.T) {
// 	fmt.Println(CreateNanoCubeFromCsvFile("test.csv", "Primary Type", 16, -1))
// }

func TestNanoCubeFromBigFile(t *testing.T) {
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "Primary Type", 20, 101000)
	PrintMemUsage()
	// fmt.Println(n.Root.Children[0])
	// fmt.Println(n.Root.Children[1].CatRoot.Summary)
	// fmt.Println(n.Root.Children[2].CatRoot.Summary)
	// fmt.Println(n.Root.Children[3].CatRoot.Summary)
	// boxes := nc.Query(n.Root, nc.Bounds{-87.9345, 42.022585817, 2, 2}, 4)
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
}
