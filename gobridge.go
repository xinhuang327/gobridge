package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"strings"
	"unicode"

	"github.com/xinhuang327/pongo2"
)

var fset = token.NewFileSet()

type GoAstVisitor struct {
	TypeName string
	Funcs    []*ast.FuncDecl
}

func (v *GoAstVisitor) GetSimpleTypeName() string {
	return strings.TrimRight(v.TypeName, "Controller")
}

func (v *GoAstVisitor) GetStructParamName(funcD *ast.FuncDecl) string {
	for _, para := range funcD.Type.Params.List {
		if v.IsStructParam(para) {
			return para.Names[0].Name
		}
	}
	return "null"
}

func (v *GoAstVisitor) GetPrimitiveParams(funcD *ast.FuncDecl) string {
	var primitiveParamNames []string
	for _, para := range funcD.Type.Params.List {
		if !v.IsStructParam(para) {
			query := para.Names[0].Name
			primitiveParamNames = append(primitiveParamNames, query)
		}
	}
	if len(primitiveParamNames) > 0 {
		return "'?'+" + strings.Join(primitiveParamNames, "+'&'+")
	} else {
		return "''"
	}
}

// Simple check if the field type is struct or primitive types: int, float64, string...
func (v *GoAstVisitor) IsStructParam(para *ast.Field) bool {
	if ident, ok := para.Type.(*ast.Ident); ok {
		if unicode.IsUpper([]rune(ident.Name)[0]) { // if the type name starts with upper case letter, it shouldn't be primitive type
			return true
		}
	}
	if _, ok := para.Type.(*ast.SelectorExpr); ok { // if the type name is a selector, we assume it's struct type
		return true
	}
	return false
}

func (v *GoAstVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		fmt.Println(reflect.TypeOf(node))
		if typeSpec, ok := node.(*ast.TypeSpec); ok {
			ast.Print(fset, typeSpec)
			v.TypeName = typeSpec.Name.Name
		}
		if funcSpec, ok := node.(*ast.FuncDecl); ok {
			ast.Print(fset, funcSpec.Recv)
			ast.Print(fset, funcSpec.Type)
			v.Funcs = append(v.Funcs, funcSpec)
			// funcSpec.Type.Params.List[0].Names[0].Name
		}
	}
	return v
}

func (v *GoAstVisitor) Render(w *os.File) error {
	tmpl, err := pongo2.FromCache("jsTemplate.js.bridge")
	if err != nil {
		return err
	}

	return tmpl.ExecuteWriter(pongo2.Context{
		"visitor": v,
	}, w)
}

func (v *GoAstVisitor) RenderFile(outFilePath string) {
	outFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("os.OpenFile() error:", err)
	}
	err = v.Render(outFile)
	if err != nil {
		fmt.Println("visitor.Render() error:", err)
	}
	err = outFile.Close()
	if err != nil {
		fmt.Println("outFile.Close() error:", err)
	}
}

func main() {
	pkgPath := "github.com/xinhuang327/web/cms/ctrls"
	srcFileName := "widget_api.go"

	pkg, err := findSourcePackage(pkgPath)
	if err != nil {
		fmt.Println("findSourcePackage error:", err)
	}
	srcFilePath := getGoFilePath(pkg, srcFileName)
	fmt.Println(srcFilePath)

	node, err := parser.ParseFile(fset, srcFilePath, nil, parser.DeclarationErrors)
	if err != nil {
		fmt.Println("parser.ParseFile error:", err)
	}
	// ast.Print(fset, node)

	visitor := &GoAstVisitor{}
	ast.Walk(visitor, node)

	fmt.Printf("%#v\n", visitor)

	outPath := "github.com/xinhuang327/web/cms/web/public/static/js/svc"
	outPathPkg, err := findSourcePackage(outPath)
	if err != nil {
		fmt.Println("findSourcePackage error:", err)
	}
	outFilePath := path.Join(outPathPkg.Dir, visitor.GetSimpleTypeName()+".js")
	visitor.RenderFile(outFilePath)
	visitor.RenderFile("output.js")
}

func findSourcePackage(pkgPath string) (*build.Package, error) {
	srcDir, err := os.Getwd()
	if err != nil {
		srcDir = ""
	}
	bp, err := build.Default.Import(pkgPath, srcDir, 0)
	if _, ok := err.(*build.NoGoError); ok {
		return bp, nil // empty directory is not an error
	}
	return bp, err
}

func getGoFilePath(pkg *build.Package, goFileName string) string {
	for _, file := range pkg.GoFiles {
		if file == goFileName {
			return path.Join(pkg.Dir, file)
		}
	}
	return ""
}
