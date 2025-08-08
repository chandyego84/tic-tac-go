package main

type GameState struct {
	GameStarted bool
	GameOver bool
	Board [9]string
	PlayerTurn string
}

// Check if a move is valid -- true if valid, false otherwise.
func (gs *GameState) validateMove(moveIndex int) bool {
	return !gs.GameOver && moveIndex >= 0 && moveIndex < 9 && gs.Board[moveIndex] == ""
}

func (gs *GameState) updateCurrentPlayer() {
	if gs.PlayerTurn == "X" {
		gs.PlayerTurn = "O"
	} else {
		gs.PlayerTurn = "X"
	}
}

// Action on game is made
func (gs *GameState) Step(moveIndex int) {
	if !gs.validateMove(moveIndex) {
        return
    }

	gs.Board[moveIndex] = gs.PlayerTurn
	if gs.isOver() {
		gs.GameOver = true
		return
	}

	gs.updateCurrentPlayer()
}

func (gs *GameState) isOver() bool {
	return gs.checkWin("X") || gs.checkWin("O") || gs.checkDraw()
}

func (gs *GameState) checkDraw() bool {
    for _, cell := range gs.Board {
        if cell == "" {
            return false
        }
    }
    return true
}

func (gs *GameState) checkWin(p string) bool {
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
		row := b[3*r : 3*r+3]
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