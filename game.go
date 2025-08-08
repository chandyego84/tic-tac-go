package main

type GameState struct {
	GameStarted bool
	GameOver bool
	Board [9]string
	PlayerTurn string
}

func (gs *GameState) updateCurrentPlayer(current string) string {
	if (current == "X") {
		return "O"
	}
	return "X"
}

// Validate move
func (gs *GameState) validateMove(moveIndex int) bool {
	valid := true

	if (gs.GameOver || moveIndex < 0 || moveIndex > 8 || gs.Board[moveIndex] != "") { valid = false }

	return valid
}

// Check for win
func (gs *GameState) checkWin() bool {
	p := gs.PlayerTurn
	b := gs.Board

	allEqual := func(slice []string, player string) bool {
		for _, v := range slice {
			if v != player {
				return false
			}
		} 

		return true
	}

	// horizontals
	for r := 0; r < 3; r++ {
		row := b[3*r : 3*r+r]
		if allEqual(row, p) {
			return true
		}
	}

	// verticals
	for c := 0; c < 3; c++ {
		col := []string {b[c], b[c + 3], b[c + 6]}
		if allEqual(col, p) {
			return true
		}
	}

	// diags
	negDiag := []string {b[0], b[4], b[8]}
	posDiag := []string {b[6], b[4], b[2]}
	if allEqual(negDiag, p) || allEqual(posDiag, p) {
		return true
	}

	return false
}