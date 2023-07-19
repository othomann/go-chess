package algorithm

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/othomann/go-chess"
)

type MateMove struct {
	Move     chess.Move
	Mobility int
}

func (move *MateMove) isCheck() bool {
	return move.Move.HasTag(chess.Check)
}

type MateNode struct {
	Children []*MateNode
	Parent   *MateNode
	Move     MateMove
	Color    chess.Color
	FEN      string
	Depth    int
}

func (mateNode *MateNode) isRoot() bool {
	return mateNode.Parent == nil
}

func (mateNode *MateNode) Root() *MateNode {
	node := mateNode
	for !node.isRoot() {
		node = node.Parent
	}
	return node
}
func reverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

func (mateNode *MateNode) String() (string, error) {
	if mateNode.isRoot() {
		fen, err := chess.FEN(mateNode.FEN)
		if err != nil {
			return "", err
		}
		options := make([]func(*chess.Game), 0)
		options = append(options, fen)
		options = append(options, chess.UseNotation(chess.UCINotation{}))
		game := chess.NewGame(options...)

		moveOutput := mateNode.uciPgn()
		moves := strings.Fields(moveOutput)
		reverseSlice(moves)
		for _, move := range moves {
			err := game.MoveStr(move)
			if err != nil {
				return "", fmt.Errorf("invalid move: %s", move)
			}
		}
		game.SetNotation(chess.AlgebraicNotation{})
		return game.String(), nil
	}
	return "", fmt.Errorf("invalid node")
}

func (mateNode *MateNode) Print() {
	if mateNode.isRoot() {
		depth := 0
		if mateNode.Children != nil && len(mateNode.Children) != 0 {
			for _, child := range mateNode.Children {
				child.printString0(depth + 1)
			}
		}
	}
}

func (mateNode *MateNode) printString0(depth int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", depth), mateNode.Move.Move.String())
	if mateNode.Children != nil && len(mateNode.Children) != 0 {
		for _, child := range mateNode.Children {
			child.printString0(depth + 1)
		}
	}
}

func (mateNode *MateNode) uciPgn() string {
	return mateNode.uciPgn0(0, mateNode.Root().Depth)
}

func (mateNode *MateNode) uciPgn0(depth, max int) string {
	if depth >= max {
		s := ""
		current := mateNode
		for current.Parent != nil {
			s += fmt.Sprintf(" %s", current.Move.Move.String())
			current = current.Parent
		}
		return s
	}
	if mateNode.Children != nil && len(mateNode.Children) != 0 {
		for _, child := range mateNode.Children {
			result := child.uciPgn0(depth+1, max)
			if result != "" {
				return result
			}
		}
	}
	return ""
}

func NewRoot(game *chess.Game, depth int) *MateNode {
	return &MateNode{
		Color: game.Position().Turn(),
		FEN:   game.FEN(),
		Depth: 1 + (depth-1)*2,
	}
}

func (node *MateNode) add(move MateMove, turn chess.Color) *MateNode {
	if node.Children == nil {
		node.Children = make([]*MateNode, 0)
	}
	newNode := &MateNode{
		Parent: node,
		Move:   move,
		Color:  turn,
	}
	node.Children = append(node.Children, newNode)
	return newNode
}

func (node *MateNode) Remove(move MateMove) {
	if node.Children == nil {
		return
	}
	found := -1
	for index, mateNode := range node.Children {
		mateMove := mateNode.Move
		if mateMove.Move.S1() == move.Move.S1() && mateMove.Move.S2() == move.Move.S2() {
			found = index
		}
	}
	if found != -1 {
		node.Children = append(node.Children[:found], node.Children[found+1:]...)
	}
}

type MateMoveSlice []*MateMove

func (a MateMoveSlice) Len() int {
	return len(a)
}

func (a MateMoveSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a MateMoveSlice) Less(i, j int) bool {
	return a[i].Mobility < a[j].Mobility
}

func createSearchableMoves(game *chess.Game) (MateMoveSlice, error) {
	nextMoves := game.ValidMoves()
	result := make([]*MateMove, len(nextMoves))

	for index, move := range nextMoves {
		err := game.Move(move)
		if err != nil {
			return nil, err
		}
		mobility := len(game.ValidMoves())
		result[index] = &MateMove{
			Mobility: mobility,
			Move:     *move,
		}
		err = game.UndoMove()
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(MateMoveSlice(result))
	return result, nil
}

func MateSearch(game *chess.Game, maximum int, mateNode *MateNode) (bool, *MateNode) {
	return mateSearch0(game, 1, maximum, mateNode)
}

func mateSearch0(game *chess.Game, depth int, maximum int, mateNode *MateNode) (bool, *MateNode) {
	next, err := createSearchableMoves(game)
	if err != nil {
		return false, mateNode.Root()
	}
	currentResult := mateNode
	/* loop */ for _, mateMove := range next {
		current := currentResult.add(*mateMove, game.Position().Turn())
		if (depth == maximum && mateMove.isCheck()) || depth != maximum {
			mobility := mateMove.Mobility
			if mobility == 0 {
				if !mateMove.isCheck() {
					// statemate
					currentResult.Remove(*mateMove)
					continue // contiue loop
				}
				// mate
				if depth != 1 || depth == maximum {
					return true, mateNode.Root()
				}
				currentResult.Remove(*mateMove)
			} else if depth < maximum {
				err := game.Move(&mateMove.Move)
				if err != nil {
					fmt.Printf("Invalid move: %s", &mateMove.Move)
				}
				moveCounter := 0
				currentResult = current
				opponentNextMoves, err := createSearchableMoves(game)
				if err != nil {
					return false, mateNode.Root()
				}
				/* opponentLoop */
				for _, opponentMateMove := range opponentNextMoves {
					opponentCurrent := currentResult.add(*opponentMateMove, game.Position().Turn())
					currentResult = opponentCurrent
					err := game.Move(&opponentMateMove.Move)
					if err != nil {
						return false, mateNode.Root()
					}
					result, _ := mateSearch0(game, depth+1, maximum, currentResult)
					if !result {
						err := game.UndoMove()
						if err != nil {
							return false, mateNode.Root()
						}
						currentResult = currentResult.Parent
						currentResult.Remove(*opponentMateMove)
						break /* break opponent loop */
					}
					currentResult = currentResult.Parent
					moveCounter++
					err = game.UndoMove()
					if err != nil {
						return false, mateNode.Root()
					}
				}

				currentResult = currentResult.Parent
				err = game.UndoMove()
				if err != nil {
					return false, mateNode.Root()
				}

				if mobility == moveCounter {
					return true, mateNode.Root()
				}
				currentResult.Remove(*mateMove)
			} else {
				currentResult.Remove(*mateMove)
				return false, mateNode
			}
		} else {
			currentResult.Remove(*mateMove)
		}
	}
	return false, mateNode
}
