package othello

import (
	"bytes"
	"io/ioutil"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"encoding/json"
	"fmt"
	//"math/rand"
	"net/http"
)

func init() {
	http.HandleFunc("/", getMove)
}

type Game struct {
	Board Board `json:board`
}

var depth = 0
var limit = 30
//var num_boards int = 0

// Provide a generic handler for move requests. If no board state is
// specified then a simple HTML form is provided to let users paste
// JSON state (which can be copy-pasted from a game running on
// http://step-reversi.appspot.com/ ).
func getMove(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var js []byte
	defer r.Body.Close()
	js, _ = ioutil.ReadAll(r.Body)
	if len(js) < 1 {
		js = []byte(r.FormValue("json"))
	}
	if len(js) < 1 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
<body><form method=get>
Paste JSON here:<p/><textarea name=json cols=80 rows=24></textarea>
<p/><input type=submit>
</form>
</body>`)
		return
	}
	var game Game
	err := json.Unmarshal(js, &game)
	if err != nil {
		fmt.Fprintf(w, "invalid json %v? %v", string(js), err)
		return
	}
	board := game.Board
	log.Infof(ctx, "got board: %v", board)
	moves := board.ValidMoves()
	if len(moves) < 1 {
		fmt.Fprintf(w, "PASS")
		return
	} else {
		// NOTE TO STUDENTS: This next line is the main line you'll want to
		// change.  Right now this is just picking a random move out of the
		// list of possible moves, but you'll want to make this choose a
		// better move (probably using some game tree traversal algorithm
		// like MinMax).
		move := getBestMove(board, moves)
		//move := moves[rand.Intn(len(moves))]
		fmt.Fprintf(w, "[%d,%d]", move.Where[0], move.Where[1])
		//fmt.Fprint(w, num_boards)
	}
}


func getScore(board Board)int{
	score := 0
	whitescore := 0
	// depth, limit: global variable (var depth = 0, var limit = 3)
	if depth != limit {
		if board.Next == Black{
			depth += 1
			score = getMaxMoveScore(board)
		} else if board.Next == White {
			depth += 1
			score = getMinMoveScore(board)
		}
	} else{
		//num_boards += 1
		x := 0
		for x >= 0 && x <= 7{
			y := 0
			for y >= 0 && y <= 7{
				if board.Pieces[y][x] == Black{
					score += 4

					if x == 1 {
						score -= 5
					} else if x == 6 {
						score -= 5
					}
					if y == 1 {
						score -= 5
					} else if y == 6 {
						score -= 5
					}

					if (x == 1 && y == 1) || (x == 1 && y == 6) || (x == 6 && y == 1) || (x == 6 && y == 6) {
						score -= 5
					}
					if (x == 0 && y == 1) || (x == 0 && y == 6) || (x == 1 && y == 0) || (x == 0 && y == 7) || (x == 6 && y == 0) || (x == 6 && y == 7) || (x == 7 && y == 1) || (x == 7 && y == 6) {
						score -= 10
					}

					if x == 0 {
						score += 5
					} else if x == 7 {
						score += 5
					}
					if y == 0 {
						score += 5
					} else if y == 7 {
						score += 5
					}

					if (x == 0 && y == 0) || (x == 0 && y == 7) || (x == 7 && y == 0) || (x == 7 && y == 7) {
						score += 50
					}
				} else if board.Pieces[y][x] == White{
					whitescore += 4

					if x == 1 {
						whitescore -= 5
					} else if x == 6 {
						whitescore -= 5
					}
					if y == 1 {
						whitescore -= 5
					} else if y == 6 {
						whitescore -= 5
					}

					if (x == 1 && y == 1) || (x == 1 && y == 6) || (x == 6 && y == 1) || (x == 6 && y == 6) {
						whitescore -= 5
					}
					if (x == 0 && y == 1) || (x == 0 && y == 6) || (x == 1 && y == 0) || (x == 0 && y == 7) || (x == 6 && y == 0) || (x == 6 && y == 7) || (x == 7 && y == 1) || (x == 7 && y == 6) {
						whitescore -= 10
					}

					if x == 0 {
						whitescore += 5
					} else if x == 7 {
						whitescore += 5
					}
					if y == 0 {
						whitescore += 5
					} else if y == 7 {
						whitescore += 5
					}

					if (x == 0 && y == 0) || (x == 0 && y == 7) || (x == 7 && y == 0) || (x == 7 && y == 7) {
						whitescore += 50
					}
				}
				score -= whitescore
				whitescore = 0
				y += 1
			}
			x += 1
		}
		if board.Next == Black{
			score += len(board.ValidMoves()) * 3
		} else if board.Next == White{
			score -= len(board.ValidMoves()) * 3
		}
	}
	return score
}

// turn black
func getMaxMoveScore(board Board)int{
	max_score := -500
	for _, move := range board.ValidMoves() {
		next_board, _ := board.After(move)
		score := getScore(next_board)
		if score >= max_score{
			max_score = score
		}
	}
	return max_score
}

// turn white
func getMinMoveScore(board Board)int{
	min_score := 500
	for _, move := range board.ValidMoves() {
		next_board, _ := board.After(move)
		score := getScore(next_board)
		if score <= min_score{
			min_score = score
		}
	}
	return min_score
}


func getBestMove(board Board, moves []Move)Move{
	score := 0
	var best_move Move
	max_score := -500
	min_score := 500
	for _, move := range moves{
		next_board, _ := board.After(move)
		score = getScore(next_board)
		if board.Next == Black{
			if score  > max_score {
				best_move = move
				max_score = score
			}
		} else if board.Next == White{
			if score  < min_score {
				best_move = move
				min_score = score
			}
		}

	}
	return best_move
}

type Piece int8

const (
	Empty Piece = iota
	Black Piece = iota
	White Piece = iota

	// Red/Blue are aliases for Black/White
	Red  = Black
	Blue = White
)

func (p Piece) Opposite() Piece {
	switch p {
	case White:
		return Black
	case Black:
		return White
	default:
		return Empty
	}
}

type Board struct {
	// Layout says what pieces are where.
	Pieces [8][8]Piece
	// Next says what the color of the next piece played must be.
	Next Piece
}

// Position represents a position on the othello board. Valid board
// coordinates are 1-8 (not 0-7)!
type Position [2]int

// Valid returns true iff this is a valid board position.
func (p Position) Valid() bool {
	ok := func(i int) bool { return 1 <= i && i <= 8 }
	return ok(p[0]) && ok(p[1])
}

// Pass returns true iff this move position represents a pass.
func (p Position) Pass() bool {
	return !p.Valid()
}

// Move describes a move on an Othello board.
type Move struct {
	// Where a piece is going to be placed. If Where is zeros, or
	// another invalid coordinate, it indicates a pass.
	Where Position
	// As is the player taking the player taking the turn.
	As Piece
}

// At returns a pointer to the piece at a given position.
func (b *Board) At(p Position) *Piece {
	return &b.Pieces[p[1]-1][p[0]-1]
}

// Get returns the piece at a given position.
func (b *Board) Get(p Position) Piece {
	return *b.At(p)
}

// Exec runs a move on a given Board, updating the given board, and
// returning it. Returns error if the move is illegal.
func (b *Board) Exec(m Move) (*Board, error) {
	if !m.Where.Pass() {
		if _, err := b.realMove(m); err != nil {
			return b, err
		}
	} else {
		// Attempting to pass.
		valid := b.ValidMoves()
		if len(valid) > 0 {
			return nil, fmt.Errorf("%v illegal move: there are valid moves available: %v", m, valid)
		}
	}
	b.Next = b.Next.Opposite()
	return b, nil
}

// Clone makes a new identical copy of an existing board and returns a
// pointer to it.
func (b *Board) Clone() *Board {
	clone := *b
	return &clone
}

// Returns the state of a new board after the given move. Returns an
// unchanged board and an error if the move is illegal.
func (b Board) After(m Move) (Board, error) {
	if _, err := b.Exec(m); err != nil {
		return b, err
	}
	return b, nil
}

// realMove executes a move that isn't a PASS. Use Exec instead to
// execute any move (include PASS moves).
func (b *Board) realMove(m Move) (*Board, error) {
	captures, err := b.tryMove(m)
	if err != nil {
		return nil, err
	}

	for _, p := range append(captures, m.Where) {
		*b.At(p) = m.As
	}
	return b, nil
}

type direction Position

var dirs []direction

func init() {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			dirs = append(dirs, direction{x, y})
		}
	}
}

// tryMove tries a non-PASS move without actually executing it.
// Returns the list of captures that would happen.
func (b *Board) tryMove(m Move) ([]Position, error) {
	if b.Get(m.Where) != Empty {
		return nil, fmt.Errorf("%v illegal move: %v is occupied by %v", m, m.Where, b.Get(m.Where))
	}

	var captures []Position
	for _, dir := range dirs {
		captures = append(captures, b.findCaptures(m, dir)...)
	}

	if len(captures) < 1 {
		return nil, fmt.Errorf("%v illegal move: no pieces were captured", m)
	}
	return captures, nil
}

func translate(p Position, d direction) Position {
	return Position{p[0] + d[0], p[1] + d[1]}
}

func (b *Board) findCaptures(m Move, dir direction) []Position {
	var caps []Position
	for p := m.Where; true; caps = append(caps, p) {
		p = translate(p, dir)
		if !p.Valid() {
			// End of board.
			return []Position{}
		}
		switch *b.At(p) {
		case m.As:
			return caps
		case Empty:
			return []Position{}
		}
	}
	panic("impossible")
}

// Returns a slice of valid moves for the given Board.
func (b *Board) ValidMoves() []Move {
	var moves []Move
	for y := 1; y <= 8; y++ {
		for x := 1; x <= 8; x++ {
			m := Move{Where: Position{x, y}, As: b.Next}
			_, err := b.tryMove(m)
			if err == nil {
				moves = append(moves, m)
			}
		}
	}
	return moves
}

// Converts a Board into a human-readable ASCII art diagram.
func (b Board) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("\n")
	buf.WriteString(" |ABCDEFGH|\n")
	buf.WriteString("-+--------+\n")
	for y := 0; y < 8; y++ {
		fmt.Fprintf(buf, "%d|", y+1)
		for x := 0; x < 8; x++ {
			p := b.Pieces[y][x]
			switch p {
			case Red:
				buf.WriteString("X")
			case Blue:
				buf.WriteString("O")
			default:
				buf.WriteString(" ")
			}
		}
		fmt.Fprintf(buf, "|%d\n", y+1)
	}
	buf.WriteString("-+--------+\n")
	buf.WriteString(" |ABCDEFGH|\n")
	return buf.String()
}
