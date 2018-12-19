package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	//"github.com/julienschmidt/httprouter"
	//"github.com/gorilla/sessions"
	//"github.com/julienschmidt/httprouter"
	//"github.com/gorilla/csrf"
	/*
		"encoding/base64"
		"net/url"
		"sync"
		"crypto/hmac"
		"crypto/sha256"
	*/)

/*
-More go routines and channel interaction with them
-CrossSiteRequestForgery protection with github.com/gorilla/csrf
*/

func BasicAuth(h httprouter.Handle, requires [][1]string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		_, err := r.Cookie("AAUEremotecredencials")
		fmt.Println("cookie check")
		if err == nil {
			id, _ := uuid.NewV4()
			expiration := time.Now().Add(24 * time.Hour) // 24h keys
			cookie := http.Cookie{Name: "AAUEremotecredencials", Value: id.String(), Expires: expiration, Domain: "localhost"}
			http.SetCookie(w, &cookie)
			//cookie := http.Cookie{Name: "AAUEremotecredencials", Domain: "localhost", Value: id.String(), Expires: expiration} // No maxAge, makes cookie ageless
			// Value:   HMAC256("AAUEremotecredencials", "AAUEremotecredencials"),  // safer than uuid
			// Domain: "aauecred.net"  // set on /etc/hosts
			//http.SetCookie(w, &cookie)
			fmt.Println("cookie set")
		}

		fmt.Println("Entering BasicAuth...")
		user, password, hasAuth := r.BasicAuth()
		var conf = false
		if hasAuth {
			for _, combo := range requires {
				if combo[0] == user+" "+password {
					conf = true
				}
			}
		}
		if conf == true {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func aboutPage(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	_, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		t, _ := template.ParseFiles("./templates/about.html")
		t.Execute(w, nil)
	}
}

func cred(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	var acessoA [8]string
	fmt.Println(r.Form)
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./templates/credform.html")
		t.Execute(w, nil)
		fmt.Println("Served credform.html file")
	} else {
		r.ParseForm()
		name := r.Form["nome"][0]
		cc := r.Form["cc"][0]
		//tipo := r.Form["tipo"]
		in, header, err := r.FormFile("photo")
		checkerr(err)
		defer in.Close()
		out, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
		checkerr(err)
		defer out.Close()
		io.Copy(out, in)
		for k, _ := range r.Form {
			if k[0] == 'z' {
				acessoA[k[1]-1] = string(k[1])
			} else {
				acessoA[k[1]-1] = "X"
			}
		}
		oldCred(header.Filename, name, cc, acessoA)
	}
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func oldCred(photo string, name string, cc string, acessoA [8]string) string {
	// Just for now so that something can be presented
	// foto.png must be named with the name or have the name an extra
	// Pcmd := "python tests.py"
	// args := ""
	// cmd := Pcmd + args
	var acessoS string
	db, err := sql.Open("sqlite3", "db.sql")
	checkerr(err)
	for _, v := range acessoA {
		acessoS += v
	}
	cmd := exec.Command("cd old; python credencias.py " + photo + " " + name + " " + cc + " " + acessoS)
	stmt, err := db.Prepare("INSERT INTO createdcreds values (?,?,?)")
	checkerr(err)
	res, err := stmt.Exec("1", time.Now(), "luis") // uuid, date, user
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
	var port, host, target string
	var hnP []string
	target = "http://" + req.URL.Host + req.URL.Path
	fmt.Println("redirecting from %s", target)
	hnP = strings.Split(req.Host, ":")
	host, port = hnP[0], hnP[1]
	if port != "9090" {
		port = "9090"
	}
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	target = "https://" + host + ":" + port + req.URL.Path
	log.Printf("to: %s", target)
	http.Redirect(w, req, target, http.StatusMovedPermanently)
}

func main() {
	db, err := sql.Open("sqlite3", "remotecreds")
	checkerr(err)
	rows, err := db.Query("SELECT * FROM user")
	checkerr(err)
	/*
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			//HostPolicy: autocert.HostWhitelist("www.checknu.de"),
			Cache:      autocert.DirCache("/home/letsencrypt/"),
		}
	*/
	var user, pass string
	var check [][1]string
	var aux [1]string
	for rows.Next() {
		err := rows.Scan(&user, &pass)
		checkerr(err)
		aux[0] = user + " " + pass
		check = append(check, aux)
	}

	router := httprouter.New()
	router.GET("/", aboutPage)
	router.GET("/about/", aboutPage)
	router.GET("/cred/", BasicAuth(cred, check))
	// CrossSiteRequestForgery protection
	//CSRF := csrf.Protect([]byte(hashKey))  // i think it is the hashKey
	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(redirTLS)) // Final version should use port 80
		if err != nil {
			panic(err)
		}
	}()
	http.ListenAndServeTLS(":9090", "cert.pem", "key.pem", router) // Final version should use port 443
}

