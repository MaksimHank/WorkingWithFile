package service

import "sync"

type Producer interface {
	Produce() ([]string, error)
}

type Presenter interface {
	Present(path []string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

type job struct {
	index int
	text  string
}

type result struct {
	index int
	value string
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

	masked := make([]string, len(data))
	var wg sync.WaitGroup

	dataChan := make(chan job, len(data))
	results := make(chan result, len(data))

	go func() {
		for i, str := range data {
			dataChan <- job{index: i, text: str}
		}
		close(dataChan)
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for str := range dataChan {
				maskedStr := s.changeTheStringToAsterisks(str.text)
				results <- result{index: str.index, value: maskedStr}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		masked[res.index] = res.value
	}

	return s.pres.Present(masked)
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
