package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var Analyzer = &analysis.Analyzer{
	Name: "checkswitch",
	Doc:  "reports checkswitch",
	//	Flags:            flag.FlagSet{},
	Run: run,
	//	RunDespiteErrors: true,
	//	Requires:         []*analysis.Analyzer{},
	//	ResultType:       nil,
	//	FactTypes:        []analysis.Fact{},
}

func main() {
	singlechecker.Main(Analyzer)
}

// перевести на мап[тип]индекс.
var orderTypes = map[string]int{
	"Document":  0,
	"Array":     1,
	"float64":   2,
	"string":    3,
	"Binary":    4,
	"ObjectID":  5,
	"bool":      6,
	"time.Time": 7,
	"NullType":  8,
	"Regex":     9,
	"int32":     10,
	"Timestamp": 11,
	"int64":     12,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			var idx int
			switch n := n.(type) {
			case *ast.CommentGroup:
			//	fmt.Printf("%+v\n", n)
			case *ast.TypeSwitchStmt:
				var name string
				for _, el := range n.Body.List {
					for _, cs := range el.(*ast.CaseClause).List {
						switch cs := cs.(type) {
						case *ast.StarExpr:
							if sexp, ok := cs.X.(*ast.SelectorExpr); ok {
								name = sexp.Sel.Name
								// name = fmt.Sprintf("%s.%s", sexp.X.(*ast.Ident).Name, sexp.X.(*ast.Ident).Name)
							}
						case *ast.SelectorExpr:
							name = fmt.Sprintf("%s.%s", cs.X, cs.Sel.Name)

						case *ast.Ident:
							name = cs.Name
						}

						iSl, ok := orderTypes[name]
						if ok && (iSl < idx) {
							pass.Reportf(n.Pos(), "non-observance of the preferred order of types")

						}
						idx = iSl
					}
				}
			}

			return true
		})
	}

	return nil, nil
}
