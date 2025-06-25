package main

import (
	"fmt"
	"os"

	sr "github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/service"
)

func main() {
	var inputFile string
	var outputFile string

	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	} else {
		fmt.Println("Error. Didn't find inputfile")
	}

	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	} else {
		outputFile = "output_default.txt"
	}

	producer := sr.NewFileProducer(inputFile)
	presenter := sr.NewFilePresenter(outputFile)

	service := sr.NewService(producer, presenter)

	if err := service.Run(); err != nil {
		fmt.Println("Error", err)
	}
}
