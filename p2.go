package main

import (
  "fmt"
  "sync"
)

type Solution struct {
  resultAdd [][]int
  resultSubstract [][]int
  resultMultiply [][]int
  resultTransposeA [][]int
  resultTransposeB [][]int
}

func NewSolution(resultAdd,
  resultSubstract,
  resultMultiply,
  resultTransposeA,
  resultTransposeB [][]int) *Solution {
  solution := Solution{
    resultAdd,
    resultSubstract,
    resultMultiply,
    resultTransposeA,
    resultTransposeB}
  return &solution
}

type TestCase struct {
  a [][]int
  b [][]int
  solution Solution
}

func NewTestCase(a, b [][]int, solution Solution) *TestCase {
  testCase := TestCase{a, b, solution}
  return &testCase
}

func addMatrices(a, b [][]int) [][]int {
  result := make([][]int, len(a))
  for i := range result {
    result[i] = make([]int, len(a[0]))
  }

  var wg sync.WaitGroup

  for i := range a {
    for j := range a[0] {
      wg.Add(1)
      go func (a, b, i, j int,result [][]int, wg *sync.WaitGroup) {
        result[i][j] = a + b
        wg.Done()
      }(a[i][j], b[i][j],
        i, j,
        result,
        &wg)
    }
  }

  wg.Wait()
  return result
}

func subtractMatrices(a, b [][]int) [][]int {
  result := make([][]int, len(a))
  for i := range result {
    result[i] = make([]int, len(a[0]))
  }

  var wg sync.WaitGroup

  for i := range a {
    for j := range a[0] {
      wg.Add(1)
      go func (a, b, i, j int,result [][]int, wg *sync.WaitGroup) {
        result[i][j] = a - b
        wg.Done()
      }(a[i][j], b[i][j],
        i, j,
        result,
        &wg)
    }
  }

  wg.Wait()
  return result
}

func multiplyMatrices(a, b [][]int) [][]int {
  if (len(a) != len(b[0])) {
    panic("Row Count for first array must match Column Count for second")
  }
  result := make([][]int, len(a))
  for i := range result {
    result[i] = make([]int, len(b[0]))
  }

  var wg sync.WaitGroup

  for i := range a {
    for j := range b[0] {
      wg.Add(1)

      go func (a, b [][]int, i, j int, result [][]int, wg *sync.WaitGroup) {
        for k := range a[0] {
            result[i][j] += a[i][k] * b[k][j]
        }
        wg.Done()
      }(a, b, i, j, result, &wg)
    }
  }

  wg.Wait()

  return result

}

func transpose(a [][]int) [][]int {
  result := make([][]int, len(a))
  for i := range result {
    result[i] = make([]int, len(a[0]))
  }

  var wg sync.WaitGroup

  for i := range a {
    for j := range a[0] {
      wg.Add(1)
      go func(a, i, j int, result [][]int, wg *sync.WaitGroup) {
        result[j][i] = a
        wg.Done()
      }(a[i][j], i, j, result, &wg)
    }
  }

  wg.Wait()
  return result

}

func isSlicesEqual(a, b [][]int) bool {
	channel := make(chan bool)

	if ((len(a) != len(b)) || (len(a[0]) != len(b[0]))) {
		return false
	}

  for i := range a {
    for j := range a[0] {
      go func (i1, i2 int) {
        channel <- (i1 == i2)
      }(a[i][j], b[i][j])
    }
  }
	for i := 0;i < len(a) * len(a[0]);i++ {
		equal := <- channel
		if (!equal) {
			return false
		}
	}

	return true

}

func runTestCase(testCase TestCase, idx int, waitGroup *sync.WaitGroup) {
  channel := make(chan bool)
  passed := true

  go func() {
    channel <- isSlicesEqual(
      addMatrices(testCase.a, testCase.b),
        testCase.solution.resultAdd)}()

  go func() {
    channel <- isSlicesEqual(
      subtractMatrices(testCase.a, testCase.b),
        testCase.solution.resultSubstract)}()

  go func() {
    channel <- isSlicesEqual(
      multiplyMatrices(testCase.a, testCase.b),
        testCase.solution.resultMultiply)}()

  go func() {
    channel <- isSlicesEqual(
      transpose(testCase.a), testCase.solution.resultTransposeA)}()

  go func() {
    channel <- isSlicesEqual(
      transpose(testCase.b), testCase.solution.resultTransposeB)}()

  for i := 0;i < 5;i++ {
    if (passed) {
      passed = <- channel
    }
  }

  if (passed) {
    fmt.Println("Test Case", idx + 1, "has passed")
  } else {
    fmt.Println("Test Case", idx + 1, "has failed")
  }

  waitGroup.Done()

}

func main() {

  testCases := make([]TestCase, 0)
  testCases = append(testCases, *NewTestCase(
      [][]int{{1, 2}, {3, 4}},
      [][]int{{2, 0}, {1, 2}},
      *NewSolution(
        [][]int {{3, 2}, {4, 6}},
        [][]int {{-1, 2}, {2, 2}},
        [][]int {{4, 4}, {10, 8}},
        [][]int {{1, 3}, {2, 4}},
        [][]int {{2, 1}, {0, 2}})))
  var waitGroup sync.WaitGroup
  for idx, testCase := range testCases {
    waitGroup.Add(1)
    go runTestCase(testCase, idx, &waitGroup)
  }
  waitGroup.Wait()
}
