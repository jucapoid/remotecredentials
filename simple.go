package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"os"

	//"io"
	"log"
	"net/http"

	//"os"
	"os/exec"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	_ "github.com/mattn/go-sqlite3"
	//"github.com/gorilla/sessions"
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
		user, password, hasAuth := r.BasicAuth()
		//var conf = false
		if hasAuth {
			for _, combo := range requires {
				if combo[0] == user+" "+password {
					expiration := time.Now().Add(24 * time.Hour)
					cookie := http.Cookie{Name: "AAUEremotecredentials", Value: user, Expires: expiration}
					http.SetCookie(w, &cookie)
					//conf = true
				}
			}
		}
		coo, _ := r.Cookie("AAUEremotecredentials")
		fmt.Println(coo)
		if coo != nil {
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
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./templates/credform.html")
		t.Execute(w, nil)
		fmt.Println("Served credform.html file")
	} else {
		r.ParseForm()
		fmt.Println("Post")
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
		acessoS += " " + v + " "
	}
	cmd := exec.Command("cd old; python credencias.py " + photo + " " + name + " " + cc + " " + acessoS)
	stmt, err := db.Prepare("INSERT INTO createdcreds values (?,?,?)")
	checkerr(err)
	res, err := stmt.Exec("1", time.Now(), name) // uuid, date, user
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
	router.POST("/cred/", cred) //BasicAuth cred check
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
