package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var numbers []int

func main() {

	var directory = flag.String("directory", "./data", "data directory")
	var prefix = flag.String("prefix", "shard", "file name prefix")
	var suffix = flag.String("suffix", ".txt", "file name suffix")
	var shard = flag.Int("shard", 1, "shard number to serve")
	var lport = flag.Int("port", 0, "port number to listen on, default 8000 + shard")
	var datafile = flag.String("datafile", "", "data file, default directory/prefix+shard+suffix")
	var number int
	flag.Parse()
	numbers = make([]int, 0)
	fname := fmt.Sprintf("%s/%s%d%s", *directory, *prefix, *shard, *suffix)
	if *datafile != "" {
		fname = *datafile
	}
	port := *lport
	if port == 0 {
		port = 8000 + *shard
	}
	fh, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	for {
		if _, err := fmt.Fscanln(fh, &number); err == nil {
			//			fmt.Println(number)
			numbers = append(numbers, number)
		} else {
			break
		}

	}
	if len(numbers) == 0 {
		panic("No data found")
	}
	sort.Ints(numbers)

	http.HandleFunc("/position", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}
			value, err := strconv.Atoi(string(body))
			if err != nil {
				log.Printf("Error reading position: %v", err)
				http.Error(w, "can't read position", http.StatusBadRequest)
				return

			}
			pos := sort.SearchInts(numbers, value) + 1 // number from 1
			fmt.Fprintf(w, strconv.Itoa(pos))

		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

		//		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	// Next number bigger than posted value
	// TODO: The case where there is no bigger number should be handled
	http.HandleFunc("/nextup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}
			value, err := strconv.Atoi(string(body))
			if err != nil {
				log.Printf("Error reading position: %v", err)
				http.Error(w, "can't read position", http.StatusBadRequest)
				return

			}
			pos := sort.SearchInts(numbers, value)
			for {
				pos++
				if pos >= len(numbers) {
					break
				}
				if numbers[pos] > value {
					value = numbers[pos]
					break
				}
			}
			fmt.Fprintf(w, strconv.Itoa(value))

		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	})
	/*
		http.HandleFunc("/nextdownorequal", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Printf("Error reading body: %v", err)
					http.Error(w, "can't read body", http.StatusBadRequest)
					return
				}
				value, err := strconv.Atoi(string(body))
				if err != nil {
					log.Printf("Error reading position: %v", err)
					http.Error(w, "can't read position", http.StatusBadRequest)
					return

				}
				pos := sort.SearchInts(numbers, value)
				for {
					if pos < 0 {
						break
					}
					if numbers[pos] <= value {
						value = numbers[pos]
						break
					}
					pos--
				}
				fmt.Fprintf(w, strconv.Itoa(value))

			} else {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}

		})
	*/
	// This provides the highest number smaller than median
	http.HandleFunc("/median", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, strconv.Itoa(numbers[(len(numbers)-1)/2]))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/size", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, strconv.Itoa(len(numbers)))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
