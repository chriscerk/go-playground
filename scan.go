package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type search struct {
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	Filename string `json:"filename"`
}

func main() {
	fmt.Println("Scanning your project...")

	var channels []<-chan string

	searches := getSearches()

	for _, s := range searches {
		c := searchFile(s.Pattern, s.Name, s.Filename)
		channels = append(channels, c)
	}

	mainChannel := fanIn(channels)

	for {
		fmt.Printf(<-mainChannel)
	}

	fmt.Println("All Scanning complete")
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
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		re := regexp.MustCompile(pattern)

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			line := scanner.Text()
			matches := re.FindStringSubmatch(line)

			if len(matches) > 0 {
				c <- fmt.Sprintf("%s | matches %s on line %d \n", scannerName, matches[0], lineNumber)
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		c <- scannerName + "-DONE \n"
		close(c)
	}()
	return c
}

func getSearches() []search {
	raw, err := ioutil.ReadFile("./searches.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []search
	json.Unmarshal(raw, &c)
	return c
}
