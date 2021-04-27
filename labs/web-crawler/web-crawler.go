// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string, height int, depthlimit int) []string {
	if height <= depthlimit {
		//fmt.Println(url)
		log.Println(url)

		tokens <- struct{}{} // acquire a token
		list, err := links.Extract(url)
		<-tokens // release the token
		//fmt.Println(list)
		if err != nil {
			log.Print(err)
		}
		return list
	} else { //Overpassed the depth limit
		var list []string
		return list
	}

}

//!-sema

//!+
func main() {
	if len(os.Args) != 4 {
		fmt.Println("Please input your arguments correctly")
		os.Exit(0)
	}
	worklist := make(chan []string)
	var n int // number of pending sends to worklist

	// Start with the command-line arguments.
	n++

	filename := strings.Split(os.Args[2], "=")
	file, err := os.OpenFile(filename[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	log.SetOutput(file)
	log.SetFlags(0)

	depthnum := strings.Split(os.Args[1], "=")

	go func() { worklist <- os.Args[3:] }()

	depthlimit, _ := strconv.Atoi(depthnum[1])

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	urlDep := make(map[string]int)
	height := 0

	//sampledata := []string

	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {

			if !seen[link] {
				seen[link] = true
				n++
				aux := height + 1
				urlDep[link] = aux
				go func(link string) {
					worklist <- crawl(link, height, depthlimit)

				}(link)
			}

		}
		height++
	}

}

//!-
