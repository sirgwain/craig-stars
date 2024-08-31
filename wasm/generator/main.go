package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/sirgwain/craig-stars/wasm/generator/generator"
	"golang.org/x/exp/maps"
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
		"BattlePlan":                     true,
		"BattleRecord":                   true,
		"BattleRecordToken":              true,
		"BattleRecordTokenAction":        true,
		"BattleRecordDestroyedToken":     true,
		"BattleRecordStats":              true,
		"BattleVector":                   true,
		"Bomb":                           true,
		"BombingResult":                  true,
		"Cargo":                          true,
		"Cost":                           true,
		"Defense":                        true,
		"DBObject":                       true,
		"Engine":                         true,
		"Fleet":                          true,
		"FleetIntel":                     true,
		"FleetOrders":                    true,
		"FleetSpec":                      true,
		"GameDBObject":                   true,
		"Hab":                            true,
		"Intel":                          true,
		"MapObject":                      true,
		"MapObjectIntel":                 true,
		"MineField":                      true,
		"MineFieldIntel":                 true,
		"MineFieldOrders":                true,
		"MineFieldSpec":                  true,
		"Mineral":                        true,
		"MineralPacketDamage":            true,
		"MineralPacketIntel":             true,
		"MysteryTrader":                  true,
		"MysteryTraderSpec":              true,
		"MysteryTraderIntel":             true,
		"MysteryTraderReward":            true,
		"Planet":                         true,
		"PlanetIntel":                    true,
		"PlanetOrders":                   true,
		"PlanetSpec":                     true,
		"PlanetStarbaseSpec":             true,
		"Player":                         true,
		"PlayerIntel":                    true,
		"PlayerIntels":                   true,
		"PlayerMessage":                  true,
		"PlayerMessageSpec":              true,
		"PlayerMessageSpecComet":         true,
		"PlayerMessageSpecMysteryTrader": true,
		"PlayerOrders":                   true,
		"PlayerPlans":                    true,
		"PlayerRelationship":             true,
		"PlayerScore":                    true,
		"PlayerSpec":                     true,
		"PlayerStats":                    true,
		"ProductionPlan":                 true,
		"ProductionPlanItem":             true,
		"ProductionQueueItem":            true,
		"QueueItemCompletionEstimate":    true,
		"Race":                           true,
		"ResearchCost":                   true,
		"SalvageIntel":                   true,
		"ScoreIntel":                     true,
		"ShipDesign":                     true,
		"ShipDesignIntel":                true,
		"ShipDesignSlot":                 true,
		"ShipDesignSpec":                 true,
		"ShipToken":                      true,
		"TechLevel":                      true,
		"Tech":                           true,
		"TechDefense":                    true,
		"TechPlanetary":                  true,
		"TechPlanetaryScanner":           true,
		"TransportPlan":                  true,
		"Vector":                         true,
		"Waypoint":                       true,
		"WaypointTransportTask":          true,
		"WaypointTransportTasks":         true,
		"WormholeIntel":                  true,
	}

	serializers := []generator.Serializer{}

	// sort all the types we loaded by their names
	keys := maps.Keys(info.Defs)
	slices.SortFunc(keys, func(a, b *ast.Ident) int { return strings.Compare(b.Name, a.Name) })

	// for each type, load a serializer if it has json fields
	for _, key := range keys {
		obj := info.Defs[key]
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
				var fields []generator.Field
				switch t := named.Underlying().(type) {
				case *types.Struct:
					fields = make([]generator.Field, t.NumFields())
					for i := 0; i < t.NumFields(); i++ {
						field := t.Field(i)
						fieldName := field.Name()
						// get the json tag name, whether it is omitted
						// or whether it is ignored
						jsonName, omitEmpty, ignore := getJsonTag(t.Tag(i))

						if !field.Exported() {
							fields[i] = generator.Field{
								Name:   fieldName,
								Ignore: true,
							}
							continue
						}
						fieldType := getTypeInfo(field.Type(), pkg)
						ignore = !field.Exported() || ignore || (!fieldType.Type.IsBasic() && fieldType.Type != generator.GeneratorTypeNamed && fieldType.Package && !typesToCheck[fieldType.TypeName])

						fields[i] = generator.Field{
							FieldType: *fieldType,
							Name:      fieldName,
							JsonName:  jsonName,
							OmitEmpty: omitEmpty,
							Ignore:    ignore,
							Exported:  field.Exported(),
						}
					}
				}

				serializers = append(serializers, generator.Serializer{
					Name:   tn.Name(),
					Fields: fields,
				})

			}
		}
	}

	for _, s := range serializers {
		fmt.Printf("  %s:\n", s.Name)
		for _, f := range s.Fields {
			fmt.Printf("    %s: json: %s, type: %v\n", f.Name, f.JsonName, f.Type)
		}
	}

	out, err := generator.RenderSerializer(pkg.Name, serializers)
	if err != nil {
		log.Fatal(err)
	}

	// format it so it looks nice
	formattedSource, err := formatGoSource(out)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 2 {
		outfile := args[1]
		f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.Write(formattedSource)
		if err != nil {
			panic(err)
		}
	}
}

