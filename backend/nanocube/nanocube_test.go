package nanocube

import (
	"fmt"
	"testing"
)

func TestStruct(t *testing.T) {
	list := make([]*SpatNode, 4)
	fmt.Println(list)
	b := Bounds{1, 1, 1, 1}
	sn := SpatNode{Bounds: b, Children: make([]*SpatNode, 4)}
	fmt.Println(sn)
}
func TestNanoCubeSetUp(t *testing.T) {
	nb := SetUpCube(10, Bounds{1, -1, 10, 10}, []string{"Android", "iphone"})
	fmt.Println("CatRoot is ", nb.Index)
}

func TestAssignIndexOnBounds(t *testing.T) {
	b := Bounds{1, -1, 3, 3}
	obj := Object{1.5, 1.5, "A", 3}
	fmt.Println(AssignIndexOnBounds(obj, b))
}

func TestBoundsIntersect(t *testing.T) {
	b1 := Bounds{0, 0, 3, 3}
	b2 := Bounds{1, -1, 4, 4}
	fmt.Println(b1.Intersect(b2))
}

func TestAddObj(t *testing.T) {
	nb := SetUpCube(16, Bounds{0, 0, 8, 8}, []string{"Android", "iPhone"})
	nb.AddObject(Object{3.1, -3.1, "Android", 50})
	fmt.Println("c:", nb.Root.CatRoot)
	fmt.Println("c:", nb.Root.Children[0].CatRoot)
	fmt.Println("c:", nb.Root.Children[0].Children[3].CatRoot)
	fmt.Println("c:", nb.Root.Children[0].Children[3].Children[3].CatRoot)
	fmt.Println("1:", nb.Root.CatRoot.Children[0])
	// fmt.Println(nb.Root.Children[0].Children[3].Summary)
	// if (nb.Root.CatRoot.Summary != nb.Root.Children[0].CatRoot.Summary) || (nb.Root.Children[0].Summary != nb.Root.Children[0].Children[3].Summary) {
	// 	t.Errorf("These three address should be equal")
	// }
	fmt.Println("Now children is:", nb.Root.Children[0].Children[3].Children)
	nb.AddObject(Object{2.51, -3.51, "iPhone", 50})
	fmt.Println(nb.Root.CatRoot.Summary)
	fmt.Println("2:", nb.Root.CatRoot.Children[0])
	fmt.Println("2 1st level:", nb.Root.Children[0].CatRoot.Children[0])
	// if (nb.Root.Summary != nb.Root.Children[0].Summary) || (nb.Root.Children[0].Summary != nb.Root.Children[0].Children[3].Summary) {
	// 	t.Errorf("These three address should be equal")
	// }
	nb.AddObject(Object{5, -5, "iPhone", 50})
	fmt.Println("3: Android", nb.Root.CatRoot.Children[0])
	fmt.Println("3: Iphone", nb.Root.CatRoot.Children[1])
	fmt.Println("3: 1st level children 0 Android", nb.Root.Children[0].CatRoot.Children[0])
	fmt.Println("3: 1st level children 3 Iphone", nb.Root.Children[3].CatRoot.Children[0])
	fmt.Println(nb.Root.Children[0].CatRoot.Summary)
	// fmt.Println(nb.Root.Children[3])
	// fmt.Println(nb.Root.Children[3].Children[0].Summary)
	if (nb.Root.CatRoot == nb.Root.Children[0].CatRoot) || (nb.Root.Children[0].CatRoot == nb.Root.Children[3].CatRoot) {
		t.Errorf("These three address should not be equal")
	}
	nb.AddObject(Object{5, -3, "Android", 50})
	fmt.Println(nb.Root.CatRoot.Summary)
	fmt.Println("4: Android", nb.Root.CatRoot.Children[0])
	fmt.Println("4: Iphone", nb.Root.CatRoot.Children[1])
	fmt.Println(QueryAll(nb.Root, 15))
	fmt.Println(Query(nb.Root, Bounds{1.5, -1.5, 5, 5}, 3))
}
