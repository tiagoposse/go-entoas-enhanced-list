package filter

import (
	"fmt"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/ogen-go/ogen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type OperationExtension struct {
	entc.DefaultExtension
	GlobalAnnotation *Annotation
}

func NewOperationExtension(opts ...MutatorOpt) *OperationExtension {
	ant := &Annotation{}
	for _, opt := range opts {
		opt(ant)
	}

	return &OperationExtension{
		GlobalAnnotation: ant,
	}
}


func (ext *OperationExtension) Hooks() []gen.Hook {
	return []gen.Hook{
		ext.generate(),
	}
}

// DisallowTypeName ensures there is no ent.Schema with the given name in the graph.
func (ext *OperationExtension) generate() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(graph *gen.Graph) error {
			for _, node := range graph.Nodes {
				ant := Annotation{}
				ant = ant.Merge(*ext.GlobalAnnotation).(Annotation)

				if ann, ok := node.Annotations[Annotation{}.Name()]; ok {
					nodeAnt := Annotation{}
					nodeAnt.Decode(ann)
					ant = ant.Merge(nodeAnt).(Annotation)
				}
				
				anns := node.Annotations
				anns[Annotation{}.Name()] = ant
				node.Annotations = anns
			}

			return next.Generate(graph)
		})
	}
}

func (ext *OperationExtension) getAnnotations(ant gen.Annotations) Annotation {
	decoded := Annotation{}
	if ann, ok := ant[Annotation{}.Name()]; ok {
		decoded.Decode(ann)
	}
	filterAnt := Annotation{}.Merge(ext.GlobalAnnotation).(Annotation).Merge(decoded)

	return filterAnt.(Annotation)
}

func (ext *OperationExtension) Mutator(graph *gen.Graph, spec *ogen.Spec) error {
	filterMods := make(map[string][]*Opt)
	opAnnotations := make(map[string]Annotation)

	for _, node := range graph.Nodes {
		filterAnt := ext.getAnnotations(node.Annotations)
		node.Annotations[Annotation{}.Name()] = filterAnt
		opName := fmt.Sprintf("list%s", node.Name)
		filterMods[opName] = filterAnt.FilterFields
		opAnnotations[opName] = filterAnt
		
		for _, edge := range node.Edges {
			edgeAnt := ext.getAnnotations(edge.Annotations)

			opName := fmt.Sprintf("list%s%s", node.Name, cases.Title(language.Und, cases.NoLower).String(edge.Name))
			filterMods[opName] = edgeAnt.FilterFields
			opAnnotations[opName] = edgeAnt
		}
	}

	for _, pathItem := range spec.Paths {
		if pathItem.Get == nil || !strings.HasPrefix(pathItem.Get.OperationID, "list") {
			continue
		}

		var ant Annotation
		if a, ok := opAnnotations[pathItem.Get.OperationID]; ok {
			ant = a
		} else {
			ant = *ext.GlobalAnnotation
		}

		newParams := make([]*ogen.Parameter, 0)
		for _, prop := range pathItem.Get.Parameters {
			prop := prop

			switch prop.Name {
			case "itemsPerPage":
				if ant.ItemsPerPage != nil {
					ant.ItemsPerPage.Set(prop)
				}
				if !ant.NoPagination {
					newParams = append(newParams, prop)
				}
			case "page":
				if ant.Page != nil {
					ant.Page.Set(prop)
				}
				if !ant.NoPagination {
					newParams = append(newParams, prop)
				}
			default:
				newParams = append(newParams, prop)
			}
		}

		if ant.Sort != nil {
			newParams = append(newParams, &ogen.Parameter{
				Name: ant.Sort.Name,
				In: ant.Sort.In,
				Schema: &ogen.Schema{Type: "string"},
			})
		}

		if _, ok := filterMods[pathItem.Get.OperationID]; ok {
			newParams = append(newParams, &ogen.Parameter{
				Name: "filter",
				In: "query",
				Schema: &ogen.Schema{Type: "string"},
			})
		}

		if ant.ReturnTotal != nil {
			for n, resp := range pathItem.Get.Responses {
				if n != "200" {
					continue
				}

				if resp.Headers == nil {
					resp.SetHeaders(make(map[string]*ogen.Parameter))
				}

				resp.Headers[ant.ReturnTotal.Name] = ogen.NewParameter().SetSchema(ogen.Int()).SetRequired(true)
			}
		}

		pathItem.Get.Parameters = newParams
	}

	return nil
}


type MutatorOpt func(*Annotation)

// func Mutator(opts ...MutatorOpt) func(graph *gen.Graph, spec *ogen.Spec) error {
// 	return func (graph *gen.Graph, spec *ogen.Spec) error {
// 		filterMods := make(map[string][]*Opt)
// 		ops := make([]string, 0)

// 		globalAnt := Annotation{}
// 		for _, opt := range opts {
// 			opt(&globalAnt)
// 		}

// 		for _, node := range graph.Nodes {
// 			filterAnt := globalAnt
// 			if ann, ok := node.Annotations[Annotation{}.Name()]; ok {
// 				filterAnt.Decode(ann)
// 			}

// 			opName := fmt.Sprintf("list%s", node.Name)
// 			ops = append(ops, opName)
// 			filterMods[opName] = filterAnt.FilterFields

// 			nodeAnt := globalAnt
// 			nodeAnt.Merge(filterAnt)
// 			anns := node.Annotations
// 			anns[Annotation{}.Name()] = nodeAnt
// 			node.Annotations = anns
// 		}

// 		for _, pathItem := range spec.Paths {
// 			if pathItem.Get == nil || slices.Index(ops, pathItem.Get.OperationID) == -1 {
// 				continue
// 			}

// 			newParams := make([]*ogen.Parameter, 0)
// 			for _, prop := range pathItem.Get.Parameters {
// 				prop := prop
// 				switch prop.Name {
// 				case "itemsPerPage":
// 					if globalAnt.ItemsPerPage != nil {
// 						globalAnt.ItemsPerPage.Set(prop)
// 					}
// 					if !globalAnt.NoPagination {
// 						newParams = append(newParams, prop)
// 					}
// 				case "page":
// 					if globalAnt.Page != nil {
// 						globalAnt.Page.Set(prop)
// 					}
// 					if !globalAnt.NoPagination {
// 						newParams = append(newParams, prop)
// 					}
// 				default:
// 					newParams = append(newParams, prop)
// 				}
// 			}

// 			if globalAnt.Sort != nil {
// 				newParams = append(newParams, &ogen.Parameter{
// 					Name: globalAnt.Sort.Name,
// 					In: globalAnt.Sort.In,
// 					Schema: &ogen.Schema{Type: "string"},
// 				})
// 			}

// 			if fs, ok := filterMods[pathItem.Get.OperationID]; ok {
// 				for _, filter := range fs {
// 					newParams = append(newParams, &ogen.Parameter{
// 						Name: filter.Name,
// 						In: filter.In,
// 						Schema: &ogen.Schema{Type: "string"},
// 					})
// 				}
// 			}

// 			pathItem.Get.Parameters = newParams
// 		}
// 		return nil
// 	}
// }
