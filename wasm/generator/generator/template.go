package generator

import (
	"bytes"
	"embed"
	"log"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed templates/*
var templatesFS embed.FS

func RenderSerializer(pkg string, serializers []Serializer) (string, error) {

	// alphabetize the serializers
	slices.SortFunc(serializers, func(a, b Serializer) int { return strings.Compare(a.Name, b.Name) })

	t := template.New("")
	t = template.Must(t.
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"include": func(name string, data interface{}) string {
				result, err := includeTemplate(t, name, data)
				if err != nil {
					log.Fatalf("failed to include %s %v", name, err)
				}
				return result
			},
		}).
		ParseFS(templatesFS, "templates/*"))

	var out bytes.Buffer
	if err := t.ExecuteTemplate(&out, "converter.go.tmpl", map[string]interface{}{
		"Pkg":         pkg,
		"Serializers": serializers,
	}); err != nil {
		return "", err
	}

	return out.String(), nil
}

// Define a custom function to include templates
func includeTemplate(tmpl *template.Template, name string, data interface{}) (string, error) {
	var result strings.Builder
	err := tmpl.ExecuteTemplate(&result, name, data)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