// New Cookie and session
/*
var store = sessions.NewCookieStore(os.Getenv("SESSION_KEY")) // export SESSION_KEY=$(bash genKey.sh)
// Useless new cookie and session stuff
var hashKey = []byte("very-secret")   // probably get genKey.sh output
var blockKey = []byte("a-lot-secret") // another genKey.sh output

var hashKey = []byte("very-secret")
var blockKey = []byte("a-lot-secret")

var c = securecookie.New(hashKey,blockKey)

func SetCookieHandler(w http.ResponseWriter, r *http.Request, user string) {
	value := map[string]string{
		"user" : user,
	}
	if encoded, err := c.Encode("AAUEremotecredencials", value); err == nil {
		cookie := &http.Cookie{
			Name:  "AAUEremotecredencials",
			Value: encoded,
			Path:  "/",
			Secure: true,
		}
		http.SetCookie(w, cookie)
	}
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request, h httprouter.Handle, ps httprouter.Params) {
	if cookie, err := r.Cookie("AAUEremotecredencials"); err == nil {
		value := make(map[string]string)
		if err = c.Decode("AAUEremotecredencials", cookie.Value, &value); err != nil {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	} else {
		fmt.Println("fuck")
	}
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "AAUEremotecredencials")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		// Maybe setcookie?
	}
	session.Values[hashKey] = blockKey
	session.Save(r, w)
}
*/
/*
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		return
		// The user isn't logged in because cookie isnt set
	}
}
*/

//var s = sessions.NewCookieStore(os.Getenv("SESSION_KEY"))
//var s = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
//var s = sessions.NewCookieStore([]byte("something-very-secret"))

/*
func MyCookie(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get req and check if it has cookie if not serve a new cookie
	_, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		id, _ := uuid.NewV4()
		// id := HMAC256() // Arguments?
		expiration := time.Now().Add(24 * time.Hour) // 24h keys
		cookie := http.Cookie{
			Name:   "AAUEremotecredencials",
			Domain: "localhost", // im guessing something like that
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
	Set(key, value interface{}) error //set custom session value
	Get(key interface{}) interface{}  //get custom session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
}

type CookieForm struct {
	Name string
	//lock     sync.Mutex // Not yet
	lifeTime int64
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
*/

// end new cookie and session
/*
func BasicAuth(h httprouter.Handle, requires [][1]string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()
		s1, err := s.Get(r, "AAUEremotecredencials")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} // if value != cookie[value]; exit
		s1.Options = &sessions.Options{
			Domain: "localhost",
			// Value:  []byte("something-very-secret"),
			// address:	ip
			// Domain:     "AAUEremotecredencials",
			Path:     "/",
			MaxAge:   86400, // a day
			HttpOnly: true,
		}
		if hasAuth {
			for _, combo := range requires {
				if combo[0] == user+" "+password {
					// conf = true
					// Save in cookie value["login"] = true
					SetCookieHandler(w, r, user)
					s1.Values["login"] = true
					s1.Values["user"] = user
					err := s1.Save(r, w)
					checkerr(err)
				}
			}
		}
		ReadCookieHandler(w,r, h, ps)
		fmt.Println(user + " " + password)
		fmt.Println(hasAuth)
	}
}
*/

