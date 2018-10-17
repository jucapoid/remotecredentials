package server

import "net/http log fmt"

func server(photoCred) {
    http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(photoCred))))
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

/*
Also serve an html and getforms to get input
Probably would be a good idea to only run it on localhost
*/
