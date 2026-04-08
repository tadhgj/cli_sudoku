// job: handle updates, key presses, and view calls
package main

import tea "charm.land/bubbletea/v2"

func (s *SudokuGameWrapperState) GameUpdate(msg tea.Msg) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.boardInteraction.MoveCursorUp()
		case "down", "j":
			s.boardInteraction.MoveCursorDown()
		case "left", "h":
			s.boardInteraction.MoveCursorLeft()
		case "right", "l":
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
	viewText := s.boardInteraction.RenderSudokuBoardState(s.styles)

	// render controls
	// TODO: based on control scheme

	viewText += "\n"
	// viewText += s.ReturnControl("↑/↓", "row")
	// viewText += "   "
	viewText += s.ReturnControl("1-9", "set num at cursor")
	viewText += "   "
	// viewText += s.ReturnControl("backspace", "clear num at cursor")
	viewText += s.ReturnControl("⌫", "clear num at cursor")
	viewText += "   "
	viewText += s.ReturnControl("q", "quit")

	v := tea.NewView(viewText)
	v.AltScreen = true
	return v
}