func BasicAuth(h httprouter.Handle, requires [][1]string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()
		//var conf = false
		if hasAuth {
			for _, combo := range requires {
				if combo[0] == user+" "+password {
					expiration := time.Now().Add(24*time.Hour)
					cookie := http.Cookie{Name: "AAUEremotecredentials", Value: user, Expires: expiration}
					http.SetCookie(w, &cookie)
					//conf = true
				}
			}
		}
		coo, _ :=  r.Cookie("AAUEremotecredentials")
		fmt.Println(coo)
		if coo != nil{
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func aboutPage(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	_, err := r.Cookie("username")
	if err != nil {
		t, _ := template.ParseFiles("./templates/about.html")
		t.Execute(w, nil)
	}
}

func cred(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	var acessoA [8]string
	fmt.Println(r.Form)
	_, err := r.Cookie("AAUEremotecredencials")
	if err != nil {
		fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("./templates/credform.html")
			t.Execute(w, nil)
		} else {
			r.ParseForm()
			//name := r.Form["nome"][0]
			//cc := r.Form["cc"][0]
			//tipo := r.Form["tipo"]

			in, header, err := r.FormFile("photo")
			checkerr(err)
			defer in.Close()
			out, err := os.OpenFile(header.Filename, os.O_WRONLY|os.O_CREATE, 0644)
			checkerr(err)
			defer out.Close()
			io.Copy(out, in)

			for k, _ := range r.Form {
				if k[0] == 'z' {
					acessoA[k[1]-1] = string(k[1])
				} else {
					acessoA[k[1]-1] = "X"
				}
			}
			//oldCred(header.Filename, name[0], cc[0], acessoA)
		}
	}
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func oldCred(photo string, name string, cc string, acessoA [8]string) string {
	// Just for now so that something can be presented
	// foto.png must be named with the name or have the name an extra
	// Pcmd := "python tests.py"
	// args := ""
	// cmd := Pcmd + args
	var acessoS string
	db, err := sql.Open("sqlite3", "db.sql")
	checkerr(err)
	for _, v := range acessoA {
		acessoS += v
	}
	cmd := exec.Command("cd old; python credencias.py " + photo + " " + name + " " + cc + " " + acessoS)
	stmt, err := db.Prepare("INSERT INTO createdcreds values (?,?,?)")
	checkerr(err)
	res, err := stmt.Exec("1", time.Now(), "luis") // uuid, date, user
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
	var port, host, target string
	var hnP []string
	target = "https://" + req.URL.Host + req.URL.Path
	hnP = strings.Split(req.Host, ":")
	host, port = hnP[0], hnP[1]
	if port != "9090" {
		port = "9090"
	}
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	target = "https://" + host + ":" + port + req.URL.Path
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusMovedPermanently)
}

func main() {
	db, err := sql.Open("sqlite3", "remotecreds")
	checkerr(err)
	rows, err := db.Query("SELECT * FROM user")
	checkerr(err)
	/*
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			//HostPolicy: autocert.HostWhitelist("www.checknu.de"),
			Cache:      autocert.DirCache("/home/letsencrypt/"),
		}
	*/
	var user, pass string
	var check [][1]string
	var aux [1]string
	for rows.Next() {
		err := rows.Scan(&user, &pass)
		checkerr(err)
		aux[0] = user + " " + pass
		check = append(check, aux)
	}
	fmt.Println(check)
	router := httprouter.New()
	router.GET("/", aboutPage)
	// router.GET("/login/", login)
	router.GET("/about/", aboutPage)
	router.GET("/cred/", BasicAuth(cred, check))
	// CrossSiteRequestForgery protection
	//CSRF := csrf.Protect([]byte(hashKey))  // i think it is the hashKey
	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(redirTLS)) // Final version should use port 80
		if err != nil {
			panic(err)
		}
	}()
	http.ListenAndServeTLS(":9090", "cert.pem", "key.pem", router) // Final version should use port 443
}

/*

func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func HMAC256(payload string, secret string) string {
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(payload))
	return b64Encode(string(sig.Sum(nil)[:]))

	key := []byte("5ebe2294ecd0e0f08eab7690d2a6ee69")
	message := "AAUEremotecredentials"
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(message))
	fmt.Println(hex.EncodeToString(sig.Sum(nil)))
}

func b64Encode(text string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(text))
}

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
