package main

import (
	"fmt"
	"os"
)

func main() {
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	producer := NewFileProducer(inputFile)
	presenter := NewFilePresenter(outputFile)

	service := NewService(producer, presenter)

	if err := service.Run(); err != nil {
		fmt.Println("Error", err)
	}
}
