package entropy

import (
	"path"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
)

type DirTree orderedmap.OrderedMap[string, *DirTree]

func NewDirTree() *DirTree {
	return (*DirTree)(orderedmap.NewOrderedMap[string, *DirTree]())
}

func (d *DirTree) Inner() *orderedmap.OrderedMap[string, *DirTree] {
	return (*orderedmap.OrderedMap[string, *DirTree])(d)
}

func (d *DirTree) AddDirs(dir string) []string {
	var result []string
	for strings.Contains(dir, "/") {
		if strings.HasSuffix(dir, "/") {
			dir = dir[:len(dir)-1]
		}
		if strings.HasPrefix(dir, "/") {
			dir = dir[1:]
		}
		result = append(result, dir)
		dir = path.Dir(dir)
	}
	if dir != "" && dir != "." {
		result = append(result, dir)
	}
	node := d.Inner()
	for i := range result {
		p := result[len(result)-i-1]
		base := path.Base(p)
		if upper, ok := node.Get(base); ok {
			node = upper.Inner()
		} else {
			newNode := NewDirTree()
			node.Set(base, newNode)
			node = newNode.Inner()
		}
	}
	return result
}
