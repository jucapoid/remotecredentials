package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"html/template"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  // parse arguments, you have to call this by yourself
	fmt.Println(r.Form)  // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	var name string
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		name += strings.Join(v," ") + " "

	}
	fmt.Fprintf(w, "hello %s\n" , name) // send data to client side
}

func login(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func main() {
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login", login) // set router
	err := http.ListenAndServe(":8081", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}