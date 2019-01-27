package bingo

import (
	"log"
	"fmt"
	"regexp"
	"net/http"
	"io/ioutil"
)

type GameSession struct {
	HostName     string
	SecretPhrase string
}

type Page struct {
    Title string
    Body  []byte
}

var sessionsMap map[string]*GameSession

var validPath = regexp.MustCompile("^/(|index|home|admin|player)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return "index", nil
	}
	fmt.Println(m[2])
	return m[2], nil // The title is the second subexpression.
}

func loadPage(title string) (*Page, error) {
	filename := "html/" + title  + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func IndexHandler(w http.ResponseWriter, r* http.Request) {
	log.Println("index URL:", r.URL)
	title, err :=  getTitle(w, r)
	if err != nil {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	if title == "" || title == "index.html" || title == "home" {
		title = "index"
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "%s", p.Body)
}

func AdminHandler(w http.ResponseWriter, r* http.Request) {
	log.Println("admin URL:", r.URL)
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
	gameSession := GameSession{HostName: r.Form.Get("hostname"), SecretPhrase: r.Form.Get("secretphrase"),}
	sessionsMap[gameSession.HostName+"-"+gameSession.SecretPhrase] = &gameSession
	log.Println("HostName:", gameSession.HostName, "SecretPhrase:", gameSession.SecretPhrase)
	fmt.Fprintf(w, "%s", "http://localhost:8080/players/add/"+gameSession.HostName + "-" + gameSession.SecretPhrase)
}

func AddPlayerHandler(w http.ResponseWriter, r* http.Request) {
	log.Println("Add Players URL:", r.URL)
	sessionId := r.URL.Path[13:]  // skip "/players/add/"
	log.Println("SessionId:", sessionId)
	if _, ok := sessionsMap[sessionId]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	p, _ := loadPage("players")
	c := http.Cookie{Raw: sessionId,}
	http.SetCookie(w, &c)
	fmt.Fprintf(w, "%s", p.Body)
}

func ListPlayersHandler(w http.ResponseWriter, r* http.Request) {
	log.Println("List Players URL:", r.URL)
	sessionId := r.URL.Path[14:]  // skip "/players/list/"
	log.Println("SessionId:", sessionId)
	if _, ok := sessionsMap[sessionId]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	c := http.Cookie{Raw: sessionId,}
	http.SetCookie(w, &c)
	fmt.Fprintf(w, "%s", "list-players")
}

func PlayerInHandler(w http.ResponseWriter, r* http.Request) {
	log.Println("PlayerIn URL:", r.URL)
	sessionId := r.URL.Path[10:]  // skip "/playerin/"
	log.Println("SessionId:", sessionId)
	if _, ok := sessionsMap[sessionId]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//p, _ := loadPage("player-sheet")
	fmt.Fprintf(w, "%s", "player-sheet") //p.Body)
}

func init() {
	sessionsMap = make(map[string]*GameSession)
}
