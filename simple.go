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
	"crypto/hmac"
	// _ "github.com/lib/pq"
	// "github.com/satori/go.uuid"
	// "github.com/nu7hatch/gouuid"
	"github.com/julienschmidt/httprouter"
)
// get libs
// go get github.com/lib/pq
// go get github.com/julienschmidt/httprouter
// go get github.com/satori/go.uuid
// or
// go get github.com/nu7hatch/gouuid

/*
Join BasicAuth and login
BasicAuth and Manager?
Change prints to log
Add protect and unprotected to protected and unprotected pages
Add go subroutines
Get a ssl crt
Cookie and not session bc cookies are pressistent
Https / http1.1
XSS protection?
*/

/*
type CookieForm struct {
	Name string
	//lock     sync.Mutex // Not yet
	lifeTime int64
}
*/

func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()
		if hasAuth && user == requiredUser && password == requiredPassword {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func MyCookie(w http.ResponseWriter, r *http.Request) {
	// get req and check if it has cookie if not serve a new cookie
	req, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		id, _ := uuid.NewV4()
		// use crypto/hmac to make sure that no cookies can be messed arround with
		expiration := time.Now().Add(24 * time.Hour) // 24h keys
		cookie := http.Cookie{
								Name: "AAUEremotecredencials",
								Domain: "aaue.",
								Value: id.String(),
								Expires: expiration
		}  // No maxAge, makes cookie ageless
			// Domain: "aauecred.net"  // set on /etc/hosts
	}
	http.SetCookie(w, &cookie)
	// force to load login
}



func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		cookie, err := r.Cookie(manager.cookieName)
		if (err != nil || cookie.Value == "") {
			sid := manager.sessionId()
			session, _ = manager.provider.SessionInit(sid)
			cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}  // HttpOnly: false
			http.SetCookie(w, &cookie)  // So i guess this replaces the MyCookie func?
		} else {
			sid, _ := url.QueryUnescape(cookie.Value)
			session, _ = manager.provider.SessionRead(sid)
		}
		return
}

//func NewManager(cookieName string, maxlifetime int64) (*CookieForm, error) {
func NewManager(provideName, cookieName string, maxlifetime int64) (*CookieForm, error){
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &CookieForm{cookieName: cookieName, maxlifetime: maxlifetime}, nil
}  // I think we need to remove this

func aboutPage(w http.ResponseWriter, r *http.Request) {
	req, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		t, _ := template.ParseFiles("./templates/about.html")
		t.Execute(w, nil)
	} else {
		login(w, r)
	}
}

func cred(w http.ResponseWriter, r *http.Request) {
	req, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("./templates/credform.html")
			t.Execute(w, nil)
		} else {
			r.ParseForm()
			name := r.Form["nome"]
			cc := r.Form["cc"]
			tipo := r.Form["tipo"]
			acessoA = [8]string
			for k, v := range r.Form {
				if (k[0] == 'z' ) {
					acessoA = append(a, r.Form(k))  // Probably this will not work
				}
			}
		}
	} else {
		login(w, r)
		http.Redirect()
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

func sayhelloName(w http.ResponseWriter, r *http.Request) {  // 
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])  // why?
	var name string
	for k, v := range r.Form {
		fmt.Println(r.Form)
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		name += strings.Join(v, " ") + " "
		// Very dangerous!! Input Validation of expected only
	}
	fmt.Fprintf(w, "Hello %s\n your input has been received", name) // send data to client side
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	in, header, err := r.FormFile("photo")
	errCount := 0
	if err != nil {
		log.Fatalf("Error: ", err)
		errCount += 1
	}
	defer in.Close()
	// you probably want to make sure header.Filename is unique and
	// use filepath.Join to put it somewhere else.
	// Should first check if file already exists
	out, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		errCount += 1
	}
	defer out.Close()
	io.Copy(out, in)
	fmt.Println()
}

func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./templates/login.html")
		t.Execute(w, sess.Get("username"))
		//w.Header().Set("Content-Type", "text/html")
		//t.Execute(w, sess.Get("username"))
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
	// Serve cookie for auth
}

func oldCred(photo string, name string, cc string) string {
	// Just for now so that something can be presented
	// foto.png must be named with the name or have the name an extra
	// Pcmd := "python tests.py"
	// args := ""
	// cmd := Pcmd + args
	// or
	cmd := exec.Command("cd old; python credencias.py " + name + " cred" + name + ".png")
	fmt.Println("Creating new credencial for " + name + " named cred" + name)
	if errV := cmd.Run(); errV != nil {
		log.Fatalf("Error: ", errV) // It's better than Start bc it waits to the command to finish
	}
	return photo + name + ".png"
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
	/*
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login/", login) // Lookup some more advanced routing
	http.HandleFunc("/about/", aboutPage)
	http.HandleFunc("/cred/", cred)
	//http.HandleFunc("/cred", cred)           // Login must always come first
	err := http.ListenAndServe(":9090", nil) // set listen port
	*/
	m := autocert.Manager{
    	Prompt:     autocert.AcceptTOS,
    	//HostPolicy: autocert.HostWhitelist("www.checknu.de"),
    	Cache:      autocert.DirCache("/home/letsencrypt/"),
	}

	user := 'root'
	password := 'toor'  // for now later be kept in a secured format

	router  := httprouter.New()
	router.GET("/", sayhelloName, BasicAuth(Protected, user, pass))
	router.GET("/login/", login, BasicAuth(Protected, user, pass))
	router.GET("/about/", aboutPage, BasicAuth(Protected, user, pass))
	router.GET("/cred/", cred, BasicAuth(Protected, user, pass))
	// Cookies must be checked

	log.Fatal(http.ListenAndServeTLS(":9090", "cert.pem", 'key.pem', router))
	fmt.Println("Server up")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
