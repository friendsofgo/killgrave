package templating

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type TemplatingData struct {
	RequestBody map[string]interface{}
	PathParams  map[string]string
	QueryParams map[string][]string
}

func ApplyTemplate(bodyStr string, templData TemplatingData) ([]byte, error) {
	tmpl, err := template.New("body").
		Funcs(template.FuncMap{
			"stringsJoin": strings.Join,
			"jsonMarshal": JsonMarshal,
			"timeNow":     TimeNow,
			"timeUTC":     TimeUTC,
			"timeAdd":     TimeAdd,
			"timeFormat":  TimeFormat,
		}).
		Parse(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, templData)
	if err != nil {
		return nil, fmt.Errorf("error applying template: %w", err)
	}

	return tpl.Bytes(), nil
}
