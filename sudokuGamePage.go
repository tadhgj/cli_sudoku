// job: handle updates, key presses, and view calls
package main

import tea "charm.land/bubbletea/v2"

func (s *SudokuGameWrapperState) GameUpdate(msg tea.Msg) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			s.boardInteraction.MoveCursorUp()
		case "down":
			s.boardInteraction.MoveCursorDown()
		case "left":
			s.boardInteraction.MoveCursorLeft()
		case "right":
			s.boardInteraction.MoveCursorRight()
		// todo: make this neater
		case "1":
			s.boardInteraction.SetNumberAtCursor(1)
		case "2":
			s.boardInteraction.SetNumberAtCursor(2)
		case "3":
			s.boardInteraction.SetNumberAtCursor(3)
		case "4":
			s.boardInteraction.SetNumberAtCursor(4)
		case "5":
			s.boardInteraction.SetNumberAtCursor(5)
		case "6":
			s.boardInteraction.SetNumberAtCursor(6)
		case "7":
			s.boardInteraction.SetNumberAtCursor(7)
		case "8":
			s.boardInteraction.SetNumberAtCursor(8)
		case "9":
			s.boardInteraction.SetNumberAtCursor(9)

		case "delete", "backspace":
			s.boardInteraction.SetNumberAtCursor(0)

			// exit back to difficulty screen
			// case "esc", "p":
			// 	s.shownPage = loadingPage
		}

	}
}

func (s *SudokuGameWrapperState) GameView() tea.View {
	v := tea.NewView(s.boardInteraction.RenderSudokuBoardState(s.styles))
	v.AltScreen = true
	return v
}
