package main

type Piece struct {
	moves         *[][][] int
	computerPiece bool
	value         int
	name          string
	imageIndex    int
}

func (p *Piece) Value() int {
	return p.value
}

func (p *Piece) SetValue(value int) {
	p.value = value
}

// Contruct a new Piece
func NewPiece(moves *[][][] int, computerPiece bool, value int, imageIndex int) *Piece {
	var p Piece
	p.moves = moves
	p.computerPiece = computerPiece
	p.value = value
	if computerPiece == false {
		p.value = -p.value
	}
	p.imageIndex = imageIndex
	return &p
}

// var board [64] *Piece
// var moves [1000] int
// Create the singletons for each piece on the board, note that some of them will be used in multiple places i.e
// be on multiple places on the board

var CPawn = NewPiece(&ComputerPawnmoves, true, 150, 6)
var CTower = NewPiece(&Towermoves, true, 850, 9)
var CBishop = NewPiece(&Bishopmoves, true, 525, 8)
var CKnight = NewPiece(&Knightmoves, true, 525, 7)
var CQueen = NewPiece(&Queenmoves, true, 1535, 10)
var CKing = NewPiece(&Kingmoves, true, 10000, 11)

// Opponent pieces
var OPawn = NewPiece(&OpponentPawnmoves, false, 150, 0)
var OTower = NewPiece(&Towermoves, false, 850, 3)
var OBishop = NewPiece(&Bishopmoves, false, 525, 2)
var OKnight = NewPiece(&Knightmoves, false, 525, 1)
var OQueen = NewPiece(&Queenmoves, false, 1535, 4)
var OKing = NewPiece(&Kingmoves, false, 10000, 5)

func getMoveArray(moveArray *[maxPossibleMoves] int, board *[64] *Piece, computer bool) int {
	var moveCount = 0
	var currentPiece, targetPiece *Piece
	for boardPosition := 0; boardPosition < 64; boardPosition++ {
		currentPiece = board[boardPosition]
		// are we playing the right piece
		if currentPiece != nil && currentPiece.computerPiece == computer {
			switch currentPiece {
			// take care of all pieces except Pawns
			case CBishop, CKing, CKnight, CQueen, CTower,
				OBishop, OKing, OKnight, OQueen, OTower:
				// Iterate over the various directions a piece can go
				for direction := 0; direction < len((*currentPiece.moves)[boardPosition]); direction++ {
					// position indicates a move in a given direction
					for position := 0; position < len((*currentPiece.moves)[boardPosition][direction]); position++ {
						targetPiece = board[(*currentPiece.moves)[boardPosition][direction][position]]
						//if target is empty, it is a possible move
						if targetPiece == nil {
							moveArray[moveCount] = boardPosition
							moveArray[moveCount+1] = (*currentPiece.moves)[boardPosition][direction][position]
							moveCount += 2
							continue
							// else if it is a piece of the opposite color
						} else if (*targetPiece).computerPiece != computer {
							moveArray[moveCount] = boardPosition
							moveArray[moveCount+1] = (*currentPiece.moves)[boardPosition][direction][position]
							moveCount += 2
						}
						//we are here when the target is occupied by a piece of the same color or a piece we have taken
						// so we stop going in that direction
						break
					}
				}
				//
			case CPawn, OPawn:
				for direction := 0; direction < 4; direction++ {
					// This is Pawn we are going to play, the first 4 direction are "take" directions
					var currentTarget = (*currentPiece.moves)[boardPosition][direction]
					if len(currentTarget) != 0 {
						if board[currentTarget[0]] != nil {
							// if we take by moving two rows only do the move if the fifth move (push by one row) is empty
							// the take by moving two rows is for direction 1 and 3
							if direction&1 == 1 {
								if board[(*currentPiece.moves)[boardPosition][5][0]] == nil {
									continue
								}
							}
							// piece is of the opposite color
							if board[currentTarget[0]].computerPiece != computer {
								moveArray[moveCount] = boardPosition
								moveArray[moveCount+1] = currentTarget[0]
								moveCount += 2
							}
						}
					}
				}
				// next two directions are push direction 5 contains the push by 1, direction 5 the push by 2
				var currentTarget = (*currentPiece.moves)[boardPosition][4]
				if len(currentTarget) != 0 && board[currentTarget[0]] == nil {
						// the target is empty
						moveArray[moveCount] = boardPosition
						moveArray[moveCount+1] = currentTarget[0]
						moveCount += 2
						// you can only move 2 if you can move 1
						currentTarget = (*currentPiece.moves)[boardPosition][5]
						// the target is empty
						if len(currentTarget) != 0 && board[currentTarget[0]] == nil {
							moveArray[moveCount] = boardPosition
							moveArray[moveCount+1] = currentTarget[0]
							moveCount += 2
					}
				}
			case nil:
				// do nothing
			default:
			}
		}
	}
	return moveCount
}

