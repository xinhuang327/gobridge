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

	"github.com/xinhuang327/pongo2"
)

var fset = token.NewFileSet()

type GoAstVisitor struct {
	TypeName string
	Funcs    []*ast.FuncDecl
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
		"ast": v,
	}, w)
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

	outFile, err := os.OpenFile("output.js", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("os.OpenFile() error:", err)
	}
	err = visitor.Render(outFile)
	if err != nil {
		fmt.Println("visitor.Render() error:", err)
	}
	err = outFile.Close()
	if err != nil {
		fmt.Println("outFile.Close() error:", err)
	}
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
