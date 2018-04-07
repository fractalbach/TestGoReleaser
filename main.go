package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/fractalbach/ninjatools/structInspector"
)

const version = "0.0.2"
const endpoint = "https://api.github.com/graphql"
const help = `
SYNOPSIS:
  Download Latest Github Release using the GRAPHQL API.

USAGE:
  TestGoReleaser -repo <repoName> -owner <ownerName> -token [OPTIONS]...

`

var verbose = false
var Query = `{"query":"{\n  repository(owner: \"{{.Owner}}\", name: \"{{.Repo}}\") {\n    releases(last: 1) {\n      edges {\n        node {\n          publishedAt\n          tag {\n            name\n          }\n          releaseAssets(first: 10) {\n            edges {\n              node {\n                name\n                downloadUrl\n              }\n            }\n          }\n        }\n      }\n    }\n  }\n}\n","variables":"{}","operationName":null}`

type Repo struct {
	Owner, Repo string
}

func ToHelp(s string) {
	if len(s) > 0 {
		fmt.Println("Input Error:", s)
	}
	fmt.Print(help)
	flag.PrintDefaults()
	os.Exit(1)
}

func checkErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func checkErrLogFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func NewClient() *http.Client {
	myTransport := &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Second * 60,
		Transport: myTransport,
	}
}

func printReqRespInfo(req *http.Request, resp *http.Response) {
	fmt.Println("\n ~~~~ Request ~~~~")
	structInspector.PrettyPrint(req, 3)

	fmt.Println("\n ~~~~ Request Headers ~~~~")
	for key, val := range req.Header {
		fmt.Printf("[%s] : %s\n", key, val)
	}
	fmt.Println("\n ~~~~ Response ~~~~")
	structInspector.PrettyPrint(resp, 3)

	fmt.Println("\n ~~~~ Response Headers ~~~~")
	for key, val := range resp.Header {
		fmt.Printf("[%s] : %s\n", key, val)
	}
}

func doFetch(reader io.Reader, url, token string) *http.Response {
	client := NewClient()
	req, err := http.NewRequest("POST", url, reader)
	checkErrLogFatal(err)
	req.Header.Add("User-Agent", "Achenbot/"+version)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "bearer "+token)

	resp, err := client.Do(req)
	checkErrLogFatal(err)

	if verbose {
		printReqRespInfo(req, resp)
	}

	return resp
}

func main() {

	repoPtr := flag.String("repo", "", "Name of the github repository. (Required)")
	ownerPtr := flag.String("owner", "", "Repository owner's name. (Required)")
	tokenPtr := flag.String("token", "", "Github OAUTH2 Token Value. (Recommended)")
	// TODO: add token path flag.
	versionBool := flag.Bool("version", false, "Prints the version of the program.")
	verboseBool := flag.Bool("verbose", false, "Prints headers and data structures of request and response.")
	flag.Parse()

	repo := *repoPtr
	owner := *ownerPtr
	token := *tokenPtr
	verbose = *verboseBool

	// Handle the flags.  Print help menus if needed.
	switch {
	case *versionBool:
		fmt.Println(version)
		return
	case repo == "":
		ToHelp("Needs a repository name.")
	case owner == "":
		ToHelp("Need to have owner name.")
	}

	// Create message body template using the data structure fields.
	myRepo := Repo{Owner: owner, Repo: repo}
	tmpl, err := template.New("test").Parse(Query)
	checkErrPanic(err)

	// Create a pipe reader/writer for the message body.
	r, w := io.Pipe()

	// Write the template using the user's input.
	go func() {
		err = tmpl.Execute(w, &myRepo)
		checkErrPanic(err)

		err = w.Close()
		checkErrPanic(err)
	}()

	// The Fetcher!
	resp := doFetch(r, endpoint, token)

	// copy response body to standard output
	if verbose {
		fmt.Println("\n ~~~~ Response Body ~~~~")
	}
	_, err = io.Copy(os.Stdout, resp.Body)
	checkErrLogFatal(err)
}

/* Potential GRAPHQL stuff.

// Means argument,
type Arg struct {
    name string
    val interface{}
}

func NewArg(name string, val interface{}) Arg {
    return Arg
}

// Converts the argument to a string in the format
// (name: value)
func (a *Arg) String() string {
    switch val.(type) {
    case string:
        return fmt.Sprintf(`%v: "%v"`, val)
    }
    return fmt.Sprintf("%v: %v", name, val)
}
*/
