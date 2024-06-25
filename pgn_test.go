package chess

import (
	"io"
	"os"
	"strings"
	"testing"
)

type pgnTest struct {
	PostPos *Position
	PGN     string
	Error   error
}

type invalidPgnTest struct {
	PGN   string
	Error error
}

func MakePGNTest(pos string, arg string) pgnTest {
	pgn, err := mustParsePGN()(arg)
	return pgnTest{
		PostPos: unsafeFEN(pos),
		PGN:     pgn,
		Error:   err,
	}
}

func MakeInvalidPGNTest(arg string) invalidPgnTest {
	pgn, err := mustParsePGN()(arg)
	return invalidPgnTest{
		PGN:   pgn,
		Error: err,
	}
}

var validPGNs = []pgnTest{
	MakePGNTest("4r3/6P1/2p2P1k/1p6/pP2p1R1/P1B5/2P2K2/3r4 b - - 0 45", "fixtures/pgns/0001.pgn"),
	MakePGNTest("4r3/6P1/2p2P1k/1p6/pP2p1R1/P1B5/2P2K2/3r4 b - - 0 45", "fixtures/pgns/0002.pgn"),
	MakePGNTest("2r2rk1/pp1bBpp1/2np4/2pp2p1/1bP5/1P4P1/P1QPPPBP/3R1RK1 b - - 0 3", "fixtures/pgns/0003.pgn"),
	MakePGNTest("r3kb1r/2qp1pp1/b1n1p2p/pp2P3/5n1B/1PPQ1N2/P1BN1PPP/R3K2R w KQkq - 1 14", "fixtures/pgns/0004.pgn"),
	MakePGNTest("r3kb1r/2qp1pp1/b1n1p2p/pp2P3/5n1B/1PPQ1N2/P1BN1PPP/R3K2R w KQkq - 1 14", "fixtures/pgns/0004.pgn"),
	MakePGNTest("rnbqkbnr/ppp2ppp/4p3/3p4/3PP3/8/PPP2PPP/RNBQKBNR w KQkq d6 0 3", "fixtures/pgns/0008.pgn"),
	MakePGNTest("r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", "fixtures/pgns/0009.pgn"),
	MakePGNTest("r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", "fixtures/pgns/0010.pgn"),
	MakePGNTest("8/8/6p1/4R3/6kQ/r2P1pP1/5P2/6K1 b - - 3 42", "fixtures/pgns/0011.pgn"),
	MakePGNTest(INITIAL_FEN_POSITION, "fixtures/pgns/0012.pgn"),
	MakePGNTest(INITIAL_FEN_POSITION, "fixtures/pgns/0013.pgn"),
	MakePGNTest("rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2", "fixtures/pgns/0016.pgn"),
}
var invalidPGNs = []invalidPgnTest{
	MakeInvalidPGNTest("fixtures/invalid/0001.pgn"),
	MakeInvalidPGNTest("fixtures/invalid/0002.pgn"),
}

func TestValidPGNs(t *testing.T) {
	for _, test := range validPGNs {
		game, err := decodePGN(nil, test.PGN)
		if err != nil {
			t.Fatalf("received unexpected pgn error %s", err.Error())
		}
		if game.Position().String() != test.PostPos.String() {
			t.Fatalf("expected board to be \n%s\nFEN:%s\n but got \n%s\n\nFEN:%s\n",
				test.PostPos.board.Draw(), test.PostPos.String(),
				game.Position().board.Draw(), game.Position().String())
		}
		if game.String() == "" {
			t.Fatalf("Unexpected pgn output")
		}
	}
}

func TestInvalidPGNs(t *testing.T) {
	for _, test := range invalidPGNs {
		_, err := decodePGN(nil, test.PGN)
		if err == nil {
			t.Fatalf("missing expected pgn error")
		}
	}
}

type commentTest struct {
	PGN         string
	Error       error
	MoveNumber  int
	CommentText string
}

func MakeCommentText(file string, moveNumber int, comment string) commentTest {
	pgn, _ := mustParsePGN()(file)
	return commentTest{
		PGN:         pgn,
		MoveNumber:  moveNumber,
		CommentText: comment,
	}
}

var (
	commentTests = []commentTest{
		MakeCommentText("fixtures/pgns/0005.pgn", 7, `(-0.25 → 0.39) Inaccuracy. cxd4 was best. [%eval 0.39] [%clk 0:05:05]`),
		MakeCommentText("fixtures/pgns/0009.pgn", 5, `This opening is called the Ruy Lopez.`),
		MakeCommentText("fixtures/pgns/0010.pgn", 5, `This opening is called the Ruy Lopez.`),
	}
)

