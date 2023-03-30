package algorithm

import (
	"testing"

	"github.com/othomann/go-chess"
)

type MateTests struct {
	fen      string
	depth    int
	expected string
}

var (
	mateTests = []MateTests{
		{fen: "rn1qkbnr/pbpp1ppp/1p6/4p3/2B1P3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 1", depth: 1, expected: "\n1. Qxf7# 1-0"},
		{fen: "3nkr2/3Rb1pp/p1B1ppn1/1p4P1/7P/6Q1/PPPNq3/1K6 w - - 0 1", depth: 2, expected: "\n1. Rxe7+ Kxe7 2. Qc7# 1-0"},
		{fen: "r4r2/pQ3ppp/2np4/2bk4/5P2/6P1/PPP5/R1B1KB1q w Q - 0 1", depth: 2, expected: "\n1. Qb3+ Ke4 2. Qd3# 1-0"},
		{fen: "r1bk3b/1pppq3/2n3n1/1p2P1BQ/3P4/8/PPP3P1/5RK1 w - - 0 1", depth: 2, expected: "\n1. Qxh8+ Nxh8 2. Rf8# 1-0"},
		{fen: "8/8/6pk/3Q4/5K2/8/8/8 w - - 0 1", depth: 3, expected: "\n1. Qg8 g5+ 2. Kf5 Kh5 3. Qh8# 1-0"},
		{fen: "r3r1k1/ppp2ppp/5n2/7q/8/2N3Qb/PPP2P1P/R1BR2K1 b - - 0 1", depth: 2, expected: "\n1... Qxd1+ 2. Nxd1 Re1# 0-1"},
		{fen: "r4rk1/ppp3pp/4b3/6K1/8/8/PB3bPP/RN1Q3R b - - 0 1", depth: 3, expected: "\n1... Rf5+ 2. Kg4 h5+ 3. Kh3 Rf3# 0-1"},
		// {fen: "1kbr3r/pp6/8/P1n2ppq/2N3n1/R3Q1P1/3B1P2/2R2BK1 w - - 0 1", depth: 11, expected: "\n1. Qf4+ gxf4 2. Bxf4+ Ne5 3. Bxe5+ Rd6 4. Bxd6+ Ka8 5. Nb6+ axb6 6. axb6+ Na6 7. Rxc8+ Rxc8 8. Rxa6+ bxa6 9. Bg2+ Qf3 10. Bxf3+ Rc6 11. Bxc6# 1-0"},
		{fen: "1r5r/p4p1k/3p1p1p/3N4/3P4/5q1P/P4P1K/4Q1R1 w - - 0 1", depth: 2, expected: "\n1. Qe4+ Qxe4 2. Nxf6# 1-0"},
		{fen: "3q1r1k/5p1p/4pb1p/3nN3/3P2Q1/6R1/5PPP/6K1 w - - 0 1", depth: 2, expected: "\n1. Qg8+ Rxg8 2. Nxf7# 1-0"},
		{fen: "7k/1p1b3p/2n3p1/5p2/4pB2/1BP1Q2P/1q3PP1/3K4 w - - 0 1", depth: 2, expected: "\n1. Qd4+ Nxd4 2. Be5# 1-0"},
		{fen: "3q1b1r/3bk2p/3p1pB1/3Pp2Q/8/4B2P/5PP1/4K3 w - - 0 1", depth: 2, expected: "\n1. Qxe5+ fxe5 2. Bg5# 1-0"},
	}
)

func TestMateSearch(t *testing.T) {
	for index, ms := range mateTests {
		fen, err := chess.FEN(ms.fen)
		if err != nil {
			t.Fatalf("error for test %d", index)
		}
		g := chess.NewGame(fen)
		result, mateNode := MateSearch(g, ms.depth, NewRoot(g, ms.depth))
		if !result {
			t.Fatalf("no mate found for test %d", index)
		}
		output, err := mateNode.String()
		if err != nil {
			t.Fatal(err)
		}
		expected := ms.expected
		if output != expected {
			t.Fatalf("wrong output: expected %s, but got %s", expected, output)
		}
	}
}
