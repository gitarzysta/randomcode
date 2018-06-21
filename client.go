package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// List of servers to query
// TODO: This should go into a config file
var servers = []string{"localhost:8001", "localhost:8002", "localhost:8003", "localhost:8004"}

func sum(numbers []int) int {
	sum := 0
	for _, i := range numbers {
		sum += i
	}
	return sum
}

func min(numbers []int) int {
	m := numbers[0]
	for _, i := range numbers {
		if m > i {
			m = i
		}
	}
	return m
}

func max(numbers []int) int {
	m := numbers[0]
	for _, i := range numbers {
		if m < i {
			m = i
		}
	}
	return m
}

//
// TODO: We should probably handle errors better

func geturls(url string) []int {
	numbers := make([]int, 0)

	for _, s := range servers {
		//		fmt.Println(s)
		resp, err := http.Get("http://" + s + url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		value, err := strconv.Atoi(string(body))
		if err != nil {
			panic(err)
		}
		numbers = append(numbers, value)

	}
	return numbers
}

func posturls(url string, pos int) []int {
	numbers := make([]int, 0)
	buf := fmt.Sprint(pos)

	for _, s := range servers {
		//		fmt.Println(s)

		resp, err := http.Post("http://"+s+url, "application/x-www-form-urlencoded", strings.NewReader(buf))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		value, err := strconv.Atoi(string(body))
		if err != nil {
			panic(err)
		}
		numbers = append(numbers, value)

	}
	return numbers
}

func main() {

	sizes := geturls("/size")
	fmt.Println("Sizes", sum(sizes), sizes)
	mpos := sum(sizes) / 2

	medians := geturls("/median")
	sort.Ints(medians)
	idx := len(medians) - 1
	// find biggest median <= global median
	for idx >= 0 {
		positions := posturls("/position", medians[idx])
		fmt.Println("Positions", sum(positions), positions)
		if sum(positions) <= mpos {
			break
		}
		idx--
	}

	// if the biggest of all medians is too small, find the next bigger number
	if idx == len(medians)-1 {
		nextups := posturls("/nextup", medians[len(medians)-1])
		sort.Ints(nextups)
		medians = append(medians, nextups[0])
	}
	// fmt.Println(idx, medians[idx], medians[idx+1], medians)

	low, high := medians[idx], medians[idx+1]
	var med int
	// find two consequitive numbers (low, high) such that median < high
	// and low exists in the data
	for (high - low) > 1 {
		med = (low + high) / 2
		positions := posturls("/position", med)
		fmt.Println("Positions", sum(positions), positions)
		if sum(positions) <= mpos {
			low = med
		} else {
			high = med
		}
		fmt.Println(med, low, high)

	}
	fmt.Println(med, low, high)

	// find the next number after low in case we need it to calculate the median
	medianHigh := posturls("/nextup", low)
	//	positionsLow := posturls("/position", max(medianLow))
	positionsHigh := posturls("/position", min(medianHigh))
	//fmt.Println("GM", medianLow, medianHigh, sum(positionsLow), sum(positionsHigh))
	answer := float64(low)
	// check if the median needs to be recalculated
	if sum(positionsHigh) == mpos+1 {
		if sum(sizes)%2 == 0 {
			answer = (answer + float64(min(medianHigh))) / 2
		} else {
			answer = float64(min(medianHigh))
		}
	}
	fmt.Println("Global median:", answer)

}
