package main

import (
	"fmt"
	"math/rand"
)

var winningPositions = [][3]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

type TicTacToe struct {
	Board []int
}

func NewTicTacToeBoard() Board {
	board := make([]int, 9)
	for idx := range board {
		board[idx] = Empty
	}
	return &TicTacToe{
		Board: board,
	}
}

func (t *TicTacToe) Clone() Board {
	board := make([]int, 9)
	copy(board, t.Board)
	return &TicTacToe{
		Board: board,
	}
}

func (t *TicTacToe) CheckWinner(player int) bool {

	for _, places := range winningPositions {
		n := 0
		for _, placeIndex := range places {
			if t.Board[placeIndex] == player {
				n += 1
			}
		}
		if n == 3 {
			return true
		}
	}
	return false
}

func (t *TicTacToe) CheckStatus() int {
	if len(t.EmptyPlaces()) == 0 {
		return Draw
	}
	if t.CheckWinner(O) {
		return O
	}
	if t.CheckWinner(X) {
		return X
	}
	return InProgress
}

func (t *TicTacToe) PerformMove(move, player int) {
	t.Board[move] = player
}

func (t *TicTacToe) EmptyPlaces() []int {
	emptyPlaces := make([]int, 0)
	for index, place := range t.Board {
		if place == Empty {
			emptyPlaces = append(emptyPlaces, index)
		}
	}
	return emptyPlaces
}

func (t *TicTacToe) Print() {
	n := 1
	for _, place := range t.Board {
		if place == O {
			fmt.Print("O ")
		} else if place == X {
			fmt.Print("X ")
		} else {
			fmt.Print("- ")
		}
		if n%3 == 0 {
			fmt.Println()
		}
		n++
	}
}

func (t *TicTacToe) RandomPlay(player int) {
	moves := t.EmptyPlaces()
	s := len(moves)
	if s == 0 {
		return
	}
	t.PerformMove(moves[rand.Intn(s)], player)
}
