package main

import tea "charm.land/bubbletea/v2"

// enumerate difficulties
type Difficulty string

const (
	easy   Difficulty = "Easy"
	normal Difficulty = "Normal"
	hard   Difficulty = "Hard"
)

var Difficulties = []Difficulty{easy, normal, hard}
var DifficultiesCount = len(Difficulties)

type DifficultyInteractionState struct {
	selectedDifficultyIndex int
}

func GenerateDifficultyState() DifficultyInteractionState {
	return DifficultyInteractionState{}
}

func (s *SudokuGameWrapperState) DifficultySelected() {
	// pull difficulty
	difficulty := Difficulties[s.difficultyInteraction.selectedDifficultyIndex]

	// generate board
	s.boardInteraction = GenerateSudokuBoardState(difficulty)

	// set viewed page to new page
	s.shownPage = gamePage
}

func (s *DifficultyInteractionState) SelectNextDifficulty() {
	s.selectedDifficultyIndex++
	if s.selectedDifficultyIndex >= DifficultiesCount {
		s.selectedDifficultyIndex = DifficultiesCount - 1
	}
}
func (s *DifficultyInteractionState) SelectPreviousDifficulty() {
	s.selectedDifficultyIndex--
	if s.selectedDifficultyIndex < 0 {
		s.selectedDifficultyIndex = 0
	}
}

func (s *SudokuGameWrapperState) DifficultyUpdate(msg tea.Msg) {
	// if keypress is s
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up":
			s.difficultyInteraction.SelectPreviousDifficulty()
		case "down":
			s.difficultyInteraction.SelectNextDifficulty()

		case "enter":
			s.DifficultySelected()
		}

	}

}

func (s *SudokuGameWrapperState) DifficultyView() tea.View {

	var difficultyText string

	gameName := "Tadhg-doku"

	difficultyText += gameName + "\n\n"

	arrow := "->"

	// arrowLength := len(arrow)

	for i := 0; i < DifficultiesCount; i++ {
		prefix := "  "

		if i == s.difficultyInteraction.selectedDifficultyIndex {
			prefix = arrow
		}

		difficultyText += prefix + " " + string(Difficulties[i]) + "\n"

	}

	// print controls

	v := tea.NewView(difficultyText)
	v.AltScreen = true
	return v
}
