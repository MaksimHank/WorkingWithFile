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

	app := &cli.App{
		Name:  "Convert some text to asterisks",
		Usage: "cli",
		Before: func(c *cli.Context) error {
			logger.SetLevel(c.String("log-level"))
			logger.Init()
			return nil
		},
		Commands: []cli.Command{
			{
				Name:      "run",
				Usage:     "my app",
				ArgsUsage: "<input path file is required, out put path file is optional>",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("must be input file in argument")
					}
					inputFile = c.Args().Get(0)

					if c.NArg() >= 2 {
						outputFile = c.Args().Get(1)
					} else {
						outputFile = "output_default.txt"
					}

					producer := prod.NewFileProducer(inputFile)
					presenter := pres.NewFilePresenter(outputFile)

					service := sr.NewService(producer, presenter)
					logger.Info("Running app... (service.Run)")
					return service.Run()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error("Application error", err)
		os.Exit(1)
	}
}
