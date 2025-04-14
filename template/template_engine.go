package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

func buildFuncMap() template.FuncMap {
	result := template.FuncMap{}
	result["env"] = os.Getenv

	return result
}

func EvaluateTemplate(text string, vars any) (string, error) {
	templateObj := template.New("template")

	templateObj = templateObj.Funcs(buildFuncMap())

	templateObj, err := templateObj.Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	output := bytes.NewBufferString("")
	err = templateObj.Execute(output, vars)
	if err != nil {
		return "", fmt.Errorf("failed to execute template using vars snapshot: %w", err)
	}
	return output.String(), nil
}
