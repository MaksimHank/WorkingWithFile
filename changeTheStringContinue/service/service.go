package service

import (
	"bufio"
	"fmt"
	"os"
)

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service {
	return &Service{
		prod: prod,
		pres: pres,
	}
}

func (s *Service) Run() error {
	data, err := s.prod.Produce()
	if err != nil {
		return err
	}

	return s.pres.Present(data)
}

type FileProducer struct {
	inputFile string
}

type FilePresenter struct {
	outputFile string
}

func NewFileProducer(inputFile string) *FileProducer {
	return &FileProducer{inputFile: inputFile}
}

func NewFilePresenter(outputFile string) *FilePresenter {
	return &FilePresenter{outputFile: outputFile}
}

func (fp *FileProducer) Produce() ([]string, error) {
	file, err := os.Open("inputFile")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	data := scanner.Text()
	sliceStr := make([]string, 0, len(data))
	for scanner.Scan() {
		sliceStr = append(sliceStr, data)
	}
	if er := scanner.Err(); er != nil {
		fmt.Println("Error reading file", er)
		return nil, er
	}

	return changeTheStringToAsterisks(sliceStr), err
}

func (fp *FilePresenter) Present(path []string) error {

	return nil
}

func (s *Service) changeTheStringToAsterisks(text string) string {
	str := []byte(text)
	prefix := []byte("http://")
	prefixlen := len(prefix)
	var result []byte
	var token []byte

	for i := 0; i < len(str); i++ {
		ch := str[i]
		if ch == ' ' || ch == ',' || ch == ';' {
			if len(token) > 0 {
				isHTTP := false
				if len(token) >= prefixlen {
					isHTTP = true
					for j := 0; j < prefixlen; j++ {
						if token[j] != prefix[j] {
							isHTTP = false
							break
						}
					}
				}
				if isHTTP {
					for j := 0; j < len(token); j++ {
						if j < prefixlen {
							result = append(result, token[j])
						} else {
							result = append(result, '*')
						}
					}
				} else {
					result = append(result, token...)
				}
				token = token[:0]
			}
			result = append(result, ch)
		} else {
			token = append(token, ch)
		}
	}

	if len(token) > 0 {
		isHTTP := false
		if len(token) >= prefixlen {
			isHTTP = true
			for j := 0; j < prefixlen; j++ {
				if token[j] != prefix[j] {
					isHTTP = false
					break
				}
			}
		}
		if isHTTP {
			for j := 0; j < len(token); j++ {
				if j < prefixlen {
					result = append(result, token[j])
				} else {
					result = append(result, '*')
				}
			}
		} else {
			result = append(result, token...)
		}
	}

	return string(result)
}
