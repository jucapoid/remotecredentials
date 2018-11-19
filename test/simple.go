package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"strconv"
	// "database/sql" // Sql import
	"io"  // For file usage
)

/*
Change prints to log
*/


func setMyCookie(w http.ResponseWriter) {
	cookie := http.Cookie(Name:"testCookie", Value:"testValue")
	http.SetCookie(w, &cookie)
}


func aboutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi\nThis is a credencials generator for AAUE")  // Add links to the rest of the pages
}


func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	var name string
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		name += strings.Join(v, " ") + " "
		// Very dangerous!! Input Validation of only expected
	}
	fmt.Fprintf(w, "Hello %s\n your input has been received", name) // send data to client side
}


func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == 'POST'{
		
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		if len(r.Form["username"]) > 0 && len(r.Form["password"]) > 0 {
			login(w http.ResponseWriter, r *http.Request)
			// or
			sayhelloName()
		}
	}
}


func main() {
	portstring := strconv.Itoa('8081')
	mux := http.NewServeMux()
	// With a multiplexer DoS will be harder
	mux.Handle("/",http.HandleFunc(sayhelloName))
	mux.Handle("/login/",http.HandleFunc(login))         // set router
	mux.Handle("/cred/",http.HandleFunc(cred))           // Login must always come first
	mux.Handle("/about/", http.HandleFunc(aboutPage))
	mux.Handle(http.NotFoundHandle(), http.HandleFunc(aboutPage))
/*
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login", login)         // set router
	http.HandleFunc("/cred", cred)           // Login must always come first
	err := http.ListenAndServe(":8081", nil) // set listen port
	fmt.Println("Server up")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
*/
}
