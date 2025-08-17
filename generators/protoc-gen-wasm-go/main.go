// protoc-gen-wasm-go is a plugin for the Protobuf compiler that generates
// Go code. To use it, build this program and make it available on your PATH as
// protoc-gen-wasm-go.
//
// The 'proto-wasm' suffix becomes part of the arguments for the Protobuf
// compiler. To generate the base Go types and wasm code using protoc:
//
//	protoc --go_out=gen --proto-wasm_out=gen path/to/file.proto
//
// With [buf], your buf.gen.yaml will look like this:
//
//	version: v2
//	plugins:
//	  - local: protoc-gen-go
//	    out: gen
//	  - local: protoc-gen-wasm-go
//	    out: gen
//
// This generates wasm service definitions for the Protobuf types and services
// defined by file.proto. If file.proto defines the foov1 Protobuf package, the
// invocations above will write output to:
//
//	gen/path/to/foov1wasm/file.wasm.go
//
// The generated code is configurable with the same parameters as the protoc-gen-go
// plugin, with the following additional parameters:
//
//   - package_suffix: To generate into a sub-package of the package containing the
//     base .pb.go files using the given suffix. An empty suffix denotes to
//     generate into the same package as the base pb.go files. Default is "wasm".
//
// For example, to generate into the same package as the base .pb.go files:
//
//	version: v2
//	plugins:
//	  - local: protoc-gen-go
//	    out: gen
//	  - local: protoc-gen-wasm-go
//	    out: gen
//	    opt: package_suffix
//
// This will generate output to:
//
//	gen/path/to/file.pb.go
//	gen/path/to/file.wasm.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/sirgwain/craig-stars/proto-wasm/grpcwasm"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	contextPackage  = protogen.GoImportPath("context")
	errorsPackage   = protogen.GoImportPath("errors")
	httpPackage     = protogen.GoImportPath("net/http")
	stringsPackage  = protogen.GoImportPath("strings")
	grpcwasmPackage = protogen.GoImportPath("github.com/sirgwain/craig-stars/proto-wasm/grpcwasm")

	generatedFilenameExtension = ".wasm.go"
	defaultPackageSuffix       = "wasm"
	packageSuffixFlagName      = "package_suffix"

	usage = "Flags:\n  -h, --help\tPrint this help and exit.\n      --version\tPrint the version and exit."

	commentWidth = 97 // leave room for "// "

	// To propagate top-level comments, we need the field number of the syntax
	// declaration and the package name in the file descriptor.
	protoSyntaxFieldNum  = 12
	protoPackageFieldNum = 2
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintln(os.Stdout, grpcwasm.Version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Fprintln(os.Stdout, usage)
		os.Exit(0)
	}
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}
	var flagSet flag.FlagSet
	packageSuffix := flagSet.String(
		packageSuffixFlagName,
		defaultPackageSuffix,
		"Generate files into a sub-package of the package containing the base .pb.go files using the given suffix. An empty suffix denotes to generate into the same package as the base pb.go files.",
	)
	protogen.Options{
		ParamFunc: flagSet.Set,
	}.Run(
		func(plugin *protogen.Plugin) error {
			plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) | uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
			plugin.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO2
			plugin.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023
			for _, file := range plugin.Files {
				if file.Generate {
					generate(plugin, file, *packageSuffix)
				}
			}
			return nil
		},
	)
}

func generate(plugin *protogen.Plugin, file *protogen.File, packageSuffix string) {
	if len(file.Services) == 0 {
		return
	}

	goImportPath := file.GoImportPath
	if packageSuffix != "" {
		if !token.IsIdentifier(packageSuffix) {
			plugin.Error(fmt.Errorf("package_suffix %q is not a valid Go identifier", packageSuffix))
			return
		}
		file.GoPackageName += protogen.GoPackageName(packageSuffix)
		generatedFilenamePrefixToSlash := filepath.ToSlash(file.GeneratedFilenamePrefix)
		file.GeneratedFilenamePrefix = path.Join(
			path.Dir(generatedFilenamePrefixToSlash),
			string(file.GoPackageName),
			path.Base(generatedFilenamePrefixToSlash),
		)
		goImportPath = protogen.GoImportPath(path.Join(
			string(file.GoImportPath),
			string(file.GoPackageName),
		))
	}
	generatedFile := plugin.NewGeneratedFile(
		file.GeneratedFilenamePrefix+generatedFilenameExtension,
		goImportPath,
	)
	if packageSuffix != "" {
		generatedFile.Import(file.GoImportPath)
	}
	generatePreamble(generatedFile, file)
	generateServiceNameConstants(generatedFile, file.Services)
	for _, service := range file.Services {
		generateService(generatedFile, file, service)
	}
}

