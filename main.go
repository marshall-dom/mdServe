package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

func main() {
	fp := os.Args[1]
	port := ":3030"
	html, err := convert(fp)
	if err != nil {
		log.Fatalf("Error converting file: %v", err)
	}

	http.Handle("/", handlePage(html))
	fmt.Printf("Spinning up server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func convert(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	md := []byte(data)
	htmlBytes := markdown.ToHTML(md, nil, nil)
	html := "<html>" + string(htmlBytes) + "</html>"
	return html, nil

}

func handlePage(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, page)
	}
}
