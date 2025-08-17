package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
	"strings"
	"text/template"

	"path/filepath"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"
)

type EnumConfig struct {
	OutputPackage      string `yaml:"output_package"`
	OutputFile         string `yaml:"output_file"`
	SourcePackage      string `yaml:"source_package"`
	SourcePackageAlias string `yaml:"source_package_alias"`
	ProtoPackage       string `yaml:"proto_package"`
	ProtoPackageAlias  string `yaml:"proto_package_alias"`
	Enums              []struct {
		Name string `yaml:"name"`
	} `yaml:"enums"`
}

type EnumValue struct {
	Name         string
	Value        string
	OriginalType string
}

type EnumData struct {
	Name   string
	Values []EnumValue
}

type TemplateEnum struct {
	Name      string
	ProtoName string
	Values    []EnumValuePair
}

type EnumValuePair struct {
	Source EnumValue
	Proto  EnumValue
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s enums.yaml", os.Args[0])
	}

	// load typenames from yaml
	yamlFile := os.Args[1]
	cfg := loadYAML(yamlFile)

	var typeNames []string
	for _, e := range cfg.Enums {
		typeNames = append(typeNames, e.Name)
	}

	// load enums from the source and proto packages
	sourceEnums := loadSourceEnumConstants(cfg.SourcePackage, typeNames)
	protoEnums := loadProtoEnumConstants(cfg.ProtoPackage, typeNames)

	// match each enum type with each other
	var enums []TemplateEnum
	for _, typeName := range typeNames {
		sourceEnum := sourceEnums[strings.ToLower(typeName)]
		protoEnum := protoEnums[strings.ToLower(typeName)]

		var pairs []EnumValuePair
		protoIndex := 0
		for _, sourceValue := range sourceEnum.Values {
			// if we have a zero value, assign it to the UNSPECIFIED proto enum (the first one)
			isZero := (sourceValue.OriginalType == "string" && sourceValue.Value == `""`) || (sourceValue.OriginalType != "string" && sourceValue.Value == "0")
			if isZero {
				pairs = append(pairs, EnumValuePair{
					Source: sourceValue,
					Proto:  protoEnum.Values[0],
				})
				continue
			}

			if protoIndex < len(protoEnum.Values) && strings.HasSuffix(protoEnum.Values[protoIndex].Name, "_UNSPECIFIED") {
				protoIndex++
			}

			if protoIndex < len(protoEnum.Values) {
				pairs = append(pairs, EnumValuePair{
					Source: sourceValue,
					Proto:  protoEnum.Values[protoIndex],
				})
				protoIndex++
			}
		}
		enums = append(enums, TemplateEnum{
			Name:      sourceEnum.Name,
			ProtoName: protoEnum.Name,
			Values:    pairs,
		})
	}

	// render the template
	data := map[string]interface{}{
		"OutputPackage":      cfg.OutputPackage,
		"SourcePackage":      cfg.SourcePackage,
		"SourcePackageAlias": cfg.SourcePackageAlias,
		"ProtoPackage":       cfg.ProtoPackage,
		"ProtoPackageAlias":  cfg.ProtoPackageAlias,
		"Enums":              enums,
	}
	tmpl := template.Must(template.New("mapper").Funcs(sprig.TxtFuncMap()).Parse(mapperTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal(err)
	}

	// format the go code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("warn: error formatting generated code: %v\n%s", err, buf.String())
	}

	// save the output
	if cfg.OutputFile != "" {
		// ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(cfg.OutputFile), 0755); err != nil {
			log.Fatalf("failed to create output directory for %s: %v", cfg.OutputFile, err)
		}
		err = os.WriteFile(cfg.OutputFile, formatted, 0644)
		if err != nil {
			log.Fatalf("failed to write to output file %s: %v", cfg.OutputFile, err)
		}
	} else {
		os.Stdout.Write(formatted)
	}
}

