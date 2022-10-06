package output

import (
	"fmt"
)

// Info represents any type that can log via Infof.
type Info interface {
	Infof(string, ...any)
}

// Error represents any type that can log via Errorf.
type Error interface {
	Errorf(string, ...any)
}

// CLI is a wrapper around fmt.Print functions for satisfying interfaces.
type CLI struct{}

func (o *CLI) print(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	fmt.Println(s)
}

func (o *CLI) Infof(msg string, args ...any) {
	o.print(msg, args...)
}

func (o *CLI) Errorf(msg string, args ...any) {
	o.print(msg, args...)
}

type Logger interface {
	Info
	Error
}
