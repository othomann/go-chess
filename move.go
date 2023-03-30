package chess

// A MoveTag represents a notable consequence of a move.
type MoveTag uint16

const (
	// KingSideCastle indicates that the move is a king side castle.
	KingSideCastle MoveTag = 1 << iota
	// QueenSideCastle indicates that the move is a queen side castle.
	QueenSideCastle
	// Capture indicates that the move captures a piece.
	Capture
	// EnPassant indicates that the move captures via en passant.
	EnPassant
	// Check indicates that the move puts the opposing player in check.
	Check
	// inCheck indicates that the move puts the moving player in check and
	// is therefore invalid.
	inCheck
)

const (
	DEFAULT       = 1
	CASTLE        = 10
	CAPTURE       = 100
	EN_PASSANT    = 100
	CHECK         = 1000
	CAPTURE_CHECK = 10000
)

// A Move is the movement of a piece from one square to another.
type Move struct {
	s1    Square
	s2    Square
	promo PieceType
	tags  MoveTag
}

// String returns a string useful for debugging.  String doesn't return
// algebraic notation.
func (m *Move) String() string {
	return m.S1().String() + m.S2().String() + m.Promo().String()
}

// S1 returns the origin square of the move.
func (m *Move) S1() Square {
	return m.s1
}

// S2 returns the destination square of the move.
func (m *Move) S2() Square {
	return m.s2
}

// Promo returns promotion piece type of the move.
func (m *Move) Promo() PieceType {
	return m.promo
}

// HasTag returns true if the move contains the MoveTag given.
func (m *Move) HasTag(tag MoveTag) bool {
	return (tag & m.tags) > 0
}

func (m *Move) addTag(tag MoveTag) {
	m.tags = m.tags | tag
}

type MoveSlice []*Move

func (a MoveSlice) find(m *Move) *Move {
	if m == nil {
		return nil
	}
	for _, move := range a {
		if move.String() == m.String() {
			return move
		}
	}
	return nil
}

func (a MoveSlice) Len() int {
	return len(a)
}

func (a MoveSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a MoveSlice) Less(i, j int) bool {
	return compare(a[i]) > compare(a[j])
}

func compare(move *Move) int {
	if move.HasTag(Capture) {
		if move.HasTag(Check) {
			return CAPTURE_CHECK
		} else {
			return CAPTURE
		}
	}
	if move.HasTag(Check) {
		return CHECK
	}
	if move.HasTag(QueenSideCastle) || move.HasTag(KingSideCastle) {
		return CASTLE
	}
	if move.HasTag(EnPassant) {
		return EN_PASSANT
	}
	return DEFAULT
}
