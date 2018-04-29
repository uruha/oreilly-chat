package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHander{filename: "chat.html"}))
	http.Handle("/room", r)

	go r.run()

	log.Println("Start web server. Port", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
