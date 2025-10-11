package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shunwuse/go-hris/ent/schema/mixin"
)

type UserRole struct {
	ent.Schema
}

func (UserRole) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.IDMixin{},
		mixin.TimeMixin{},
	}
}

func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("user_id").
			Positive(),
		field.Uint("role_id").
			Positive(),
	}
}

func (UserRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			Required().
			Field("user_id"),
		edge.To("role", Role.Type).
			Unique().
			Required().
			Field("role_id"),
	}
}

func (UserRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "role_id").
			Unique(),
	}
}

func (UserRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "user_role",
		},
	}
}
