package entropy

import (
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/gabotechs/dep-tree/internal/graph"
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
	for strings.Contains(dir, string(os.PathSeparator)) {
		if strings.HasSuffix(dir, string(os.PathSeparator)) {
			dir = dir[:len(dir)-1]
		}
		if strings.HasPrefix(dir, string(os.PathSeparator)) {
			dir = dir[1:]
		}
		result = append(result, dir)
		dir = filepath.Dir(dir)
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
		result[i] = filepath.Base(fullPaths[len(fullPaths)-i-1])
	}
	return result
}

func (d *DirTree) AddDirs(dirs []string) {
	node := d.inner()
	for _, dir := range dirs {
		if upper, ok := node.Get(dir); ok {
			node = upper.entry.inner()
		} else {
			newNode := NewDirTree()
			node.Set(dir, DirTreeEntry{newNode, node.Len()})
			node = newNode.inner()
		}
	}
}

func (d *DirTree) AddDirsFromDisplay(display graph.DisplayResult) {
	dirs := splitBaseNames(filepath.Dir(display.Name))
	if display.Group != "" {
		d.AddDirs(utils.AppendFront(display.Group, dirs))
	} else {
		d.AddDirs(dirs)
	}
}

type colorFormat int

const (
	HSV colorFormat = iota
	RGB
)

// ColorForDir smartly assigns a color for the specified dir based on all the dir tree that
// the codebase has. Files in the same folder will receive the same color, and colors for
// each sub folder will be assigned evenly following a radial distribution in an HSV wheel.
// As it goes deeper into more sub folders, colors fade, but the distribution rules are
// the same.
func (d *DirTree) ColorForDir(dirs []string, format colorFormat) []float64 {
	node := d.inner()
	h, s, v := 0., 0., 1.
	depth := 0
	for depth < len(dirs) {
		el, ok := node.Get(dirs[depth])
		if !ok {
			return []float64{0, 0, 0}
		}

		// It might happen that all the nodes have some common folders, like src/,
		// so if literally all of them have the same common folders, we do not want to take
		// them into account for reducing the saturation, as they will appear very faded.
		if node.Len() > 1 {
			h = float64(int(h+360*float64(el.index)/float64(node.Len())) % 360)
			if s == 0 {
				s = 1
			}
			s -= .2
			s = utils.Scale(s, 0, 1, .2, .9)
		}

		depth += 1
		node = el.entry.inner()
	}
	if format == RGB {
		r, g, b := HSVToRGB(h, s, v)
		return []float64{float64(r), float64(g), float64(b)}
	} else {
		return []float64{h, s, v}
	}
}

func (d *DirTree) ColorForDisplay(display graph.DisplayResult, format colorFormat) []float64 {
	dirs := splitBaseNames(filepath.Dir(display.Name))
	if display.Group != "" {
		return d.ColorForDir(utils.AppendFront(display.Group, dirs), format)
	} else {
		return d.ColorForDir(dirs, format)
	}
}

func (d *DirTree) GroupingsForDir(dirs []string) []string {
	depth := 0
	var result []string

	node := d.inner()
	acc := ""
	for depth < len(dirs) {
		acc = filepath.Join(acc, dirs[depth])
		el, ok := node.Get(dirs[depth])
		if !ok {
			return result
		}
		if node.Len() > 1 {
			result = append(result, acc)
		}

		depth += 1
		node = el.entry.inner()
	}
	return result
}

func (d *DirTree) GroupingsForDisplay(display graph.DisplayResult) []string {
	dirs := splitBaseNames(filepath.Dir(display.Name))
	if display.Group != "" {
		return d.GroupingsForDir(utils.AppendFront(display.Group, dirs))
	} else {
		return d.GroupingsForDir(dirs)
	}
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
