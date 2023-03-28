package algorithm

import "github.com/othomann/go-chess"

const (
	INF = 1000000
	MAX = 1
	MIN = 0
)

func AlphaBeta(g *chess.Game, depth int, alpha int, beta int, maximizingPlayer bool) int {
	if depth == 0 {
		return g.Evaluate()
	}

	moves := g.ValidMoves()
	if len(moves) == 0 {
		switch g.Method() {
		case chess.InCheck:
			fallthrough
		case chess.Checkmate:
			return -INF
		case chess.Stalemate:
			fallthrough
		case chess.DrawOffer:
			fallthrough
		case chess.FiftyMoveRule:
			fallthrough
		case chess.FivefoldRepetition:
			fallthrough
		case chess.SeventyFiveMoveRule:
			fallthrough
		case chess.InsufficientMaterial:
			return 0
		}
	}

	if maximizingPlayer {
		bestValue := -INF
		for _, move := range moves {
			g.Move(move)
			value := AlphaBeta(g, depth-1, alpha, beta, !maximizingPlayer)
			g.UndoMove()
			bestValue = max(bestValue, value)
			alpha = max(alpha, bestValue)
			if beta <= alpha {
				break
			}
		}
		return bestValue
	} else {
		bestValue := INF
		for _, move := range moves {
			g.Move(move)
			value := AlphaBeta(g, depth-1, alpha, beta, !maximizingPlayer)
			g.UndoMove()
			bestValue = min(bestValue, value)
			beta = min(beta, bestValue)
			if beta <= alpha {
				break
			}
		}
		return bestValue
	}
}

func Minimax(g *chess.Game, depth int, alpha int, beta int) (int, *chess.Move) {
	if depth == 0 || g.Method() == chess.Checkmate {
		return INF, &chess.Move{}
	}

	bestScore := -99999
	var bestMove *chess.Move

	for _, move := range g.ValidMoves() {
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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
