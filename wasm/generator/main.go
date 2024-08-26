package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"os"
	"regexp"
	"strings"

	"github.com/sirgwain/craig-stars/wasm/generator/generator"
	"golang.org/x/tools/go/packages"
)

const cmdUsage = `
Usage : generator [options] <package> <outfile>
Examples:
generate source for serializing js.Value to/from structs in a package`

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println(cmdUsage)
		return
	}

	pkg, err := loadPackage(args[0])
	if err != nil {
		panic(err)
	}

	info, err := getPackageInfo(*pkg)
	if err != nil {
		panic(err)
	}

	// TODO: do them all, eventually
	typesToCheck := map[string]bool{
		"DBObject":     true,
		"Hab":          true,
		"Race":         true,
		"ResearchCost": true,
		"TechLevel":    true,
		"Cost":         true,
	}

	serializers := []generator.Serializer{}

	// Iterate over types in the package
	for _, obj := range info.Defs {
		if obj == nil {
			continue
		}

		if tn, ok := obj.(*types.TypeName); ok {
			if !tn.Exported() {
				continue
			}

			if ok := typesToCheck[tn.Name()]; !ok {
				continue
			}
			if named, ok := tn.Type().(*types.Named); ok {
				hasJsonFields := false
				var fields []generator.Field
				switch t := named.Underlying().(type) {
				case *types.Struct:
					fmt.Printf("Type: %s\n", tn.Name())
					fields = make([]generator.Field, t.NumFields())
					for i := 0; i < t.NumFields(); i++ {
						field := t.Field(i)
						fieldName := field.Name()
						fieldType := field.Type()
						basicType := getBasicType(fieldType)
						underlyingType := getUnderlyingType(fieldType)
						jsonName, omitEmpty := getJsonTag(t.Tag(i))
						if jsonName != "" {
							hasJsonFields = true
						}
						fmt.Printf("  Field: %s, Type: %s, Underlying: %s, Tag: %s\n",
							fieldName,
							strings.ReplaceAll(fieldType.String(), pkg.PkgPath+".", ""),
							underlyingType,
							jsonName,
						)

						packageType := false
						fullFieldType := fieldType.String()
						if strings.Index(fullFieldType, pkg.ID) == 0 {
							// this is a type internal to our package
							packageType = true
						}

						fields[i] = generator.Field{
							Name:        fieldName,
							JsonName:    jsonName,
							ObjectType:  underlyingType,
							OmitEmpty:   omitEmpty,
							JSType:      getJsType(fieldType),
							BasicType:   basicType,
							IsBasicType: basicType != "",
							PackageType: packageType,
							Ignore:      (basicType == "") && packageType && !typesToCheck[underlyingType],
						}
					}
				case *types.Interface:
					// ignore
				case *types.Basic:
					fmt.Printf("Type: %s => %v\n", tn.Name(), named.Underlying())
				}

				if hasJsonFields {
					serializers = append(serializers, generator.Serializer{
						Name:   tn.Name(),
						Fields: fields,
					})
				}
			}
		}
	}

	fmt.Printf("found %d serialzers\n", len(serializers))

	for _, s := range serializers {
		fmt.Printf("  %s:\n", s.Name)
		for _, f := range s.Fields {
			fmt.Printf("    %s: json: %s, jsType: %v\n", f.Name, f.JsonName, f.JSType)
		}
	}

	out, err := generator.RenderSerializer(pkg.Name, serializers)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\nConverter code\n\n")
	fmt.Print(out)

	if len(args) == 2 {
		outfile := args[1]
		f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.Write([]byte(out))
		if err != nil {
			panic(err)
		}

		// err = format.Node(w, fset, f)
		// if err != nil {
		// 	fmt.Printf("Error formating file %s", err)
		// 	return
		// }
		// w.Flush()
	}

}

func loadPackage(path string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedDeps | packages.NeedImports,
		Fset:  token.NewFileSet(),
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %v", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found")
	}

	// Assuming only one package is loaded, but we can handle multiple packages as well.
	pkg := pkgs[0]
	if pkg.Errors != nil {
		for _, e := range pkg.Errors {
			fmt.Println(e)
		}
		return nil, fmt.Errorf("found errors in package")
	}

	return pkg, nil
}

func getPackageInfo(pkg packages.Package) (*types.Info, error) {
	// Create a types.Config to type-check the loaded syntax
	conf := types.Config{Importer: importer.ForCompiler(pkg.Fset, "source", nil)}
	info := &types.Info{
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Types:      make(map[ast.Expr]types.TypeAndValue),
	}

	// Type-check the package
	_, err := conf.Check(pkg.PkgPath, pkg.Fset, pkg.Syntax, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func getUnderlyingType(t types.Type) string {
	switch v := t.(type) {
	case *types.Basic:
		return v.Name()
	case *types.Slice:
		return "[] " + getUnderlyingType(v.Elem())
	case *types.Array:
		return fmt.Sprintf("[%d]%s", v.Len(), getUnderlyingType(v.Elem()))
	case *types.Map:
		return fmt.Sprintf("map[%s]%s", getUnderlyingType(v.Key()), getUnderlyingType(v.Elem()))
	case *types.Pointer:
		return "*" + getUnderlyingType(v.Elem())
	case *types.Struct:
		return "struct"
	case *types.Named:
		return v.Obj().Name()
	default:
		return "unknown"
	}
}

func getBasicType(t types.Type) string {
	t, ok := t.Underlying().(*types.Basic)
	if ok {
		return t.String()
	}
	return ""
}

// get a jsonName and omitEmpty from a tag
func getJsonTag(tag string) (name string, omitEmpty bool) {
	regex := regexp.MustCompile("json:\"(?P<jsonName>[a-zA-Z0-9]+)(?P<omitEmpty>.*)\"")
	matches := regex.FindStringSubmatch(tag)
	if len(matches) == 3 {
		name = matches[1]
		omitEmpty = matches[2] == ",omitempty"
	}

	return name, omitEmpty
}

func getJsType(t types.Type) generator.JSType {
	switch v := t.(type) {
	case *types.Basic:
		return generator.JSTypeFromBasicType(v.Name())
	case *types.Named:
		fullName := v.String()
		switch fullName {
		case "time.Time":
			return generator.JSTime
		default:
			if basic, ok := v.Underlying().(*types.Basic); ok {
				return generator.JSTypeFromBasicType(basic.Name())
			}
			fmt.Printf("found unknown named type %s\n", fullName)
		}
	case *types.Struct:
		return generator.JSObject
	}

	return generator.JSString
}
