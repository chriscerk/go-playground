package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type search struct {
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	Filename string `json:"filename"`
}

func main() {
	fmt.Println("Scanning your project...")

	var channels []<-chan string

	fileSearches := getFileSearches()

	for _, s := range fileSearches {
		c := searchFile(s.Pattern, s.Name, s.Filename)
		channels = append(channels, c)
	}

	mainChannel := fanIn(channels)

	i := 1
	for i < 100 {
		fmt.Printf(<-mainChannel)
		time.Sleep(1 * time.Millisecond)
		i++
	}

	fmt.Println("All Scanning complete")
	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func fanIn(cSet []<-chan string) <-chan string {
	c := make(chan string)
	for i := range cSet {
		go func(in <-chan string) {
			for {
				x := <-in
				c <- x
			}
		}(cSet[i])
	}
	return c
}

func searchFile(pattern, scannerName, filename string) <-chan string {
	c := make(chan string)
	go func() {
		if filename == "all" {
			allFilesChannel := searchAllFiles(pattern, scannerName, c)
			go func() {
				for {
					c <- <-allFilesChannel
				}
			}()
		} else {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			re := regexp.MustCompile(pattern)

			scanner := bufio.NewScanner(file)
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				matches := re.FindStringSubmatch(line)

				if len(matches) > 0 {
					c <- fmt.Sprintf("%s | matches %s in %s on line %d \n", scannerName, matches[0], filename, lineNumber)
				}
				lineNumber++
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			c <- scannerName + "-DONE \n"
			close(c)
		}
	}()
	return c
}

func searchAllFiles(pattern, scannerName string, c <-chan string) <-chan string {
	searchDir := "C:\\Code\\go-playground\\valigator"

	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		isDirectory := f.IsDir()
		if !isDirectory {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	var channels []<-chan string

	for _, file := range fileList {
		c := searchFile(pattern, scannerName, file)
		channels = append(channels, c)
	}

	return fanIn(channels)
}

func getFileSearches() []search {
	raw, err := ioutil.ReadFile("./searches.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []search
	json.Unmarshal(raw, &c)
	return c
}
