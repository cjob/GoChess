package main

import (
    "math"
    "sync"
)

const maxDepth= 6
const maxPossibleMoves=1000
// const minDepth= 6
const maxNodeCount = 10000000

type Board struct {
	BoardArray [64] *Piece // An array of pointers to the Pieces, the Pieces are singletons
	maxDepthReached int // maximum depth reached by the engine
	moveArray[maxDepth+1][maxPossibleMoves] int //
	BestMoveStart[maxDepth] int
	BestMoveEnd[maxDepth] int
	numLeaf int // number of leaf nodes explored during one move
	currentPositionScore int // current score of the board as seen by the computer
	sync.Mutex
}

func (b *Board) move(src int, dest int) {

    // fmt.Printf("Move from %d to %d", src, dest)
    b.BoardArray[dest] = b.BoardArray[src]
    b.BoardArray[src] = nil

    // Promote the pawn to a queen (BUG!: it could be a knight as well)
    if (b.BoardArray[dest] == CPawn) && (dest > 55) {
        b.BoardArray[dest] = CQueen
    } else
    // Promote the pawn to a queen (BUG!: it could be a knight as well)
    if (b.BoardArray[dest] == OPawn) && (dest <8)  {
        b.BoardArray[dest] = OQueen
    }
}



func (b *Board) initBoard() {
    // I am sure there is a better way of clearing an array of pointers
    for i := 0 ; i < 63 ; i++ {
        b.BoardArray[i] = nil
    }

    b.BoardArray[0] = CTower
    b.BoardArray[1] = CKnight
    b.BoardArray[2] = CBishop
    b.BoardArray[4] = CKing
    b.BoardArray[3] = CQueen
    b.BoardArray[5] = CBishop
    b.BoardArray[6] = CKnight
    b.BoardArray[7] = CTower
    b.BoardArray[8] = CPawn
    b.BoardArray[9] = CPawn
    b.BoardArray[10] = CPawn
    b.BoardArray[11] = CPawn
    b.BoardArray[12] = CPawn
    b.BoardArray[13] = CPawn
    b.BoardArray[14] = CPawn
    b.BoardArray[15] = CPawn


    b.BoardArray[56] = OTower
    b.BoardArray[57] = OKnight
    b.BoardArray[58] = OBishop
    b.BoardArray[60] = OKing
    b.BoardArray[59] = OQueen
    b.BoardArray[61] = OBishop
    b.BoardArray[62] = OKnight
    b.BoardArray[63] = OTower
    b.BoardArray[48] = OPawn
    b.BoardArray[49] = OPawn
    b.BoardArray[50] = OPawn
    b.BoardArray[51] = OPawn
    b.BoardArray[52] = OPawn
    b.BoardArray[53] = OPawn
    b.BoardArray[54] = OPawn
    b.BoardArray[55] = OPawn

}


