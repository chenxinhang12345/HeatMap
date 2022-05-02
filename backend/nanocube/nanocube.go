package nanocube

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
	B     Bounds
	Count int64
}

/*Bounds encode spatial information for each node
|---------Lng ++
|  0    1
|
|  2    3
Lat --
*/
type Bounds struct {
	Lng    float64
	Lat    float64
	Width  float64
	Height float64
}

//Object represent event
type Object struct {
	Lng       float64
	Lat       float64
	Type      string
	TimeStamp int
}

//Summary a summary for a bunch of nodes
type Summary struct {
	Count             int64
	TimeStampedCounts []int64
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

//Copy return a deep copy of a summary
func (s *Summary) Copy() *Summary {
	return &Summary{Count: s.Count}
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
				count := s.CatRoot.Summary.Count
				// fmt.Println("original count", count)
				s.CatRoot = &CatNode{Summary: &Summary{Count: count}, Children: cpy} //summary is new, children is old
				if s.CatRoot.Children[index] == nil {
					s.CatRoot.Children[index] = &Summary{Count: 1}
				} else {
					s.CatRoot.Children[index] = s.CatRoot.Children[index].Copy() //only make a new copy on this cat children
					s.CatRoot.Children[index].Count++
				}
				s.CatRoot.Summary.Count++
			}
		}
	} else {
		if s.CatRoot == nil {
			s.CatRoot = &CatNode{Summary: &Summary{Count: 1}, Children: make([]*Summary, len(nc.Types))}
			index := nc.getIndex(obj.Type)
			s.CatRoot.Children[index] = s.CatRoot.Summary
		} else { //need update
			index := nc.getIndex(obj.Type)
			if s.CatRoot.Children[index] != nil {
				s.CatRoot.Children[index] = s.CatRoot.Children[index].Copy()
				s.CatRoot.Children[index].Count++
				// s.CatRoot.Summary.Count++
			} else {
				s.CatRoot.Children[index] = &Summary{Count: 1}
				// for i := 0; i < len(s.CatRoot.Children); i++ { //all categories unaffected
				// 	if i != index {
				// 		if s.CatRoot.Children[i] != nil {
				// 			s.CatRoot.Children[i] = &Summary{Count: s.CatRoot.Children[i].Count} //deep copy
				// 		}
				// 	}
				// }
			}
			s.CatRoot.Summary = &Summary{Count: s.CatRoot.Summary.Count + 1}
		}
	}
}

//AddObject Add an object
func (nc *Nanocube) AddObject(obj Object) {
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
		currentNode.UpdateSummary(obj, levels, nc)
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
