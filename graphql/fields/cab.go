package fields

import (
	"bookCab/booker"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/graphql-go/graphql"
)

var cab = graphql.NewObject(graphql.ObjectConfig{
	Name: "cab",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.Int},
		"type":     &graphql.Field{Type: graphql.String},
		"driverId": &graphql.Field{Type: graphql.String},
		"number":   &graphql.Field{Type: graphql.String}},

	Description: "cab details"})

//SetLocation request booking
func SetLocation(db *sql.DB, cache *redis.Client) *graphql.Field {
	return &graphql.Field{
		Type: booking,
		Args: graphql.FieldConfigArgument{
			"cabId":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"locationId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return rslvCabLocation(p, db, cache)
		},
		Description: "cab",
	}
}

func rslvCabLocation(p graphql.ResolveParams, db *sql.DB, cache *redis.Client) (i interface{}, e error) {

	l, c, err := parseCabLocationArgs(p, db, cache)

	if err != nil {
		return nil, err
	}

	if err = <-l.Load(); err != nil {
		return nil, err
	}

	if err := <-c.Load(); err != nil {
		return nil, err
	}

	if err := l.AddCab(c.ID); err != nil {
		return nil, err
	}

	//Return the cab details after booking
	return c, nil
}

func parseCabLocationArgs(p graphql.ResolveParams, db *sql.DB, cache *redis.Client) (*booker.Location, *booker.Cab, error) {
	var locID, cabID int64

	if err := validateInt64("locationId", p.Args["locationId"], &locID); err != nil {
		return nil, nil, err
	}
	if err := validateInt64("cabId", p.Args["cabId"], &cabID); err != nil {
		return nil, nil, err
	}

	l := booker.NewLocation(locID, "", -1, -1, db, cache)
	c := booker.NewCab(cabID, -1, -1, "", db)
	return l, c, nil

}
