package httpservice

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
)

type ResponseDetail struct {
	Code         int
	MessageQuery string
}

// RuleCOnfig matches the struct of each entry in route json file
type RuleConfig struct {
	Method    string                    `json:"method"`
	Filters   map[string]string         `json:"filters"`
	Responses map[string]ResponseDetail `json:"responses"`
}

// CompiledRule holds the pre-compiled assets for a route
type CompiledRuleEntry struct {
	PathPattern string
	Regex       *regexp.Regexp
	Method      string
	// Filters: category -> compiled JQ query
	Filters map[string]*gojq.Query
	// Responses: category -> code and message
	Responses map[string]ResponseDetail
}

type RuleCollector struct {
	activeRules []CompiledRuleEntry
}

func NewRuleColector() *RuleCollector {
	return &RuleCollector{
		activeRules: make([]CompiledRuleEntry, 0),
	}
}

func (rC *RuleCollector) GetActiveRules() []CompiledRuleEntry {
	return rC.activeRules
}

func (rC *RuleCollector) CollectRules(configData []byte) {
	// 1. unmarshal into a temporary map
	var rawRules map[string]RuleConfig
	if err := json.Unmarshal(configData, &rawRules); err != nil {
		panic(fmt.Sprintf("failed to parse config JSON: %v", err))
	}

	// 2. Compile assets and store in the activeRules slice
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

		rC.activeRules = append(rC.activeRules, CompiledRuleEntry{
			PathPattern: path,
			Regex:       re,
			Method:      ruleConfig.Method,
			Filters:     compiledFilters,
			Responses:   ruleConfig.Responses,
		})
	}

	fmt.Sprintf("Successfully loaded %d validation rules from config json\n", len(rC.activeRules))
}
