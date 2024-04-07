package vrahos

import (
	"bytes"
	"net/http"
	"text/template"
)

type MetaData struct {
	Authenticated bool
	Sse           *Sse
	Template      *template.Template
}

type TemplateData struct {
	Request     *http.Request
	Props       any
	PayloadJson string
	Page        string
	Meta        MetaData
}

type Component interface {
	Name() string
	URL() string
	Template() string
	Functions() *map[string]any
	Props(r *http.Request, meta *MetaData) (any, map[string]string)
}

type ComponentFull interface {
	Component
}

type BasicComponent struct{}

func (p BasicComponent) Template() string {
	return ""
}

func (p BasicComponent) URL() string {
	return ""
}

func (p BasicComponent) Props(r *http.Request, meta *MetaData) (any, map[string]string) {
	return nil, nil
}

func (p BasicComponent) Functions() *map[string]any {
	return nil
}

func RenderComponentToString(tmp *template.Template, component Component, meta *MetaData, props any) (string, error) {

	data := TemplateData{
		Props:   props,
		Page:    "",
		Request: nil,
		Meta:    *meta,
	}

	tmp, err := tmp.Parse(component.Template())

	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	tmp.Execute(&buff, &data)
	return buff.String(), nil
}
