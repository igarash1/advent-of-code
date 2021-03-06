package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"aoc"
)

type INTERVAL struct {
	left, right int
}

func (ivl INTERVAL) Within(p int) bool {
	return ivl.left <= p && p <= ivl.right
}

func splitToInts(fields string) []int {
	var nums []int
	ss := strings.Split(fields, ",")
	for _, snum := range ss {
		nums = append(nums, aoc.ToInt(snum))
	}
	return nums
}

func computeErrors(tickets []int, ivls []INTERVAL) int {
	ret := 0
	for _, t := range tickets {
		valid := false
		for _, ivl := range ivls {
			if ivl.Within(t) {
				valid = true
				break
			}
		}
		if !valid {
			ret += t
		}
	}
	return ret
}

func main() {
	var ivls []INTERVAL

	// scan the rules
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			log.Print("BLANK LINE")
			break
		}
		st := strings.Index(line, ":")
		if st == -1 {
			log.Print("INVALID RULE LINE", line)
		}

		ss := bufio.NewScanner(strings.NewReader(line[st+1:]))
		ss.Split(bufio.ScanWords)
		for ss.Scan() {
			sivl := ss.Text()
			if sivl == "or" {
				continue
			}
			var l, r int
			_, err := fmt.Sscanf(sivl, "%d-%d", &l, &r)
			if err != nil {
				log.Fatal(err)
			}
			ivls = append(ivls, INTERVAL{l, r})
		}
	}

	// scan your ticket
	scanner.Scan()
	if scanner.Text() != "your ticket:" {
		log.Fatal("expected 'your ticket:', but ", scanner.Text())
	}
	scanner.Scan()
	// ignore your ticket for now
	splitToInts(scanner.Text())

	// scan nearby tickets
	scanner.Scan()
	log.Print("BLANK LINE")
	if scanner.Scan(); scanner.Text() != "nearby tickets:" {
		log.Fatal("expected 'nearby tickets:', but ", scanner.Text())
	}

	// compute errors
	result := 0
	for scanner.Scan() {
		nearbyTicket := splitToInts(scanner.Text())
		result += computeErrors(nearbyTicket, ivls)
	}

	fmt.Println(result)
}
