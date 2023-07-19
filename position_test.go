package chess

import (
	"testing"
)

func TestPositionBinary(t *testing.T) {
	for _, fen := range validFENs {
		pos, err := decodeFEN(fen)
		if err != nil {
			t.Fatal(err)
		}
		b, err := pos.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		cp := &Position{}
		if err := cp.UnmarshalBinary(b); err != nil {
			t.Fatal(err)
		}
		if pos.String() != cp.String() {
			t.Fatalf("expected %s but got %s", pos.String(), cp.String())
		}
	}
}

func TestPositionUpdate(t *testing.T) {
	for _, fen := range validFENs {
		pos, err := decodeFEN(fen)
		if err != nil {
			t.Fatal(err)
		}

		np := pos.Update(pos.ValidMoves()[0])
		if pos.Turn().Other() != np.turn {
			t.Fatal("expected other turn")
		}
		if pos.board.String() == np.board.String() {
			t.Fatal("expected board update")
		}

		np = pos.Update(nil)
		if pos.Turn() != np.turn {
			t.Fatal("expected other turn")
		}
		if pos.halfMoveClock != np.halfMoveClock {
			t.Fatalf("expected half move clock increment: %d - expected %d", np.halfMoveClock, pos.halfMoveClock)
		}
		if pos.board.String() != np.board.String() {
			t.Fatal("expected same board")
		}
	}
}

func TestPositionMarshalling(t *testing.T) {
	pos := &Position{}
	expected := "r4r2/1b2bppk/ppq1p3/2pp3n/5P2/1P2P3/PBPPQ1PP/R4RK1 w - - 0 2"
	err := pos.UnmarshalText([]byte(expected))
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := pos.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	output := string(bytes)
	if output != expected {
		t.Fatalf("Expected: %s, but got %s", expected, output)
	}
	if pos.HalfMoveClock() != 0 {
		t.Fatalf("Wrong half-move clock. Expected: %d, but got %d", 2, pos.HalfMoveClock())
	}
	if pos.EnPassantSquare() != -1 {
		t.Fatalf("Wrong en passant square. Expected: %d, but got %d", 2, pos.EnPassantSquare())
	}
	if pos.CastleRights() != "-" {
		t.Fatalf("Wrong castling rights. Expected: %s, but got %s", "-", pos.CastleRights())
	}
}

func TestPositionMarshalling2(t *testing.T) {
	pos := &Position{}
	expected := "r4rk1/1b2bppp/ppq1p3/2pp3n/5P2/1P1BP3/PBPPQ1PP/R4RK1 w e4 - 0 1"
	err := pos.UnmarshalText([]byte(expected))
	if err == nil {
		t.Fatal("Missing error")
	}
}
