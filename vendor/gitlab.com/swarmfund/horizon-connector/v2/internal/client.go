package internal

import "io"

// Client exists for testing purpose only
//go:generate mockery -case underscore -name Client
type Client interface {
	Get(string) ([]byte, error)
	Put(string,io.Reader) ([]byte, error)
}
