package event

import (
	"errors"
	"log"

	"github.com/fatih/color"
	"github.com/r3labs/terraform/builtin/providers/azurerm"
)

// Base : common Event method container
type Base struct {
}

// Log : ...
func (ev *Base) Log(level, message string) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	switch level {
	case "error":
		level = red("[ERROR]")
	case "warn":
		level = yellow("[WARNING]")
	case "info":
		level = blue("[INFO]")
	case "debug":
		level = green("[DEBUG]")
	}

	log.Println(level, message)

}

// Client : not implemented
func (ev *Base) Client() (*azurerm.ArmClient, error) {
	return nil, errors.New("Not implemented")
}
