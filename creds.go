package credsmain

import ("fmt" "image" "os/exec")

func main() {
	fmt.Println("Hello this does basicly nothing yet...")	
}


func oldCred() {
	// Just for now so that something can be presented
	// Add variables so that it works with multiple stuff
	// Vars like foto.png cred.png
	// foto.png must be named with the name or have the name an extra
	// Pcmd := "python tests.py"
	// args := ""
	// cmd := Pcmd + args
	// ou
	cmd := exec.Command("cd old; python test.py " + name + " cred.png")
	fmt.Println("Creating new credencial for " + name + " named cred" + name)
	if errV := cmd.Run(); errV != nil {
		log.Fatalf("Error: ", errV)  // It's better than Start bc it waits to the command to finish
	}
	else {
		fmt.Println("Done.")
	}
}


func pdfCreateCred() {
	fmt.Fprintf("Create pdf with cred")
	/*
	"gofpdf"  - go get github.com/jung-kurt/gofpdf
	"wktml" gihub.com/SebastiaanKlippert/go-wkhtmltopdf
	*/
}
