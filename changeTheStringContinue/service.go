package main

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

type FilePath struct {
	inputPath  string
	outputPath string
}

func (fp *FilePath) Producer() ([]string, error) {

}

func (fp *FilePath) Presenter(path []string) error {

}

func changeTheStringToAsterisks(text string) string {
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
