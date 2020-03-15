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

	// The following code will open a browser window, but
	// requires additional permissions:

	//err := openBrowser("localhost:3030")
	//if err != nil {
	//	log.Fatal(err)
	//}
	wg.Wait()
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

func handlePage(f string) http.HandlerFunc {
	page, err := convert(f)
	if err != nil {
		log.Fatalf("Error converting file: %v", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, page)
	}
}

func openBrowser(url string) error {
	fox := "/Applications/Firefox.app"
	cmd := exec.Command(fox, url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Print("attempting to open browser...\n")
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
