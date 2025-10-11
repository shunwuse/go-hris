package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type IDMixin struct {
	mixin.Schema
}

func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Positive().
			Unique().
			Immutable(),
	}
}