/*
    func main() {
        var b Board
        var startTime time.Time
        var elapsed time.Duration
        b.initBoard()
        b.numLeaf = 0
        reader := bufio.NewReader(os.Stdin)
        var start, end int


       // fmt.Printf("number of move %d", getMoveArray(&b.moveArray[0], &b.BoardArray, true))
        // fmt.Printf("Moves %v", b.moveArray[0]);
    for true  {
            b.numLeaf = 0
            fmt.Printf("Thinking ...%#U ♔ starts at \n", 2654)
            // do not evaluate more than maxNodeCount nodes
            startTime = time.Now()
            b.currentPositionScore = b.Evaluate(0, math.MinInt32, math.MaxInt32, maxNodeCount)
            fmt.Printf("Move From %d to %d \n", b.BestMoveStart[0], b.BestMoveEnd[0])
            b.move(b.BestMoveStart[0], b.BestMoveEnd[0])
            elapsed = time.Since(startTime)
            fmt.Printf("Looked at : %d nodes in %s seconds \n ", b.numLeaf, elapsed )
            fmt.Printf("Current Score : %d \n", b.currentPositionScore)
            startS, _ := reader.ReadString('\n')
            endS, _ := reader.ReadString('\n')
            start, _ = strconv.Atoi(startS)
            end, _ = strconv.Atoi(endS)
            b.move(start,end)
        }
    }

*/


    // Evaluate a leaf node by summing up the value of the different pieces
    // subtracting the value of the opponent pieces
    // note a big optimization is to instead of computing from scratch update incrementally the value when you execute the move
    func (b *Board) EvaluateLeaf() int {
        var tempTotal  = 0
        var p *Piece
        for i := 0; i < 64; i++ {
            p = b.BoardArray[i]
            if p != nil {
                tempTotal += p.Value()
                // value more the pawns that are further down the board
                if p == CPawn {
                    tempTotal += i >> 4
                } else {
                    if p == OPawn {
                        tempTotal -= (63 - i) >> 4
                    }
                }

            }
        }
        b.numLeaf++
        // fmt.Println("Leaf evaluation : %d" ,tempTotal)
        return tempTotal
    }

    func  max(a int, b int) int {
        if a  > b {
                return a
        }
        return b
    }


    func min(a int, b int) int {
        if a < b {
            return a
        }
        return b
    }


    func (b *Board) Evaluate( depth int, A int, B int, currentNodeCount int) int {

           var lastTarget,lastSource *Piece
           var moveCount int
 //          var checkMoveCount int
           var alpha,beta,val,src,dest int
//           var check bool;
 //          var finalNodeCount int
 //          var checkMoveArray[1000]int
 //           fmt.Println(" depth %d Alpha %d Beta %B)", depth, A, B)
           computer := (depth & 1) == 0

           // if we have reached the maximum depth of the tree
           if depth == maxDepth {

                // ||
              // or of we are past the minimum and there is not enough nodes left to go another round
              // be careful that we always need to evaluate at the same depth
              // (depth >= minDepth) && (currentNodeCount < 2500) && ((depth & 1) == (minDepth & 1)))  {
                b.maxDepthReached = max(b.maxDepthReached,depth)
                    return b.EvaluateLeaf()
           } else {

               alpha = math.MinInt32
               beta = math.MaxInt32


               // finalNodeCount = b.numLeaf + currentNodeCount;
               // get all the moves
               moveCount = getMoveArray(&b.moveArray[depth], &b.BoardArray, computer )
               // play the moves one by one
               for i :=0 ; i < moveCount ; i+=2 {

                   dest= b.moveArray[depth][i+1]
                   src=b.moveArray[depth][i]
                   // save the the piece at the target to be able to undo the move
                   lastTarget = b.BoardArray[dest]
                   lastSource = b.BoardArray[src]


                   // if one of the available move is to take the King,
                   // no need to explore other moves
                   if lastTarget== OKing {
                       // substract the depth so that the engine chooses the path to the shortest checkmate
                       return math.MaxInt32 - depth
                   }
                   if lastTarget== CKing {
                       return math.MinInt32 + depth
                   }

                   // make the move
                   b.move(src,dest)


                   // if it is a MAX node
                   if computer {
                       val = b.Evaluate(depth+1,max(A,alpha),B,currentNodeCount- b.numLeaf)

/*
                       // if the opponent is going to loose his King at the next move, make sure that he is in check
                       // otherwise it is a draw
                       if (val == math.MaxInt32 - depth - 2) {
                           // get all the moves
                           checkMoveCount = getMoveArray(&checkMoveArray, &b.BoardArray, computer );
                           // can the computer take it already
                           check = false;
                           for j:=1; j < checkMoveCount ; j+=2 {

                               if (b.BoardArray[checkMoveArray[j]] == OKing) {
                                   check = true;
                                   break;
                               }
                           }
                           // if the opponent was not in check already, it is a pat
                           if (!check) {
                               val = math.MaxInt32 / 2
                           }
                       }
*/
                       // undo the move
                       b.BoardArray[src] = lastSource
                       b.BoardArray[dest]= lastTarget


                       // if it found a batter move
                       if val > alpha {
                           //fmt.Printf("Found a better move %d %d ", src, dest)
                           alpha = val
                           b.BestMoveStart[depth] = src
                           b.BestMoveEnd[depth] = dest
                       }
                       // if it is too good, return. We know the opponent wont take us down that path
                       if alpha >= B {
                           return alpha
                       }
                   } else {
                       // playing the opponent
                       val = b.Evaluate(depth+1,A,min(B,beta),currentNodeCount- b.numLeaf)
/*
                       // if the computer is going to loose his King, make sure that he is in check
                       if (val == math.MinInt32 + depth + 2 ) {

                           // get all the moves
                           checkMoveCount = getMoveArray(&checkMoveArray, &b.BoardArray, computer )
                           // can the computer take it already
                           check= false;
                           for j:=1; j < checkMoveCount ; j+=2 {
                           if (b.BoardArray[checkMoveArray[j]] == CKing) {
                               check = true;
                               break;
                           }
                           }

                           // if the computer was not in check already, it is a pat
                           if (!check) {
                               val = math.MinInt32 / 2
                           }
                       }
*/
                       // undo the move
                       b.BoardArray[src] = lastSource
                       b.BoardArray[dest]= lastTarget

                       if val < beta {
                           // fmt.Printf("Found a better move %d %d ", src, dest)
                           beta = val
                           b.BestMoveStart[depth] = src
                           b.BestMoveEnd[depth] = dest
                       }
                       // if the move is bad (i.e worse than A return. We assume the opponent will not take us down that path
                       if A >= beta  {
                           return beta
                       }
                   }

               }
           }
        if computer {
            // fmt.Printf("Board evaluation returns (computer/alpha) %d", alpha);
            return alpha

        }
        // fmt.Printf("Board evaluation returns (Opponent/beta) %d", beta);
        return beta
    }





