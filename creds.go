package creds

import ("fmt" "image" "os/exec" "bufio" "os" "database/sql")
import "github.com/lib/pq"
// for using postgresql

/*
-Create DB
-Finish unwritten func's
*/


func main() {
	// var name, cc, photo string  // For normal usage with fmt.Scanf
	var name, cc, photo string = 'Ze Teste', '11227788', 'photo.png'
	fmt.Println("Hello this does basicly nothing yet...")
	//reader := bufio.NewReader(os.Stdin)  // Use this if fmt.Scanf still dosen't wait for input, this or the workaround for the fmt.Scanf
	//name := reader.ReadString('\n')
/*
	fmt.Println("Insert your photo path: ")
	fmt.Scanf('%s', &photo)
	fmt.Println("Insert your name: ")
	fmt.Scanf('%s', &name)
	fmt.Println("Insert your CC number: ")
	fmt.Scanf('%s', &cc)
*/
	oldCred(photo, name, cc)
}


func oldCred(photo string, name string, cc string) {
	// Just for now so that something can be presented
	// foto.png must be named with the name or have the name an extra
	// Pcmd := "python tests.py"
	// args := ""
	// cmd := Pcmd + args
	// ou

	cmd := exec.Command("cd old; python test.py " + photo + name + cc)
	fmt.Println("Creating new credencial: " + name + ".png")
	if errV := cmd.Run(); errV != nil {
		// It's better than Start bc it waits to the command to finish
		log.Fatalf("Error: ", errV)
	}
	else {
		fmt.Println("Done.")
	}
}


func pdfCreateCred() {
	fmt.Fprintf("Create pdf with cred from html")
	/*
	"gofpdf"  - go get github.com/jung-kurt/gofpdf
	"wktml" gihub.com/SebastiaanKlippert/go-wkhtmltopdf  // Probably this one
	*/
}


func dbManager() {
	fmt.Println("postgresql interface")
	/*
	"pq" - go get github.com/lib/pq
	*/
}