func generatePreamble(g *protogen.GeneratedFile, file *protogen.File) {
	syntaxPath := protoreflect.SourcePath{protoSyntaxFieldNum}
	syntaxLocation := file.Desc.SourceLocations().ByPath(syntaxPath)
	for _, comment := range syntaxLocation.LeadingDetachedComments {
		leadingComments(g, protogen.Comments(comment), false /* deprecated */)
	}
	g.P()
	leadingComments(g, protogen.Comments(syntaxLocation.LeadingComments), false /* deprecated */)
	g.P()

	programName := filepath.Base(os.Args[0])
	// Remove .exe suffix on Windows so that generated code is stable, regardless
	// of whether it was generated on a Windows machine or not.
	if ext := filepath.Ext(programName); strings.ToLower(ext) == ".exe" {
		programName = strings.TrimSuffix(programName, ext)
	}
	g.P("// Code generated by ", programName, ". DO NOT EDIT.")
	g.P("//")
	if file.Proto.GetOptions().GetDeprecated() {
		wrapComments(g, file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// Source: ", file.Desc.Path())
	}
	g.P()

	pkgPath := protoreflect.SourcePath{protoPackageFieldNum}
	pkgLocation := file.Desc.SourceLocations().ByPath(pkgPath)
	for _, comment := range pkgLocation.LeadingDetachedComments {
		leadingComments(g, protogen.Comments(comment), false /* deprecated */)
	}
	g.P()
	leadingComments(g, protogen.Comments(pkgLocation.LeadingComments), false /* deprecated */)

	g.P("package ", file.GoPackageName)
	g.P()
}

func generateServiceNameConstants(g *protogen.GeneratedFile, services []*protogen.Service) {
	var numMethods int
	g.P("const (")
	for _, service := range services {
		constName := fmt.Sprintf("%sName", service.Desc.Name())
		wrapComments(g, constName, " is the fully-qualified name of the ",
			service.Desc.Name(), " service.")
		g.P(constName, ` = "`, service.Desc.FullName(), `"`)
		numMethods += len(service.Methods)
	}
	g.P(")")
	g.P()

	if numMethods == 0 {
		return
	}
	wrapComments(g, "These constants are the fully-qualified names of the RPCs defined in this package. ",
		"They're exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.")
	g.P("//")
	wrapComments(g, "Note that these are different from the fully-qualified method names used by ",
		"google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to ",
		"reflection-formatted method names, remove the leading slash and convert the ",
		"remaining slash to a period.")
	g.P("const (")
	for _, service := range services {
		for _, method := range service.Methods {
			// The runtime exposes this value as Spec.Procedure, so we should use the
			// same term here.
			wrapComments(g, procedureConstName(method), " is the fully-qualified name of the ",
				service.Desc.Name(), "'s ", method.Desc.Name(), " RPC.")
			g.P(procedureConstName(method), ` = "`, fmt.Sprintf("/%s/%s", service.Desc.FullName(), method.Desc.Name()), `"`)
		}
	}
	g.P(")")
	g.P()
}

func generateService(g *protogen.GeneratedFile, file *protogen.File, service *protogen.Service) {
	names := newNames(service)
	generateServerInterface(g, service, names)
	generateServerConstructor(g, file, service, names)
}

func generateServerInterface(g *protogen.GeneratedFile, service *protogen.Service, names names) {
	wrapComments(g, names.Server, " is the WASM service interface for the ", service.Desc.FullName(), " service.")
	if isDeprecatedService(service) {
		g.P("//")
		deprecated(g)
	}
	g.AnnotateSymbol(names.Server, protogen.Annotation{Location: service.Location})
	g.P("type ", names.Server, " interface {")
	for _, method := range service.Methods {
		leadingComments(
			g,
			method.Comments.Leading,
			isDeprecatedMethod(method),
		)
		g.AnnotateSymbol(names.Server+"."+method.GoName, protogen.Annotation{Location: method.Location})
		g.P(serverSignature(g, method))
	}
	g.P("}")
	g.P()
}

func generateServerConstructor(g *protogen.GeneratedFile, file *protogen.File, service *protogen.Service, names names) {
	wrapComments(g, names.ServerConstructor, " builds a WASM handler from the service implementation.")
	if isDeprecatedService(service) {
		g.P("//")
		deprecated(g)
	}
	g.P("func ", names.ServerConstructor, "(svc ", names.Server, ") ", grpcwasmPackage.Ident("Handler"), " {")
	g.P(`return `, grpcwasmPackage.Ident("HandlerFunc"), `(func(ctx `, contextPackage.Ident("Context"), `, method string, reqBytes []byte) ([]byte, error) {`)
	g.P("switch method {")
	for _, method := range service.Methods {
		g.P("case ", procedureConstName(method), ":")
		g.P("respBytes, err := ", grpcwasmPackage.Ident("HandleUnary"), `(ctx, reqBytes, &`, method.Input.GoIdent, `{}, svc.`, method.GoName, `)`)
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return respBytes, nil")
	}
	g.P("default:")
	g.P("return nil, ", grpcwasmPackage.Ident("ErrMethodNotFound"))
	g.P("}")
	g.P("})")
	g.P("}")
	g.P()
}

func serverSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	return method.GoName + serverSignatureParams(g, method, false /* named */)
}

