package main

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
)

const boardSizeGeneric int = 9
const boardWidth int = boardSizeGeneric
const boardHeight int = boardSizeGeneric
const boardBlockCount int = boardSizeGeneric
const blockHeight = 3
const blockWidth = 3

const cursorHorizMin int = 0
const cursorHorizMax int = boardWidth
const cursorVertMin int = 0
const cursorVertMax int = boardHeight

const perSquareMin int = 0
const perSquareMax int = 9

const numberOfPossibleDigits int = perSquareMax - perSquareMin

type CellDigit int

const blankDigitInt int = 0
const blankDigit CellDigit = CellDigit(blankDigitInt)

func IsValidCellDigit(input int) bool {
	return input <= perSquareMax && input >= perSquareMin
}

func CastCellDigit(input int) CellDigit {
	if IsValidCellDigit(input) {
		return CellDigit(input)
	}
	// error
	return CellDigit(0)
}

type BoardGridErrors struct {
	board [boardHeight][boardWidth]bool
}

type BoardGridNotes struct {
	board [boardHeight][boardWidth][numberOfPossibleDigits]int
}

type BoardGrid struct {
	board [boardHeight][boardWidth]int
}

func (inputBoard *BoardGrid) SetValueAt(inputPosition BoardPosition, newValue int) {
	// check if value is good
	if newValue < perSquareMin || newValue > perSquareMax {
		// todo: throw error
		return
	}
	inputBoard.board[inputPosition.vert][inputPosition.horiz] = newValue
}

func (inputErrorBoard *BoardGridErrors) SetValueAt(inputPosition BoardPosition, newValue bool) {
	inputErrorBoard.board[inputPosition.vert][inputPosition.horiz] = newValue
}

func (inputBoard *SudokuBoard) GetAnyValueAt(inputPosition BoardPosition) int {
	if inputBoard.givenBoard.GetValueAt(inputPosition) != 0 {
		return inputBoard.givenBoard.GetValueAt(inputPosition)
	} else {
		return inputBoard.userEntries.GetValueAt(inputPosition)
	}
}

func (inputErrorBoard *BoardGridErrors) GetValueAt(inputPosition BoardPosition) bool {
	return inputErrorBoard.board[inputPosition.vert][inputPosition.horiz]
}

func (inputBoard *BoardGrid) GetValueAt(inputPosition BoardPosition) int {
	return inputBoard.board[inputPosition.vert][inputPosition.horiz]
}

func (inputBoard *SudokuBoard) GetRowArray(inputRow int) [boardWidth]int {
	// check bounds
	if inputRow <= cursorHorizMin {
		// throw error
		// for now, just trim
		inputRow = cursorHorizMin
	}
	if inputRow > cursorHorizMax {
		// throw error
		// for now, just trim
		inputRow = (cursorHorizMax - 1)
	}

	var outRow [boardWidth]int
	for colIndex := 0; colIndex < boardWidth; colIndex++ {
		cPos := BoardPosition{horiz: colIndex, vert: inputRow}
		if inputBoard.givenBoard.GetValueAt(cPos) != 0 {
			outRow[colIndex] = inputBoard.givenBoard.GetValueAt(cPos)
		} else {
			outRow[colIndex] = inputBoard.userEntries.GetValueAt(cPos)
		}
	}
	return outRow
}

func (inputBoard *SudokuBoard) GetColumnArray(inputColumn int) [boardHeight]int {
	// TODO: check bounds
	var outColumn [boardHeight]int
	for rowIndex := 0; rowIndex < boardHeight; rowIndex++ {
		cPos := BoardPosition{horiz: inputColumn, vert: rowIndex}
		if inputBoard.givenBoard.GetValueAt(cPos) != 0 {
			outColumn[rowIndex] = inputBoard.givenBoard.GetValueAt(cPos)
		} else {
			outColumn[rowIndex] = inputBoard.userEntries.GetValueAt(cPos)
		}
	}
	return outColumn
}

