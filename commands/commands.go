// Package commands contains the methods and helpers for plugins to be exported
// for use by bogon.
// it also exports a set of basic plugins for use
package commands

import (
	"regexp"

	"github.com/go-kit/kit/log"
)

// RegisterFunc is exported by a plugin
// name, regular expression, method and options
type RegisterFunc func() (string, *regexp.Regexp, CommandFunc, Options)

// CommandFunc is a function type that all plugins register to be called
type CommandFunc func(Message, MessageFunc) string

// MessageFunc is a return function to send data back inside of a command
type MessageFunc func(string, ...interface{})

// Options con tains options that describe how a command works
type Options struct {
	AdminOnly bool
	Raw       bool
}

// Message represents a message sent to a command
type Message struct {
	Params []string
	To     string
	From   string

	Logger log.Logger
}