func TestCommentsDetection(t *testing.T) {
	for _, test := range commentTests {
		game, err := decodePGN(nil, test.PGN)
		if err != nil {
			t.Fatal(err)
		}
		comment := strings.Join(game.Comments()[test.MoveNumber], " ")
		if comment != test.CommentText {
			t.Fatalf("expected pgn comment to be %s but got %s", test.CommentText, comment)
		}
	}
}

func TestNewGameComments(t *testing.T) {
	for _, test := range commentTests {
		pgn, err := PGN(NewInput(strings.NewReader(test.PGN)))
		if err != nil {
			t.Fatal(err)
		}
		game := NewGame(pgn)
		comment := strings.Join(game.Comments()[test.MoveNumber], " ")
		if comment != test.CommentText {
			t.Fatalf("expected pgn comment to be %s but got %s", test.CommentText, comment)
		}
	}
}

func TestWriteComments(t *testing.T) {
	pgn, err := mustParsePGN()("fixtures/pgns/0005.pgn")
	if err != nil {
		t.Fatal(err)
	}
	game, err := decodePGN(nil, pgn)
	if err != nil {
		t.Fatal(err)
	}
	game, err = decodePGN(nil, game.String())
	if err != nil {
		t.Fatal(err)
	}
	if len(game.Comments()[7]) != 2 {
		t.Fatalf("expected %d comments for move 7 but got %d", 2, len(game.Comments()[7]))
	}
}

func TestScanner(t *testing.T) {
	for _, fname := range []string{"fixtures/pgns/0006.pgn", "fixtures/pgns/0007.pgn", "fixtures/pgns/0014.pgn"} {
		f, err := os.Open(fname)
		if err != nil {
			t.Fatal("could not open file")
		}
		defer f.Close()
		scanner := NewScanner(f)
		games := []*Game{}
		for scanner.Scan() {
			game := scanner.Next()
			games = append(games, game)
		}
		if len(games) != 5 {
			t.Fatalf(fname+" expected 5 games but got %d", len(games))
		}
	}
}

func BenchmarkPGN(b *testing.B) {
	pgn, _ := mustParsePGN()("fixtures/pgns/0001.pgn")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		opt, _ := PGN(NewInput(strings.NewReader(pgn)))
		NewGame(opt)
	}
}

func mustParsePGN() func(n string) (string, error) {
	return func(n string) (string, error) {
		f, err := os.Open(n)
		if err != nil {
			return "", err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}
func TestGamesFromPGN(t *testing.T) {
	for _, test := range validPGNs {
		reader := strings.NewReader(test.PGN)
		games, err := GamesFromPGN(reader)
		if err != nil {
			t.Fatalf("fail to read games from valid pgn: %s", err.Error())
		}
		if len(games) != 1 {
			t.Fatalf("expected to get 1 game from pgn, got %d", len(games))
		}
	}
}

func TestScannerWithFromPosFENs(t *testing.T) {
	finalPositions := []string{
		"rnbqkbnr/pp2pppp/2p5/3p4/3PP3/5P2/PPP3PP/RNBQKBNR b KQkq - 0 3",
		"r2qkb1r/pp1n1ppp/2p2n2/4p3/2BPP1b1/2P2N2/PP4PP/RNBQ1RK1 b kq - 0 8",
		"rnbqk2r/pp2nppp/2p1p3/3p4/1b1PP3/2NB1P2/PPPB2PP/R2QK1NR b KQkq - 5 6",
		"rnbqk1nr/pp2ppbp/2p3p1/3p4/3PP3/2N1BP2/PPP3PP/R2QKBNR b KQkq - 3 5",
		"rnb1kbnr/pp3ppp/1qp5/8/3NP3/2N5/PPP3PP/R1BQKB1R b KQkq - 0 7",
	}
	fname := "fixtures/pgns/0015.pgn"
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := NewScanner(f)
	games := []*Game{}
	for idx := 0; scanner.Scan(); {
		game := scanner.Next()
		if len(game.moves) == 0 {
			continue
		}
		finalPos := game.Position().String()
		if finalPos != finalPositions[idx] {
			t.Fatalf(fname+" game %v expected final pos %v but got %v", idx,
				finalPositions[idx], finalPos)
		}
		games = append(games, game)
		idx++
	}
	if len(games) != len(finalPositions) {
		t.Fatalf(fname+" expected %v games but got %v", len(finalPositions),
			len(games))
	}
}
