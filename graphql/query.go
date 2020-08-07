package graphql

import (
	"bookCab/graphql/fields"
	"database/sql"

	"github.com/go-redis/redis/v8"

	"github.com/graphql-go/graphql"
)

func newQuery(db *sql.DB, cache *redis.Client) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "query",
		Fields: graphql.Fields{
			"user": fields.User(db),
		},
	})
}
