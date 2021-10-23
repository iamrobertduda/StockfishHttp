package internal

import (
	internal "StockfishHttp/internal/stockfish"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func unmarshalGame(r *http.Request) (game internal.Game, valid bool) {
	body, err := ioutil.ReadAll(r.Body)
	if len(body) < 1 {
		return internal.Game{}, false
	}
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &game)
	if err != nil {
		log.Fatal(err)
	}
	return game, true
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	game, valid := unmarshalGame(r)
	if !valid {
		log.Println("No game found")
		return
	}
	i, err := strconv.ParseInt(game.Depth, 10, 8)
	if err != nil {
		log.Println("Invalid depth provided")
		return
	}
	if i > 20 {
		log.Println("A depth above 20 is invalid")
		return
	}

	internal.PlayPlayer(&game)
	internal.FetchComputerMove(&game)
	internal.PlayComputerMove(&game)

	js, err := json.Marshal(game)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		log.Fatal(err)
	}
}

func InitHandler() {
	fmt.Println("Now listening on :8080")
	http.HandleFunc("/move", handleMove)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
