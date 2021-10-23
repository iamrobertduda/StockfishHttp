package internal

import (
	"io"
	"log"
	"os/exec"
	"strings"
)

type Game struct {
	Fen        string
	Playermove string
	Answer     string
	Depth      string
}

func getStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

func fetchFen(out string) (result string) {
	result, found := getStringInBetweenTwoString(out, "Fen: ", "\nKey:")
	if found != true {
		log.Fatal("Not able to find FEN.")
	}
	return strings.ReplaceAll(result, "\r", "")
}

func fetchBestMove(out string) (result string) {
	result, found := getStringInBetweenTwoString(out, "bestmove ", " ponder")
	if found != true {
		return out[strings.Index(out, "bestmove ")+len("bestmove "):]
	}
	return result
}

func outputAvailable(search string, out string) bool {
	if strings.Index(out, search) == -1 {
		return false
	}
	return true
}

func writeString(cmd string, stdin io.WriteCloser) {
	_, err := io.WriteString(stdin, cmd+"\n")
	if err != nil {
		log.Fatal(err)
	}
}

func PlayPlayer(game *Game) {
	out := useStockfish([]string{"position " + fetchStartingPos(game) + " moves " + game.Playermove, "d"}, "Fen: ")
	game.Fen = fetchFen(out)
}

func FetchComputerMove(game *Game) {
	out := useStockfish([]string{"position " + fetchStartingPos(game), "go depth " + game.Depth}, "bestmove ")
	game.Answer = fetchBestMove(out)[0:4]
}

func PlayComputerMove(game *Game) {
	out := useStockfish([]string{"position " + fetchStartingPos(game) + " moves " + game.Answer, "d"}, "Fen: ")
	game.Fen = fetchFen(out)
}

func fetchStartingPos(game *Game) string {
	if len(game.Fen) != 0 {
		return strings.ReplaceAll("fen "+game.Fen, "\r", "")
	}
	return "startpos"
}

func useStockfish(commands []string, search string) string {
	stdin, stdout, err := startEngine("stockfish")
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range commands {
		writeString(cmd, stdin)
	}

	b := make([]byte, 10240)

	for {
		n, err := stdout.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		if outputAvailable(search, string(b)) {
			break
		}
		if n > 1000000 {
			break
		}
	}
	return string(b)
}

func startEngine(enginePath string) (stdin io.WriteCloser, stdout io.ReadCloser, err error) {
	cmd := exec.Command(enginePath)

	stdin, err = cmd.StdinPipe()
	if err != nil {
		return
	}
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start() // start command - but don't wait for it to complete
	return
}
