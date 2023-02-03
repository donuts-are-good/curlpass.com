package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"strings"
)

const helpText = `
curlpass.com - a password generator

Usage:
curl curlpass.com/[passwordType]

passwordType:
0 - just hex
1 - just letters
2 - letters + numbers
3 - letters + numbers + symbols

Contributing:
github.com/curlpass

Please keep it standard library and free from dependencies and bullshit :)

Donations: 
If you would like to support the development of this project, you can donate to the following addresses:

	Bitcoin: bc1qg72tguntckez8qy2xy4rqvksfn3qwt2an8df2n
	Monero: 42eCCGcwz5veoys3Hx4kEDQB2BXBWimo9fk3djZWnQHSSfnyY2uSf5iL9BBJR5EnM7PeHRMFJD5BD6TRYqaTpGp2QnsQNgC

Thank you for your support!
`

func generatePassword(passwordType string) string {

	// make a string variable for the character set
	var characters string

	// switch because if/else isn't cool anymore
	switch passwordType {

	// if "0", make a hex password
	case "0":
		characters = "1234567890ABCDEF"

	// if "1", upper and lowercase letters
	case "1":
		characters = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

	// if "2", use letters and numbers
	case "2":
		characters = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

	// if "3", user letters and numbers and symbols
	case "3":
		characters = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM!@#$%^&*(){}|:<>?-=[]\\;,./"

	// if the passwordType is "help", return an empty string as a trigger.
	// this seems out of place to be in this switch/case, however, you see,
	// the way we're parsing the urls for password type conflicts with anything
	// else under the root path, meaning we'd have been better off doing like
	// curlpass.com/password/1 or curlpass.com/p/1 but that'd been a stumbling
	// point when typing all those extra characters, so we sacrifice a little
	// bit of developer comfort for the sake of user comfort.
	// it's fine.
	case "help":
		return ""

	// if the passwordType is anything else, use the default character set of numbers and letters
	default:
		characters = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	}

	// slice for our password
	password := make([]byte, 12)

	// define the max
	max := big.NewInt(int64(len(characters)))

	// loop through the password
	for i := range password {

		// store our random int, return errors
		ourInt, err := rand.Int(rand.Reader, max)
		if err != nil {
			return err.Error()
		}

		// convert the big int to int64
		ourInt64 := ourInt.Int64()

		// make the chosen character this character
		password[i] = characters[ourInt64]
	}

	// return the password
	return string(password)

}

func passwordHandler(w http.ResponseWriter, r *http.Request) {

	// get the type from the number string (i know) after the /
	passwordType := strings.TrimPrefix(r.URL.Path, "/")

	// make the correct password and .... call it 'password'
	password := generatePassword(passwordType)

	// if it's help, send them the help page
	if passwordType == "help" {
		http.Redirect(w, r, "/help", http.StatusSeeOther)
		return
	}

	// check the user agent for curl
	// in retrospect, it's a good thing we're not evil.
	if strings.Contains(r.Header.Get("User-Agent"), "curl") {

		// if it's curl, just give them the password in the response with no html
		fmt.Fprintln(w, password)
	} else {

		// if it's not curl, it's probably a browser or something that appreciates
		// the madness that is html parsing
		tmpl, err := template.ParseFiles("./template/index.html")
		if err != nil {

			// if there's an error, report the appropriate error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// render the template with the password passed in
		tmpl.Execute(w, password)
	}
}

func helpHandler(w http.ResponseWriter, r *http.Request) {

	// et tu curl-ay?
	if strings.Contains(r.Header.Get("User-Agent"), "curl") {

		// if it's curl, just return the help text
		// ignore the linter, it's just being snooty.
		// if the linter bothers you though, do this
		// w.Write([]byte(helpText))
		fmt.Fprintln(w, helpText)
	} else {

		// if it's not curl, parse the help.html template and insert the help text
		tmpl, err := template.ParseFiles("./template/help.html")
		if err != nil {

			// if there's an error, report the appropriate error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// render the template with the helpText string passed in
		tmpl.Execute(w, helpText)
	}
}
func main() {

	// define the password handler
	http.HandleFunc("/", passwordHandler)

	// send help requests to the help handler
	http.HandleFunc("/help", helpHandler)

	// port 9945
	err := http.ListenAndServe(":9945", nil)

	// shit the bed if we're not happy
	if err != nil {
		panic(err)
	}
}
