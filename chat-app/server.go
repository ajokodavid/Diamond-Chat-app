package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var tmpl *template.Template

func chat(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("template/chat.html"))

	tmpl.Execute(w,nil)
}

func friends(w http.ResponseWriter, r *http.Request) {
	tmpl = template.Must(template.ParseFiles("template/friends.html"))

	tmpl.Execute(w,nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var savedsocketreader []*socketReader

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Socket Request")

	if savedsocketreader == nil {
		savedsocketreader = make([]*socketReader,0)
	}

	defer func() {
		err := recover()

		if err != nil {
			fmt.Println(err)
		}

		r.Body.Close()
	}()

	con, _ := upgrader.Upgrade(w,r,nil)

	ptrSocketReader := &socketReader{
		con: con,
	}

	savedsocketreader = append(savedsocketreader, ptrSocketReader)

	ptrSocketReader.startThread()
}

type socketReader struct{
	con *websocket.Conn
	mode int
	name string
}

func(i *socketReader) broadcast(str string) {
	for _, g := range savedsocketreader {
		if g == i {
			continue
		}

		if g.mode == 1 {
			continue
		}

		g.writeMsg(i.name, str)
	}
}

func(i *socketReader) read() {
	_, msg, err := i.con.ReadMessage()

	if err != nil{
		panic(err)
	}

	fmt.Println(i.name + " " + string(msg))
	fmt.Println(i.mode)

	if i.mode == 1 {
		i.name = string(msg)
		i.writeMsg("Sys", " Welcome " + i.name + " Feel free to chat with other users")
		i.mode = 2

		return
	}

	i.broadcast(string(msg))

	fmt.Println(i.name + " " + string(msg))
}

func(i *socketReader) writeMsg(name string, str string) {
	i.con.WriteMessage(websocket.TextMessage, []byte("<b>" + name + ": </b>" + str))
}

func(i *socketReader) startThread() {
	i.writeMsg("Sys", "Please write your name")
	i.mode = 1

	go func() {
		defer func() {
			err := recover()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("Thread socketreader finish")
		}()

		for {
			i.read()
		}
	}()
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/echo", echo)

	mux.HandleFunc("/", chat)

	mux.HandleFunc("/friends", friends)

	fs := http.FileServer(http.Dir("./static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	fmt.Println("Server Started...")
	
	port := os.Getenv("PORT")

	log.Print("Listening on :" + port)

	log.Fatal(http.ListenAndServe(":"+port, mux))

}