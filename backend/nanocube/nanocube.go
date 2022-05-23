package nanocube

import (
	"encoding/json"
	"sort"
)

/*
Nanocube ...
*/
type Nanocube struct {
	Root     *SpatNode
	MaxLevel int //maximum level allowed for spatial attribute
	Types    []string
	Index    map[string]int //the map stores categorical index
}

//SpatNode for encoding spatial attribute
type SpatNode struct {
	Bounds   Bounds
	Children []*SpatNode
	CatRoot  *CatNode
	Level    int //current level
}

//CatNode for encoding categorical attribute
type CatNode struct {
	Children []*Summary
	Summary  *Summary
	// Type     string //the category
}

//HeatMapGrid encode information for each grid of heatmap
type HeatMapGrid struct {
	B     Bounds `json:"bounds"`
	Count int64  `json:"count"`
}

/*Bounds encode spatial information for each node
|---------Lng ++
|  0    1
|
|  2    3
Lat --
*/
type Bounds struct {
	Lng    float64 `json:"lng"`
	Lat    float64 `json:"lat"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

//Object represent event
type Object struct {
	Lng       float64
	Lat       float64
	Type      string
	TimeStamp int64
}

//Summary a summary for a bunch of nodes
type Summary struct {
	Count           int64
	TimeStampCounts []TemporalCount
}

//TemporalCount tuple in timestampcounts, provide seconds and number of spatial points at this timestamp (culmulative)
type TemporalCount struct {
	TimeStamp int64
	Count     int64
}

func PrintTimestampCounts(data []TemporalCount) {
	for _, tc := range data {
		print("(", tc.TimeStamp, ",", tc.Count, ") ")
	}
}

//InsertAt insert the temporal count at a specific index in temporal count array
func InsertAt(data []TemporalCount, i int, v TemporalCount) []TemporalCount {
	if i == len(data) {
		return append(data, v)
	}
	data = append(data[:i+1], data[i:]...)
	data[i] = v
	return data
}

//AddTemporalCount add a temporal count in temporal count array make the array sorted and with culmulative count
func AddTemporalCount(tc TemporalCount, data []TemporalCount) []TemporalCount {
	f := func(i int) bool {
		return data[i].TimeStamp >= tc.TimeStamp
	}
	i := sort.Search(len(data), f)
	if i < len(data) && data[i].TimeStamp == tc.TimeStamp {
		data[i].Count += tc.Count
	} else {
		data = InsertAt(data, i, tc)
		if i != 0 {
			data[i].Count += data[i-1].Count
		}
	}
	for j := i + 1; j < len(data); j++ {
		data[j].Count += 1
	}
	return data
}

func TemporalCountRangeQuery(data []TemporalCount, startTime int64, endTime int64) int64 {
	if len(data) == 0 {
		return 0
	}
	fstart := func(i int) bool {
		return data[i].TimeStamp >= startTime
	}
	fend := func(i int) bool {
		return data[i].TimeStamp >= endTime
	}
	startIndex := sort.Search(len(data), fstart)
	endIndex := sort.Search(len(data), fend)
	if endIndex == len(data) {
		endIndex -= 1
	}
	if startIndex == len(data) {
		startIndex -= 1
	}
	// println(startIndex, endIndex)
	// println(endTime, data[endIndex].TimeStamp)
	if endTime < data[endIndex].TimeStamp {
		endIndex -= 1
	}
	if startTime > data[startIndex].TimeStamp {
		startIndex += 1
	}
	// println(startIndex, endIndex)
	if startIndex > endIndex {
		return 0
	}
	if startIndex == endIndex {
		if startIndex > 0 {
			return data[endIndex].Count - data[endIndex-1].Count
		} else {
			return data[startIndex].Count
		}
	}
	var res int64
	if startIndex > 0 {
		res = data[endIndex].Count - data[startIndex-1].Count
	} else {
		res = data[endIndex].Count
	}
	// if res != int64(len(data)) {
	// println("s, e", startIndex, endIndex)
	// println(data[endIndex].Count, data[startIndex].Count)
	// println("res", res, "len", len(data))
	// PrintTimestampCounts(data)
	// println("")
	// }
	return res
}

//SetUpCube Initialize the cube
func SetUpCube(MaxLevel int, MaxBounds Bounds, Types []string) *Nanocube {
	nc := &Nanocube{Root: &SpatNode{Bounds: MaxBounds, Children: make([]*SpatNode, 4), Level: 1}, MaxLevel: MaxLevel, Types: Types}
	m := make(map[string]int)
	for i := 0; i < len(Types); i++ {
		m[Types[i]] = i
	}
	nc.Index = m
	return nc
}

//AssignIndexOnBounds helper function for assigning index on specific bounds for an object
func AssignIndexOnBounds(obj Object, b Bounds) (int, Bounds) {
	HalfWidth := b.Width / 2
	HalfHeight := b.Height / 2
	MidLng := b.Lng + HalfWidth
	MidLat := b.Lat - HalfHeight
	// fmt.Println("func AssignIndexBounds ", obj, " ", b, "MidLng ", MidLng, "MidLat ", MidLat)
	if obj.Lng <= MidLng && obj.Lat >= MidLat {
		return 0, Bounds{b.Lng, b.Lat, HalfWidth, HalfHeight}
	} else if obj.Lng > MidLng && obj.Lat >= MidLat {
		return 1, Bounds{MidLng, b.Lat, HalfWidth, HalfHeight}
	} else if obj.Lng <= MidLng && obj.Lat < MidLat {
		return 2, Bounds{b.Lng, MidLat, HalfWidth, HalfHeight}
	} else if obj.Lng > MidLng && obj.Lat < MidLat {
		return 3, Bounds{MidLng, MidLat, HalfWidth, HalfHeight}
	} else {
		return -1, Bounds{}
	}
}

func (nc *Nanocube) getIndex(t string) int {
	return nc.Index[t]
}

//HasOnlyOneChild check if the SpatNode has only one child
func (s *SpatNode) HasOnlyOneChild() (bool, *SpatNode) {
	counter := 0
	var retptr *SpatNode = nil
	for i := 0; i < 4; i++ {
		if s.Children[i] != nil {
			retptr = s.Children[i]
			counter++
		}
	}
	return (counter == 1), retptr
}

func Copy(tc TemporalCount) TemporalCount {
	return TemporalCount{TimeStamp: tc.TimeStamp, Count: tc.Count}
}

//Copy return a deep copy of a summary
func (s *Summary) Copy() *Summary {
	ts := []TemporalCount{}
	for _, t := range s.TimeStampCounts {
		ts = append(ts, Copy(t))
	}
	return &Summary{Count: s.Count, TimeStampCounts: ts}
}

//HasOnlyOneChild check if the cat node has only one child
func (c *CatNode) HasOnlyOneChild() bool {
	counter := 0
	for i := 0; i < len(c.Children); i++ {
		if c.Children[i] != nil {
			counter++
		}
		if counter > 1 { //more than one child
			return false
		}
	}
	return true
}

//UpdateSummary update current summary when adding an object to current SpatNode
func (s *SpatNode) UpdateSummary(obj Object, maxLevel int, nc *Nanocube) {
	hasOnlyOneChild, child := s.HasOnlyOneChild()
	if s.Level < maxLevel {
		if s.CatRoot == nil { //if it doesn't have categorical root
			s.CatRoot = child.CatRoot
		} else { //if it has
			if hasOnlyOneChild { //only one child
				s.CatRoot = child.CatRoot
			} else { //need update
				index := nc.getIndex(obj.Type) //update categorical node
				cpy := make([]*Summary, len(s.CatRoot.Children))
				copy(cpy, s.CatRoot.Children)
				// count := s.CatRoot.Summary.Count
				// temporalCounts := s.CatRoot.Summary.TimeStampCounts
				// fmt.Println("original count", count)
				s.CatRoot = &CatNode{Summary: s.CatRoot.Summary.Copy(), Children: cpy} //summary is new, children is old
				if s.CatRoot.Children[index] == nil {
					s.CatRoot.Children[index] = &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}} //init summary
				} else {
					s.CatRoot.Children[index] = s.CatRoot.Children[index].Copy() //only make a new copy on this cat children and then update
					s.CatRoot.Children[index].Count++
					s.CatRoot.Children[index].TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Children[index].TimeStampCounts)
				}
				s.CatRoot.Summary.Count++
				s.CatRoot.Summary.TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Summary.TimeStampCounts)
			}
		}
	} else {
		if s.CatRoot == nil {
			s.CatRoot = &CatNode{Summary: &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}, Children: make([]*Summary, len(nc.Types))}
			index := nc.getIndex(obj.Type)
			s.CatRoot.Children[index] = s.CatRoot.Summary
		} else { //need update
			index := nc.getIndex(obj.Type)
			if s.CatRoot.Children[index] != nil {
				s.CatRoot.Children[index] = s.CatRoot.Children[index].Copy()
				s.CatRoot.Children[index].Count++
				s.CatRoot.Children[index].TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Children[index].TimeStampCounts)
			} else {
				s.CatRoot.Children[index] = &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}
			}
			s.CatRoot.Summary = s.CatRoot.Summary.Copy()
			s.CatRoot.Summary.Count++
			s.CatRoot.Summary.TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Summary.TimeStampCounts)
		}
	}
}

func (c *CatNode) Copy() *CatNode {
	length := len(c.Children)
	childrenCopy := make([]*Summary, length)
	for i := 0; i < len(c.Children); i++ {
		if c.Children[i] != nil {
			childrenCopy[i] = c.Children[i].Copy()
		} else {
			childrenCopy[i] = nil
		}
	}
	return &CatNode{Summary: c.Summary.Copy(), Children: childrenCopy}
}

func (s *SpatNode) UpdateSummaryWithoutSharing(obj Object, maxLevel int, nc *Nanocube) {
	// _, child := s.HasOnlyOneChild()
	if s.Level < maxLevel {
		if s.CatRoot == nil {
			s.CatRoot = &CatNode{Summary: &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}, Children: make([]*Summary, len(nc.Types))}
			index := nc.getIndex(obj.Type) //update categorical node
			if s.CatRoot.Children[index] != nil {
				s.CatRoot.Children[index].Count++
				s.CatRoot.Children[index].TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Children[index].TimeStampCounts)
			} else {
				s.CatRoot.Children[index] = &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}
				s.CatRoot.Summary.TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Summary.TimeStampCounts)
			}
		} else {
			index := nc.getIndex(obj.Type) //update categorical node
			if s.CatRoot.Children[index] != nil {
				s.CatRoot.Children[index].Count++
				s.CatRoot.Children[index].TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Children[index].TimeStampCounts)
			} else {
				s.CatRoot.Children[index] = &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}
			}
			s.CatRoot.Summary.Count++
			s.CatRoot.Summary.TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Summary.TimeStampCounts)
		}
	} else {
		if s.CatRoot == nil {
			s.CatRoot = &CatNode{Summary: &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}, Children: make([]*Summary, len(nc.Types))}
			index := nc.getIndex(obj.Type)
			s.CatRoot.Children[index] = &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}
		} else { //need update
			index := nc.getIndex(obj.Type)
			if s.CatRoot.Children[index] != nil {
				s.CatRoot.Children[index].Count++
				s.CatRoot.Children[index].TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Children[index].TimeStampCounts)
			} else {
				s.CatRoot.Children[index] = &Summary{Count: 1, TimeStampCounts: []TemporalCount{{TimeStamp: obj.TimeStamp, Count: 1}}}
			}
			s.CatRoot.Summary.Count++
			s.CatRoot.Summary.TimeStampCounts = AddTemporalCount(TemporalCount{TimeStamp: obj.TimeStamp, Count: 1}, s.CatRoot.Summary.TimeStampCounts)
		}
	}
}

//AddObject Add an object
func (nc *Nanocube) AddObject(obj Object, isWithSharing bool) {
	stack := make([]*SpatNode, 0)
	levels := nc.MaxLevel
	currentNode := nc.Root
	currentLevel := 1
	for currentLevel < levels {
		// fmt.Println("currentLevel: ", currentLevel)
		index, b := AssignIndexOnBounds(obj, currentNode.Bounds)
		// fmt.Println(index, b)
		// fmt.Println("Assignindex: ", index)
		if currentNode.Children[index] == nil { //no nodes on current index
			currentNode.Children[index] = &SpatNode{Bounds: b, Children: make([]*SpatNode, 4)} //create a new node on current index
		}
		currentNode.Level = currentLevel
		stack = append(stack, currentNode)
		currentNode = currentNode.Children[index] //next level node
		currentLevel++
	}
	currentNode.Level = currentLevel
	// fmt.Println("leave level:", currentLevel)
	stack = append(stack, currentNode) //update leaves
	for i := len(stack) - 1; i >= 0; i-- {
		currentNode = stack[i]
		if isWithSharing {
			currentNode.UpdateSummary(obj, levels, nc)
		} else {
			currentNode.UpdateSummaryWithoutSharing(obj, levels, nc)
		}
	}
}

//Intersect Decide whether these two bounds are interseted or not
func (b1 *Bounds) Intersect(b2 Bounds) bool {
	return !((b1.Lat-b1.Height >= b2.Lat) || (b2.Lat-b2.Height >= b1.Lat) || (b1.Lng+b1.Width <= b2.Lng) || (b2.Lng+b2.Width <= b1.Lng))
}

//Equal Decide whether these two bounds are equal
func (b1 *Bounds) Equal(b2 Bounds) bool {
	return b1.Lat == b2.Lat && b1.Lng == b2.Lng && b2.Width == b1.Width && b1.Height == b2.Height
}

//QueryAll return all the grids within current node at specifying level
func QueryAll(s *SpatNode, level int) []HeatMapGrid {
	if s == nil {
		return []HeatMapGrid{}
	}
	s1 := s.Children[0]
	s2 := s.Children[1]
	s3 := s.Children[2]
	s4 := s.Children[3]
	if s.Level < level {
		c1 := QueryAll(s1, level)
		c2 := QueryAll(s2, level)
		c3 := QueryAll(s3, level)
		c4 := QueryAll(s4, level)
		res := append(c1, c2...)
		res = append(res, c3...)
		res = append(res, c4...)
		return res
	}
	return []HeatMapGrid{{s.Bounds, s.CatRoot.Summary.Count}}

}

//Query basic function for query a heatmap without specifying type
func Query(s *SpatNode, b Bounds, level int) []HeatMapGrid {
	s1 := s.Children[0]
	s2 := s.Children[1]
	s3 := s.Children[2]
	s4 := s.Children[3]
	c1 := []HeatMapGrid{}
	c2 := []HeatMapGrid{}
	c3 := []HeatMapGrid{}
	c4 := []HeatMapGrid{}
	if s.Level < level {
		if s1 != nil {
			b1 := s1.Bounds
			if b1.Intersect(b) { //Intersect
				c1 = Query(s1, b, level)
			}
		}
		if s2 != nil {
			b1 := s2.Bounds
			if b1.Intersect(b) { //Intersect
				c2 = Query(s2, b, level)
			}
		}
		if s3 != nil {
			b1 := s3.Bounds
			if b1.Intersect(b) { //Intersect
				c3 = Query(s3, b, level)
			}
		}
		if s4 != nil {
			b1 := s4.Bounds
			if b1.Intersect(b) { //Intersect
				c4 = Query(s4, b, level)
			}
		}
		res := append(c1, c2...)
		res = append(res, c3...)
		res = append(res, c4...)
		return res
	}
	return []HeatMapGrid{{s.Bounds, s.CatRoot.Summary.Count}}
}

func JsonQuery(s *SpatNode, b Bounds, level int) string {
	var grids []HeatMapGrid = Query(s, b, level)
	binStr, err := json.Marshal(grids)
	if err != nil {
		return "error with deserialize"
	}
	return string(binStr)
}

func QueryType(typeIndex int, s *SpatNode, b Bounds, level int) []HeatMapGrid {
	s1 := s.Children[0]
	s2 := s.Children[1]
	s3 := s.Children[2]
	s4 := s.Children[3]
	c1 := []HeatMapGrid{}
	c2 := []HeatMapGrid{}
	c3 := []HeatMapGrid{}
	c4 := []HeatMapGrid{}
	if s.Level < level {
		if s1 != nil {
			b1 := s1.Bounds
			if b1.Intersect(b) { //Intersect
				c1 = QueryType(typeIndex, s1, b, level)
			}
		}
		if s2 != nil {
			b1 := s2.Bounds
			if b1.Intersect(b) { //Intersect
				c2 = QueryType(typeIndex, s2, b, level)
			}
		}
		if s3 != nil {
			b1 := s3.Bounds
			if b1.Intersect(b) { //Intersect
				c3 = QueryType(typeIndex, s3, b, level)
			}
		}
		if s4 != nil {
			b1 := s4.Bounds
			if b1.Intersect(b) { //Intersect
				c4 = QueryType(typeIndex, s4, b, level)
			}
		}
		res := append(c1, c2...)
		res = append(res, c3...)
		res = append(res, c4...)
		return res
	}
	if s.CatRoot.Children[typeIndex] == nil {
		return []HeatMapGrid{}
	}
	return []HeatMapGrid{{s.Bounds, s.CatRoot.Children[typeIndex].Count}}
}

func QueryTypeTime(startTime int64, endTime int64, typeIndex int, s *SpatNode, b Bounds, level int) []HeatMapGrid {
	s1 := s.Children[0]
	s2 := s.Children[1]
	s3 := s.Children[2]
	s4 := s.Children[3]
	c1 := []HeatMapGrid{}
	c2 := []HeatMapGrid{}
	c3 := []HeatMapGrid{}
	c4 := []HeatMapGrid{}
	if s.Level < level {
		if s1 != nil {
			b1 := s1.Bounds
			if b1.Intersect(b) { //Intersect
				c1 = QueryTypeTime(startTime, endTime, typeIndex, s1, b, level)
			}
		}
		if s2 != nil {
			b1 := s2.Bounds
			if b1.Intersect(b) { //Intersect
				c2 = QueryTypeTime(startTime, endTime, typeIndex, s2, b, level)
			}
		}
		if s3 != nil {
			b1 := s3.Bounds
			if b1.Intersect(b) { //Intersect
				c3 = QueryTypeTime(startTime, endTime, typeIndex, s3, b, level)
			}
		}
		if s4 != nil {
			b1 := s4.Bounds
			if b1.Intersect(b) { //Intersect
				c4 = QueryTypeTime(startTime, endTime, typeIndex, s4, b, level)
			}
		}
		res := append(c1, c2...)
		res = append(res, c3...)
		res = append(res, c4...)
		return res
	}

	var resCount int64
	if typeIndex != -1 {
		if s.CatRoot.Children[typeIndex] == nil {
			return []HeatMapGrid{}
		}
		resCount = TemporalCountRangeQuery(s.CatRoot.Children[typeIndex].TimeStampCounts, startTime, endTime) //query on specific category
	} else {
		resCount = TemporalCountRangeQuery(s.CatRoot.Summary.TimeStampCounts, startTime, endTime) //query all regardless categories
	}
	if resCount == 0 {
		return []HeatMapGrid{} //Do not need to render
	}
	return []HeatMapGrid{{s.Bounds, resCount}}
}

func JsonQueryType(typeIndex int, s *SpatNode, b Bounds, level int) string {
	var grids []HeatMapGrid = QueryType(typeIndex, s, b, level)
	binStr, err := json.Marshal(grids)
	if err != nil {
		return "error with deserialize"
	}
	return string(binStr)
}
