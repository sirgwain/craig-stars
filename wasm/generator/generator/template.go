package generator

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

var funcMap = template.FuncMap{
	"JSString": func() JSType { return JSString },
	"JSBool":   func() JSType { return JSBool },
	"JSInt":    func() JSType { return JSInt },
	"JSFloat":  func() JSType { return JSFloat },
	"JSObject": func() JSType { return JSObject },
	"JSArray":  func() JSType { return JSArray },
	"JSTime":   func() JSType { return JSTime },
}

//go:embed templates/converter.go.tmpl
var converterTemplate embed.FS

func RenderSerializer(pkg string, serializers []Serializer) (string, error) {

	tmpl, err := converterTemplate.ReadFile("templates/converter.go.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to load embedded template %v", err)
	}

	t := template.Must(template.New("serializerTemplate").Funcs(funcMap).Parse(string(tmpl)))

	var out bytes.Buffer
	t.Execute(&out, map[string]interface{}{
		"Pkg":         pkg,
		"Serializers": serializers,
	})

	return out.String(), nil
}