// formatGoSource formats the go source we generate to make it look pretty
func formatGoSource(source string) ([]byte, error) {
	// Convert the source string to a byte slice
	sourceBytes := []byte(source)

	// Format the source code
	formattedBytes, err := format.Source(sourceBytes)
	if err != nil {
		return nil, err
	}

	// Convert the formatted bytes back to a string
	return formattedBytes, nil
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

// getTypeInfo returns a generator.FieldType from a go/types field Variable
func getTypeInfo(fieldType types.Type, pkg *packages.Package) *generator.FieldType {
	fullType := fieldType.String()

	// is this type defined in the package we're scanning?
	isPackageType := strings.Contains(fullType, pkg.ID)

	// get the actual go type used in  generation, like int or []cs.Planet
	goType := getGoType(fieldType, pkg)
	typeName := goType

	var underlyingType *generator.FieldType
	var keyType *generator.FieldType
	var valueType *generator.FieldType
	var isPointer bool
	var arrayLength int64

	underlyingName := fieldType.Underlying().String()
	isStruct := strings.Index(underlyingName, "struct") == 0

	// if the underlying type differs, fill it in
	if !isStruct && fieldType.Underlying().String() != fullType {
		underlyingType = getTypeInfo(fieldType.Underlying(), pkg)
	}

	var generatorType generator.GeneratorType

	switch t := fieldType.(type) {
	case *types.Basic:
		generatorType = generator.GeneratorTypeFromBasicType(goType)
	case *types.Pointer:
		isPointer = true
		valueType = getTypeInfo(t.Elem(), pkg)
		typeName = valueType.TypeName
		generatorType = valueType.Type
	case *types.Named:
		if isPackageType {
			if isStruct {
				generatorType = generator.GeneratorTypeObject
				typeName = t.Obj().Name()
			} else {
				generatorType = generator.GeneratorTypeNamed
			}
		}
	case *types.Map:
		generatorType = generator.GeneratorTypeMap
		keyType = getTypeInfo(t.Key(), pkg)
		valueType = getTypeInfo(t.Elem(), pkg)
	case *types.Array:
		generatorType = generator.GeneratorTypeArray
		valueType = getTypeInfo(t.Elem(), pkg)
		typeName = valueType.TypeName
		arrayLength = t.Len()
	case *types.Slice:
		generatorType = generator.GeneratorTypeSlice
		valueType = getTypeInfo(t.Elem(), pkg)
		typeName = valueType.TypeName
	}

	return &generator.FieldType{
		TypeName:       typeName,
		FullType:       fullType,
		GoType:         goType,
		Type:           generatorType,
		Package:        isPackageType,
		Pointer:        isPointer,
		UnderlyingType: underlyingType,
		KeyType:        keyType,
		ValueType:      valueType,
		ArrayLength:    arrayLength,
	}

}

func getGoType(t types.Type, pkg *packages.Package) string {
	fullType := t.String()
	pkgPrefix := ""
	if strings.Contains(fullType, pkg.ID) {
		pkgPrefix = pkg.Name + "."
	}

	switch v := t.(type) {
	case *types.Basic:
		return pkgPrefix + v.Name()
	case *types.Slice:
		return fmt.Sprintf("[]%s", getGoType(v.Elem(), pkg))
	case *types.Array:
		return fmt.Sprintf("[%d]%s", v.Len(), getGoType(v.Elem(), pkg))
	case *types.Map:
		return fmt.Sprintf("map[%s]%s", getGoType(v.Key(), pkg), getGoType(v.Elem(), pkg))
	case *types.Pointer:
		return fmt.Sprintf("*%s", getGoType(v.Elem(), pkg))
	case *types.Struct:
		return "struct"
	case *types.Named:
		return pkgPrefix + v.Obj().Name()
	default:
		log.Fatalf("unknown type %#v", v)
	}
	return ""
}

// get a jsonName and omitEmpty from a tag
func getJsonTag(tag string) (name string, omitEmpty bool, ignore bool) {
	if strings.Contains(tag, "json:\"-\"") {
		return "", false, true
	}
	regex := regexp.MustCompile("json:\"(?P<jsonName>[a-zA-Z0-9]+)(?P<omitEmpty>.*)\"")
	matches := regex.FindStringSubmatch(tag)
	if len(matches) == 3 {
		name = matches[1]
		omitEmpty = matches[2] == ",omitempty"
	}

	return name, omitEmpty, false
}
