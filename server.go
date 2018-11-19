package creds

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var counter int
var mutex = &sync.Mutex{}

func server() {
	photoCred := "old/cred.png" // later put this as an argument
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(photoCred))))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// move this to main

/*
Also serve an html and getforms to get input
Probably would be a good idea to only run it on localhost
*/

func saveInput(w http.ResponseWriter, r *http.Request) (string, string) {
	cc := r.FormValue("CC")
	name := r.FormValue("Name")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return "1", "1"
	}

	// Not sure what to do here, probably upload it and save it on the db
	if err != nil {
		http.Error(w, err.Error(), 500)
		return "1", "1"
	}
	return cc, name
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
}

func HandleUpload(w http.ResponseWriter, req *http.Request) {
	in, header, err := req.FormFile("file")
	errCount := 0
	if err != nil {
		log.FatalF("Error: ", err)
		errCount += 1
	}
	defer in.Close()
	//you probably want to make sure header.Filename is unique and
	// use filepath.Join to put it somewhere else.
	out, err := os.OpenFile(header.Filename, os.O_WRONLY, 0644)
	if err != nil {
		errCount += 1
		out := os.OpenFile(header.Filename+strconv.Itoa(errCount), os.O_WRONLY, 0644)
	}
	defer out.Close()
	io.Copy(out, in)
	//do other stuff
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		incrementCounter()
	})
	http.HandleFunc("/save", HandleUpload, saveInput)

	// http.HandleFunc("/increment", incrementCounter)  // Increment per cookie

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
			http.ServeFile(w, r, r.URL.Path[1:])
			})
		// Serves html files
	*/

	http.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi\nThis is a credentials generator for AAUE")
		incrementCounter()
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
