package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func main() {
	type Result struct {
		input  string
		output uuid.UUID
		err    error
	}
	uuids := []Result{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		output, err := uuid.Parse(input)
		uuids = append(uuids, Result{
			input:  input,
			output: output,
			err:    err,
		})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("scanner.Err() failed: %+v", err))
	}

	failures := []Result{}
	passed := []Result{}
	for _, result := range uuids {
		if result.err != nil {
			failures = append(failures, result)
		} else {
			passed = append(passed, result)
		}
	}

	fmt.Printf("Failed:\n")
	for _, failed := range failures {
		fmt.Printf("\tuuid=%+v, err=%+v\n", failed.input, failed.err)
	}

	fmt.Printf("Passed:\n")
	for _, pass := range passed {
		fmt.Printf("\tuuid=%+v, result=%+v\n", pass.input, pass.output)
	}
}
