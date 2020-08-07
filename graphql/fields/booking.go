package fields

import (
	"bookCab/booker"
	"database/sql"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/graphql-go/graphql"
)

var booking = graphql.NewObject(graphql.ObjectConfig{
	Name: "booking",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.Int},
		"userId":      &graphql.Field{Type: graphql.String},
		"pickUpLocId": &graphql.Field{Type: graphql.String},
		"dropLocId":   &graphql.Field{Type: graphql.String},
		"time":        &graphql.Field{Type: graphql.DateTime},
		"cab":         &graphql.Field{Type: cab}},

	Description: "Booking data"})

//BookCab request booking
func BookCab(db *sql.DB, cache *redis.Client) *graphql.Field {
	return &graphql.Field{
		Type: booking,
		Args: graphql.FieldConfigArgument{
			"userId":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"pickUpLocId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"dropLocId":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return rslvBookCab(p, db, cache)
		},
		Description: "booking",
	}
}

func rslvBookCab(p graphql.ResolveParams, db *sql.DB, cache *redis.Client) (i interface{}, e error) {

	b, err := parseBookingArgs(p, db, cache)

	if err != nil {
		return nil, err
	}

	if err = <-b.Save(); err != nil {
		return nil, err
	}

	//load the cab details for user to look at
	if <-b.Cab.Load(); err != nil {
		return nil, err
	}

	return b, nil
}

func parseBookingArgs(p graphql.ResolveParams, db *sql.DB, cache *redis.Client) (*booker.Booking, error) {

	var userID, pickUpLocID, dropLocID int64

	if err := validateInt64("userId", p.Args["userId"], &userID); err != nil {
		return nil, err
	}
	if err := validateInt64("pickUpLocId", p.Args["pickUpLocId"], &pickUpLocID); err != nil {
		return nil, err
	}
	if err := validateInt64("dropLocId", p.Args["dropLocId"], &dropLocID); err != nil {
		return nil, err
	}

	b := booker.NewBooking(-1, userID, pickUpLocID, dropLocID, time.Now(), db, cache)

	return b, nil
}
