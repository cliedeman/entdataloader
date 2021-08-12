package entdataloader

import (
	"embed"
	"entgo.io/ent/entc/gen"
	"fmt"
	"text/template"
)

//go:embed template/*.tmpl
var templates embed.FS

var (
	DataLoaderTemplate = parse("template/entdataloader.tmpl")

	// AllTemplates holds all templates for extending ent to support GraphQL.
	AllTemplates = []*gen.Template{
		DataLoaderTemplate,
	}
)

func hasEdgeField(edge *gen.Edge) bool {
	if edge == nil {
		return false
	}
	f := edge.Field()
	if f != nil {
		return f.IsEdgeField()
	}
	return false
}

func dataloaderEdge(edge *gen.Edge) (bool, error) {
	ant := &Annotation{}
	if edge.Annotations != nil && edge.Annotations[ant.Name()] != nil {
		if err := ant.Decode(edge.Annotations[ant.Name()]); err != nil {
			return false, err
		}
		switch edge.Rel.Type {
		case gen.O2M:
			{
				isEdgeField := hasEdgeField(edge.Ref)
				if isEdgeField {
					return true, nil
				} else {
					return false, fmt.Errorf("edge %s on Type %s does not have an edge field", edge.Name, edge.Type.Name)
				}
			}
		case gen.M2O:
			{
				return true, nil
			}
		}

	}
	return false, nil
}

func dataloaderFieldName(edge *gen.Edge) string {
	switch edge.Rel.Type {
	case gen.O2M:
		return edge.Ref.Field().Name
	case gen.M2O:
		return edge.Type.ID.Name
	}
	return "invalid"
}

func dataloaderPlural(edge *gen.Edge) bool {
	switch edge.Rel.Type {
	case gen.O2M:
		return true
	case gen.M2O:
		return false
	}
	return false
}

func parse(name string) *gen.Template {
	bytes, err := templates.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return gen.MustParse(gen.NewTemplate(name).
		Funcs(gen.Funcs).
		Funcs(template.FuncMap{
			"dataloaderEdge":      dataloaderEdge,
			"dataloaderPlural":    dataloaderPlural,
			"dataloaderFieldName": dataloaderFieldName,
		}).
		Parse(string(bytes)))
}
