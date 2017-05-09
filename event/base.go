package event

import (
	"log"

	"github.com/fatih/color"
)

// Base : common Event method container
type Base struct {
	Name          string `json:"name" validate:"required"`
	ProviderType  string `json:"_provider"`
	ComponentType string `json:"_component"`
	ComponentID   string `json:"_component_id"`
	State         string `json:"_state"`
	Action        string `json:"_action"`
	Service       string `json:"service"`
	ErrorMessage  string `json:"error,omitempty"`
	Subject       string `json:"-"`
	Body          []byte `json:"-"`
	CryptoKey     string `json:"-"`
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
