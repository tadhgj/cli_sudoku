package main

import (
	"strconv"

	"charm.land/lipgloss/v2"
)

const boardSizeGeneric int = 9
const boardWidth int = boardSizeGeneric
const boardHeight int = boardSizeGeneric

const cursorHorizMin int = 0
const cursorHorizMax int = boardWidth
const cursorVertMin int = 0
const cursorVertMax int = boardHeight

const perSquareMin int = 0
const perSquareMax int = 9

type BoardGridNotes struct {
	board [boardHeight][boardWidth][(perSquareMax - perSquareMin)]int
}

type BoardGrid struct {
	board [boardHeight][boardWidth]int
}

type SudokuBoard2 struct {
	// solutionBoard ?
	givenBoard  BoardGrid
	userEntries BoardGrid
	userNotes   BoardGridNotes
}

type BoardPosition struct {
	horiz int
	vert  int
}

type SudokuBoardInteractionState struct {
	board          SudokuBoard2
	selectedNumber int
	cursor         BoardPosition
	wrap           bool // should cursor wrap over
	toggle         bool // should entering 5 in same square clear square
}

func internalGenerateBlankSudokuBoard() BoardGrid {
	return BoardGrid{
		board: [9][9]int{
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
}

func internalGenerateBlankSudokuBoardNotes() BoardGridNotes {
	return BoardGridNotes{
		board: [9][9][9]int{
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
}

func internalGenerateSudokuBoard(difficulty string) BoardGrid {
	if difficulty == "" {
		return BoardGrid{
			board: [9][9]int{
				{5, 3, 0, 0, 7, 0, 0, 0, 0},
				{6, 0, 0, 1, 9, 5, 0, 0, 0},
				{0, 9, 8, 0, 0, 0, 0, 6, 0},
				{8, 0, 0, 0, 6, 0, 0, 0, 3},
				{4, 0, 0, 8, 0, 3, 0, 0, 1},
				{7, 0, 0, 0, 2, 0, 0, 0, 6},
				{0, 6, 0, 0, 0, 0, 2, 8, 0},
				{0, 0, 0, 4, 1, 9, 0, 0, 5},
				{0, 0, 0, 0, 8, 0, 0, 7, 9},
			},
		}
	} else {
		return BoardGrid{
			board: [9][9]int{
				{5, 3, 0, 0, 7, 0, 0, 0, 0},
				{6, 0, 0, 1, 9, 5, 0, 0, 0},
				{0, 9, 8, 0, 0, 0, 0, 6, 0},
				{8, 0, 0, 0, 6, 0, 0, 0, 3},
				{4, 0, 0, 8, 0, 3, 0, 0, 1},
				{7, 0, 0, 0, 2, 0, 0, 0, 6},
				{0, 6, 0, 0, 0, 0, 2, 8, 0},
				{0, 0, 0, 4, 1, 9, 0, 0, 5},
				{0, 0, 0, 0, 8, 0, 0, 7, 9},
			},
		}
	}
}

func GenerateSudokuBoard(difficulty string) SudokuBoard2 {
	return SudokuBoard2{
		givenBoard:  internalGenerateSudokuBoard(difficulty),
		userEntries: internalGenerateBlankSudokuBoard(),
		userNotes:   internalGenerateBlankSudokuBoardNotes(),
	}
}

func GenerateSudokuBoardState(difficulty string) SudokuBoardInteractionState {
	return SudokuBoardInteractionState{
		board:          GenerateSudokuBoard(difficulty),
		selectedNumber: -1,
		cursor: BoardPosition{
			horiz: 0,
			vert:  0,
		},
		wrap: false,
	}
}

func RenderSudokuBoardState(state SudokuBoardInteractionState, style lipgloss.Style) string {

	lockedNumberStyle := style.Foreground(lipgloss.Color("#000"))
	userNumberStyle := style.Foreground(lipgloss.Color("#00f"))
	blankStyle := style.Foreground(lipgloss.Color("#ccc"))

	board := state.board
	cursor := state.cursor

	// get number under cursor
	var cursorSelectedNum int
	if board.givenBoard.board[cursor.vert][cursor.horiz] != 0 {
		cursorSelectedNum = board.givenBoard.board[cursor.vert][cursor.horiz]
	} else if board.userEntries.board[cursor.vert][cursor.horiz] != 0 {
		cursorSelectedNum = board.userEntries.board[cursor.vert][cursor.horiz]
	} else {
		cursorSelectedNum = -1
	}

	var result string
	for i := 0; i < boardHeight; i++ {
		for j := 0; j < boardWidth; j++ {

			// look at given board
			var currentNum int
			var currentChar string
			// currentNumForegroundStyle := lockedNumberStyle
			if board.givenBoard.board[i][j] != 0 {
				// 'locked' number
				currentNum = board.givenBoard.board[i][j]
				userNumberString := strconv.Itoa(currentNum)
				currentChar = lockedNumberStyle.Render(userNumberString)
			} else if board.userEntries.board[i][j] != 0 {
				// user generated number
				currentNum = board.userEntries.board[i][j]
				userNumberString := strconv.Itoa(currentNum)
				// currentNumForegroundStyle = userNumberStyle
				currentChar = userNumberStyle.Render(userNumberString)
			} else {
				// blank
				currentChar = blankStyle.Render(".")
			}
			// check if special highlights need to be applied to this number
			// cases: (in order of precedence)
			// userHighlight
			// numberSelected
			if cursor.vert == i && cursor.horiz == j {
				result += style.Background(lipgloss.Color("#000")).Foreground(lipgloss.Color("#fff")).Render(currentChar)
			} else if currentNum == cursorSelectedNum {
				result += style.Background(lipgloss.Color("#888")).Render(currentChar)
			} else {
				result += style.Render(currentChar)
			}
			// spacer
			if j%3 == 2 {
				result += "|"
			} else {
				result += " "
			}
		}
		result += "\n"

		// spacer
		if i%3 == 2 {
			result += "- - -+- - -+- - - \n"
		}
	}
	return result
}

func (cursor *BoardPosition) setCursor(pos BoardPosition, wrap bool) {
	// enforce bounds
	if wrap {
		if pos.horiz < cursorHorizMin {
			pos.horiz = cursorHorizMax - 1
		}
		if pos.vert < cursorVertMin {
			pos.vert = cursorVertMax - 1
		}
		if pos.horiz >= cursorHorizMax {
			pos.horiz = cursorHorizMin
		}
		if pos.vert >= cursorVertMax {
			pos.vert = cursorVertMin
		}
	} else {
		if pos.horiz < cursorHorizMin {
			pos.horiz = cursorHorizMin
		}
		if pos.vert < cursorVertMin {
			pos.vert = cursorVertMin
		}
		if pos.horiz >= cursorHorizMax {
			pos.horiz = cursorHorizMax - 1
		}
		if pos.vert >= cursorVertMax {
			pos.vert = cursorVertMax - 1
		}
	}
	// set
	cursor.horiz = pos.horiz
	cursor.vert = pos.vert
}

// public
func (board *SudokuBoardInteractionState) SetCursor(pos BoardPosition) {
	board.cursor.setCursor(pos, board.wrap)
}

func (board *BoardGrid) setNumberAtPos(input int, pos BoardPosition, toggle bool) {
	// check if int is valid is valid
	if input < perSquareMin || input > perSquareMax {
		return
	}
	// check if cursor position is valid
	// TODO
	if input != 0 && toggle && board.board[pos.vert][pos.horiz] == input {
		board.board[pos.vert][pos.horiz] = 0
	} else {
		board.board[pos.vert][pos.horiz] = input
	}
}

func (board *SudokuBoardInteractionState) SetNumberAtCursor(input int) {
	// check if givenBoard is blank there; we can't edit 'set' positions
	if board.board.givenBoard.board[board.cursor.vert][board.cursor.horiz] != 0 {
		return
	}
	board.board.userEntries.setNumberAtPos(input, board.cursor, board.toggle)
}

func (board *SudokuBoardInteractionState) MoveCursorLeft() {
	board.SetCursor(
		BoardPosition{
			horiz: board.cursor.horiz - 1,
			vert:  board.cursor.vert,
		},
	)
}

func (board *SudokuBoardInteractionState) MoveCursorRight() {
	board.SetCursor(BoardPosition{
		horiz: board.cursor.horiz + 1,
		vert:  board.cursor.vert,
	})
}

func (board *SudokuBoardInteractionState) MoveCursorUp() {
	board.SetCursor(BoardPosition{
		horiz: board.cursor.horiz,
		vert:  board.cursor.vert - 1,
	})
}

func (board *SudokuBoardInteractionState) MoveCursorDown() {
	board.SetCursor(BoardPosition{
		horiz: board.cursor.horiz,
		vert:  board.cursor.vert + 1,
	})
}
