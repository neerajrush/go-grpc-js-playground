package main

import (
	"log"
	"net/http"
	bingo "github.com/neerajrush/go-grpc-js-playground/bingo"
)

func main() {
	http.HandleFunc("/admin/", bingo.AdminHandler)
	http.HandleFunc("/players/add/", bingo.AddPlayerHandler)
	http.HandleFunc("/players/list/", bingo.ListPlayersHandler)
	http.HandleFunc("/playerin/", bingo.PlayerInHandler)
	http.HandleFunc("/", bingo.IndexHandler)
	http.HandleFunc("/index", bingo.IndexHandler)
	http.HandleFunc("/home", bingo.IndexHandler)
	//log.Fatalf("failed to listen http:", http.ListenAndServe(":80", nil))
	log.Fatalf("failed to listen http:", http.ListenAndServeTLS(":443", "/usr/local/certs/cert.pem", "/usr/local/certs/key.pem", nil))
}
