package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"sync"

	"github.com/gomarkdown/markdown"
)

func main() {
	fp := os.Args[1]
	port := ":3030"

	http.Handle("/", handlePage(fp))
	fmt.Printf("Spinning up server on port %s...\n", port)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Fatal(http.ListenAndServe(port, nil))
		wg.Done()
	}()

	err := openBrowser("localhost:3030")
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}

// convert returns an html string from the specified markdown file.
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

// handlePage takes a file path and returns a handler func that serves
// an html file.
func handlePage(f string) http.HandlerFunc {
	page, err := convert(f)
	if err != nil {
		log.Fatalf("Error converting file: %v", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, page)
	}
}

// openBrowser opens a browser window to a specified url.
func openBrowser(url string) error {
	fox := "/Applications/Firefox.app/Contents/MacOS/firefox-bin"
	cmd := exec.Command(fox, url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
		if err != nil {
			return err
		}
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		if err != nil {
			return err
		}
	case "darwin":
		err = cmd.Run()
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("unsupported platform")
		if err != nil {
			return err
		}
	}
	return nil
}
