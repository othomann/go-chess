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

func MakePGNTest(pos string, arg string) pgnTest {
	pgn, err := mustParsePGN()(arg)
	return pgnTest{
		PostPos: unsafeFEN(pos),
		PGN:     pgn,
		Error:   err,
	}
}

var (
	validPGNs = []pgnTest{
		MakePGNTest("4r3/6P1/2p2P1k/1p6/pP2p1R1/P1B5/2P2K2/3r4 b - - 0 45", "fixtures/pgns/0001.pgn"),
		MakePGNTest("4r3/6P1/2p2P1k/1p6/pP2p1R1/P1B5/2P2K2/3r4 b - - 0 45", "fixtures/pgns/0002.pgn"),
		MakePGNTest("2r2rk1/pp1bBpp1/2np4/2pp2p1/1bP5/1P4P1/P1QPPPBP/3R1RK1 b - - 0 3", "fixtures/pgns/0003.pgn"),
		MakePGNTest("r3kb1r/2qp1pp1/b1n1p2p/pp2P3/5n1B/1PPQ1N2/P1BN1PPP/R3K2R w KQkq - 1 14", "fixtures/pgns/0004.pgn"),
		MakePGNTest("r3kb1r/2qp1pp1/b1n1p2p/pp2P3/5n1B/1PPQ1N2/P1BN1PPP/R3K2R w KQkq - 1 14", "fixtures/pgns/0004.pgn"),
		MakePGNTest("rnbqkbnr/ppp2ppp/4p3/3p4/3PP3/8/PPP2PPP/RNBQKBNR w KQkq d6 0 3", "fixtures/pgns/0008.pgn"),
		MakePGNTest("r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", "fixtures/pgns/0009.pgn"),
		MakePGNTest("r1bqkbnr/1ppp1ppp/p1n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 4", "fixtures/pgns/0010.pgn"),
		MakePGNTest("8/8/6p1/4R3/6kQ/r2P1pP1/5P2/6K1 b - - 3 42", "fixtures/pgns/0011.pgn"),
		MakePGNTest("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "fixtures/pgns/0012.pgn"),
		MakePGNTest("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "fixtures/pgns/0013.pgn"),
	}
)

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
