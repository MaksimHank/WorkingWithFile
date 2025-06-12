package main

type Producer interface {
	Produce() ([]string, error)
}

type Presenter interface {
	Present(path []string) error
}
