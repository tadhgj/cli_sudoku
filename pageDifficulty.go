package main

import tea "charm.land/bubbletea/v2"

// enumerate difficulties
type Difficulty string

// please update "Difficulties" array after editing this const list
const (
	testeasy Difficulty = "(FOR TESTING) Extremely easy"
	easy     Difficulty = "Easy"
	normal   Difficulty = "Normal"
	hard     Difficulty = "Hard"
)

var Difficulties = []Difficulty{testeasy, easy, normal, hard}

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

	s.AbsoluteDifficultySelected(difficulty)
}

func (s *SudokuGameWrapperState) AbsoluteDifficultySelected(difficulty Difficulty) {
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

		case "up", "k":
			s.difficultyInteraction.SelectPreviousDifficulty()
		case "down", "j":
			s.difficultyInteraction.SelectNextDifficulty()

		// // quick-picks
		// case "e":
		// 	s.AbsoluteDifficultySelected(easy)
		// case "n":
		// 	s.AbsoluteDifficultySelected(normal)
		// case "h":
		// 	s.AbsoluteDifficultySelected(hard)

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
	difficultyArrowControl := Control{
		keys: "↑/↓",
		desc: "difficulty",
	}
	difficultySelectControl := Control{
		keys: "enter",
		desc: "select",
	}
	difficultyText += s.RenderControlList([]Control{difficultyArrowControl, difficultySelectControl, QuitControl()})

	v := tea.NewView(difficultyText)
	v.AltScreen = true
	return v
}
