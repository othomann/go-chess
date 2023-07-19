package chess

import (
	"log"
	"strings"
	"testing"
)

func TestCheckmate(t *testing.T) {
	fenStr := "rn1qkbnr/pbpp1ppp/1p6/4p3/2B1P3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 1"
	fen, err := FEN(fenStr)
	if err != nil {
		t.Fatal(err)
	}
	g := NewGame(fen)
	if err := g.MoveStr("Qxf7#"); err != nil {
		t.Fatal(err)
	}
	if g.Method() != Checkmate {
		t.Fatalf("expected method %s but got %s", Checkmate, g.Method())
	}
	if g.Outcome() != WhiteWon {
		t.Fatalf("expected outcome %s but got %s", WhiteWon, g.Outcome())
	}

	// Checkmate on castle
	fenStr = "Q7/5Qp1/3k2N1/7p/8/4B3/PP3PPP/R3K2R w KQ - 0 31"
	fen, err = FEN(fenStr)
	if err != nil {
		t.Fatal(err)
	}
	g = NewGame(fen)
	if err := g.MoveStr("O-O-O"); err != nil {
		t.Fatal(err)
	}
	if g.Method() != Checkmate {
		t.Fatalf("expected method %s but got %s", Checkmate, g.Method())
	}
	if g.Outcome() != WhiteWon {
		t.Fatalf("expected outcome %s but got %s", WhiteWon, g.Outcome())
	}
}

func TestCheckmateFromFen(t *testing.T) {
	fenStr := "rn1qkbnr/pbpp1Qpp/1p6/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 1"
	fen, err := FEN(fenStr)
	if err != nil {
		t.Fatal(err)
	}
	g := NewGame(fen)
	if g.Method() != Checkmate {
		t.Error(g.Position().Board().Draw())
		t.Fatalf("expected method %s but got %s", Checkmate, g.Method())
	}
	if g.Outcome() != WhiteWon {
		t.Fatalf("expected outcome %s but got %s", WhiteWon, g.Outcome())
	}
}

func TestStalemate(t *testing.T) {
	fenStr := "k1K5/8/8/8/8/8/8/1Q6 w - - 0 1"
	fen, err := FEN(fenStr)
	if err != nil {
		t.Fatal(err)
	}
	g := NewGame(fen)
	if err := g.MoveStr("Qb6"); err != nil {
		t.Fatal(err)
	}
	if g.Method() != Stalemate {
		t.Fatalf("expected method %s but got %s", Stalemate, g.Method())
	}
	if g.Outcome() != Draw {
		t.Fatalf("expected outcome %s but got %s", Draw, g.Outcome())
	}
}

// position shouldn't result in stalemate because pawn can move http://en.lichess.org/Pc6mJDZN#138
func TestInvalidStalemate(t *testing.T) {
	fenStr := "8/3P4/8/8/8/7k/7p/7K w - - 2 70"
	fen, err := FEN(fenStr)
	if err != nil {
		t.Fatal(err)
	}
	g := NewGame(fen)
	if err := g.MoveStr("d8=Q"); err != nil {
		t.Fatal(err)
	}
	if g.Outcome() != NoOutcome {
		t.Fatalf("expected outcome %s but got %s", NoOutcome, g.Outcome())
	}
}

