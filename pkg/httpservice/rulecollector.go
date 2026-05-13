package httpservice

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/itchyny/gojq"
	"github.com/shamhub/exercise/pkg/config"
)

var activeRules []CompiledRuleEntry = make([]CompiledRuleEntry, 0)

func init() {

	// 1. Read environment config
	configData := config.ReadRuleConfig()

	if configData == nil {
		panic("configuration file mising")
	}

	// 2. Collect activeRules from rules json file
	activeRules = CollectRules(configData)
}

func GetActiveRules() []CompiledRuleEntry {
	if len(activeRules) == 0 {
		panic("rule config file is missing")
	}
	return activeRules
}

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
	Regex       *regexp.Regexp // URL regex
	Method      string
	// Filters: category -> compiled JQ query
	Filters map[string]*gojq.Query
	// Responses: category -> code and message
	TemplateFilePath string
	Transform        *gojq.Query
}

func CollectRules(configData []byte) []CompiledRuleEntry {
	// 1. unmarshal into a temporary map
	var rawRules map[string]RuleConfig
	if err := json.Unmarshal(configData, &rawRules); err != nil {
		panic(fmt.Sprintf("failed to parse config JSON: %v", err))
	}

	var activeRules = make([]CompiledRuleEntry, 0)

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

		activeRules = append(activeRules, CompiledRuleEntry{
			PathPattern:      path,
			Regex:            re,
			Method:           ruleConfigEntry.Method,
			Filters:          compiledFilters,
			TemplateFilePath: templateFilePath,
			Transform:        transformQueryTree,
		})
	}

	fmt.Printf("Successfully loaded %d validation rules from config json\n", len(activeRules))
	return activeRules

}
