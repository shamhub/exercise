package httpservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"

	"github.com/itchyny/gojq"
	"github.com/shamhub/exercise/pkg/errorlib"
)

type TemplateData struct {
	TemplateHandle *template.Template
	Data           interface{}
	TemplateName   string
}

type TemplateHandler func(*RequestContextForTemplate) (interface{}, error)

func (h TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	matchedRule, httpErr := validateRequest(r)
	if httpErr != nil {
		processError(w, r.URL.Path, httpErr)
		return
	}

	resourcePath := r.URL.Path

	injector := newTemplateContextInjector()

	injector.injectRequestContextWithTemplate(r, matchedRule.TemplateFilePath)

	data, err := h(injector)
	if err != nil {
		processError(w, resourcePath, err)
		return
	}

	processDataWithTemplates(w, resourcePath, data)
}

func validateRequest(r *http.Request) (*CompiledRuleEntry, errorlib.HttpResponseError) {

	// 1. Identify matching rule based on Method and Path Regex Pattern
	activeRules := GetActiveRules()
	if len(activeRules) == 0 {
		panic("rule config is missing and not loaded")
	}
	var matchedRule *CompiledRuleEntry
	for i := range activeRules {
		rule := activeRules[i]
		if (r.Method == rule.Method) && rule.Regex.MatchString(r.URL.Path) {
			matchedRule = &rule
			break
		}
	}

	fmt.Println("matchedrule: ", matchedRule) // map[userId:johndoe] for "/api/v1/user/(?P<userId>[^/]+)"

	// If no rule exists for this path and method, return error response
	if matchedRule == nil {
		return nil, errorlib.NewResponseError(http.StatusBadRequest, "request route is not valid")
	}

	// 2. Evaluate 'route' Filter Category (expects parameters context)
	fmt.Println("Evaluating route filter")
	if filter, ok := matchedRule.Filters["route"]; ok {
		params := extractNamedMatches(matchedRule.Regex, r.URL.Path)
		fmt.Println("extractNamedMatches for route params:", params)
		routeCtx := map[string]any{"params": params}
		if err := executeFilter(filter, routeCtx); err != nil {
			return nil, errorlib.NewResponseError(http.StatusBadRequest, "route validation failed: "+err.Error())
		}
	}

	// 3. Evaluate 'query_params' Filter Category
	fmt.Println("Evaluating query_params filter")
	if filter, ok := matchedRule.Filters["query_params"]; ok {
		queryCtx := make(map[string]any)
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				queryCtx[key] = values[0] // Simplify multi-values to standard string for basic match
			}
		}

		if err := executeFilter(filter, queryCtx); err != nil {
			return nil, errorlib.NewResponseError(http.StatusBadRequest, "query parameters validation failed: "+err.Error())
		}
	}

	// 4. Evaluate 'headers' Filter Category
	fmt.Println("Evaluating headers filter")
	if filter, ok := matchedRule.Filters["headers"]; ok {
		headerCtx := make(map[string]any)
		headersMap := make(map[string]any)
		for key, values := range r.Header {
			if len(values) > 0 {
				headersMap[key] = values[0]
			}
		}
		headerCtx["headers"] = headersMap

		if err := executeFilter(filter, headerCtx); err != nil {
			return nil, errorlib.NewResponseError(http.StatusBadRequest, "header validation failed: "+err.Error())
		}
	}

	// 5. Evaluate 'payload' Filter Category
	fmt.Println("Evaluating payload filter")
	if filter, ok := matchedRule.Filters["payload"]; ok {
		if r.Body == nil || r.Body == http.NoBody {
			return nil, errorlib.NewResponseError(http.StatusBadRequest, "payload validation failed: request body is empty")
		}

		// Read and preserve body buffer
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, errorlib.NewResponseError(http.StatusInternalServerError, "failed to read request body: "+err.Error())
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var payloadCtx any
		if err := json.Unmarshal(bodyBytes, &payloadCtx); err != nil {
			return nil, errorlib.NewResponseError(http.StatusBadRequest, "payload validation failed: body is invalid JSON: "+err.Error())
		}

		if err := executeFilter(filter, payloadCtx); err != nil {
			return nil, errorlib.NewResponseError(http.StatusBadRequest, "payload validation failed: "+err.Error())
		}
	}

	return matchedRule, nil
}

// ExtractNamedMatches converts path regex capture groups into a key-value map
func extractNamedMatches(re *regexp.Regexp, path string) map[string]any {
	matches := re.FindStringSubmatch(path)
	if len(matches) == 0 {
		return nil
	}

	result := make(map[string]any)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}
	return result
}

// executeFilter runs the gojq pipeline and ensures it resolves explicitly to a true assertion
func executeFilter(query *gojq.Query, input any) error {
	code, err := gojq.Compile(query)
	if err != nil {
		return errorlib.NewResponseError(400, "compile error: "+err.Error())
	}

	iter := code.Run(input)
	v, ok := iter.Next()
	if !ok {
		return fmt.Errorf("filter evaluated to empty sequence")
	}

	if err, isErr := v.(error); isErr {
		return fmt.Errorf("execution error: %w", err)
	}

	booleanResult, isBool := v.(bool)
	if !isBool {
		return fmt.Errorf("filter rules must evaluate to a boolean expression, got %T", v)
	}

	if !booleanResult {
		return fmt.Errorf("rule constraint rejected incoming input data")
	}

	return nil
}

func processDataWithTemplates(w http.ResponseWriter, resourcePath string, data interface{}) {
	fmt.Println("processing response data to render template")

	switch v := data.(type) {
	case TemplateData:
		err := v.TemplateHandle.ExecuteTemplate(w, v.TemplateName, v.Data)
		if err != nil {
			processError(w, resourcePath, err)
		}
	case *TemplateData:
		err := v.TemplateHandle.ExecuteTemplate(w, v.TemplateName, v.Data)
		if err != nil {
			processError(w, resourcePath, err)
		}
	default:
		processData(w, data)
	}
}
