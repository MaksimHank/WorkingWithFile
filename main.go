package main

import (
	"fmt"
	"github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/logger"
	pres "github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/presenter"
	prod "github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/producer"
	sr "github.com/MaksimHank/WorkingWithFile/changeTheStringContinue/service"
	"github.com/urfave/cli"
	"os"
)

func main() {
	var inputFile string
	var outputFile string

	producer := prod.NewFileProducer(inputFile)
	presenter := pres.NewFilePresenter(outputFile)

	service := sr.NewService(producer, presenter)

	app := &cli.App{
		Name:  "",
		Usage: "",
		Commands: []cli.Command{
			{
				Name:      "Run",
				Usage:     "My app",
				ArgsUsage: "<input path file, out put path file is optional>",
				Action: func(c *cli.Context) error {
					logger.Info("Running app... (service.Run)")
					if c.NArg() < 1 {
						return fmt.Errorf("must be input file in argument")
					}
					inputFile = c.Args().Get(0)
					return service.Run()
				},
			},
		},
	}

	/*if len(os.Args) < 2 {
		fmt.Println("go run main.go inputfile.txt(required) outputfile.txt(optional)")
		return
	}
	inputFile = os.Args[1]

	if len(os.Args) == 3 {
		outputFile = os.Args[2]
	} else {
		outputFile = "output_default.txt"
	}*/

	if err := app.Run(os.Args); err != nil {
		logger.Error("Application error", err)
		os.Exit(1)
	}
}
