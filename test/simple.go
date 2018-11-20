package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os" // For file storing
	"strings"
	"time" // For cookie but could also serve timestamp on pages
	// "database/sql"
)

/*
Change prints to log
*/

func setMyCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "testCookie", Value: "testValue", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi\nThis is a credencials generator for AAUE") // Add links to the rest of the pages
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	var name string
	for k, v := range r.Form {
		fmt.Println(r.Form)
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		name += strings.Join(v, " ") + " "
		// Very dangerous!! Input Validation of only expected
	}
	fmt.Fprintf(w, "Hello %s\n your input has been received", name) // send data to client side
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	in, header, err := r.FormFile("file")
	errCount := 0
	if err != nil {
		log.Fatalf("Error: ", err)
		errCount += 1
	}
	defer in.Close()
	//you probably want to make sure header.Filename is unique and
	// use filepath.Join to put it somewhere else.
	out, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		errCount += 1
	}
	defer out.Close()
	io.Copy(out, in)
	//do other stuff
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
			sayhelloName(w, r)
		}
	}
}

func main() {
	/*
		portstring := strconv.Itoa('8081')
		mux := http.NewServeMux()
		mux.Handle("/",http.HandleFunc(sayhelloName))
		mux.Handle("/login/",http.HandleFunc(login))         // set router
		mux.Handle("/cred/",http.HandleFunc(cred))           // Login must always come first
		mux.Handle("/about/", http.HandleFunc(aboutPage))
		mux.Handle(http.NotFoundHandle(), http.HandleFunc(aboutPage))
	*/
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login/", login) // set router
	http.HandleFunc("/about/", aboutPage)
	//http.HandleFunc("/cred", cred)           // Login must always come first
	err := http.ListenAndServe(":8081", nil) // set listen port
	fmt.Println("Server up")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
