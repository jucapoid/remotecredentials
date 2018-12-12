package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os" // For file storing
	"os/exec"
	"strings"
	"sync"
	"time" // For cookie but could also serve timestamp on pages

	//"github.com/satori/go.uuid"
	//"crypto/hmac"
	// _ "github.com/lib/pq"
	//"github.com/nu7hatch/gouuid"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

// get libs
// go get github.com/mattn/go-sqlite3
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

/*
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
*/

func HMAC256(payload string, secret string) string {
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(payload))
	return b64Encode(string(sig.Sum(nil)[:]))
	/*
		key := []byte("5ebe2294ecd0e0f08eab7690d2a6ee69")
		message := "AAUEremotecredentials"
		sig := hmac.New(sha256.New, key)
		sig.Write([]byte(message))
		fmt.Println(hex.EncodeToString(sig.Sum(nil)))
	*/
}

func b64Encode(text string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(text))
}

func MyCookie(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get req and check if it has cookie if not serve a new cookie
	_, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		id, _ := uuid.NewV4()
		// id := HMAC256() // Arguments?
		// use crypto/hmac to make sure that no cookies can be messed arround with
		expiration := time.Now().Add(24 * time.Hour) // 24h keys
		cookie := http.Cookie{
			Name:   "AAUEremotecredencials",
			Domain: "aaue.",
			Value:  id.String(),
			// Value:   HMAC256("AAUEremotecredencials", "AAUEremotecredencials"),  // safer than uuid
			Expires: expiration} // No maxAge, makes cookie ageless
		// Domain: "aauecred.net"  // set on /etc/hosts
		http.SetCookie(w, &cookie)
	}
	// force to load login
}

type Manager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxlifetime int64
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := session.SessionID()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)} // HttpOnly: false
		http.SetCookie(w, &cookie)                                                                                                                // So i guess this replaces the MyCookie func?
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func aboutPage(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	_, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		t, _ := template.ParseFiles("./templates/about.html")
		t.Execute(w, nil)
	} else {
		login(w, r, h)
	}
}

func cred(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	_, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("./templates/credform.html")
			t.Execute(w, nil)
		} else {
			r.ParseForm()
			/*name := r.Form["nome"]  // Uncomment this
			cc := r.Form["cc"]
			tipo := r.Form["tipo"]*/
			var acessoA [8]string
			for k, _ := range r.Form {
				if k[0] == 'z' {
					acessoA[k[1]-1] = string(k[1]) // Probably this will not work
				} else {
					acessoA[k[1]-1] = "X"
				}
			}
			//oldCred()
		}
	} else {
		login(w, r, h)
		//http.Redirect()
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

func sayhelloName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { //
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"]) // why?
	var name string
	for k, v := range r.Form { // why would we parse the form if there's no POST here?
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
	if err != nil {
		log.Fatalf("Error: ", err)
	}
	defer in.Close()
	// you probably want to make sure header.Filename is unique and
	// use filepath.Join to put it somewhere else.
	// Should first check if file already exists
	out, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	checkerr(err)
	defer out.Close()
	io.Copy(out, in)
	fmt.Println()
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//sess := globalSessions.SessionStart(w, r)
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./templates/login.html")
		t.Execute(w, nil)
		//t.Execute(w, sess.Get("username"))
		//w.Header().Set("Content-Type", "text/html")
		//t.Execute(w, sess.Get("username"))
	} else {
		r.ParseForm()
		db, err := sql.Open("sqlite3", "remotecreds")
		checkerr(err)
		stmt, err := db.Prepare("SELECT * FROM user WHERE username =?")
		checkerr(err)
		rows, err := stmt.Query(r.Form["username"])
		fmt.Println(rows)
		checkerr(err)
		var user string
		var password string
		if rows.Next() {
			err := rows.Scan(&user, &password)
			checkerr(err)
			if user == r.Form["username"][0] && password == r.Form["password"][0] {
				expiration := time.Now().Add(24 * time.Hour)
				cookie := http.Cookie{Name: "username", Value: user, Expires: expiration}
				http.SetCookie(w, &cookie)
				//sess.Set("username", user)
			} else {
				http.Redirect(w, r, "/", 302)
			}
		}
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		//sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
	// Serve cookie for auth
}

func oldCred(photo string, name string, cc string) string {
	db, err := sql.Open("sqlite3", "db.sql")
	checkerr(err)
	// Just for now so that something can be presented
	// foto.png must be named with the name or have the name an extra
	// Pcmd := "python tests.py"
	// args := ""
	// cmd := Pcmd + args
	// or
	cmd := exec.Command("cd old; python credencias.py " + name + " cred" + name + ".png")
	stmt, err := db.Prepare("INSERT INTO createdcreds values (?,?,?)")
	checkerr(err)
	res, err := stmt.Exec("1", time.Now(), "luis")
	checkerr(err)
	affect, err := res.RowsAffected()
	checkerr(err)
	fmt.Println(affect)
	fmt.Println("Creating new credencial for " + name + " named cred" + name)
	if errV := cmd.Run(); errV != nil {
		log.Fatalf("Error: ", errV) // It's better than Start bc it waits to the command to finish
	}
	return photo + name + ".png"
}

func redirTLS(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		http.StatusTemporaryRedirect)
}

func main() {
	/*
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			//HostPolicy: autocert.HostWhitelist("www.checknu.de"),
			Cache:      autocert.DirCache("/home/letsencrypt/"),
		}
	*/

	user := "root"
	password := "toor" // this can now be removed

	router := httprouter.New()
	router.GET("/", sayhelloName)
	router.GET("/login/", login)
	router.GET("/about/", aboutPage)
	router.GET("/cred/", BasicAuth(cred, user, password))
	// Cookies must be checked
	go func() { // a go routine so that can start multiple threads
		err := http.ListenAndServe(":9090", http.HandlerFunc(redirTLS)) // This may fail if so try using router
		if err != nil {
			panic(err)
		}
	}() // idk if this is required
	// log.Fatal(http.ListenAndServeTLS(":9090", "cert.pem", "key.pem", router))
	http.ListenAndServeTLS(":"+Config.String("9090"), Config.Key("https").String("cert.pem"), Config.Key("https").String("key.pem"), router)
}
