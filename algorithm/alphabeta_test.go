package algorithm

import (
	"math"
	"testing"

	"github.com/othomann/go-chess"
)

func TestCheckmate(t *testing.T) {
	fenStr := "rn1qkbnr/pbpp1ppp/1p6/4p3/2B1P3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 1"
	fen, err := chess.FEN(fenStr)
	if err != nil {
		t.Fatal(err)
	}
	g := chess.NewGame(fen)
	score, move, err := Minimax(g, 2, math.MinInt32, math.MaxInt32)
	if err != nil {
		t.Fatal(err)
	}
	if score != 2147483647 {
		t.Fatalf("Wrong score; expected %d, but got %d", -100, score)
	}
	if move.String() != "f3f7" {
		t.Fatalf("Wrong score; expected %s, but got %s", "f3f7", move.String())
	}
}

func TestCheckmate2(t *testing.T) {
	fenStr := "rn1qkbnr/pbpp1ppp/1p6/4p3/2B1P3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 1"
	fen, err := chess.FEN(fenStr)
	if err != nil {
		t.Fatal(err)
	}
	g := chess.NewGame(fen)
	score, err := AlphaBeta(g, 2, math.MinInt32, math.MaxInt32, true)
	if err != nil {
		t.Fatal(err)
	}
	if score != 2147483647 {
		t.Fatalf("Wrong score; expected %d, but got %d", -100, score)
	}
}
