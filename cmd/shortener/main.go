package main

import (
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
	name := params.ByName("name")
	writer.WriteHeader(http.StatusCreated)
	_, err := writer.Write([]byte(host + "/" + MakeShortURL(name)))
	if err != nil {
		panic(err)
	}
}

func MakeShortURL(name string) string {
	id := randSeq(6)
	Dct[id] = name
	return id

}

func main() {
	Dct["ya.ru"] = "yandex.ru"
	router := httprouter.New()
	router.POST("/:name", PostHandle)
	router.GET("/:id", GetHandle)
	listen, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler: router,
	}

	log.Fatalln(server.Serve(listen))
}
