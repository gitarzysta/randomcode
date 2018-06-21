package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

func main() {
	var shards = flag.Int("shards", 4, "number of shards")
	var size = flag.Int("size", 10000, "shard size")
	var directory = flag.String("directory", "./data", "data directory")
	var prefix = flag.String("prefix", "shard", "file name prefix")
	var suffix = flag.String("suffix", ".txt", "file name suffix")
	var force = flag.Bool("force", false, "override files in data directory")
	var min = flag.Int("min", 0, "minimum number")
	var max = flag.Int("max", 1000, "maximum number")

	flag.Parse()

	numbers := make([]int, 0)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	fmt.Println(*shards, *size, *directory, *prefix, *suffix)
	if *shards <= 2 {
		fmt.Println("invalid number of shards", *shards)
		os.Exit(1)
	}
	if *size <= 0 {
		fmt.Println("invalid shard size", *size)
		os.Exit(1)
	}
	if *max < *min {
		fmt.Println("invalid range", *min, "-", *max)
		os.Exit(1)
	}
	for i := 1; i <= *shards; i++ {
		fmt.Println("Shard:", i)
		fname := fmt.Sprintf("%s/%s%d%s", *directory, *prefix, i, *suffix)
		if _, err := os.Stat(fname); os.IsNotExist(err) || *force {

			fh, err := os.Create(fname)
			if err != nil {
				panic(err)
			}
			for j := 0; j < *size; j++ {
				v := *min + r.Intn(*max-*min)
				fmt.Fprintln(fh, v)
				numbers = append(numbers, v)

			}
			err = fh.Close()
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("File", fname, "exists and no --force used")
			os.Exit(1)
		}

	}

	fname := fmt.Sprintf("%s/%s", *directory, "median")
	fh, err := os.Create(fname)
	sort.Ints(numbers)
	if len(numbers)%2 == 0 {
		fmt.Fprintln(fh, float64(numbers[(len(numbers)-1)/2]+numbers[len(numbers)/2])/2.0)
	} else {
		fmt.Fprintln(fh, numbers[(len(numbers)-1)/2])
	}
	err = fh.Close()
	if err != nil {
		panic(err)
	}

}
