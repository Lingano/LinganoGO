package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Annotations(entgql.OrderField("ID")),
		field.String("name").
			NotEmpty().
			Annotations(entgql.OrderField("NAME")),
		field.String("email").
			Unique().
			NotEmpty().
			Annotations(entgql.OrderField("EMAIL")),
		field.String("password").
			Sensitive().
			NotEmpty(),
		field.Enum("role").
			Values("ADMIN", "USER").
			Default("USER").
			Annotations(entgql.OrderField("ROLE")),
		field.Bool("is_verified").
			Default(false),
		field.String("verification_token").
			Optional().
			Nillable(),
		field.String("reset_password_token").
			Optional().
			Nillable(),
		field.Time("reset_password_expires").
			Optional().
			Nillable(),
		field.JSON("profile", map[string]interface{}{}).
			Optional(),
		field.JSON("preferences", map[string]interface{}{}).
			Optional(),
		field.JSON("saved_words", map[string]interface{}{}).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Annotations(entgql.OrderField("CREATED_AT")),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entgql.OrderField("UPDATED_AT")),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("readings", Reading.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("flashcards", Flashcard.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("posts", Post.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
