package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
)

type templateHander struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "address of application")
	flag.Parse()
	gomniauth.SetSecurityKey("gochat")
	gomniauth.WithProviders(
		github.New("089ffd6d320420fb6774", "a424d77abc3cae5eb0959df50132c23d4aa2c6ac", "http://localhost:3308/auth/callback/github"),
	)
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHander{filename: "chat.html"}))
	http.Handle("/login", &templateHander{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	log.Println("Start web server. Port", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
