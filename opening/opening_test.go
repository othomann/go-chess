package opening_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/othomann/go-chess"
	"github.com/othomann/go-chess/opening"
)

func TestFind(t *testing.T) {
	g := chess.NewGame()
	g.MoveStr("e4")
	g.MoveStr("e6")

	// print French Defense
	book, err := opening.NewBookECO()
	if err != nil {
		t.Fatal(err)
	}
	o := book.Find(g.Moves())
	if o == nil || o.Title() != "French Defense" {
		t.Fatalf("Wrong opening found for %s", g.String())
	}
}

func TestPossible(t *testing.T) {
	g := chess.NewGame()
	g.MoveStr("e4")
	g.MoveStr("d5")

	// print all variantions of the Scandinavian Defense
	book, err := opening.NewBookECO()
	if err != nil {
		t.Fatal(err)
	}
	openings := book.Possible(g.Moves())
	if len(openings) != 40 {
		t.Fatalf("Wrong number of openings for %s: expected: %d but got %d", g.String(), 40, len(openings))
	}
	o := openings[0]
	if o.Code() == "" {
		t.Fatalf("Wrong code for opening: %s", o.Title())
	}
	if o.PGN() == "" {
		t.Fatalf("Wrong pgn for opening: %s", o.Title())
	}
}

func TestFind2(t *testing.T) {
	g := chess.NewGame()
	if err := g.MoveStr("e4"); err != nil {
		t.Fatal(err)
	}
	if err := g.MoveStr("d5"); err != nil {
		t.Fatal(err)
	}
	book, err := opening.NewBookECO()
	if err != nil {
		t.Fatal(err)
	}
	o := book.Find(g.Moves())
	expected := "Scandinavian Defense"
	if o == nil || o.Title() != expected {
		t.Fatalf("expected to find opening %s but got %s", expected, o.Title())
	}
}

func TestPossible2(t *testing.T) {
	g := chess.NewGame()
	if err := g.MoveStr("g3"); err != nil {
		t.Fatal(err)
	}
	book, err := opening.NewBookECO()
	if err != nil {
		t.Fatal(err)
	}
	openings := book.Possible(g.Moves())
	actual := len(openings)
	if actual != 22 {
		t.Fatalf("expected %d possible openings but got %d", 22, actual)
	}
}

func TestDraw(t *testing.T) {
	g := chess.NewGame()
	if err := g.MoveStr("g3"); err != nil {
		t.Fatal(err)
	}
	book, err := opening.NewBookECO()
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if err := book.Draw(buf); err != nil {
		t.Fatalf("Draw failed on ECO book")
	}
}

func TestParsing(t *testing.T) {
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	book, err := opening.NewBookECO()
	if err != nil {
		t.Fatal("Error in opening book")
	}
	moves := "e2e4 g8f6 e4e5 f6d5 f1c4"
	ml := strings.Split(moves, " ")

	for _, move := range ml {
		if err := game.MoveStr(move); err != nil {
			break
		}
	}
	s := game.Moves()
	p := book.Possible(s)
	g2 := p[0]
	m := g2.Game().Moves()
	bestMove := m[5].String()
	if bestMove != "d5b6" {
		t.Fatalf("Wrong move; expected: %s, but got: %s", "d5b6", bestMove)
	}
}

func BenchmarkNewBookECO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opening.NewBookECO()
	}
}
