package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

type solution struct {
	union               []int
	intersection        []int
	difference_1_2      []int
	difference_2_1      []int
	symmetricDifference []int
}

func NewSolution(
	union []int,
	intersection []int,
	difference_1_2 []int,
	difference_2_1 []int,
	symmetricDifference []int) *solution {
	sol := solution{
		union:               union,
		intersection:        intersection,
		difference_1_2:      difference_1_2,
		difference_2_1:      difference_2_1,
		symmetricDifference: symmetricDifference}

	return &sol

}

type TestCase struct {
	array1 []int
	array2 []int
	sol    solution
}

func NewTestCase(array1 []int, array2 []int, sol solution) *TestCase {
	test_case := TestCase{
		array1: array1,
		array2: array2,
		sol:    sol}
	return &test_case
}

func findMin(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func makeSet(array []int) []int {
	set := make([]int, 0)
	hashMap := make(map[int]bool)
	for _, element := range array {
		if !hashMap[element] {
			hashMap[element] = true
			set = append(set, element)
		}
	}

	return set
}

func makeSuperSet(array_1, array_2 []int) []int {
	superSet := make([]int, 0)
	for _, element := range array_1 {
		superSet = append(superSet, element)
	}

	for _, element := range array_2 {
		superSet = append(superSet, element)
	}
	return superSet
}

func makeUnion(array_1, array_2 []int) []int {
	return makeSet(makeSuperSet(array_1, array_2))
}

func makeIntersection(array_1, array_2 []int) []int {
	array_1 = makeSet(array_1)
	array_2 = makeSet(array_2)

	intersectionArray := make([]int, 0)
	hashMap := make(map[int]bool)

	for _, element := range array_1 {
		hashMap[element] = true
	}

	for _, element := range array_2 {
		if hashMap[element] {
			intersectionArray = append(intersectionArray, element)
		}
	}

	return intersectionArray
}

func makeDifference(array_1, array_2 []int) []int {

	array_1 = makeSet(array_1)
	array_2 = makeSet(array_2)

	differenceArray := make([]int, 0)
	hashMap := make(map[int]bool)

	for _, element := range array_2 {
		hashMap[element] = true
	}

	for _, element := range array_1 {
		if !hashMap[element] {
			differenceArray = append(differenceArray, element)
		}
	}

	return differenceArray
}

func makeSymmetricDifference(array_1, array_2 []int) []int {
	return makeSuperSet(
		makeDifference(array_1, array_2),
		makeDifference(array_2, array_1))
}

func isSlicesEqual(a, b []int) bool {
	channel := make(chan bool)

	if len(a) != len(b) {
		return false
	}

	sort.Ints(a)
	sort.Ints(b)

	for i := 0; i < len(a); i++ {
		go func(i1, i2 int) {
			channel <- (i1 == i2)
		}(a[i], b[i])
	}

	for i := 0; i < len(a); i++ {
		equal := <-channel
		if !equal {
			return false
		}
	}

	return true

}

func runTestCase(idx int, tCase TestCase, wg *sync.WaitGroup) {
	passed := true
	channel := make(chan bool)

	go func(array_1, array_2 []int) {
		channel <- isSlicesEqual(tCase.sol.union, makeUnion(array_1, array_2))
	}(tCase.array1, tCase.array2)

	go func(array_1, array_2 []int) {
		channel <- isSlicesEqual(
			tCase.sol.intersection, makeIntersection(array_1, array_2))
	}(tCase.array1, tCase.array2)

	go func(array_1, array_2 []int) {

		channel <- isSlicesEqual(
			tCase.sol.difference_1_2, makeDifference(array_1, array_2))
	}(tCase.array1, tCase.array2)

	go func(array_1, array_2 []int) {
		channel <- isSlicesEqual(
			tCase.sol.difference_2_1, makeDifference(array_2, array_1))
	}(tCase.array1, tCase.array2)

	go func(array_1, array_2 []int) {
		channel <- isSlicesEqual(
			tCase.sol.symmetricDifference, makeSymmetricDifference(
				array_1, array_2))
	}(tCase.array1, tCase.array2)
	for i := 0; i < 5; i++ {
		if passed {
			passed = <-channel
		}
	}
	if passed {
		fmt.Println(
			"TestCase " + strconv.Itoa(idx) + " has passed")
	} else {
		fmt.Println(
			"TestCase " + strconv.Itoa(idx) + " has failed")
	}

	wg.Done()
}

func main() {

	var wg sync.WaitGroup
	testCases := make([]TestCase, 0)

	testCases = append(testCases, *NewTestCase(
		[]int{1, 2, 3},
		[]int{2, 3, 4},
		*NewSolution(
			[]int{1, 2, 3, 4},
			[]int{2, 3},
			[]int{1},
			[]int{4},
			[]int{1, 4})))
	testCases = append(testCases, *NewTestCase(
		[]int{1, 2, 3, 4},
		[]int{2, 3, 4},
		*NewSolution(
			[]int{1, 2, 3, 4},
			[]int{2, 3, 4},
			[]int{1},
			[]int{},
			[]int{1})))

	for idx, testCase := range testCases {
		wg.Add(1)
		go runTestCase(idx, testCase, &wg)
	}

	wg.Wait()

}
