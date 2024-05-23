package utils

type IFiles interface {
	Open(path string) error
}
