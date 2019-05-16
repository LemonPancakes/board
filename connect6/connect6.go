package connect6

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Move struct {
	x, y int
}

func (m *Move) Add(a *Move) {
	m.x += a.x
	m.y += a.y
}

func (m *Move) Sub(a *Move) {
	m.x -= a.x
	m.y -= a.y
}

type Connect6 struct {
	board         [19][19]int
	CurrentPlayer int
	firstMove     bool
	Finished      bool
}

func (game *Connect6) NewGame() {
	var tmp [19][19]int
	game.board = tmp
	game.CurrentPlayer = 1
	game.firstMove = false
	game.Finished = false
}

func (game *Connect6) MakeMove(m Move) (int, error) {
	if game.Finished {
		return -1, errors.New("Game is already finished")
	}

	if game.board[m.x][m.y] != 0 {
		return -1, errors.New("Space is already taken")
	}

	player := game.CurrentPlayer
	game.board[m.x][m.y] = player
	if game.CheckWin(m) {
		game.Finished = true
		return player, nil
	}

	if !game.firstMove {
		game.CurrentPlayer = 3 - game.CurrentPlayer
	}
	game.firstMove = !game.firstMove
	return player, nil
}

func (game *Connect6) Print() {
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			switch curr := game.board[i][j]; curr {
			case 0:
				fmt.Print(".")
			case 1:
				fmt.Print("X")
			case 2:
				fmt.Print("O")
			default:
				fmt.Println(errors.New("WTF???"))
			}
		}
		fmt.Println("")
	}
}

func (game *Connect6) CheckWin(m Move) bool {
	var directions = [4]Move{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, dir := range directions {
		currLength := 1
		tmp := m
		tmp.Add(&dir)
		// Check in forward direction first
		for currLength < 6 && tmp.x < 19 && tmp.x >= 0 && tmp.y < 19 && tmp.y >= 0 && game.board[tmp.x][tmp.y] == game.CurrentPlayer {
			tmp.Add(&dir)
			currLength++
		}

		tmp = m
		tmp.Sub(&dir)
		// Check in opposite direction
		for currLength < 6 && tmp.x < 19 && tmp.x >= 0 && tmp.y < 19 && tmp.y >= 0 && game.board[tmp.x][tmp.y] == game.CurrentPlayer {
			tmp.Sub(&dir)
			currLength++
		}

		if currLength == 6 {
			return true
		}
	}
	return false
}

func (game *Connect6) GetState() string {
	gameState := ""

	gameState += strconv.Itoa(game.CurrentPlayer) + ","

	if game.firstMove {
		gameState += "1"
	} else {
		gameState += "0"
	}

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			gameState += "," + strconv.Itoa(game.board[i][j])
		}
	}

	return gameState
}

func ParseMove(s string) Move {
	inputs := strings.Split(strings.TrimSpace(s), ",")
	row, _ := strconv.Atoi(inputs[0])
	col, _ := strconv.Atoi(inputs[1])
	fmt.Printf("Move at (%d, %d)\n", row, col)
	return Move{row, col}
}