func TestThreeFoldRepetition(t *testing.T) {
	g := NewGame()
	moves := []string{
		"Nf3", "Nf6", "Ng1", "Ng8",
		"Nf3", "Nf6", "Ng1", "Ng8",
	}
	for _, m := range moves {
		if err := g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.Draw(ThreefoldRepetition); err != nil {
		for _, pos := range g.Positions() {
			log.Println(pos.String())
		}
		t.Fatalf("%s - %d reps", err.Error(), g.numOfRepetitions())
	}
}

func TestInvalidThreeFoldRepetition(t *testing.T) {
	g := NewGame()
	moves := []string{
		"Nf3", "Nf6", "Ng1", "Ng8",
	}
	for _, m := range moves {
		if err := g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.Draw(ThreefoldRepetition); err == nil {
		t.Fatal("should require three repeated board states")
	}
}

func TestFiveFoldRepetition(t *testing.T) {
	g := NewGame()
	moves := []string{
		"Nf3", "Nf6", "Ng1", "Ng8",
		"Nf3", "Nf6", "Ng1", "Ng8",
		"Nf3", "Nf6", "Ng1", "Ng8",
		"Nf3", "Nf6", "Ng1", "Ng8",
	}
	for _, m := range moves {
		if err := g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	if g.Outcome() != Draw || g.Method() != FivefoldRepetition {
		t.Fatal("should automatically draw after five repetitions")
	}
}

func TestFiftyMoveRule(t *testing.T) {
	fen, _ := FEN("2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 100 60")
	g := NewGame(fen)
	if err := g.Draw(FiftyMoveRule); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidFiftyMoveRule(t *testing.T) {
	fen, _ := FEN("2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 99 60")
	g := NewGame(fen)
	if err := g.Draw(FiftyMoveRule); err == nil {
		t.Fatal("should require fifty moves")
	}
}

func TestDrawMethod(t *testing.T) {
	fen, _ := FEN("2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 99 60")
	g := NewGame(fen)
	if err := g.Draw(Resignation); err == nil {
		t.Fatal("should require fifty moves")
	}
}

func TestSeventyFiveMoveRule(t *testing.T) {
	fen, _ := FEN("2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 149 80")
	g := NewGame(fen)
	if err := g.MoveStr("Kf8"); err != nil {
		t.Fatal(err)
	}
	if g.Outcome() != Draw || g.Method() != SeventyFiveMoveRule {
		t.Fatal("should automatically draw after seventy five moves w/ no pawn move or capture")
	}
}

func TestInsufficientMaterial(t *testing.T) {
	fens := []string{
		"8/2k5/8/8/8/3K4/8/8 w - - 1 1",
		"8/2k5/8/8/8/3K1N2/8/8 w - - 1 1",
		"8/2k5/8/8/8/3K1B2/8/8 w - - 1 1",
		"8/2k5/2b5/8/8/3K1B2/8/8 w - - 1 1",
		"4b3/2k5/2b5/8/8/3K1B2/8/8 w - - 1 1",
	}
	for _, f := range fens {
		fen, err := FEN(f)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGame(fen)
		if g.Method().String() != "InsufficientMaterial" {
			t.Fatalf("Method InsufficientMaterial String() failed; expected %s, not got %s", "InsufficientMaterial", g.Method().String())
		}
		if g.Outcome().String() != "1/2-1/2" {
			t.Fatalf("Outcome Draw String() failed; expected %s, not got %s", "1/2-1/2", g.Outcome().String())
		}
		if g.Outcome() != Draw || g.Method() != InsufficientMaterial {
			log.Println(g.Position().Board().Draw())
			t.Fatalf("%s should automatically draw by insufficient material", f)
		}
	}
}

func TestInsufficientMaterial2(t *testing.T) {
	fens := []string{
		"8/8/8/8/8/3K4/8/8 w - - 1 1",
		"8/2k5/8/8/8/8/8/8 w - - 1 1",
	}
	for _, f := range fens {
		fen, err := FEN(f)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGame(fen)
		if g.Method().String() != "InsufficientMaterial" {
			t.Fatalf("Method InsufficientMaterial String() failed; expected %s, not got %s", "InsufficientMaterial", g.Method().String())
		}
	}
}

func TestSufficientMaterial(t *testing.T) {
	fens := []string{
		"8/2k5/8/8/8/3K1B2/4N3/8 w - - 1 1",
		"8/2k5/8/8/8/3KBB2/8/8 w - - 1 1",
		"8/2k1b3/8/8/8/3K1B2/8/8 w - - 1 1",
		"8/2k5/8/8/4P3/3K4/8/8 w - - 1 1",
		"8/2k5/8/8/8/3KQ3/8/8 w - - 1 1",
		"8/2k5/8/8/8/3KR3/8/8 w - - 1 1",
	}
	for _, f := range fens {
		fen, err := FEN(f)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGame(fen)
		if g.Outcome() != NoOutcome {
			log.Println(g.Position().Board().Draw())
			t.Fatalf("%s should not find insufficient material", f)
		}
	}
}

func TestSerializationCycle(t *testing.T) {
	g := NewGame()
	g.MoveStr("e4")
	g.MoveStr("e5")
	pgn, err := PGN(NewInput(strings.NewReader(g.String())))
	if err != nil {
		t.Fatal(err)
	}
	cp := NewGame(pgn)
	if cp.String() != g.String() {
		t.Fatalf("expected %s but got %s", g.String(), cp.String())
	}
}

func TestInitialNumOfValidMoves(t *testing.T) {
	g := NewGame()
	if len(g.ValidMoves()) != 20 {
		t.Fatal("should find 20 valid moves from the initial position")
	}
}

func TestTagPairs(t *testing.T) {
	g := NewGame()
	override := g.AddTagPair("Draw Offer", "White")
	if override {
		t.Fatalf("TagPair was overriden")
	}
	tagPair := g.GetTagPair("Draw Offer")
	if tagPair == nil {
		t.Fatalf("expected %s but got %s", "White", "nil")
	}
	if tagPair.Value != "White" {
		t.Fatalf("expected %s but got %s", "White", tagPair.Value)
	}
	override = g.AddTagPair("Draw Offer", "Black")
	if !override {
		t.Fatalf("TagPair was not overriden")
	}
	g.RemoveTagPair("Draw Offer")
	tagPair = g.GetTagPair("Draw Offer")
	if tagPair != nil {
		t.Fatalf("expected %s but got %s", "nil", "not nil")
	}
}

func TestPositionHash(t *testing.T) {
	g1 := NewGame()
	for _, s := range []string{"Nc3", "e5", "Nf3"} {
		g1.MoveStr(s)
	}
	g2 := NewGame()
	for _, s := range []string{"Nf3", "e5", "Nc3"} {
		g2.MoveStr(s)
	}
	if g1.Position().Hash() != g2.Position().Hash() {
		t.Fatalf("expected position hashes to be equal but got %s and %s", g1.Position().Hash(), g2.Position().Hash())
	}
}

func TestMoveHistory(t *testing.T) {
	lens := []int{89, 89, 5, 26}
	for i, test := range validPGNs[0:4] {
		pgn, err := PGN(NewInput(strings.NewReader(test.PGN)))
		if err != nil {
			t.Fatal(err)
		}
		game := NewGame(pgn)
		l := len(game.MoveHistory())
		if lens[i] != l {
			t.Fatalf("expected history length to be %d but got %d", lens[i], l)
		}
	}
}

func TestMoveHistory2(t *testing.T) {
	game := NewGame()
	game.MoveStr("e4")
	game.MoveStr("e5")
	game.Resign(Black)
	history := game.MoveHistory()
	if len(history) != 2 {
		t.Fatal("Didn't retrieve full history")
	}
	output := game.Position().Board().Draw()
	if output == "" {
		t.Fatalf("Wrong board output; expected %s, but got %s", "", output)
	}
	fen := game.FEN()
	expected := "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"
	if fen != expected {
		t.Fatalf("Wrong fen output; expected %s, but got %s", expected, fen)
	}
	err := game.Move(nil)
	if err == nil {
		t.Fatalf("Error should be reported")
	}
	err = game.Move(&Move{
		s1: A1,
		s2: A2,
	})
	if err == nil {
		t.Fatalf("Error should be reported")
	}
}

func TestMoveHistory3(t *testing.T) {
	game := NewGame()
	game.MoveStr("e4")
	game.MoveStr("e5")
	game.Resign(Black)
	history := game.MoveHistory()
	if len(history) != 2 {
		t.Fatal("Didn't retrieve full history")
	}
	game.UndoMove()
	history = game.MoveHistory()
	if len(history) != 1 {
		t.Fatal("Didn't undo move")
	}
	game.UndoMove()
	history = game.MoveHistory()
	if len(history) != 0 {
		t.Fatal("Didn't undo move")
	}
	err := game.UndoMove()
	if err == nil {
		t.Fatal("error should be returned as there is no move to undo")
	}
}

func TestMoveHistory4(t *testing.T) {
	game := NewGame()
	game.MoveStr("e4")
	game.MoveStr("e5")
	game.Resign(Black)
	history := game.MoveHistory()
	if len(history) != 2 {
		t.Fatal("Didn't retrieve full history")
	}
	game.UndoMoves(2)
	history = game.MoveHistory()
	if len(history) != 0 {
		t.Fatal("Didn't undo last 2 moves")
	}
	err := game.UndoMove()
	if err == nil {
		t.Fatal("error should be returned as there is no move to undo")
	}
}

func TestMarshalling(t *testing.T) {
	game := NewGame()
	game.MoveStr("e4")
	game.MoveStr("e5")
	game.Resign(Black)
	bytes, err := game.MarshalText()
	if err != nil {
		t.Fatalf("Marshalling failed with %s", err)
	}
	output := string(bytes)
	expected := "\n1. e4 e5 1-0"
	if output != expected {
		t.Fatalf("Wrong marshalling; expected %s, but got %s", expected, output)
	}
	err = game.UnmarshalText(bytes)
	if err != nil {
		t.Fatalf("Marshalling failed with %s", err)
	}
}

func TestMarshalling2(t *testing.T) {
	game := NewGame()
	bytes := []byte("1. e4 e3")
	err := game.UnmarshalText(bytes)
	if err == nil {
		t.Fatalf("Missing error while marshalling")
	}
}

func BenchmarkStalemateStatus(b *testing.B) {
	fenStr := "k1K5/8/8/8/8/8/8/1Q6 w - - 0 1"
	fen, err := FEN(fenStr)
	if err != nil {
		b.Fatal(err)
	}
	g := NewGame(fen)
	if err := g.MoveStr("Qb6"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Position().Status()
	}
}

func BenchmarkInvalidStalemateStatus(b *testing.B) {
	fenStr := "8/3P4/8/8/8/7k/7p/7K w - - 2 70"
	fen, err := FEN(fenStr)
	if err != nil {
		b.Fatal(err)
	}
	g := NewGame(fen)
	if err := g.MoveStr("d8=Q"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Position().Status()
	}
}

func BenchmarkPositionHash(b *testing.B) {
	fenStr := "8/3P4/8/8/8/7k/7p/7K w - - 2 70"
	fen, err := FEN(fenStr)
	if err != nil {
		b.Fatal(err)
	}
	g := NewGame(fen)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Position().Hash()
	}
}

type PerfTest struct {
	fen      string
	depth    int
	expected int
}

func TestPerft(t *testing.T) {

	mateTests := []PerfTest{
		{fen: INITIAL_FEN_POSITION, depth: 1, expected: 20},
		{fen: INITIAL_FEN_POSITION, depth: 2, expected: 400},
		{fen: INITIAL_FEN_POSITION, depth: 3, expected: 8902},
		{fen: INITIAL_FEN_POSITION, depth: 4, expected: 197281},
		{fen: INITIAL_FEN_POSITION, depth: 5, expected: 4865609},
	}

	for _, perfTest := range mateTests {
		result, duration, err := Perft(perfTest.fen, perfTest.depth)

		if err != nil {
			t.Fatalf("Perft failed for depth %d with err %s", perfTest.depth, err)
		}
		if result != perfTest.expected {
			t.Fatalf("Perft failed for depth %d: expected %d, got %d", 4, perfTest.expected, result)
		}
		t.Logf("Time spend on perf %s with depth %d: %f\n", perfTest.fen, perfTest.depth, duration.Seconds())
	}
}
