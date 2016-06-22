package client

import "fmt"

type speedTest interface {
	printableResult() string
}

// Main entry function
func Run() {
	tests := []speedTest{
		newPinger(),
		newDownloader(),
	}

	for _, t := range tests {
		fmt.Println(t.printableResult())
	}
}
