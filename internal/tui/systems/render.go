package systems

import (
	"strings"

	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/board/graphics"
	"dep-tree/internal/dep_tree"
	"dep-tree/internal/utils"
)

type RenderState struct {
	Cells  [][]graphics.CellStack
	Errors map[string][]error
}

func computeCursor(s *State, rs *RenderState) {
	if s.Cursor.Y < len(rs.Cells) {
		for j := range rs.Cells[s.Cursor.Y] {
			if nodeId := rs.Cells[s.Cursor.Y][j].Tag(dep_tree.NodeIdTag); nodeId != "" {
				s.Cursor.X = j
				s.SelectedId = nodeId
				return
			}
		}
	}
	s.SelectedId = ""
}

var defaultStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var primaryStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
var errorStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)
var errorSelectedStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)
var secondaryStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDarkCyan)

func forEachCell(
	s *State,
	rs *RenderState,
	ss *SpatialState,
) {
	from := utils.Vec(ss.Offset.X, ss.Offset.Y)
	to := utils.Vec(ss.Offset.X+ss.ScreenSize.X, ss.Offset.Y+ss.ScreenSize.Y)

	for i := from.Y; i < to.Y; i++ {
		if i >= len(rs.Cells) || i < 0 {
			break
		}
		for j := from.X; j < to.X; j++ {
			if j >= len(rs.Cells[i]) || j < 0 {
				break
			}
			cell := rs.Cells[i][j]
			priorityTags := map[string]string{}
			style := defaultStyle
			if errors := rs.Errors[cell.Tag(dep_tree.NodeIdTag)]; len(errors) > 0 {
				style = errorStyle
			}
			switch {
			case s.SelectedId == "":
				// nothing here.
			case cell.Is(dep_tree.NodeIdTag, s.SelectedId):
				if style == errorStyle {
					style = errorSelectedStyle
				} else {
					style = primaryStyle
				}
			case cell.Is(dep_tree.ConnectorOriginNodeIdTag, s.SelectedId):
				style = primaryStyle
				priorityTags[dep_tree.ConnectorOriginNodeIdTag] = s.SelectedId
			case strings.Contains(cell.Tag(dep_tree.NodeParentsTag), s.SelectedId):
				style = secondaryStyle
			}

			s.Screen.SetContent(
				j-from.X,
				i-from.Y,
				cell.Render(priorityTags),
				nil,
				style,
			)
		}
	}
}

const renderErrorMargin = 40

func extractWords(err error, maxLength int) []string {
	words := []string{"-"}
	for _, word := range strings.Split(err.Error(), " ") {
		// If the word is too long break it.
		for len(word) > maxLength {
			brokenWord := word[:maxLength]
			word = word[maxLength:]
			words = append(words, brokenWord)
		}
		words = append(words, word)
	}
	return words
}

func renderError(
	s *State,
	rs *RenderState,
	ss *SpatialState,
) {
	w := ss.ScreenSize.X
	availableSpace := utils.Clamp(renderErrorMargin, renderErrorMargin, w-renderErrorMargin)
	wordsStart := w - availableSpace

	seen := make(map[string]bool)

	// first, retrieve lines.
	lines := make([][]string, 1)
	for _, err := range rs.Errors[s.SelectedId] {
		if _, ok := seen[err.Error()]; ok {
			continue
		}
		seen[err.Error()] = true
		x := wordsStart
		// get words from error message.
		words := extractWords(err, availableSpace)
		// make lines with those words.
		for _, word := range words {
			if x+len(word)+1 >= w {
				x = wordsStart
				lines = append(lines, []string{})
			}
			x += len(word) + 1 // +1 because a space is added after each word.
			lines[len(lines)-1] = append(lines[len(lines)-1], word)
		}
		lines = append(lines, []string{})
	}

	// second, render each line.
	for y, line := range lines {
		x := wordsStart
		for _, word := range line {
			for _, letter := range word + " " {
				s.Screen.SetContent(
					x,
					y,
					letter,
					nil,
					errorStyle,
				)
				x++
			}
		}
	}
}

func RenderSystem(s *State, rs *RenderState, ss *SpatialState) {
	computeCursor(s, rs)
	forEachCell(s, rs, ss)
	renderError(s, rs, ss)
}
