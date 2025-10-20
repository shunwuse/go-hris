package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shunwuse/go-hris/ent/schema/mixin"
	"github.com/shunwuse/go-hris/internal/constants"
)

// Approval holds the schema definition for the Approval entity.
type Approval struct {
	ent.Schema
}

func (Approval) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.IDMixin{},
		mixin.TimeMixin{},
	}
}

func (Approval) Fields() []ent.Field {
	return []ent.Field{
		field.String("status").
			NotEmpty().
			GoType(constants.ApprovalStatus("")),
		field.Uint("creator_id").
			Positive(),
		field.Uint("approver_id").
			Positive().
			Optional().
			Nillable(),
	}
}

func (Approval) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("creator", User.Type).
			Required().
			Unique().
			Field("creator_id"),
		edge.To("approver", User.Type).
			Unique().
			Field("approver_id"),
	}
}

func (Approval) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "approvals",
		},
	}
}
