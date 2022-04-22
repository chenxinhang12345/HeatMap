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
	n := CreateNanoCubeFromCsvFile("crime2020.csv", "Primary Type", 16, 1000)
	PrintMemUsage()
	fmt.Println(n.Root.Children[0])
	fmt.Println(n.Root.Children[1].CatRoot.Summary)
	fmt.Println(n.Root.Children[2].CatRoot.Summary)
	fmt.Println(n.Root.Children[3].CatRoot.Summary)
	// boxes := nc.Query(n.Root, nc.Bounds{-87.9345, 42.022585817, 2, 2}, 4)
	boxes := nc.QueryAll(n.Root, 14)
	sum := 0
	for _, box := range boxes {
		sum += int(box.Count)
	}
	fmt.Println("Total sum is ", sum)

}
