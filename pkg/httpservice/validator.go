package httpservice

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
)

type TemplateConfig struct {
	TemplateFile string `jaon:"template_file"`
}

// RuleCOnfig matches the struct of each entry in route json file
type RuleConfig struct {
	Method    string            `json:"method"`
	Filters   map[string]string `json:"filters"`
	Responses map[string]string `json:"responses"`
	Transform string            `json:"transform"`
}

// CompiledRule holds the pre-compiled assets for a route
type CompiledRuleEntry struct {
	PathPattern string
	Regex       *regexp.Regexp
	Method      string
	// Filters: category -> compiled JQ query
	Filters map[string]*gojq.Query
	// Responses: category -> code and message
	TemplateFilePath string
	Transform        *gojq.Query
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
	for path, ruleConfigEntry := range rawRules {
		// Compile regex (anchored to start and end for exact matching)
		re, err := regexp.Compile("^" + path + "$")
		if err != nil {
			panic(fmt.Sprintf("Invalid regex in path %s, %v", path, err))
		}

		compiledFilters := make(map[string]*gojq.Query)
		// Compile each JQ filter string for this route
		for category, jqFilterString := range ruleConfigEntry.Filters {
			queryStruct, err := gojq.Parse(jqFilterString)
			if err != nil {
				panic(fmt.Sprintf("invalid JQ filter in %s (%s): %v ", path, jqFilterString, err))
			}
			compiledFilters[category] = queryStruct
		}

		// Collect templatefile path data
		var templateFilePath string
		for _, pathString := range ruleConfigEntry.Responses {
			templateFilePath = pathString
		}

		// Compile Transform Logic
		tranformString := ruleConfigEntry.Transform
		transformQueryTree, err := gojq.Parse(tranformString)
		if err != nil {
			panic(fmt.Sprintf("invalid JQ filter in %s (%s): %v ", path, tranformString, err))
		}

		rC.activeRules = append(rC.activeRules, CompiledRuleEntry{
			PathPattern:      path,
			Regex:            re,
			Method:           ruleConfigEntry.Method,
			Filters:          compiledFilters,
			TemplateFilePath: templateFilePath,
			Transform:        transformQueryTree,
		})
	}

	fmt.Printf("Successfully loaded %d validation rules from config json\n", len(rC.activeRules))
}
