package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/shamhub/exercise/pkg/errorlib"
	"github.com/shamhub/exercise/pkg/httpservice"
)

func PostUserId(ctx *httpservice.RequestContextForTemplate) (interface{}, error) {

	queryParams, _ := ctx.GetQueryParams()
	qMap := make(map[string][]string)
	for k, v := range queryParams {
		if len(v) > 0 {
			qMap[k] = v
		}
	}

	// 3. Prepare Payload for gojq
	var pMap map[string]any
	json.NewDecoder(ctx.GetRequestPayload()).Decode(&pMap)

	// 5. Render Go Template
	templateFilePath, _ := ctx.GetTemplatePath()
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return nil, errorlib.NewResponseError(http.StatusInternalServerError, "template file not found")
	}
	return httpservice.TemplateData{
		TemplateHandle: tmpl,
		Data:           &pMap,
	}, nil
}
