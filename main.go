package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"text/template"
)

const verbose = true
const endpoint = "https://api.github.com/graphql"

const help = `
SYNOPSIS:
  Download Latest Github Release using the GRAPHQL API.

USAGE:
  TestGoReleaser -repo <repoName> -owner <ownerName> [OPTIONS]...

`

var Query = `
query {
  repository(owner: {{.Owner}}, name: {{.Repo}}) {
    releases(last: 1) {
      edges {
        node {
          publishedAt
          tag {
            name
          }
          releaseAssets(first: 10) {
            edges {
              node {
                name
                downloadUrl
              }
            }
          }
        }
      }
    }
  }
}
`

// var repo_owner = "fractalbach"
// var repo_name  = "TestGoReleaser"

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

func main() {

	// Flags
	repoPtr := flag.String("repo", "", "Name of the github repository. (Required)")
	ownerPtr := flag.String("owner", "", "Repository owner's name. (Required)")
	flag.Parse()

	var repo string = *repoPtr
	var owner string = *ownerPtr

	// Display the help message if the required flags aren't filled.
	switch {
	case repo == "":
		ToHelp("Needs a repository name.")
	case owner == "":
		ToHelp("Need to have owner name.")
	}

	// Store the user input.
	myRepo := Repo{Owner: owner, Repo: repo}

	// Create message body template using the data structure fields.
	tmpl, err := template.New("test").Parse(Query)
	if err != nil {
		panic(err)
	}

	// Create a reader to store the message body.
	var b bytes.Buffer
	var body = bufio.NewReadWriter(bufio.NewReader(&b), bufio.NewWriter(&b))

	// Write the template using the user's input.
	err = tmpl.Execute(body, &myRepo)
	checkErrPanic(err)

	err = body.Flush()
	checkErrPanic(err)

	// The Fetcher!
	resp, err := http.Post(endpoint, "application/json", body)
	checkErrLogFatal(err)

	// This is for looking at the headers and info about the request.
	if verbose == true {
		s := reflect.ValueOf(resp).Elem()
		t := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fmt.Fprintf(os.Stdout, "%d: %s %s = %v\n", i, t.Field(i).Name, f.Type(), f.Interface())
		}
	}

	n, err := io.Copy(os.Stdout, resp.Body)
	checkErrLogFatal(err)

	fmt.Println(n)

}

// func Test() {
// 	ownerArg := Arg("owner", "fractalbach")
// 	nameArg := Arg("name", "TestGoReleaser")
// 	last1 := Arg("last", 1)
// }

/*
   Example Commands:
   myRepo := repository{owner: "fractalbach", name: "TestGoReleaser"}
*/

/* Example Query


{
  repository(owner: "fractalbach", name: "TestGoReleaser") {
    releases(last: 1) {
      edges {
        node {
          publishedAt
          tag {
            name
          }
          releaseAssets(first: 10) {
            edges {
              node {
                name
                downloadUrl
              }
            }
          }
        }
      }
    }
  }
}

*/

/* Example Response


{
  "data": {
    "repository": {
      "releases": {
        "edges": [
          {
            "node": {
              "publishedAt": "2018-04-04T10:50:33Z",
              "tag": {
                "name": "HelloWorld_Build_39"
              },
              "releaseAssets": {
                "edges": [
                  {
                    "node": {
                      "name": "HelloWorld_windows_amd64.tar.gz",
                      "downloadUrl": "https://github.com/fractalbach/TestGoReleaser/releases/download/HelloWorld_Build_39/HelloWorld_windows_amd64.tar.gz"
                    }
                  },
                  {
                    "node": {
                      "name": "HelloWorld_linux_amd64.tar.gz",
                      "downloadUrl": "https://github.com/fractalbach/TestGoReleaser/releases/download/HelloWorld_Build_39/HelloWorld_linux_amd64.tar.gz"
                    }
                  },
                  {
                    "node": {
                      "name": "HelloWorld_darwin_amd64.tar.gz",
                      "downloadUrl": "https://github.com/fractalbach/TestGoReleaser/releases/download/HelloWorld_Build_39/HelloWorld_darwin_amd64.tar.gz"
                    }
                  },
                  {
                    "node": {
                      "name": "HelloWorld_linux_arm6.tar.gz",
                      "downloadUrl": "https://github.com/fractalbach/TestGoReleaser/releases/download/HelloWorld_Build_39/HelloWorld_linux_arm6.tar.gz"
                    }
                  },
                  {
                    "node": {
                      "name": "HelloWorld_linux_arm7.tar.gz",
                      "downloadUrl": "https://github.com/fractalbach/TestGoReleaser/releases/download/HelloWorld_Build_39/HelloWorld_linux_arm7.tar.gz"
                    }
                  }
                ]
              }
            }
          }
        ]
      }
    }
  }
}

*/

/*
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
