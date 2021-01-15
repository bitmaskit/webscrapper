package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var url = "https://bitmask.it/"

func main() {
	// Create HTTP client with timout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create out file
	dirName := makeDir("download")

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "Not Firefox")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Get the response body as string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	pageTitle, err := FindTitle(pageContent)
	if err != nil {
		fmt.Println("Couldn't find title: ", err)
	} else {
		fmt.Printf("Page title: %s\n", pageTitle)
	}

	outFile, err := os.Create(dirName + "/output.html")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Copy data from the response to standard output
	n, err := io.Copy(outFile, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Number of bytes copied to file: ", n)
}

func FindTitle(content string) ([]byte, error) {
	titleStartIdx := strings.Index(content, "<title>")
	if titleStartIdx == - 1 {
		return []byte(""), errors.New("no title")
	}
	titleEndIdx := strings.Index(content, "</title>")
	if titleEndIdx == - 1 {
		return []byte(""), errors.New("no closing tag")
	}
	titleStartIdx += 7 // offset the html tag

	return []byte(content[titleStartIdx:titleEndIdx]), nil
}

func makeDir(dirName string) string {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		// try to create
		err = os.Mkdir(dirName, os.ModePerm)
	} else if os.IsExist(err) {
		err = os.Chmod(dirName, os.ModePerm)
	} else {
		dirName = "."
	}
	return dirName
}
