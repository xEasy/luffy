package consistent

import (
	"fmt"
	"strconv"
	"testing"
)

func TestConsistent(t *testing.T) {
	cs := NewConsistent()
	cs.Add("one", 150)
	cs.Add("two", 100)

	err := cs.Add("one", 1)
	if err == nil {
		t.Fatal("add existed node!")
	}

	node := cs.GetNode("hello")
	if node == "" {
		t.Fatal("get empty node info!")
	}
	t.Logf("get node success, node = %s, key = %s", node, "hello")

	err = cs.Remove("one")
	if err != nil {
		t.Logf("remove node <one> fail!")
	}

	if len(cs.hashSortNodes) != 100 {
		t.Logf("remove node wrong")
	}
	if len(cs.circle) != 100 {
		t.Logf("remove node wrong")
	}

	node = cs.GetNode("hello")
	if node != "two" {
		t.Logf("Get node wrong")
	}

}

func BenchmarkConsistent(t *testing.B) {
	// test 10 nodes with 100 vitual nodes
	cs := NewConsistent()
	for i := 0; i < 100; i++ {
		cs.Add(fmt.Sprintf("node:%d", i), 100)
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		cs.GetNode("call:" + strconv.Itoa(i))
	}
}
