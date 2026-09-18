package main

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type configFile struct {
	Symlinks []symLinkFile
	Files    []writeToFile
}

func parseConfigFile(file string) ([]appearanceChangeAction, error) {
	content, err := os.ReadFile(file) //nolint:gosec // File being read came from command line argument
	if err != nil {
		return nil, err
	}

	var parsed configFile
	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return nil, err
	}

	var actions []appearanceChangeAction

	for _, a := range parsed.Files {
		actions = append(actions, appearanceChangeAction(a))
	}

	for _, a := range parsed.Symlinks {
		actions = append(actions, appearanceChangeAction(a))
	}

	return actions, nil
}
