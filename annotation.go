package entdataloader

import (
	"encoding/json"
	"entgo.io/ent/schema"
)

type Annotation struct {
	// TODO
}

func (Annotation) Name() string {
	return "entdataloader"
}

func (a *Annotation) Decode(annotation interface{}) error {
	buf, err := json.Marshal(annotation)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	// _ schema.Merger     = (*Annotation)(nil)
)
