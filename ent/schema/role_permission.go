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

type RolePermission struct {
	ent.Schema
}

func (RolePermission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.IDMixin{},
		mixin.TimeMixin{},
	}
}

func (RolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("role_id").
			Positive(),
		field.Uint("permission_id").
			Positive(),
	}
}

func (RolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role", Role.Type).
			Unique().
			Required().
			Field("role_id"),
		edge.To("permission", Permission.Type).
			Unique().
			Required().
			Field("permission_id"),
	}
}

func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "permission_id").
			Unique(),
	}
}

func (RolePermission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "role_permission",
		},
	}
}
