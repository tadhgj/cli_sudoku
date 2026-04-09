package main

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/jhunters/sudoku"
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
	toggle         bool // should entering 5 in same square clear square
}

func internalGenerateBlankSudokuBoardErrors() BoardGridErrors {
	return BoardGridErrors{} // fills with falses. false treated as no error for this square/cell
}

func internalGenerateBlankSudokuBoard() BoardGrid {
	return BoardGrid{} // fills with 0's; 0 treated as blank
}

func internalGenerateBlankSudokuBoardNotes() BoardGridNotes {
	return BoardGridNotes{} // fills grid with arrays of length 'numofpossibledigits' with 0's; 0 treated as no notes
}

func internalGenerateSudokuBoard(difficulty Difficulty) BoardGrid {
	// I am outsourcing my sudoku board generation because I fear that learning how to create the boards will take the fun out of solving the boards

	var misscount int

	switch difficulty {
	case testeasy:
		misscount = 3
	case easy:
		misscount = 20
	case normal:
		misscount = 30
	case hard:
		misscount = 40
	}
	// note that these numbers are not a very good determinor of the 'difficulty'
	// of a puzzle but are much better than nothing for the time being

	// random stub
	// misscount = 50

	// misscount limit is the max cell value squared, minus 5
	// so with a max cell value of 9, the max misscount value is 81-5 which is 76
	// the minimum is not enforced
	// a miss count of 76 leaves literally 5 filled squares on the board
	// so you'd probably need at least 9 filled squares to ensure a unique solution
	// default value was 40
	// will change this value with difficulty selection

	// NewSudokuGenX(maxcellvalue, 'misscount')
	sg, _ := sudoku.NewSudokuGenX(9, misscount)
	// GenSudoku -> result, answer, err
	result, _, _ := sg.GenSudoku()

	// because I declared the board type to be [9][9]int and the generator returns [][]int
	var resultCopy [9][9]int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			resultCopy[i][j] = result[i][j]
		}
	}

	return BoardGrid{board: resultCopy}

	// // placeholder until true generation
	// return BoardGrid{
	// 	board: [9][9]int{
	// 		{5, 3, 0, 0, 7, 0, 0, 0, 0},
	// 		{6, 0, 0, 1, 9, 5, 0, 0, 0},
	// 		{0, 9, 8, 0, 0, 0, 0, 6, 0},
	// 		{8, 0, 0, 0, 6, 0, 0, 0, 3},
	// 		{4, 0, 0, 8, 0, 3, 0, 0, 1},
	// 		{7, 0, 0, 0, 2, 0, 0, 0, 6},
	// 		{0, 6, 0, 0, 0, 0, 2, 8, 0},
	// 		{0, 0, 0, 4, 1, 9, 0, 0, 5},
	// 		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	// 	},
	// }
}

