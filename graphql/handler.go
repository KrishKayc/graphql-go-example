package graphql

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// NewHandler is for the graphql
// @Summary The GraphQL end point of the app
// @Description Open 'GraphiQL' tool and type localhost:3000/graphql
// @Router /graphql [post]
func NewHandler(db *sql.DB, cache *redis.Client) (*handler.Handler, error) {
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    newQuery(db, cache),
			Mutation: newMutation(db, cache),
		},
	)
	if err != nil {
		return nil, err
	}

	return handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	}), nil
}
