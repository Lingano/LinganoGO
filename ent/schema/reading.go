package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Reading holds the schema definition for the Reading entity.
type Reading struct {
	ent.Schema
}

// Fields of the Reading.
func (Reading) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Annotations(entgql.OrderField("ID")),
		field.String("title").
			NotEmpty().
			Annotations(entgql.OrderField("TITLE")),
		field.UUID("user_id", uuid.UUID{}).
			Annotations(entgql.Skip()),
		field.Bool("finished").
			Default(false).
			Annotations(entgql.OrderField("FINISHED")),
		field.Bool("public").
			Default(false).
			Annotations(entgql.OrderField("PUBLIC")),
	}
}

// Edges of the Reading.
func (Reading) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("readings").
			Field("user_id").
			Required().
			Unique(),
	}
}
