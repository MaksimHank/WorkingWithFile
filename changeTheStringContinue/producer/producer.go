package producer

import (
	"bufio"
	"fmt"
	"os"
)

type FileProducer struct {
	inputFile string
}

func NewFileProducer(inputFile string) *FileProducer {
	return &FileProducer{inputFile: inputFile}
}

func (fp *FileProducer) Produce() ([]string, error) {
	file, err := os.Open(fp.inputFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var res []string
	for scanner.Scan() {
		data := scanner.Text()
		res = append(res, data)
	}
	if err = scanner.Err(); err != nil {
		fmt.Println("Error reading file", err)
		return nil, err
	}

	return res, nil
}