func serverSignatureParams(g *protogen.GeneratedFile, method *protogen.Method, named bool) string {
	ctxName := "ctx "
	reqName := "req "
	if !named {
		ctxName, reqName = "", ""
	}
	// unary
	return "(" + ctxName + g.QualifiedGoIdent(contextPackage.Ident("Context")) +
		", " + reqName + "*" +
		g.QualifiedGoIdent(method.Input.GoIdent) + ") " +
		"(*" +
		g.QualifiedGoIdent(method.Output.GoIdent) + ", error)"
}

func procedureConstName(m *protogen.Method) string {
	return fmt.Sprintf("%s%sProcedure", m.Parent.GoName, m.GoName)
}

func isDeprecatedService(service *protogen.Service) bool {
	serviceOptions, ok := service.Desc.Options().(*descriptorpb.ServiceOptions)
	return ok && serviceOptions.GetDeprecated()
}

func isDeprecatedMethod(method *protogen.Method) bool {
	methodOptions, ok := method.Desc.Options().(*descriptorpb.MethodOptions)
	return ok && methodOptions.GetDeprecated()
}

// stolen from connect-go
func wrapComments(g *protogen.GeneratedFile, elems ...any) {
	text := &bytes.Buffer{}
	for _, el := range elems {
		switch el := el.(type) {
		case protogen.GoIdent:
			fmt.Fprint(text, g.QualifiedGoIdent(el))
		default:
			fmt.Fprint(text, el)
		}
	}
	words := strings.Fields(text.String())
	text.Reset()
	var pos int
	for _, word := range words {
		numRunes := utf8.RuneCountInString(word)
		if pos > 0 && pos+numRunes+1 > commentWidth {
			g.P("// ", text.String())
			text.Reset()
			pos = 0
		}
		if pos > 0 {
			text.WriteRune(' ')
			pos++
		}
		text.WriteString(word)
		pos += numRunes
	}
	if text.Len() > 0 {
		g.P("// ", text.String())
	}
}

func leadingComments(g *protogen.GeneratedFile, comments protogen.Comments, isDeprecated bool) {
	if comments.String() != "" {
		g.P(strings.TrimSpace(comments.String()))
	}
	if isDeprecated {
		if comments.String() != "" {
			g.P("//")
		}
		deprecated(g)
	}
}

func deprecated(g *protogen.GeneratedFile) {
	g.P("// Deprecated: do not use.")
}

type names struct {
	Base              string
	Server            string
	ServerConstructor string
}

func newNames(service *protogen.Service) names {
	base := service.GoName
	return names{
		Base:              base,
		Server:            fmt.Sprintf("%sHandler", base),
		ServerConstructor: fmt.Sprintf("New%sHandler", base),
	}
}
