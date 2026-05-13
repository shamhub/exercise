package httpservice

import (
	"context"
	"errors"
	"net/http"
)

type templateFilePathType string

type RequestContextForTemplate struct {
	*RequestContext
	templateFilePathKey templateFilePathType
}

func newTemplateContextInjector() *RequestContextForTemplate {
	return &RequestContextForTemplate{
		RequestContext:      newContextInjector(),
		templateFilePathKey: "filepathKey",
	}
}

func (c *RequestContextForTemplate) injectRequestContextWithTemplate(r *http.Request, templatePath string) {
	if c == nil {
		panic("initialize request context")
	}
	c.injectQueryParams(r)
	c.injectRouteVar(r)
	c.injectRequestPayLoad(r)
	c.injectTemplatePath(templatePath)
}

func (c *RequestContextForTemplate) injectTemplatePath(templatePath string) {
	c.ctx = context.WithValue(c.ctx, c.templateFilePathKey, templatePath)
}

func (c *RequestContextForTemplate) GetTemplatePath() (string, error) {
	if c == nil {
		panic("initialize request context")
	}
	val, ok := c.ctx.Value(c.templateFilePathKey).(string)
	if !ok {
		return "", errors.New("route vars not found in context")
	}
	return val, nil
}
