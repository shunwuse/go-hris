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

// RefreshToken holds the schema definition for the RefreshToken entity.
type RefreshToken struct {
	ent.Schema
}

func (RefreshToken) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.IDMixin{},
		mixin.TimeMixin{},
	}
}

func (RefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("token_hash").
			NotEmpty().
			Unique(),
		field.Uint("user_id").
			Positive(),
		field.Time("expires_at"),
		field.Bool("revoked").
			Default(false),
		field.Time("revoked_at").
			Nillable().
			Optional(),
	}
}

func (RefreshToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("refresh_tokens").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (RefreshToken) Indexes() []ent.Index {
	return []ent.Index{
		// fast token lookup
		index.Fields("token_hash").
			Unique(),
		// finding user's tokens
		index.Fields("user_id"),
		// cleaning up expired tokens
		index.Fields("expires_at"),
		// finding user's active tokens
		index.Fields("user_id", "revoked"),
	}
}

func (RefreshToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "refresh_tokens",
		},
	}
}
