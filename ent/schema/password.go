package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shunwuse/go-hris/ent/schema/mixin"
)

// Password holds the schema definition for the Password entity.
type Password struct {
	ent.Schema
}

func (Password) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.IDMixin{},
		mixin.TimeMixin{},
	}
}

func (Password) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("user_id").
			Positive(),
		field.String("hash").
			NotEmpty(),
	}
}

func (Password) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("password").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (Password) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "passwords",
		},
	}
}
