package algorithm

import (
	"math"
	"sort"

	"github.com/othomann/go-chess"
)

func AlphaBeta(game *chess.Game, depth, alpha, beta int, maximizingPlayer bool) (int, error) {
	// check if we've reached the maximum depth or if the game is over
	if game.Method() == chess.Checkmate {
		return math.MaxInt32, nil
	}
	if depth == 0 {
		return game.Evaluate(), nil
	}

	// generate all possible moves
	moves := game.ValidMoves()
	sort.Sort(chess.MoveSlice(moves))

	// if the maximizing player is playing
	if maximizingPlayer {
		bestScore := math.MinInt32
		for _, move := range moves {
			// apply the move to the game state
			err := game.Move(move)
			if err != nil {
				return 0, err
			}

			// calculate the score of this move using alpha-beta pruning
			score, err := AlphaBeta(game, depth-1, alpha, beta, false)
			if err != nil {
				return 0, err
			}

			// update the best score and alpha value
			bestScore = max(bestScore, score)
			alpha = max(alpha, bestScore)

			// undo the move
			err = game.UndoMove()
			if err != nil {
				return 0, err
			}

			// check if we can prune the remaining moves
			if beta <= alpha {
				break
			}
		}
		return bestScore, nil
	} else { // if the minimizing player is playing
		bestScore := math.MaxInt32
		for _, move := range moves {
			// apply the move to the game state
			err := game.Move(move)
			if err != nil {
				return 0, err
			}

			// calculate the score of this move using alpha-beta pruning
			score, err := AlphaBeta(game, depth-1, alpha, beta, true)
			if err != nil {
				return 0, err
			}

			// update the best score and beta value
			bestScore = min(bestScore, score)
			beta = min(beta, bestScore)

			// undo the move
			err = game.UndoMove()
			if err != nil {
				return 0, err
			}

			// check if we can prune the remaining moves
			if beta <= alpha {
				break
			}
		}
		return bestScore, nil
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

func Minimax(g *chess.Game, depth int, alpha int, beta int) (int, *chess.Move, error) {
	if depth == 0 {
		if g.Method() == chess.Checkmate {
			return math.MaxInt32, &chess.Move{}, nil
		}
		return g.Evaluate(), &chess.Move{}, nil
	} else if g.Method() == chess.Checkmate {
		return math.MaxInt32, &chess.Move{}, nil
	} else if g.Method() == chess.Stalemate {
		return 0, &chess.Move{}, nil
	}
	bestScore := -99999
	var bestMove *chess.Move

	validMoves := g.ValidMoves()
	sort.Sort(chess.MoveSlice(validMoves))
	for _, move := range validMoves {
		err := g.Move(move)
		if err != nil {
			return 0, nil, err
		}
		score, _, err := Minimax(g, depth-1, alpha, beta)
		if err != nil {
			return 0, nil, err
		}
		err = g.UndoMove()
		if err != nil {
			return 0, nil, err
		}

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

	return bestScore, bestMove, nil
}
