package chess

import (
	"testing"
)

const ERROR_MESSAGE = "received unexpected error"
const BOARD_STRING_ERROR_MESSAGE = "expected board string %s but got %s"

func TestBoardTextSerialization(t *testing.T) {
	fen := INITIAL_POSITION
	b := &Board{}
	if err := b.UnmarshalText([]byte(fen)); err != nil {
		t.Fatal(ERROR_MESSAGE, err)
	}
	txt, err := b.MarshalText()
	if err != nil {
		t.Fatal(ERROR_MESSAGE, err)
	}
	if fen != string(txt) {
		t.Fatalf("fen expected board string %s but got %s", fen, string(txt))
	}
}

func TestBoardBinarySerialization(t *testing.T) {
	g := NewGame()
	board := g.Position().Board()
	b, err := board.MarshalBinary()
	if err != nil {
		t.Fatal(ERROR_MESSAGE, err)
	}
	cpBoard := &Board{}
	err = cpBoard.UnmarshalBinary(b)
	if err != nil {
		t.Fatal(ERROR_MESSAGE, err)
	}
	s := INITIAL_POSITION
	if s != cpBoard.String() {
		t.Fatalf(BOARD_STRING_ERROR_MESSAGE, s, cpBoard.String())
	}
}

func TestBoardRotation(t *testing.T) {
	fens := []string{
		"RP4pr/NP4pn/BP4pb/QP4pq/KP4pk/BP4pb/NP4pn/RP4pr",
		"RNBKQBNR/PPPPPPPP/8/8/8/8/pppppppp/rnbkqbnr",
		"rp4PR/np4PN/bp4PB/kp4PK/qp4PQ/bp4PB/np4PN/rp4PR",
		INITIAL_POSITION,
	}
	g := NewGame()
	board := g.Position().Board()
	for i := 0; i < 4; i++ {
		board = board.Rotate()
		if fens[i] != board.String() {
			t.Fatalf(BOARD_STRING_ERROR_MESSAGE, fens[i], board.String())
		}
	}
}

func TestBoardFlip(t *testing.T) {
	g := NewGame()
	board := g.Position().Board()
	board = board.Flip(UpDown)
	b := "RNBQKBNR/PPPPPPPP/8/8/8/8/pppppppp/rnbqkbnr"
	if b != board.String() {
		t.Fatalf(BOARD_STRING_ERROR_MESSAGE, b, board.String())
	}
	board = board.Flip(UpDown)
	b = INITIAL_POSITION
	if b != board.String() {
		t.Fatalf(BOARD_STRING_ERROR_MESSAGE, b, board.String())
	}
	board = board.Flip(LeftRight)
	b = "rnbkqbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR"
	if b != board.String() {
		t.Fatalf(BOARD_STRING_ERROR_MESSAGE, b, board.String())
	}
	board = board.Flip(LeftRight)
	b = INITIAL_POSITION
	if b != board.String() {
		t.Fatalf(BOARD_STRING_ERROR_MESSAGE, b, board.String())
	}
}

func TestBoardTranspose(t *testing.T) {
	g := NewGame()
	board := g.Position().Board()
	board = board.Transpose()
	b := "rp4PR/np4PN/bp4PB/qp4PQ/kp4PK/bp4PB/np4PN/rp4PR"
	if b != board.String() {
		t.Fatalf(BOARD_STRING_ERROR_MESSAGE, b, board.String())
	}
}
