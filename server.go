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

func saveInput(w http.ResponseWriter, r *http.Request) {
	cc := r.FormValue("CC")
	name, err := r.FormValue("Name")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return 1
	}

	// Not sure what to do here, probably upload it and save it on the db
	f, err := os.Open("somefile.json")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return 1
	}

	f.Write(b)
	f.Close()
}

func init() {
	http.HandleFunc("/save", save)
}
