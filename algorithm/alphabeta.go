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
