package entropy

import (
	"math"
	"path"
	"strings"

	"github.com/elliotchance/orderedmap/v2"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type DirTreeEntry struct {
	entry *DirTree
	index int
}
type DirTree orderedmap.OrderedMap[string, DirTreeEntry]

func NewDirTree() *DirTree {
	return (*DirTree)(orderedmap.NewOrderedMap[string, DirTreeEntry]())
}

func (d *DirTree) inner() *orderedmap.OrderedMap[string, DirTreeEntry] {
	return (*orderedmap.OrderedMap[string, DirTreeEntry])(d)
}

func splitFullPaths(dir string) []string {
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
	return result
}

func splitBaseNames(dir string) []string {
	fullPaths := splitFullPaths(dir)
	result := make([]string, len(fullPaths))
	for i := range fullPaths {
		result[i] = path.Base(fullPaths[len(fullPaths)-i-1])
	}
	return result
}

func (d *DirTree) AddDirs(dir string) {
	node := d.inner()
	for _, p := range splitBaseNames(dir) {
		base := path.Base(p)
		if upper, ok := node.Get(base); ok {
			node = upper.entry.inner()
		} else {
			newNode := NewDirTree()
			node.Set(base, DirTreeEntry{newNode, node.Len()})
			node = newNode.inner()
		}
	}
}

func (d *DirTree) ColorFor(dir string) []int {
	baseNames := splitBaseNames(dir)
	depth := 0
	node := d.inner()
	h, s, v := float64(0), 0., 1.
	for depth < len(baseNames) {
		el, ok := node.Get(baseNames[depth])
		if !ok {
			return []int{0, 0, 0}
		}
		h = float64(int(h+360*float64(el.index)/float64(node.Len())) % 360)
		s = utils.Scale(1-float64(depth)/float64(len(baseNames)), 0, 1, .2, .9)

		depth += 1
		node = el.entry.inner()
	}
	r, g, b := HSVToRGB(h, s, v)
	return []int{int(r), int(g), int(b)}
}

// HSVToRGB converts an HSV triple to an RGB triple.
// taken from https://github.com/Crazy3lf/colorconv/blob/master/colorconv.go
//
//nolint:gocyclo
func HSVToRGB(h, s, v float64) (r, g, b uint8) {
	if h < 0 || h >= 360 ||
		s < 0 || s > 1 ||
		v < 0 || v > 1 {
		return 0, 0, 0
	}
	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ v ≤ 1:
	C := v * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - C
	var Rnot, Gnot, Bnot float64
	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r = uint8(math.Round((Rnot + m) * 255))
	g = uint8(math.Round((Gnot + m) * 255))
	b = uint8(math.Round((Bnot + m) * 255))
	return r, g, b
}