func loadYAML(path string) EnumConfig {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg EnumConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func loadSourceEnumConstants(pkgPath string, typeNames []string) map[string]EnumData {
	return loadEnumConstants(pkgPath, typeNames, func(c *types.Const) string {
		return c.Name()
	})
}

func loadProtoEnumConstants(pkgPath string, typeNames []string) map[string]EnumData {
	return loadEnumConstants(pkgPath, typeNames, func(c *types.Const) string {
		return strings.TrimPrefix(c.Name(), c.Type().(*types.Named).Obj().Name()+"_")
	})
}

// load all enum constants from a package
// takes a load function to extract the name of the constant to use when mapping source to proto enums
func loadEnumConstants(pkgPath string, typeNames []string, nameFromConst func(*types.Const) string) map[string]EnumData {
	cfg := &packages.Config{Mode: packages.LoadSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		log.Fatalf("failed to load package %s: %v", pkgPath, err)
	}
	if len(pkgs) == 0 {
		log.Fatalf("no packages found for %s", pkgPath)
	}

	typeNamesSet := make(map[string]string, len(typeNames))
	for _, name := range typeNames {
		typeNamesSet[strings.ToLower(name)] = name
	}

	constants := make(map[string][]*types.Const)
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				decl, ok := n.(*ast.GenDecl)
				if !ok || decl.Tok != token.CONST {
					return true
				}

				for _, spec := range decl.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}

					for _, name := range valueSpec.Names {
						obj := pkg.TypesInfo.Defs[name]
						if obj == nil {
							continue
						}

						c, ok := obj.(*types.Const)
						if !ok {
							continue
						}

						named, ok := c.Type().(*types.Named)
						if !ok {
							continue
						}

						typeName := named.Obj().Name()
						if _, ok := typeNamesSet[strings.ToLower(typeName)]; !ok {
							continue
						}

						// this constant is one of the enum types we care about
						constants[strings.ToLower(typeName)] = append(constants[strings.ToLower(typeName)], c)
					}
				}
				return false
			})
		}
	}

	values := make(map[string]EnumData)
	for typeName, consts := range constants {
		var enumValues []EnumValue
		for _, c := range consts {
			enumValues = append(enumValues, EnumValue{
				Name:         nameFromConst(c),
				Value:        c.Val().String(),
				OriginalType: c.Type().Underlying().String(),
			})
		}
		values[typeName] = EnumData{
			Name:   consts[0].Type().(*types.Named).Obj().Name(),
			Values: enumValues,
		}
	}

	if len(values) != len(typeNames) {
		var missing []string
		for _, typeName := range typeNames {
			if _, ok := values[strings.ToLower(typeName)]; !ok {
				missing = append(missing, typeName)
			}
		}
		log.Fatalf("could not find all constants, found %d, expected %d. Missing: %s", len(values), len(typeNames), strings.Join(missing, ", "))
	}
	return values
}

const mapperTemplate = `// Code generated by enum-mapper-gen. DO NOT EDIT.

package {{ .OutputPackage }}

import (
	{{ .SourcePackageAlias }} "{{ .SourcePackage }}"
	{{ .ProtoPackageAlias }} "{{ .ProtoPackage }}"
)

{{ range .Enums }}
{{- $protoName := .ProtoName -}}
func {{ .Name }}ToCS{{ .Name }}(m {{ $.ProtoPackageAlias }}.{{ .ProtoName }}) {{ $.SourcePackageAlias }}.{{ .Name }} {
	switch m {
	{{- range .Values }}
	case {{ $.ProtoPackageAlias }}.{{ $protoName }}_{{ .Proto.Name }}:
		return {{ $.SourcePackageAlias }}.{{ .Source.Name }}
	{{- end }}
	default:
		return {{ $.SourcePackageAlias }}.{{ .Name }}({{- if eq (index .Values 0).Source.OriginalType "string" -}}""{{- else -}}0{{- end -}})
	}
}

func CS{{ .Name }}To{{ .Name }}(m {{ $.SourcePackageAlias }}.{{ .Name }}) {{ $.ProtoPackageAlias }}.{{ .ProtoName }} {
	switch m {
	{{- range .Values }}
	case {{ $.SourcePackageAlias }}.{{ .Source.Name }}:
		return {{ $.ProtoPackageAlias }}.{{ $protoName }}_{{ .Proto.Name }}
	{{- end }}
	default:
		return {{ $.ProtoPackageAlias }}.{{ $protoName }}_{{ (index .Values 0).Proto.Name }}
	}
}
{{ end -}}
`
