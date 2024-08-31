package generator

type GeneratorType string

const (
	GeneratorTypeNone = ""
	// a basic type like an int
	GeneratorTypeBasicBool   = "BasicBool"
	GeneratorTypeBasicFloat  = "BasicFloat"
	GeneratorTypeBasicInt    = "BasicInt"
	GeneratorTypeBasicString = "BasicString"
	// an object like a MapObject or Planet
	GeneratorTypeObject = "Object"
	// a named type like an enum, with an underlying type
	GeneratorTypeNamed = "Named"

	// collection types
	GeneratorTypeMap   = "Map"
	GeneratorTypeArray = "Array"
	GeneratorTypeSlice = "Slice"
)

func (t GeneratorType) IsBasic() bool {
	return t == GeneratorTypeBasicBool || t == GeneratorTypeBasicFloat || t == GeneratorTypeBasicInt || t == GeneratorTypeBasicString
}

type Serializer struct {
	Name   string
	Fields []Field
}

// A Field represents a field in a struct we are
// serializing/deserializing
// It can represent
// * A basic type, like an int -- obj.Ironium = GetInt[int](o, "ironium")
// * An object type, like a Hab -- obj.Hab = GetHab(o.Get("hab"))
// * A slice of objects -- obj.Relations = GetPlayerRelationshipArray(o.Get("relations"))
// * A named type, like an enum -- obj.PrimaryTarget = cs.BattleTarget(GetString(o, "primaryTarget"))
// * A named type pointing to a map -- obj.Tags = cs.Tags{}, then populate map
// * A map with a basic key and value type that could be one of the above
// * An array with a length - obj.FuelUsage = [11]int{}
type Field struct {
	FieldType
	Name      string
	JsonName  string
	OmitEmpty bool
	Ignore    bool
}

// TODO: make a TypeType with struct, map, basic, named (type name mapping to a basic or map type)
// Named types should use the ValueType as their underlying type, or myabe an underlying type pointer?
type FieldType struct {
	TypeName string
	FullType string
	GoType   string
	Type     GeneratorType
	Package  bool
	Pointer  bool
	// the underlying type, if different
	UnderlyingType *FieldType
	// the value type for maps, slices, or arrays
	ValueType *FieldType
	// the key type for maps
	KeyType     *FieldType
	ArrayLength int64
}

func (t FieldType) IsBasic() bool {
	return t.Type.IsBasic()
}

func GeneratorTypeFromBasicType(basic string) GeneratorType {
	switch basic {
	case "int":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		fallthrough
	case "uint":
		fallthrough
	case "uint32":
		fallthrough
	case "uint64":
		return GeneratorTypeBasicInt
	case "float":
		fallthrough
	case "float32":
		fallthrough
	case "float64":
		return GeneratorTypeBasicFloat
	case "bool":
		return GeneratorTypeBasicBool
	case "string":
		return GeneratorTypeBasicString
	}

	return GeneratorTypeNone
}
