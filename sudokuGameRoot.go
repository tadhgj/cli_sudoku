// job: determine which 'page' is visible at the moment and render and update it
package main

import (
	// "fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type page = int

const (
	loadingPage    page = iota // unused
	splashPage                 // unused
	difficultyPage             // starting page
	gamePage
	pausePage // unused
)

type SudokuGameWrapperState struct {
	boardInteraction      SudokuBoardInteractionState
	difficultyInteraction DifficultyInteractionState
	shownPage             page
	styles                styles
}

func NewWrapper() SudokuGameWrapperState {
	return SudokuGameWrapperState{
		boardInteraction:      SudokuBoardInteractionState{},
		difficultyInteraction: DifficultyInteractionState{},
		shownPage:             difficultyPage,
		styles:                newStyles(true), // assume dark terminal until we are told otherwise
	}
}

func (s *SudokuGameWrapperState) SetToFreshWrapper() {
	s.boardInteraction = GenerateSudokuBoardState("")
	s.difficultyInteraction = GenerateDifficultyState()
	s.shownPage = gamePage
}

// func (s *SudokuGameWrapperState) SetWrapper(...) {

// }

type styles struct {
	primary             lipgloss.Style
	secondary           lipgloss.Style
	tertiary            lipgloss.Style
	accent              lipgloss.Style
	errorForeground     lipgloss.Style
	invertHighlight     lipgloss.Style
	primaryInvert       lipgloss.Style
	secondaryInvert     lipgloss.Style
	tertiaryInvert      lipgloss.Style
	accentInvert        lipgloss.Style
	errorInvert         lipgloss.Style
	lightHighlight      lipgloss.Style
	darkHighlight       lipgloss.Style
	errorHighlightUser  lipgloss.Style
	errorHighlightGiven lipgloss.Style
	keybindKey          lipgloss.Style
	keybindText         lipgloss.Style
}

func newStyles(bgIsDark bool) styles {

	// fmt.Println("newStyles called: is dark? ")
	// fmt.Println(bgIsDark)

	lightDark := lipgloss.LightDark(!bgIsDark)

	return styles{
		primary: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#ddd"), // on dark background
			lipgloss.Color("#000"), // on light background
		)),
		secondary: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#bbb"),
			lipgloss.Color("#555"),
		)),
		tertiary: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#555"),
			lipgloss.Color("#ccc"),
		)),
		accent: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#77f"),
			lipgloss.Color("#33f"),
		)),
		errorForeground: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#f00"),
			lipgloss.Color("#f00"),
		)),
		invertHighlight: lipgloss.NewStyle().Background(lightDark(
			lipgloss.Color("#fff"),
			lipgloss.Color("#000"),
		)),
		primaryInvert: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#000"),
			lipgloss.Color("#fff"),
		)),
		secondaryInvert: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#555"),
			lipgloss.Color("#bbb"),
		)),
		tertiaryInvert: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#ccc"),
			lipgloss.Color("#555"),
		)),
		accentInvert: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#33f"),
			lipgloss.Color("#77f"),
		)),
		errorInvert: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#f00"),
			lipgloss.Color("#f00"),
		)),
		lightHighlight: lipgloss.NewStyle().Background(lightDark(
			lipgloss.Color("#333"),
			lipgloss.Color("#efefef"),
		)),
		darkHighlight: lipgloss.NewStyle().Background(lightDark(
			lipgloss.Color("#777"),
			lipgloss.Color("#bbb"),
		)),
		errorHighlightGiven: lipgloss.NewStyle().Background(lightDark(
			lipgloss.Color("#f00"),
			lipgloss.Color("#f00"),
		)),
		errorHighlightUser: lipgloss.NewStyle().Background(lightDark(
			lipgloss.Color("#f00"),
			lipgloss.Color("#f00"),
		)).Foreground(lightDark(
			lipgloss.Color("#900"),
			lipgloss.Color("#900"),
		)),
		keybindKey: lipgloss.NewStyle().Bold(true).Foreground(lightDark(
			lipgloss.Color("#fff"),
			lipgloss.Color("#000"),
		)),
		keybindText: lipgloss.NewStyle().Foreground(lightDark(
			lipgloss.Color("#ccc"),
			lipgloss.Color("#666"),
		)),
	}
}

func (s SudokuGameWrapperState) ReturnControl(controlKeyText string, controlKeyDesc string) string {
	returnText := s.styles.keybindKey.Render(controlKeyText)
	returnText += " "
	returnText += s.styles.keybindText.Render(controlKeyDesc)
	return returnText
}

func (s SudokuGameWrapperState) WrapperUpdate(msg tea.Msg) (SudokuGameWrapperState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// set styles for sudokugamewrapperstate
		s.styles = newStyles(msg.IsDark())
	case tea.KeyMsg:
		// todo make nice keymap
		// todo: handle these controls based on a 'control scheme' option
		// that can be changed by the user
		switch s.shownPage {
		case difficultyPage:
			s.DifficultyUpdate(msg)
		case gamePage:
			s.GameUpdate(msg)
		}
		switch msg.String() {
		case "q":
			return s, tea.Quit
		}

	}
	return s, nil
}

func (s SudokuGameWrapperState) WrapperView() tea.View {
	switch s.shownPage {
	// case loadingPage:
	// 	return LoadingView()
	// case splashPage:
	// 	break
	case difficultyPage:
		return s.DifficultyView()
	case gamePage:
		return s.GameView()
		// case pausePage:
		// 	break
	}

	var contentString string
	contentString = "Unknown page shown"
	v := tea.NewView(contentString)
	v.AltScreen = true
	return v
}
