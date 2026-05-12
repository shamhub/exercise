package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
	"github.com/shamhub/exercise/pkg/config"
)

// RuleCOnfig matches the struct of each entry in route json file
type RuleConfig struct {
	Method  string            `json:"method"`
	Filters map[string]string `json:"filters"`
}

// CompiledRule holds the pre-compiled assets for a route
type CompiledRule struct {
	PathPattern string
	Regex       *regexp.Regexp
	Method      string
	Filters     map[string]*gojq.Query
}

var activeRules []CompiledRule

func main() {

	// 1. Read environment config
	configData := config.ReadConfig()

	// 2. unmarshal into a temporary map
	var rawRules map[string]RuleConfig
	if err := json.Unmarshal(configData, &rawRules); err != nil {
		panic(fmt.Sprintf("failed to parse config JSON: %v", err))
	}

	// 3. Compile assets and store in the activeRules slice
	for path, ruleConfig := range rawRules {
		// Compile regex (anchored to start and end for exact matching)
		re, err := regexp.Compile("^" + path + "$")
		if err != nil {
			panic(fmt.Sprintf("Invalid regex in path %s, %v", path, err))
		}

		compiledFilters := make(map[string]*gojq.Query)
		// Compile each JQ filter string for this route
		for category, jqFilterString := range ruleConfig.Filters {
			queryStruct, err := gojq.Parse(jqFilterString)
			if err != nil {
				panic(fmt.Sprintf("invalid JQ filter in %s (%s): %v ", path, jqFilterString, err))
			}
			compiledFilters[category] = queryStruct
		}

		activeRules = append(activeRules, CompiledRule{
			PathPattern: path,
			Regex:       re,
			Method:      ruleConfig.Method,
			Filters:     compiledFilters,
		})
	}

	fmt.Sprintf("Successfully loaded %d validation rules from config json\n", len(activeRules))
}
