package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var serviceName string
var dashBoardTitle string
var namespace string

const structName = "metricConfig"

func main() {
	flag.StringVar(&serviceName, "service_name", "", "name of your service")
	flag.StringVar(&dashBoardTitle, "title", "", "title of your dashboard")
	flag.StringVar(&namespace, "namespace", "", "namespace")
	flag.Parse()
	if serviceName == "" {
		panic("service name can not be nil.Usage: pd_gen -service_name api")
	}
	args := flag.Args()

	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}
	dir := args[0]

	pkg, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", dir, err)
	}
	//fmt.Printf("%+v\n", pkg)
	fileList := pkg.GoFiles
	//
	panels := []string{}
	gen := false
	id := 0
	for _, fl := range fileList {
		src, err := os.Open(dir + "/" + fl)
		if err != nil {
			log.Fatalf("cannot process directory %s: %s", dir, err)
		}

		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, src.Name(), src, 0)
		if err != nil {
			panic(err)
		}
		ast.FilterFile(f, func(s string) bool {
			return s == structName
			//return true
		})
		ast.Inspect(f, func(node ast.Node) bool {

			switch val := node.(type) {
			case *ast.StructType:
				if gen {
					gen = false
					for _, field := range val.Fields.List {
						if field.Tag != nil {
							v := strings.TrimRight(field.Tag.Value[1:], "`") // 去掉头尾的`
							tag := reflect.StructTag(v)
							//fmt.Println(field.Names[0].Name, v, tag, )
							//fmt.Println(tag.Lookup("trace_val"))
							if dimension, b := tag.Lookup("trace_val"); b && dimension != "" {
								panel := panelMeta{
									namespace:   namespace,
									serviceName: serviceName,
									dimension:   dimension,
									endpoint:    "{{endpoint}}",
								}
								id += 2
								panel.id = id
								panels = append(panels, panel.Qps())
								id += 2
								panel.id = id
								panels = append(panels, panel.Duration())
							}
						}
					}
				}

			case *ast.GenDecl:
				for _, spe := range val.Specs {
					if spe.(*ast.TypeSpec).Name.Name == structName {
						gen = true
					}
				}
			default:
				//if val != nil && reflect.TypeOf(val) != nil {
				//	fmt.Println(reflect.TypeOf(val), val)
				//}
			}
			return true
		})
	}

	outputName := filepath.Join(dir, fmt.Sprintf("%s.json", serviceName))

	// Write to file.
	outputFile, err := os.Create(outputName)
	if err != nil {
		log.Fatalf("create file output error: %s", err)
	}

	var renderData = map[string]string{
		"panels":         strings.Join(panels, ","),
		"dashBoardTitle": dashBoardTitle,
	}

	err = tem.Execute(outputFile, renderData)
	if err != nil {
		log.Fatalf("render template error: %s", err)
	}

}