func GenerateSudokuBoard(difficulty Difficulty) SudokuBoard {
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

func GenerateSudokuBoardState(difficulty Difficulty) SudokuBoardInteractionState {
	return SudokuBoardInteractionState{
		board:          GenerateSudokuBoard(difficulty),
		selectedNumber: -1,
		cursor: BoardPosition{
			horiz: 0,
			vert:  0,
		},
	}
}

func (s SudokuGameWrapperState) RenderSudokuBoardState(styles styles) string {

	// in a console, some of the most basic and available styling you have available includes:
	// underline
	// highlighting
	// a 'light' text color

	state := s.boardInteraction
	board := state.board
	cursor := state.cursor

	// get number under cursor
	cursorSelectedNum := board.GetAnyValueAt(cursor)
	if cursorSelectedNum == blankDigitInt {
		cursorSelectedNum = -1 // set to -1 so it doesn't match empty value of 0
	}

	var result string

	var heightOffset int
	var widthOffset int
	var rowStartOffset int
	var columnStartOffset int
	// IF we have 'center of universe' setting on...
	// calculate offset based off character
	// heightOffset = 10
	// widthOffset = 10

	// x123|456|789x
	// x 1 2 3 | 4 5 6 | 7 8 9 x
	// x-------+
	// cursor.horiz = 8 widthOffset = 0
	// cursor.horiz = 7 widthOffset = 2
	// 25? wide?

	//     ,
	//  123 cursor.vert = 0 (heightOffset = 10)
	//  123 cursor.vert = 1 (heightOffset = 9)
	//  123 cursor.vert = 2 (heightOffset = 8)
	//  ---
	//  123 (heightOffset = 6)
	//  123 (heightOffset = 5)
	//  123 (heightOffset = 4)
	//  ---
	//  123 (heightOffset = 2)
	//  123 (heightOffset = 1)
	//  123 cursor.vert = 8 (heightOffset = 0)
	//     '
	// 13 tall

	if s.userOptions.selectedRenderingStyleIndex == int(centeredCursor) {
		heightOffset = (10 - cursor.vert) - (cursor.vert / 3)
		widthOffset = (10-(cursor.horiz))*2 - (cursor.horiz/3)*2
	} else if s.userOptions.selectedRenderingStyleIndex == int(infiniteBoard) {
		heightOffset = 0
		widthOffset = 0

		columnStartOffset = (6 + ((cursor.vert / 3) * 3)) % boardHeight
		rowStartOffset = (6 + ((cursor.horiz / 3) * 3)) % boardWidth
	}
	// else if selrenderstyle is the infinite centered one, set the height/width offset fixed and add the startingRow/startingColumn offsets

	var heightOffsetString string
	var widthOffsetString string
	var debugwidthOffsetString string
	// calc height offset
	for i := 0; i < heightOffset; i++ {
		// heightOffsetString += "↓\n"
		heightOffsetString += "\n"
	}
	// calc width offset
	for i := 0; i < widthOffset; i++ {
		if i%2 == 0 {
			// debugwidthOffsetString += "→"
			debugwidthOffsetString += " "
		}
		if i%2 == 1 {
			debugwidthOffsetString += " "
		}
		widthOffsetString += " "
	}

	result += heightOffsetString

	// top decoration

	result += debugwidthOffsetString
	result += "       ,       ,       \n"

	// for each row...
	for ri := 0; ri < boardHeight; ri++ {

		i := (columnStartOffset + ri) % boardHeight

		isSameRowAsCursor := i == cursor.vert

		// apply width offset
		result += widthOffsetString
		// spacer at beginning of row
		result += " "

		// for each item in row...
		for rj := 0; rj < boardWidth; rj++ {

			j := (rowStartOffset + rj) % boardHeight

			cPos := BoardPosition{horiz: j, vert: i}

			cNumValue := board.GetAnyValueAt(cPos)
			isEmpty := cNumValue == blankDigitInt
			var currentChar string
			if isEmpty {
				currentChar = " "
			} else {
				currentChar = strconv.Itoa(cNumValue)
			}

			isSameColumnAsCursor := cPos.horiz == cursor.horiz
			isSameSquareAsCursor := (cPos.horiz/3 == cursor.horiz/3) && (cPos.vert/3 == cursor.vert/3)

			isUnderCursor := isSameColumnAsCursor && isSameRowAsCursor
			isErroneousUserNumber := board.invalidEntries.board[cPos.vert][cPos.horiz]

			// TODO:
			// IF DARK MODE

			// look at given board
			currentNumForegroundStyle := styles.tertiary
			currentNumInvertedForegroundStyle := styles.tertiary
			if isEmpty {
				// blank
			} else if board.givenBoard.GetValueAt(cPos) != 0 {
				// 'locked' number
				currentNumForegroundStyle = styles.primary
				currentNumInvertedForegroundStyle = styles.primaryInvert
			} else if board.userEntries.GetValueAt(cPos) != 0 {
				// user generated number
				currentNumForegroundStyle = styles.accent
				currentNumInvertedForegroundStyle = styles.accentInvert
				// currentNumIsUserEntered = true
			}

			// check if special highlights need to be applied to this number
			// cases: (in order of precedence)
			// userHighlight
			// numberSelected
			if isUnderCursor {
				if isErroneousUserNumber {
					result += styles.invertHighlight.Render(styles.errorInvert.Render(currentChar))
				} else {
					result += styles.invertHighlight.Render(currentNumInvertedForegroundStyle.Render(currentChar))
				}
			} else if cNumValue == cursorSelectedNum {
				matchingCausingError := (isErroneousUserNumber || (isSameColumnAsCursor || isSameRowAsCursor) || isSameSquareAsCursor)
				if matchingCausingError {
					if isErroneousUserNumber {
						result += styles.errorHighlightUser.Render(currentChar)
						// result += style.Background(lipgloss.Color("#f00")).Foreground(lipgloss.Color("#900")).Render(currentChar)
					} else {
						result += styles.errorHighlightGiven.Render(currentChar)
					}
				} else {
					result += styles.darkHighlight.Render(currentChar)
				}
			} else if isSameColumnAsCursor || isSameRowAsCursor || isSameSquareAsCursor {
				if isErroneousUserNumber {
					result += styles.lightHighlight.Render(styles.errorForeground.Render(currentChar))
				} else {
					result += styles.lightHighlight.Render(currentNumForegroundStyle.Render(currentChar))
				}
			} else {
				if isErroneousUserNumber {
					result += styles.errorForeground.Render(currentChar)
				} else {
					result += currentNumForegroundStyle.Render(currentChar)
				}
			}

			// spacer
			if rj != boardWidth-1 {
				if j%3 == 2 {
					result += " | "
				} else if isSameRowAsCursor || isSameSquareAsCursor {
					result += styles.lightHighlight.Render(" ")
				} else {
					result += " "
				}
			}
		}
		result += "\n"

		// spacer
		if ri != boardHeight-1 && ri%3 == 2 {
			result += widthOffsetString
			result += "-------+-------+-------\n"
		}
	}

	// bottom
	result += widthOffsetString
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
			if cursorSelectedNum == number {
				numbStr = lipgloss.NewStyle().Underline(true).Render(numbStr)
			}
			thisNumberFinished := board.numberOfDigitsTotal[numberIndex] >= 9
			numbStrRendered := numbStr
			if thisNumberFinished {
				numbStrRendered = styles.tertiary.Render(numbStr)
			} else {
				numbStrRendered = styles.secondary.Render(numbStr)
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
func (s *SudokuGameWrapperState) SetCursor(pos BoardPosition) {
	s.boardInteraction.cursor.setCursor(pos, s.userOptions.loopCursorAroundEdges)
}

func (board *BoardGrid) setNumberAtPos(input int, pos BoardPosition, toggle bool) bool {
	// check if input number is valid
	if input < perSquareMin || input > perSquareMax {
		return false
	}

	// check if already set to this value
	if board.GetValueAt(pos) == input {
		return false
	}

	// check if cursor position is valid
	// TODO: do something different based on toggle setting
	// if input != 0 && toggle && board.GetValueAt(pos) == input {
	// 	board.SetValueAt(pos, input)
	// } else {
	// 	board.SetValueAt(pos, input)
	// }

	board.SetValueAt(pos, input)
	return true
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

			// clear error at position (since we have removed the number, the script below will not remove it)
			board.board.invalidEntries.SetValueAt(at, false)

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

func (s *SudokuGameWrapperState) checkForWin() {
	// if user errors is empty and numbers are full call it a day
	board := s.boardInteraction
	anyError := false
	for rowIndex := 0; rowIndex < boardHeight; rowIndex++ {
		for colIndex := 0; colIndex < boardWidth; colIndex++ {

			if board.board.invalidEntries.board[rowIndex][colIndex] {
				anyError = true
			}
			if anyError {
				break
			}

			if board.board.GetAnyValueAt(BoardPosition{horiz: colIndex, vert: rowIndex}) == 0 {
				anyError = true
			}
			if anyError {
				break
			}
		}
	}

	if !anyError {
		// call win stuff
		s.shownPage = winPage
	}
}

// returns true if number at cursor has changed value, necessitating a check for new error states
func (s *SudokuGameWrapperState) SetNumberAtCursor(input int) {
	// check if givenBoard is blank there; we can't edit 'set' positions
	if s.boardInteraction.board.givenBoard.GetValueAt(s.boardInteraction.cursor) != 0 {
		return
	}

	// check if it's the same as the cursor. don't do any logic for 'change 3 to 3'
	if s.boardInteraction.board.userEntries.GetValueAt(s.boardInteraction.cursor) == input {
		return
	}

	prevInput := s.boardInteraction.board.userEntries.GetValueAt(s.boardInteraction.cursor)
	wasChanged := s.boardInteraction.board.userEntries.setNumberAtPos(input, s.boardInteraction.cursor, s.boardInteraction.toggle)
	if !wasChanged {
		return
	}

	s.boardInteraction.CheckForErrorsForMove(prevInput, input, s.boardInteraction.cursor)

	numberIndex := input - 1
	if numberIndex >= 0 {
		s.boardInteraction.board.numberOfDigitsTotal[numberIndex]++
	}

	prevNumberIndex := prevInput - 1
	if prevNumberIndex >= 0 {
		s.boardInteraction.board.numberOfDigitsTotal[prevNumberIndex]--
	}

	// check for win state
	s.checkForWin()
}

func (s *SudokuGameWrapperState) MoveCursorLeft() {
	s.SetCursor(BoardPosition{
		horiz: s.boardInteraction.cursor.horiz - 1,
		vert:  s.boardInteraction.cursor.vert,
	},
	)
}

func (s *SudokuGameWrapperState) MoveCursorRight() {
	s.SetCursor(BoardPosition{
		horiz: s.boardInteraction.cursor.horiz + 1,
		vert:  s.boardInteraction.cursor.vert,
	})
}

func (s *SudokuGameWrapperState) MoveCursorUp() {
	s.SetCursor(BoardPosition{
		horiz: s.boardInteraction.cursor.horiz,
		vert:  s.boardInteraction.cursor.vert - 1,
	})
}

func (s *SudokuGameWrapperState) MoveCursorDown() {
	s.SetCursor(BoardPosition{
		horiz: s.boardInteraction.cursor.horiz,
		vert:  s.boardInteraction.cursor.vert + 1,
	})
}
