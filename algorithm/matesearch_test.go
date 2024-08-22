package algorithm

import (
	"strings"
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
		{fen: "rn1qkbnr/pbpp1ppp/1p6/4p3/2B1P3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 1", depth: 1, expected: "1. Qxf7# 1-0"},
		{fen: "3nkr2/3Rb1pp/p1B1ppn1/1p4P1/7P/6Q1/PPPNq3/1K6 w - - 0 1", depth: 2, expected: "1. Rxe7+ Kxe7 2. Qc7# 1-0"},
		{fen: "r4r2/pQ3ppp/2np4/2bk4/5P2/6P1/PPP5/R1B1KB1q w Q - 0 1", depth: 2, expected: "1. Qb3+ Ke4 2. Qd3# 1-0"},
		{fen: "r1bk3b/1pppq3/2n3n1/1p2P1BQ/3P4/8/PPP3P1/5RK1 w - - 0 1", depth: 2, expected: "1. Qxh8+ Nxh8 2. Rf8# 1-0"},
		{fen: "8/8/6pk/3Q4/5K2/8/8/8 w - - 0 1", depth: 3, expected: "1. Qg8 g5+ 2. Kf5 Kh5 3. Qh8# 1-0"},
		{fen: "r3r1k1/ppp2ppp/5n2/7q/8/2N3Qb/PPP2P1P/R1BR2K1 b - - 0 1", depth: 2, expected: "1... Qxd1+ 2. Nxd1 Re1# 0-1"},
		{fen: "r4rk1/ppp3pp/4b3/6K1/8/8/PB3bPP/RN1Q3R b - - 0 1", depth: 3, expected: "1... Rf5+ 2. Kg4 h5+ 3. Kh3 Rf3# 0-1"},
		// {fen: "1kbr3r/pp6/8/P1n2ppq/2N3n1/R3Q1P1/3B1P2/2R2BK1 w - - 0 1", depth: 11, expected: "1. Qf4+ gxf4 2. Bxf4+ Ne5 3. Bxe5+ Rd6 4. Bxd6+ Ka8 5. Nb6+ axb6 6. axb6+ Na6 7. Rxc8+ Rxc8 8. Rxa6+ bxa6 9. Bg2+ Qf3 10. Bxf3+ Rc6 11. Bxc6# 1-0"},
		{fen: "1r5r/p4p1k/3p1p1p/3N4/3P4/5q1P/P4P1K/4Q1R1 w - - 0 1", depth: 2, expected: "1. Qe4+ Qxe4 2. Nxf6# 1-0"},
		{fen: "3q1r1k/5p1p/4pb1p/3nN3/3P2Q1/6R1/5PPP/6K1 w - - 0 1", depth: 2, expected: "1. Qg8+ Rxg8 2. Nxf7# 1-0"},
		{fen: "7k/1p1b3p/2n3p1/5p2/4pB2/1BP1Q2P/1q3PP1/3K4 w - - 0 1", depth: 2, expected: "1. Qd4+ Nxd4 2. Be5# 1-0"},
		{fen: "3q1b1r/3bk2p/3p1pB1/3Pp2Q/8/4B2P/5PP1/4K3 w - - 0 1", depth: 2, expected: "1. Qxe5+ fxe5 2. Bg5# 1-0"},
		{fen: "1q2r1k1/2p1bpp1/8/1r6/8/1B6/1B6/1K2Q2R w - - 0 1", depth: 4, expected: "1. Qe6 Rxb3 2. Rh8+ Kxh8 3. Qh6+ Kg8 4. Qxg7# 1-0"},
		{fen: "kb5Q/p7/Pp6/1P6/4p3/4R3/4P1p1/6K1 w - - 0 1", depth: 3, expected: "1. Rh3 e3 2. Rh1 gxh1=Q+ 3. Qxh1# 1-0"},
		{fen: "r5k1/2p2Rpp/8/p3N3/2B3P1/q1P5/7P/6K1 w - - 0 1", depth: 3, expected: "1. Rf3+ Kh8 2. Ng6+ hxg6 3. Rh3# 1-0"},
		{fen: "2r2r1k/6pp/2NN4/5p2/2Q2nq1/8/6PP/2R4K w - - 0 1", depth: 5, expected: "1. Qg8+ Kxg8 2. Ne7+ Kh8 3. Nf7+ Rxf7 4. Rxc8+ Rf8 5. Rxf8# 1-0"},
		{fen: "k7/2K1R3/6bR/8/8/8/8/8 w - - 0 1", depth: 2, expected: "1. Re4 Bxe4 2. Ra6# 1-0"},
		{fen: "kr5R/rp6/6K1/8/4Q3/8/8/R7 w - - 0 1", depth: 2, expected: "1. Qh1 Rg8+ 2. Rxg8# 1-0"},
		{fen: "r5k1/6p1/2p2p1p/N3p3/1P2n3/2P2nbP/2QB1qB1/4RR1K b - - 0 1", depth: 2, expected: "1... Qg1+ 2. Rxg1 Nf2# 0-1"},
		{fen: "8/8/8/8/5N2/6n1/1np5/1rk1K2Q w - - 0 1", depth: 2, expected: "1. Qa8 Nd3+ 2. Nxd3# 1-0"},
		{fen: "r3r1k1/pp3Rpp/2pp4/8/2P1P2B/5Q2/P1q3PP/5R1K w - - 0 1", depth: 4, expected: "1. Rxg7+ Kxg7 2. Qf6+ Kg8 3. Qf7+ Kh8 4. Bf6# 1-0"},
		{fen: "r2qrk2/ppp3pQ/5p2/3Np1b1/2B1P1P1/P2P4/1PP5/1K5R w - - 0 1", depth: 3, expected: "1. Qg8+ Kxg8 2. Ne7+ Kf8 3. Ng6# 1-0"},
		{fen: "3k4/8/3K2b1/8/8/8/7Q/8 w - - 0 1", depth: 2, expected: "1. Qa2 Bc2 2. Qg8# 1-0"},
		{fen: "5rk1/7p/b2pNp2/3P1N2/1p2P3/6KP/P2Q1P2/2B1nq2 w - - 0 1", depth: 4, expected: "1. Qg5+ fxg5 2. Nh6+ Kh8 3. Bb2+ Rf6 4. Bxf6# 1-0"},
		{fen: "5R2/5Kbk/R7/8/8/8/8/8 w - - 0 1", depth: 3, expected: "1. Rf6 Bxf6 2. Kxf6 Kh6 3. Rh8# 1-0"},
		{fen: "8/5P1k/5B2/8/8/3K4/8/8 w - - 0 1", depth: 4, expected: "1. Ke4 Kg6 2. f8=R Kh5 3. Kf5 Kh6 4. Rh8# 1-0"},
		{fen: "5R2/6pp/p4pk1/4Pb2/8/4QPKP/3B4/r2q4 w - - 0 1", depth: 4, expected: "1. Qh6+ gxh6 2. Rxf6+ Kg7 3. Bxh6+ Kg8 4. Rf8# 1-0"},
		{fen: "2r4k/p5p1/3Qq2p/4N3/Pp3n2/6P1/PK5P/3R4 b - - 0 1", depth: 3, expected: "1... Rc2+ 2. Kxc2 Qxa2+ 3. Kc1 Ne2# 0-1"},
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
		output = strings.TrimSpace(output)
		expected := strings.TrimSpace(ms.expected)
		if output != expected {
			t.Fatalf("wrong output: expected %s, but got %s", expected, output)
		}
	}
}
