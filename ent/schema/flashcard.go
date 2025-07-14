package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Flashcard struct {
	ent.Schema
}

// Fields of the Flashcard.
func (Flashcard) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Annotations(entgql.OrderField("ID")),
		field.String("question").
			NotEmpty().
			Annotations(entgql.OrderField("QUESTION")),
		field.String("answer").
			NotEmpty().
			Annotations(entgql.OrderField("ANSWER")),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Annotations(entgql.OrderField("CREATED_AT")),
		field.Time("last_reviewed_at").
			Optional().
			Annotations(entgql.OrderField("LAST_REVIEWED_AT")),
		field.UUID("user_id", uuid.UUID{}).
			Annotations(entgql.Skip()),
	}
}

// Edges of the Flashcard.
func (Flashcard) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("flashcards").
			Field("user_id").
			Required().
			Unique(),
	}
}