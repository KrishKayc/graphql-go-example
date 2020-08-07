package graphql

import (
	"bookCab/graphql/fields"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/graphql-go/graphql"
)

func newMutation(db *sql.DB, cache *redis.Client) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "mutation",
		Fields: graphql.Fields{
			"createUser":  fields.CreateUser(db),
			"bookCab":     fields.BookCab(db, cache),
			"setLocation": fields.SetLocation(db, cache),
		},
	})
}
