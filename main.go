package main

import (
	"fmt"
	"os"

	sr "github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/service"
)

func main() {
	var inputFile string
	var outputFile string

	if len(os.Args) < 2 {
		fmt.Println("go run main.go inputfile.txt(required) outputfile.txt(optional)")
		return
	}
	inputFile = os.Args[1]

	if len(os.Args) == 3 {
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