func (inputBoard *SudokuBoard) GetSquareForPosition(at BoardPosition) [blockHeight][blockWidth]int {
	// TODO: check bounds

	var outArray [blockHeight][blockWidth]int

	// input: BoardPosition 0,0
	// output: row 012 col 012
	// input: BoardPosition 1,1
	// output: row 012 col 012
	// input: BoardPosition 2,2
	// output: row 012 col 012
	// input: BoardPosition 3,3
	// output: row 345 col 345
	fmt.Println("Input: ")
	fmt.Println(at)

	startingRow := blockHeight * (at.vert / blockHeight)
	fmt.Println("startingRow: ")
	fmt.Println(startingRow)
	for rowOffset := 0; rowOffset < blockHeight; rowOffset++ {
		rowIndex := startingRow + rowOffset
		startingColumn := blockWidth * (at.horiz / blockWidth)
		if rowIndex == startingRow {
			fmt.Println("startingColumn: ")
			fmt.Println(startingColumn)
		}
		for colOffset := 0; colOffset < blockWidth; colOffset++ {
			colIndex := startingColumn + colOffset
			fmt.Println("Pulling ")
			cPos := BoardPosition{horiz: colIndex, vert: rowIndex}
			fmt.Println(cPos)
			if inputBoard.givenBoard.GetValueAt(cPos) != 0 {
				outArray[rowOffset][colOffset] = inputBoard.givenBoard.GetValueAt(cPos)
			} else {
				outArray[rowOffset][colOffset] = inputBoard.userEntries.GetValueAt(cPos)
			}
		}
	}

	return outArray
}

type SudokuBoard struct {
	// solutionBoard ?
	givenBoard     BoardGrid
	userEntries    BoardGrid
	invalidEntries BoardGridErrors
	userNotes      BoardGridNotes

	// number of 1s, 2s, 3s, ... present in each row
	numberOfDigitsPerRow [boardHeight][numberOfPossibleDigits]int
	// number of 1s, 2s, 3s, ... present in each col
	numberOfDigitsPerCol [boardWidth][numberOfPossibleDigits]int
	// number of 1s, 2s, 3s, ... present in each square
	numberOfDigitsPerSquare [boardBlockCount][numberOfPossibleDigits]int
	// number of 1s, 2s, 3s, ... present on board
	numberOfDigitsTotal [numberOfPossibleDigits]int
}

type BoardPosition struct {
	horiz int
	vert  int
}

func (pos BoardPosition) String() string {
	return fmt.Sprintf("horiz: %d, vert: %d", pos.horiz, pos.vert)
}

type SudokuBoardInteractionState struct {
	board          SudokuBoard
	selectedNumber int
	cursor         BoardPosition
	wrap           bool // should cursor wrap over
	toggle         bool // should entering 5 in same square clear square
}

func internalGenerateBlankSudokuBoardErrors() BoardGridErrors {
	return BoardGridErrors{} // fills with falses. false treated as no error for this square/cell
}

