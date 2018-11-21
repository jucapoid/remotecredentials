package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os" // For file storing
	"os/exec"
	"strings"
	"time" // For cookie but could also serve timestamp on pages

	_ "github.com/lib/pq"
)

/*
Change prints to log
Put everything working with mutexes
*/

type CookieForm struct {
	Name string
	//lock     sync.Mutex // Not yet
	lifeTime int64
}

func setMyCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "AAUEremotecredencials", Value: "testValue", Expires: expiration}
	http.SetCookie(w, &cookie)
}

/*
//func NewManager(cookieName string, maxlifetime int64) (*CookieForm, error) {
func NewManager(provideName, cookieName string, maxlifetime int64) (*CookieForm, error){
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &CookieForm{cookieName: cookieName, maxlifetime: maxlifetime}, nil
}
*/

func aboutPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./test/templates/about.html")
	if err != nil {
		fmt.Println(err)                              // Ugly debug output
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	t.Execute(w, nil)
	//fmt.Fprintf(w, "Hi\nThis is a credencials generator for AAUE") // Add links to the rest of the pages
}

func credform(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, err := template.ParseFiles("./test/templates/credform.html")
		if err != nil {
			fmt.Println(err)                              // Ugly debug output
			w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
			return
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()  // acesso should be a map
	}
}

/*
func dbManager(q string) {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
}
*/

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
	// you probably want to make sure header.Filename is unique and
	// use filepath.Join to put it somewhere else.
	out, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		errCount += 1
	}
	defer out.Close()
	io.Copy(out, in)
	// do other stuff
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./test/templates/login.html")
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

func oldCred(photo string, name string, cc string) string {
	// Just for now so that something can be presented
	// foto.png must be named with the name or have the name an extra
	// Pcmd := "python tests.py"
	// args := ""
	// cmd := Pcmd + args
	// or
	cmd := exec.Command("cd ../old; python credencias.py " + name + " cred" + name + ".png")
	fmt.Println("Creating new credencial for " + name + " named cred" + name)
	if errV := cmd.Run(); errV != nil {
		log.Fatalf("Error: ", errV)  // It's better than Start bc it waits to the command to finish
	}
	return (photo + name)
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
	http.HandleFunc("/cred/", credform)
	//http.HandleFunc("/cred", cred)           // Login must always come first
	err := http.ListenAndServe(":9090", nil) // set listen port
	fmt.Println("Server up")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
