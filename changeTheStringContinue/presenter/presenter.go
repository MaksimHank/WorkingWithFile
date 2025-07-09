package presenter

import (
	"bufio"
	"fmt"
	"os"
)

type FilePresenter struct {
	OutputFile string
}

func NewFilePresenter(outputFile string) *FilePresenter {
	return &FilePresenter{OutputFile: outputFile}
}

func (fp *FilePresenter) Present(data []string) error {
	file, err := os.Create(fp.OutputFile)
	if err != nil {
		fmt.Println("Error creating file!")
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range data {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing data to file!")
			return err
		}
	}

	return writer.Flush()
}