func internalGenerateBlankSudokuBoard() BoardGrid {
	return BoardGrid{} // fills with 0's; 0 treated as blank
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

func GenerateSudokuBoard(difficulty string) SudokuBoard {
	givenBoard := internalGenerateSudokuBoard(difficulty)

	var numberOfDigitsPerRow [boardHeight][numberOfPossibleDigits]int
	var numberOfDigitsPerCol [boardWidth][numberOfPossibleDigits]int
	var numberOfDigitsPerSquare [boardBlockCount][numberOfPossibleDigits]int
	var numberOfDigitsTotal [numberOfPossibleDigits]int
	for rowIndex := 0; rowIndex < boardHeight; rowIndex++ {
		for colIndex := 0; colIndex < boardWidth; colIndex++ {
			squareIndex := ((rowIndex / 3) * 3) + (colIndex / 3)
			cDigit := givenBoard.GetValueAt(BoardPosition{vert: rowIndex, horiz: colIndex})
			numberIndex := cDigit - 1
			if numberIndex >= 0 {
				numberOfDigitsPerRow[rowIndex][numberIndex]++
				numberOfDigitsPerCol[colIndex][numberIndex]++
				numberOfDigitsPerSquare[squareIndex][numberIndex]++
				numberOfDigitsTotal[numberIndex]++
			}
		}
	}

	return SudokuBoard{
		givenBoard:              givenBoard,
		userEntries:             internalGenerateBlankSudokuBoard(),
		invalidEntries:          internalGenerateBlankSudokuBoardErrors(),
		userNotes:               internalGenerateBlankSudokuBoardNotes(),
		numberOfDigitsPerRow:    numberOfDigitsPerRow,
		numberOfDigitsPerCol:    numberOfDigitsPerCol,
		numberOfDigitsPerSquare: numberOfDigitsPerSquare,
		numberOfDigitsTotal:     numberOfDigitsTotal,
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

func RenderSudokuBoardState(state SudokuBoardInteractionState, style lipgloss.Style, bg string) string {

	// lockedNumberStyle := style.Foreground(lipgloss.Color("#000"))
	// userNumberStyle := style.Foreground(lipgloss.Color("#00f"))
	// blankStyle := style.Foreground(lipgloss.Color("#ccc"))
	// types of numbers:
	// given numbers (primary)
	primaryForeground := style.Foreground(lipgloss.Color("#000"))
	invertPrimaryForeground := style.Foreground(lipgloss.Color("#fff"))

	secondaryForeground := style.Foreground(lipgloss.Color("#555"))
	// user numbers (accent)
	accentForeground := style.Foreground(lipgloss.Color("#00f"))
	invertAccentForeground := style.Foreground(lipgloss.Color("#aaf"))
	// user numbers that violate rules (error, bold)
	errorForeground := style.Foreground(lipgloss.Color("#f00")).Bold(true)
	invertErrorForeground := style.Foreground(lipgloss.Color("#f00"))
	// blank numbers (tertiary)
	tertiaryForeground := style.Foreground(lipgloss.Color("#ccc"))
	invertTertiaryForeground := style.Foreground(lipgloss.Color("#ccc"))

	// types of highlights:
	// under cursor (black highlight)
	invertHighlight := style.Background(lipgloss.Color("#000"))
	// - if given number, use opposite color scheme primary color
	// - if user number, use opposite color scheme accent color
	// - - if isErroneousUserNumber, use accent
	// - if blank number, use opposite color scheme tertiary color

	// same row/col as cursor (light highlight)
	lightHighlight := style.Background(lipgloss.Color("#eee"))
	// same row/col as cursor AND same number as number under cursor (causing error, red highlight)
	// errorHighlight := style.Background(lipgloss.Color("#000"))
	// same square as cursor AND same number as number under cursor (causing error, red highlight)
	// same number as number under cursor (dark highlight)
	darkHighlight := style.Background(lipgloss.Color("#aaa"))
	// IF LIGHT MODE
	if bg == "light" {
		primaryForeground = style.Foreground(lipgloss.Color("#000"))
		invertPrimaryForeground = style.Foreground(lipgloss.Color("#fff"))
		secondaryForeground = style.Foreground(lipgloss.Color("#555"))
		accentForeground = style.Foreground(lipgloss.Color("#00f"))
		invertAccentForeground = style.Foreground(lipgloss.Color("#aaf"))
		errorForeground = style.Foreground(lipgloss.Color("#f00")).Bold(true)
		invertErrorForeground = style.Foreground(lipgloss.Color("#f00"))
		tertiaryForeground = style.Foreground(lipgloss.Color("#ccc"))
		invertTertiaryForeground = style.Foreground(lipgloss.Color("#ccc"))

		// types of highlights:
		invertHighlight = style.Background(lipgloss.Color("#000"))
		lightHighlight = style.Background(lipgloss.Color("#eee"))
		darkHighlight = style.Background(lipgloss.Color("#aaa"))

	} else if bg == "dark" {
		primaryForeground = style.Foreground(lipgloss.Color("#fff"))
		invertPrimaryForeground = style.Foreground(lipgloss.Color("#000"))
		secondaryForeground = style.Foreground(lipgloss.Color("#bbb"))
		accentForeground = style.Foreground(lipgloss.Color("#55f"))
		invertAccentForeground = style.Foreground(lipgloss.Color("#33f"))
		errorForeground = style.Foreground(lipgloss.Color("#f00")).Bold(true)
		invertErrorForeground = style.Foreground(lipgloss.Color("#f00"))
		tertiaryForeground = style.Foreground(lipgloss.Color("#555"))
		invertTertiaryForeground = style.Foreground(lipgloss.Color("#666"))

		// types of highlights:
		invertHighlight = style.Background(lipgloss.Color("#fff"))
		lightHighlight = style.Background(lipgloss.Color("#333"))
		darkHighlight = style.Background(lipgloss.Color("#777"))
	} else {
		// unhandled
	}

	board := state.board
	cursor := state.cursor

	// get number under cursor
	var cursorSelectedNum int
	if board.givenBoard.GetValueAt(cursor) != 0 {
		cursorSelectedNum = board.givenBoard.GetValueAt(cursor)
	} else if board.userEntries.GetValueAt(cursor) != 0 {
		cursorSelectedNum = board.userEntries.GetValueAt(cursor)
	} else {
		cursorSelectedNum = -1 // set to -1 so it doesn't match empty value of 0
	}

	var result string
	// top
	result += "       ,       ,       \n"
	for i := 0; i < boardHeight; i++ {
		isSameRowAsCursor := i == cursor.vert

		// spacer
		result += " "

		for j := 0; j < boardWidth; j++ {

			cPos := BoardPosition{horiz: j, vert: i}

			numberUnderCursor := board.GetAnyValueAt(cPos)
			isEmpty := numberUnderCursor == 0
			var currentChar string
			if isEmpty {
				currentChar = "."
			} else {
				currentChar = strconv.Itoa(numberUnderCursor)
			}

			isSameColumnAsCursor := cPos.horiz == cursor.horiz
			isSameSquareAsCursor := (cPos.horiz/3 == cursor.horiz/3) && (cPos.vert/3 == cursor.vert/3)

			isUnderCursor := isSameColumnAsCursor && isSameRowAsCursor
			isErroneousUserNumber := board.invalidEntries.board[cPos.vert][cPos.horiz]

			// TODO:
			// IF DARK MODE

			// look at given board
			currentNumForegroundStyle := tertiaryForeground
			currentNumInvertedForegroundStyle := invertTertiaryForeground
			if isEmpty {
				// blank
			} else if board.givenBoard.GetValueAt(cPos) != 0 {
				// 'locked' number
				currentNumForegroundStyle = primaryForeground
				currentNumInvertedForegroundStyle = invertPrimaryForeground
			} else if board.userEntries.GetValueAt(cPos) != 0 {
				// user generated number
				currentNumForegroundStyle = accentForeground
				currentNumInvertedForegroundStyle = invertAccentForeground
				// currentNumIsUserEntered = true
			}

			// check if special highlights need to be applied to this number
			// cases: (in order of precedence)
			// userHighlight
			// numberSelected
			if isUnderCursor {
				if isErroneousUserNumber {
					result += invertHighlight.Render(invertErrorForeground.Render(currentChar))
				} else {
					result += invertHighlight.Render(currentNumInvertedForegroundStyle.Render(currentChar))
				}
			} else if numberUnderCursor == cursorSelectedNum {
				matchingCausingError := (isErroneousUserNumber || (isSameColumnAsCursor || isSameRowAsCursor) || isSameSquareAsCursor)
				if matchingCausingError {
					if isSameColumnAsCursor || isSameRowAsCursor {
						result += style.Background(lipgloss.Color("#f00")).Foreground(lipgloss.Color("#900")).Render(currentChar)
					} else {
						result += style.Background(lipgloss.Color("#f00")).Render(currentChar)
					}
				} else {
					result += darkHighlight.Render(currentChar)
				}
			} else if isSameColumnAsCursor || isSameRowAsCursor || isSameSquareAsCursor {
				if isErroneousUserNumber {
					result += lightHighlight.Render(errorForeground.Render(currentChar))
				} else {
					result += lightHighlight.Render(currentNumForegroundStyle.Render(currentChar))
				}
			} else {
				if isErroneousUserNumber {
					result += errorForeground.Render(currentChar)
				} else {
					result += currentNumForegroundStyle.Render(currentChar)
				}
			}

			// spacer
			if j != boardWidth-1 {
				if j%3 == 2 {
					result += " | "
				} else if isSameRowAsCursor || isSameSquareAsCursor {
					result += lightHighlight.Render(" ")
				} else {
					result += " "
				}
			}
		}
		result += "\n"

		// spacer
		if i != boardHeight-1 && i%3 == 2 {
			result += "-------+-------+-------\n"
		}
	}

	// bottom
	result += "       '       '       \n"

	// print remaining numbers
	// if space
	tallEnoughForBottomNumberCount := true

	displayBottomNumberCount := tallEnoughForBottomNumberCount

	if displayBottomNumberCount {
		result += "\n"
		result += "   "
		for numberIndex := (perSquareMin); numberIndex < perSquareMax; numberIndex++ {
			number := numberIndex + 1
			numbStr := strconv.Itoa(number)
			thisNumberFinished := board.numberOfDigitsTotal[numberIndex] >= 9
			numbStrRendered := "x"
			if thisNumberFinished {
				numbStrRendered = tertiaryForeground.Render(numbStr)
			} else {
				numbStrRendered = secondaryForeground.Render(numbStr)
			}
			result += numbStrRendered + " "
		}
		result += "  \n"
		result += "\n"
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
	if input != 0 && toggle && board.GetValueAt(pos) == input {
		board.SetValueAt(pos, input)
	} else {
		board.SetValueAt(pos, input)
	}
}

// func (board *SudokuBoardInteractionState) CheckForAllInvalidPositions {
//
// }

func (board *SudokuBoardInteractionState) CheckForErrorsForMove(existingValue int, newValue int, at BoardPosition) {
	// if we change a value on the board
	//    then we can always remove the error at the cursor position
	//    then the only possible change in erroneous user inputs is along the rows, columns, and square
	//    FOR USER ENTRIES EQUAL TO THE NEW VALUE
	//    then we check if there is over 1 instance of an X
	//    - in the row
	//    - in the column
	//    - in the square
	//    and any position that is user entered is marked as an error

	// get new value
	fmt.Print(fmt.Sprintf("Switching from %d to %d\n", existingValue, newValue))
	if newValue == existingValue {
		// really shouldn't be getting this
		fmt.Println("No change, not checking for errors")
		return // no change
	} else {
		fmt.Println("Yes change, checking for errors")
		if existingValue != 0 {

			// check for remaining errors for existingValue
			// a full check is needed because the error determination for items along the row/column
			// is not singly determined by the removal of the current position.
			// eg. removing this 5
			// . . 5 | . . 5 | . . .
			// . . . | . . . | . . .
			// . . 5 | . . . | . . .
			// to a new state
			// . . 5 | . . . | . . .
			// . . . | . . . | . . .
			// . . 5 | . . . | . . .
			// does not mean we can clear the error value of the 'top' 5
			// it has both a same-square error and a same-column error
			// so we need to check the whole board for errors with 5's
			// and if a row/column/square has more than one 5
			// then every user entered 5 in this row/column/square should be marked as an error
			valueToCheck := existingValue
			fmt.Print("Need to check what errors are left after removing: ")
			fmt.Println(valueToCheck)

			// clear board of errors for this number
			for boardIndex := 0; boardIndex < (boardWidth * boardHeight); boardIndex++ {
				cPos := BoardPosition{vert: (boardIndex / boardWidth), horiz: (boardIndex % boardWidth)}
				if board.board.userEntries.GetValueAt(cPos) == valueToCheck || board.board.userEntries.GetValueAt(cPos) == 0 {
					board.board.invalidEntries.SetValueAt(cPos, false)
				}
			}

			// check all rows
			for rowIndex := 0; rowIndex < boardHeight; rowIndex++ {
				digitFoundInThisRow := false
				for colIndex := 0; colIndex < boardWidth; colIndex++ {
					cPos := BoardPosition{vert: rowIndex, horiz: colIndex}
					if board.board.GetAnyValueAt(cPos) == valueToCheck {
						if digitFoundInThisRow {
							// fill all user entered valueToCheck instances in this row with errors
							for colIndex2 := 0; colIndex2 < boardWidth; colIndex2++ {
								cPos2 := BoardPosition{vert: at.vert, horiz: colIndex2}
								if board.board.userEntries.GetValueAt(cPos2) == valueToCheck {
									board.board.invalidEntries.SetValueAt(cPos2, true)
								}
							}
							break
						}
						digitFoundInThisRow = true
					}
				}
			}

			// check all columns
			for colIndex := 0; colIndex < boardHeight; colIndex++ {
				digitFoundInThisColumn := false
				for rowIndex := 0; rowIndex < boardWidth; rowIndex++ {
					cPos := BoardPosition{vert: rowIndex, horiz: colIndex}
					if board.board.GetAnyValueAt(cPos) == valueToCheck {
						if digitFoundInThisColumn {
							// fill all user entered valueToCheck instances in this column with errors
							for rowIndex2 := 0; rowIndex2 < boardWidth; rowIndex2++ {
								cPos2 := BoardPosition{vert: rowIndex2, horiz: at.horiz}
								if board.board.userEntries.GetValueAt(cPos2) == valueToCheck {
									board.board.invalidEntries.SetValueAt(cPos2, true) // it's fine if we write existing true to true again
								}
							}
							break
						}
						digitFoundInThisColumn = true
					}
				}
			}

			// check all squares
			for squareNumber := 0; squareNumber < 9; squareNumber++ {
				squareOffsetHoriz := squareNumber / 3
				squareOffsetVert := squareNumber % 3
				digitFoundInThisSquare := false
				for squareIndex := 0; squareIndex < 9; squareIndex++ {
					cPos := BoardPosition{
						horiz: (squareOffsetHoriz * 3) + (squareIndex % 3),
						vert:  (squareOffsetVert * 3) + (squareIndex / 3),
					}
					fmt.Print("Checking position ")
					fmt.Print(cPos)
					fmt.Println("...")
					if board.board.GetAnyValueAt(cPos) == valueToCheck {
						if digitFoundInThisSquare {
							fmt.Print("This square has an error. Setting all user entries in this square for this digit to error")
							for squareIndex2 := 0; squareIndex2 < 9; squareIndex2++ {
								cPos2 := BoardPosition{
									horiz: (squareOffsetHoriz * 3) + (squareIndex2 % 3),
									vert:  (squareOffsetVert * 3) + (squareIndex2 / 3),
								}
								if board.board.userEntries.GetValueAt(cPos2) == valueToCheck {
									board.board.invalidEntries.SetValueAt(cPos2, true)
								}
							}
							break
						}
						digitFoundInThisSquare = true
					}
				}
			}
		}
		if newValue != 0 {

			// (slim) check for new errors for newValue
			// eg. adding this 5
			// could only *add* an error among its row, column, or square
			// so we don't need to check the whole board
			valueToCheck := newValue
			fmt.Print("Need to check new errors are introduced by adding: ")
			fmt.Println(valueToCheck)

			// row
			fmt.Print("Checking row ")
			fmt.Print(at.vert)
			fmt.Println("...")

			for colIndex := 0; colIndex < boardWidth; colIndex++ {
				if colIndex == at.horiz {
					fmt.Print("Skipping row item number ")
					fmt.Print(colIndex)
					fmt.Print(" because we know ")
					fmt.Print(valueToCheck)
					fmt.Println(" is there.")
					continue
				} // skip self
				fmt.Print("Checking row item number ")
				fmt.Print(colIndex)
				fmt.Println("...")
				cPos := BoardPosition{vert: at.vert, horiz: colIndex}
				if board.board.GetAnyValueAt(cPos) == valueToCheck {
					fmt.Print("This row has an error. Setting all user entries in this row for this digit to error")
					fmt.Print(colIndex)
					fmt.Println("...")
					for colIndex2 := 0; colIndex2 < boardWidth; colIndex2++ {
						cPos2 := BoardPosition{vert: at.vert, horiz: colIndex2}
						if board.board.userEntries.GetValueAt(cPos2) == valueToCheck {
							board.board.invalidEntries.SetValueAt(cPos2, true) // it's fine if we write existing true to true again
						}
					}
					break
				}
			}
			// could be much more efficient if we keep a few counts of 'how many of digit X is in this row'

			// col
			fmt.Print("Checking column ")
			fmt.Print(at.horiz)
			fmt.Println("...")
			for rowIndex := 0; rowIndex < boardWidth; rowIndex++ {
				if rowIndex == at.vert {
					fmt.Print("Skipping column item number ")
					fmt.Print(rowIndex)
					fmt.Print(" because we know ")
					fmt.Print(valueToCheck)
					fmt.Println(" is there.")
					continue
				} // skip self
				fmt.Print("Checking column item number ")
				fmt.Print(rowIndex)
				fmt.Println("...")
				cPos := BoardPosition{vert: rowIndex, horiz: at.horiz}
				if board.board.GetAnyValueAt(cPos) == valueToCheck {
					fmt.Print("This column has an error. Setting all user entries in this column for this digit to error")
					fmt.Print(rowIndex)
					fmt.Println("...")
					for rowIndex2 := 0; rowIndex2 < boardWidth; rowIndex2++ {
						cPos2 := BoardPosition{vert: rowIndex2, horiz: at.horiz}
						if board.board.userEntries.GetValueAt(cPos2) == valueToCheck {
							board.board.invalidEntries.SetValueAt(cPos2, true) // it's fine if we write existing true to true again
						}
					}
					break
				}
			}

			// square
			squareOffsetHoriz := at.horiz / 3
			squareOffsetVert := at.vert / 3
			fmt.Print("Checking square ")
			fmt.Print((squareOffsetVert * 3) + squareOffsetHoriz)
			fmt.Println("...")
			for squareIndex := 0; squareIndex < 9; squareIndex++ {
				cPos := BoardPosition{
					horiz: (squareOffsetHoriz * 3) + (squareIndex % 3),
					vert:  (squareOffsetVert * 3) + (squareIndex / 3),
				}
				// skip self
				if at == cPos {
					fmt.Print("Skipping position ")
					fmt.Print(cPos)
					fmt.Print(" because we know ")
					fmt.Print(valueToCheck)
					fmt.Println(" is there.")
					continue
				}
				fmt.Print("Checking position ")
				fmt.Print(cPos)
				fmt.Println("...")
				if board.board.GetAnyValueAt(cPos) == valueToCheck {
					fmt.Print("This square has an error. Setting all user entries in this square for this digit to error")
					for squareIndex2 := 0; squareIndex2 < 9; squareIndex2++ {
						cPos2 := BoardPosition{
							horiz: (squareOffsetHoriz * 3) + (squareIndex2 % 3),
							vert:  (squareOffsetVert * 3) + (squareIndex2 / 3),
						}
						if board.board.userEntries.GetValueAt(cPos2) == valueToCheck {
							board.board.invalidEntries.SetValueAt(cPos2, true)
						}
					}
					break
				}
			}
		}
	}
}

func (board *SudokuBoardInteractionState) SetNumberAtCursor(input int) {
	// check if givenBoard is blank there; we can't edit 'set' positions
	if board.board.givenBoard.GetValueAt(board.cursor) != 0 {
		return
	}
	// check if it's the same as the cursor. don't do any logic for 'change 3 to 3'
	if board.board.userEntries.GetValueAt(board.cursor) == input {
		return
	}

	prevInput := board.board.userEntries.GetValueAt(board.cursor)
	board.board.userEntries.setNumberAtPos(input, board.cursor, board.toggle)
	board.CheckForErrorsForMove(prevInput, input, board.cursor)

	numberIndex := input - 1
	if numberIndex >= 0 {
		board.board.numberOfDigitsTotal[numberIndex]++
	}

	prevNumberIndex := prevInput - 1
	if prevNumberIndex >= 0 {
		board.board.numberOfDigitsTotal[prevNumberIndex]--
	}
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
