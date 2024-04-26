package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"math/rand"
	"net"
	"net/http"
)

const (
	host = "localhost:8080"
)

var Dct = make(map[string]string)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func MakeShortURL(name string) string {
	id := randSeq(6)
	Dct[id] = name
	return id

}

func GetHandle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//writer.Write([]byte("Hi"))
	id := params.ByName("id")
	if value, exist := Dct[id]; exist {
		http.Redirect(writer, request, value, http.StatusTemporaryRedirect)
		return
	}
	writer.WriteHeader(400)
}

func PostHandle(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var url string

	if err := request.ParseForm(); err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	for k, v := range request.Form {
		url += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	writer.WriteHeader(201)

	_, err := writer.Write([]byte(host + "/" + MakeShortURL(url)))
	if err != nil {
		panic(err)
	}
}

func DefaultHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	writer.WriteHeader(400)
}

func main() {
	router := httprouter.New()
	router.POST("/", PostHandle)
	router.GET("/:id", GetHandle)
	router.GET("/", DefaultHandler)
	listen, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler: router,
	}

	log.Fatalln(server.Serve(listen))
}
