package generator

import "fmt"

type JSType int

const (
	JSString JSType = iota
	JSBool
	JSInt
	JSFloat
	JSObject
	JSArray
	JSTime
)

type Serializer struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name        string
	JsonName    string
	JSType      JSType
	ObjectType  string
	OmitEmpty   bool
	BasicType   string
	IsBasicType bool
	PackageType bool
	Ignore      bool
}

func (t JSType) String() string {
	switch t {
	case JSString:
		return "String"
	case JSBool:
		return "Bool"
	case JSInt:
		return "Int"
	case JSFloat:
		return "Float"
	case JSObject:
		return "Object"
	case JSArray:
		return "Array"
	case JSTime:
		return "Time"
	}

	panic(fmt.Sprintf("unknown jstype %d", t))
}

func JSTypeFromBasicType(basic string) JSType {
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
		return JSInt
	case "float":
		fallthrough
	case "float32":
		fallthrough
	case "float64":
		return JSFloat
	case "bool":
		return JSBool
	default:
		return JSString

	}
}
