// job: handle updates, key presses, and view calls
package main

import tea "charm.land/bubbletea/v2"

func (s *SudokuGameWrapperState) BackToDifficulty() {
	s.shownPage = difficultyPage
}

func (s *SudokuGameWrapperState) GameToPause() {
	s.shownPage = pausePage
}

func (s *SudokuGameWrapperState) GameUpdate(msg tea.Msg) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			s.MoveCursorUp()
		case "down", "j":
			s.MoveCursorDown()
		case "left", "h":
			s.MoveCursorLeft()
		case "right", "l":
			s.MoveCursorRight()
		// case "shift+up", "K": // doesn't work in Terminal.app, the default mac terminal
		// todo: make this neater
		case "1":
			s.SetNumberAtCursor(1)
		case "2":
			s.SetNumberAtCursor(2)
		case "3":
			s.SetNumberAtCursor(3)
		case "4":
			s.SetNumberAtCursor(4)
		case "5":
			s.SetNumberAtCursor(5)
		case "6":
			s.SetNumberAtCursor(6)
		case "7":
			s.SetNumberAtCursor(7)
		case "8":
			s.SetNumberAtCursor(8)
		case "9":
			s.SetNumberAtCursor(9)

		case "delete", "backspace":
			s.SetNumberAtCursor(0)

		// exit back to difficulty screen
		case "esc":
			s.BackToDifficulty()

		// pause
		case "p":
			s.GameToPause()
		}

	}
}

func (s *SudokuGameWrapperState) GameView() tea.View {
	viewText := s.RenderSudokuBoardState(s.styles)

	// render controls
	// TODO: based on control scheme
	viewText += "\n"

	gameArrowControl := Control{
		keys: "↑/↓/←/→",
		desc: "cursor",
	}
	gameSetNumControl := Control{
		keys: "1-9",
		desc: "set",
	}
	gameClearNumControl := Control{
		keys: "backsp",
		desc: "clear",
	}
	gameEscToDiffControl := Control{
		keys: "esc",
		desc: "back to diff.",
	}
	gameEscToPauseControl := Control{
		keys: "p",
		desc: "pause",
	}

	viewText += s.RenderControlList([]Control{gameArrowControl, gameSetNumControl, gameClearNumControl, gameEscToDiffControl, gameEscToPauseControl, QuitControl()})

	v := tea.NewView(viewText)
	v.AltScreen = true
	return v
}
