package systems

import "github.com/gdamore/tcell/v2"

func PrintScreen(s tcell.SimulationScreen) string {
	result := ""
	w, h := s.Size()
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			char, _, _, _ := s.GetContent(j, i)
			result += string(char)
		}
		result += "\n"
	}
	return result
}
