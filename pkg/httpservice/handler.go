package httpservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shamhub/exercise/pkg/config"
	"github.com/shamhub/exercise/pkg/errorlib"
)

var ruleCollection *RuleCollector

func init() {
	// 1. Read environment config
	configData := config.ReadConfig()

	// 2. Collect activeRules from rules json file
	ruleCollection = NewRuleColector()
	ruleCollection.CollectRules(configData)
}

type MyHandler func(*RequestContext) (interface{}, error)

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	httpErr := validateRequest(r)
	if httpErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpErr.GetStatusCode())
		json.NewEncoder(w).Encode(map[string][]string{"error": httpErr.ProvideReason()})
		return
	}

	resourcePath := r.URL.Path

	injector := newContextInjector()

	injector.injectRequestContext(r)

	data, err := h(injector)
	if err != nil {
		processError(w, resourcePath, err)
		return
	}

	processData(w, data)
}

func validateRequest(r *http.Request) errorlib.HttpResponseError {

	// 1. Attempt to match the route and extract variables
	matchedRule, params := findMatch(r)

	// If no rule exists for this path, return error response
	if matchedRule == nil {
		return errorlib.NewResponseError(http.StatusNotFound, "Route not found")
	}

	// 2. Execute validation logic using the matchedRule and params
	isValid, errorMsg := validateWithRule(r, matchedRule, params)
	if !isValid {
		return errorlib.NewResponseError(http.StatusBadRequest, errorMsg)
	}
	return nil
}

func findMatch(r *http.Request) (*CompiledRule, map[string]string) {
	params := make(map[string]string)

	var activeRuleCollection []CompiledRule
	if ruleCollection != nil {
		activeRuleCollection = ruleCollection.GetActiveRules()
	}
	for _, rule := range activeRuleCollection {
		// 1. Verify HTTP Method (skip if rule specifies a method and it doesn't match)
		if rule.Method != "" && rule.Method != r.Method {
			continue
		}

		// 2. Check if the URL path matches the regex pattern
		match := rule.Regex.FindStringSubmatch(r.URL.Path)
		if match == nil {
			continue
		}

		// 3. Extract named captures into the params map
		// index 0 is the full match, submatches start at index 1
		groupNames := rule.Regex.SubexpNames()
		for i, value := range match {
			if i > 0 && groupNames[i] != "" {
				params[groupNames[i]] = value
			}
		}

		return &rule, params
	}

	return nil, nil
}

func validateWithRule(r *http.Request, rule *CompiledRule, params map[string]string) (bool, string) {
	// 1. Extract Body (Handle empty bodies gracefully)
	var body interface{}
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return false, "Invalid JSON payload"
		}
	}

	type RequestData struct {
		Params  map[string]string   `json:"params"`
		Query   map[string][]string `json:"query"`
		Headers map[string][]string `json:"headers"`
		Body    interface{}         `json:"body"`
	}

	// 2. Build the JQ input object
	input := RequestData{
		Params:  params,
		Query:   r.URL.Query(),
		Headers: r.Header,
		Body:    body,
	}

	// Convert struct to map for gojq compatibility
	var inputMap map[string]interface{}
	data, _ := json.Marshal(input)
	json.Unmarshal(data, &inputMap)

	// 3. Run each filter category (route, headers, payload, etc.)
	for category, query := range rule.Filters {
		iter := query.Run(inputMap)
		v, ok := iter.Next()

		// If JQ returns an error or anything other than 'true'
		if !ok {
			return false, fmt.Sprintf("Validation failed: %s logic error", category)
		}
		if err, ok := v.(error); ok {
			return false, fmt.Sprintf("%s error: %v", category, err)
		}
		if v != true {
			return false, fmt.Sprintf("Request failed %s validation", category)
		}
	}

	return true, ""
}

func processData(w http.ResponseWriter, data interface{}) {
	fmt.Println("processing data")
	// encodedData, err := json.Marshal(data)
	// if err != nil {
	// 	marshalErrorResponse(w, err)
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	// w.Write(encodedData)
}

func processError(w http.ResponseWriter, resourcePath string, err error) {
	switch err := err.(type) {
	case *errorlib.ResponseError:
		c := newCustomErrorForSingleErrorResponse(resourcePath, err)
		SendErrorResponse(w, c)
	case *errorlib.MultiErrors:
		c := newCustomErrorForMultiErrorResponse(resourcePath, err)
		SendErrorResponse(w, c)
	}
}
