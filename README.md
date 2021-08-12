# Ent Dataloader

## Introduction

Ent dataloader generates datalodaers for edges.

Only `O2M` (with edge field) and `M2O` edges are currently supported.


```golang
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/cliedeman/entdataloader"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type).
			Annotations(entdataloader.Annotation{}),
	}
}

type Post struct {
    ent.Schema
}

func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.Int("author_id").
            Optional(),
    }
}

func (Post) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("author", User.Type).
            Ref("posts").
            Field("author_id").
            Unique().
            Annotations(entdataloader.Annotation{}),
    }
}
```