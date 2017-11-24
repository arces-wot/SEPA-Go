package sepa

import (
	"encoding/json"
	"io"
)

type PortsType struct {
	Http int
	Ws   int
}

type PathsType struct {
	Query     string
	Update    string
	Subscribe string
}

type Configuration struct {
	Host     string
	Ports    PortsType
	Paths    PathsType
	Extended map[string]interface{}
}

type parameters struct {
	Parameters Configuration
}

var default_config = parameters{
	Parameters: Configuration{
		Ports:    PortsType{Http: 8000, Ws: 9000},
		Paths:    PathsType{Query: "/query", Update: "/update", Subscribe: "/subscribe"},
		Extended: make(map[string]interface{}),
	},
}

func DefaultConfig() Configuration {
	return default_config.Parameters
}

func ConfigFromJson(reader io.Reader) (Configuration, error) {
	config := default_config
	err := json.NewDecoder(reader).Decode(&config)
	return config.Parameters, err
}
