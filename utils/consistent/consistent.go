package consistent

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Hash hashing
type Hash struct {
	sync.RWMutex                          // map RWMutex lock
	hashSortNodes     []uint32            // sorted hash virtual nodes
	circle            map[uint32]string   // virtual nodes info
	nodes             map[string]bool     //binding nodes
	nodeVirtuals      map[string][]uint32 // node's virtual nodes
	virtualNodesCount int
}

func NewConsistent() *Hash {
	return &Hash{}
}

func (c *Hash) Add(node string, virtualNodeCount int) error {
	if node == "" {
		return nil
	}
	c.Lock()
	defer c.Unlock()

	if c.circle == nil {
		c.circle = make(map[uint32]string)
	}

	if c.nodes == nil {
		c.nodes = make(map[string]bool)
	}

	if c.nodeVirtuals == nil {
		c.nodeVirtuals = make(map[string][]uint32)
	}

	if _, ok := c.nodes[node]; ok {
		return errors.New("node existed!")
	}
	c.nodes[node] = true
	c.nodeVirtuals[node] = make([]uint32, virtualNodeCount)

	// add virtual nodes
	for i := 0; i < virtualNodeCount; i++ {
		virtualKey := c.hashKey(strings.Join([]string{node, strconv.Itoa(i)}, "-"))
		c.circle[virtualKey] = node
		c.nodeVirtuals[node][i] = virtualKey
		c.hashSortNodes = append(c.hashSortNodes, virtualKey)
	}

	// sort virtual nodes
	sort.Slice(c.hashSortNodes, func(x, y int) bool {
		return c.hashSortNodes[x] < c.hashSortNodes[y]
	})

	return nil
}

func (c *Hash) GetNode(key string) string {
	c.RLock()
	defer c.RUnlock()

	hash := c.hashKey(key)
	i := c.getPosition(hash)

	return c.circle[i]
}

func (c *Hash) Remove(node string) error {
	if node == "" {
		return nil
	}

	c.Lock()
	defer c.Unlock()

	if _, ok := c.nodes[node]; !ok {
		return errors.New("node not existed!")
	} else {
		delete(c.nodes, node)
	}

	// remove virtual nodes from sortedSlice
	nodeVirtuals := c.nodeVirtuals[node]
	index := 0
	for i := 0; i < len(c.hashSortNodes); i++ {
		nodeKey := c.hashSortNodes[i]
		if c.circle[nodeKey] != node {
			continue
		}
		c.hashSortNodes[index] = c.hashSortNodes[i]
		index++
	}
	c.hashSortNodes = c.hashSortNodes[:index]

	// delete virtual nodes form circle
	for _, v := range nodeVirtuals {
		delete(c.circle, v)
	}

	delete(c.nodeVirtuals, node)

	return nil
}

func (c *Hash) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Hash) getPosition(hash uint32) uint32 {
	nodesLen := len(c.hashSortNodes)
	i := sort.Search(nodesLen, func(i int) bool {
		return c.hashSortNodes[i] >= hash
	})
	if i < nodesLen {
		if i == nodesLen-1 {
			i = 0
		}
	} else {
		i = nodesLen - 1
	}

	return c.hashSortNodes[i]
}
