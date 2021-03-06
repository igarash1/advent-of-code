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
	for _, s := range ss {
		nums = append(nums, aoc.ToInt(s))
	}
	return nums
}

func isValidTicket(ticket int, ivls []INTERVAL) bool {
	if ticket == -1 {
		return true
	}
	for _, ivl := range ivls {
		if ivl.Within(ticket) {
			return true
		}
	}
	return false
}

func getValidTickets(tickets []int, ivlsMap map[string][]INTERVAL) []int {
	var validTickets []int
	for _, t := range tickets {
		valid := false
		for _, ivls := range ivlsMap {
			if isValidTicket(t, ivls) {
				valid = true
				validTickets = append(validTickets, t)
				break
			}
		}
		// the discarded values are represented by -1
		if !valid {
			validTickets = append(validTickets, -1)
		}
	}
	return validTickets
}

func removeFromInts(nums []int, target int) []int {
	var ret []int
	for _, v := range nums {
		if v != target {
			ret = append(ret, v)
		}
	}
	return ret
}

func main() {
	ivlsMap := make(map[string][]INTERVAL)

	// scan the rules
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			log.Print("BLANK LINE")
			break
		}
		colonPos := strings.Index(line, ":")
		if colonPos == -1 {
			log.Print("INVALID RULE LINE", line)
		}
		fieldName := line[:colonPos]

		ss := bufio.NewScanner(strings.NewReader(line[colonPos+1:]))
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
			ivlsMap[fieldName] = append(ivlsMap[fieldName], INTERVAL{l, r})
		}
	}

	// scan your ticket
	scanner.Scan()
	if scanner.Text() != "your ticket:" {
		log.Fatal("expected 'your ticket:', but ", scanner.Text())
	}

	scanner.Scan()
	yourTickets := splitToInts(scanner.Text())
	validYTs := getValidTickets(yourTickets, ivlsMap)

	// scan nearby tickets
	scanner.Scan()
	log.Print("BLANK LINE")
	if scanner.Scan(); scanner.Text() != "nearby tickets:" {
		log.Fatal("expected 'nearby tickets:', but ", scanner.Text())
	}

	var nearbyTickets [][]int
	for scanner.Scan() {
		nearbyTicket := splitToInts(scanner.Text())
		nearbyTickets = append(nearbyTickets, getValidTickets(nearbyTicket, ivlsMap))
	}

	// enumerate possible indices for each field
	// (which corresponds to making a bipartite graph)
	possibleIdx := make(map[string][]int)
	for i, yt := range validYTs {
		for fName, ivls := range ivlsMap {
			if !isValidTicket(yt, ivls) {
				continue
			}
			possible := true
			for _, nearbyTicket := range nearbyTickets {
				if !isValidTicket(nearbyTicket[i], ivls) {
					possible = false
				}
			}
			if possible {
				log.Printf("the field of %d-th index can be %s", i, fName)
				possibleIdx[fName] = append(possibleIdx[fName], i)
			}
		}
	}

	// determine which field is which
	// (which corresponds to determining a perfect matching in a bipartite graph)
	fieldIdx := make(map[string]int)
	for len(fieldIdx) < len(ivlsMap) {
		dIdx := -1
		for fName, idx := range possibleIdx {
			if len(idx) == 1 {
				dIdx = idx[0]
				log.Printf("the field of %d-th index is determined to be %s", dIdx, fName)
				fieldIdx[fName] = dIdx
				break
			}
		}
		if dIdx == -1 {
			log.Fatal("something went wrong")
		}
		// remove the determined index from candidates
		for fName := range possibleIdx {
			possibleIdx[fName] = removeFromInts(possibleIdx[fName], dIdx)
		}
	}

	// compute the multiple of the 'departure' values on your ticket
	result := 1
	for fName, i := range fieldIdx {
		if strings.HasPrefix(fName, "departure") {
			result *= yourTickets[i]
		}
	}
	fmt.Println(result)
}
