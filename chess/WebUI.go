package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// Note: we only support one client
var b Board

// the javascript and html used for the UI
const htmlTemplate = "<!DOCTYPE html><html><head>\n" +
"<script>\n " +
	" var editmode = false;" +

	"function onclickreset_handler(ev) {\n" +
	"ev.preventDefault();\n"+
	"console.log(\"onclickreset is working\");\n" +
	"window.location.replace( window.location.origin + \"?reset=yes\");\n" +
	"}\n " +



	"function dragover_handler(ev) {\n" +
	"ev.preventDefault();\n"+
	"console.log(\"Drag over is working\");\n" +
	//	" console.log(\"source \" + ev.srcElement.id + \" destination \" + ev.target.id );\n" +
	"}\n " +

	"function startdrag_handler(ev) {\n" +
	"ev.preventDefault();\n"+
	"console.log(\"startDrag is working\");\n" +
	" ev.dataTransfer.setData(\"start\", ev.target.id);\n" +
//	" console.log(\"source \" + ev.target.id );\n" +
	"}\n " +

	"function dropover_handler(ev) {\n" +
	"ev.preventDefault();\n"+
	"console.log(\"Drop is working\");\n" +
//	"console.log(\"source \" + ev.dataTransfer.getData(\"start\") + \" destination \" + ev.target.id );\n" +
//	"console.log( window.location.origin + \"?start=\" + ev.dataTransfer.getData(\"start\") + \"&dest=\" + ev.target.id );\n" +
	"window.location.replace( window.location.origin + \"?start=\" + ev.dataTransfer.getData(\"start\") + \"&dest=\" + ev.target.id );\n" +
	"}\n  " +
	"</script> \n" +
	"\n</head><html>\n<body>\n"


// There is a bug, since the http calls can be aynch's the UI might ask for a refresh as we gernerate the board
func GenerateHTMLBoard(b *Board, elapsed float64 ) string {
	var s = htmlTemplate
	for i := 0 ; i < 64 ; i++ {
		// change line
		if i & 7 == 0 {
			s = s + "<br>"
		}
		s = s + "<img height=50 width=50 id=\"" + strconv.Itoa(i) + "\" ondrop=\"dropover_handler(event)\"" +" ondragover=\"dragover_handler(event)\" " +
			  " ondragstart=\"event.dataTransfer.setData('start','" + strconv.Itoa(i) +"')\" "
		// if the position is not empty put the piece in it
		if b.BoardArray[i] != nil {
			s = s + "src=/gifs/"+ strconv.Itoa(b.BoardArray[i].imageIndex)+".gif"
			if b.BoardArray[i].computerPiece {
				s = s +  " draggable=false"
			} else {
				s = s +  " draggable=true"
			}
		}
		// alternate color
		if  (i + i >> 3)  & 1 == 1 {
			s = s + " style=background-color:grey"
		}
		s = s + ">"
	}
	s = s + "<BR><BR> Analyzed : " + strconv.Itoa(b.numLeaf) + " positions in " + fmt.Sprintf("%.2f seconds %.0f positions/seconds", elapsed , float64(b.numLeaf)/elapsed)
	s = s +  "<BR> Score: \n" + strconv.Itoa(b.currentPositionScore)
	s = s + "<BR><button onclick=\"onclickreset_handler(event)\">reset board</button>"

	s = s + "<body>\n<html>"
	return s
}

//  handler for the browser request
func handler(w http.ResponseWriter, r *http.Request) {
	var startTime time.Time
	var elapsed time.Duration
	b.Mutex.Lock()
	defer b.Unlock()
	if r.URL.Query()["start"] != nil && r.URL.Query()["dest"] != nil {
		start, _ := strconv.Atoi(r.URL.Query()["start"][0])
		dest, _ :=  strconv.Atoi(r.URL.Query()["dest"][0])
		fmt.Println("Request for page")
		fmt.Println("start ",r.URL.Query()["start"][0])
		fmt.Println("dest ",r.URL.Query()["dest"][0])
		b.move(start, dest)
		startTime = time.Now()
		b.numLeaf = 0
		b.currentPositionScore = b.Evaluate(0, math.MinInt32, math.MaxInt32, maxNodeCount)
		fmt.Printf("Move From %d to %d \n", b.BestMoveStart[0], b.BestMoveEnd[0])
		b.move(b.BestMoveStart[0], b.BestMoveEnd[0])
		elapsed = time.Since(startTime)
	} else {
		if r.URL.Query()["reset"] != nil {
			b.initBoard()
		}
	}
	_, _ = fmt.Fprintf(w, GenerateHTMLBoard(&b, elapsed.Seconds()))

}


func main() {
//	b.initBoard()

	b.BoardArray[6] = CKing
	b.BoardArray[0] = CTower
	b.BoardArray[7] = CTower
	b.BoardArray[57] = OKing


	// set the handler for the static loading of the images
	http.Handle("/gifs/", http.StripPrefix("/gifs/", http.FileServer(http.Dir("./gifs"))))
	// set the handler for the browser URL request
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
