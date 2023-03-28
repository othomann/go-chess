package algorithm

import (
	"math"
	"sort"

	"github.com/othomann/go-chess"
)

func AlphaBeta(game *chess.Game, depth, alpha, beta int, maximizingPlayer bool) int {
	// check if we've reached the maximum depth or if the game is over
	if game.Method() == chess.Checkmate {
		return math.MaxInt32
	}
	if depth == 0 {
		return game.Evaluate()
	}

	// generate all possible moves
	moves := game.ValidMoves()
	sort.Sort(chess.MoveSlice(moves))

	// if the maximizing player is playing
	if maximizingPlayer {
		bestScore := math.MinInt32
		for _, move := range moves {
			// apply the move to the game state
			game.Move(move)

			// calculate the score of this move using alpha-beta pruning
			score := AlphaBeta(game, depth-1, alpha, beta, false)

			// update the best score and alpha value
			bestScore = max(bestScore, score)
			alpha = max(alpha, bestScore)

			// undo the move
			game.UndoMove()

			// check if we can prune the remaining moves
			if beta <= alpha {
				break
			}
		}
		return bestScore
	} else { // if the minimizing player is playing
		bestScore := math.MaxInt32
		for _, move := range moves {
			// apply the move to the game state
			game.Move(move)

			// calculate the score of this move using alpha-beta pruning
			score := AlphaBeta(game, depth-1, alpha, beta, true)

			// update the best score and beta value
			bestScore = min(bestScore, score)
			beta = min(beta, bestScore)

			// undo the move
			game.UndoMove()

			// check if we can prune the remaining moves
			if beta <= alpha {
				break
			}
		}
		return bestScore
	}
}

// helper function to get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Minimax(g *chess.Game, depth int, alpha int, beta int) (int, *chess.Move) {
	if depth == 0 {
		if g.Method() == chess.Checkmate {
			return math.MaxInt32, &chess.Move{}
		}
		return g.Evaluate(), &chess.Move{}
	} else if g.Method() == chess.Checkmate {
		return math.MaxInt32, &chess.Move{}
	} else if g.Method() == chess.Stalemate {
		return 0, &chess.Move{}
	}
	bestScore := -99999
	var bestMove *chess.Move

	validMoves := g.ValidMoves()
	sort.Sort(chess.MoveSlice(validMoves))
	for _, move := range validMoves {
		g.Move(move)
		score, _ := Minimax(g, depth-1, alpha, beta)
		g.UndoMove()

		if g.Position().Turn() == chess.White && score > bestScore {
			bestScore = score
			bestMove = move
			if score > alpha {
				alpha = score
			}
		} else if g.Position().Turn() == chess.Black && score < bestScore {
			bestScore = score
			bestMove = move
			if score < beta {
				beta = score
			}
		}

		if alpha >= beta {
			break
		}
	}

	return bestScore, bestMove
}
