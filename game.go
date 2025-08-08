package main

type GameState struct {
	GameStarted bool
	GameWon bool
	Board [9]string
	PlayerTurn string
}

func nextPlayer(current string) string {
	if (current == "X") {
		return "O"
	}
	return "X"
}