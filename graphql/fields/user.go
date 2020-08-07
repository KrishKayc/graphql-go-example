package fields

import (
	"bookCab/booker"
	"database/sql"

	"github.com/graphql-go/graphql"
	ast "github.com/graphql-go/graphql/language/ast"
)

var user = graphql.NewObject(graphql.ObjectConfig{
	Name: "user",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.Int},
		"name":  &graphql.Field{Type: graphql.String},
		"locID": &graphql.Field{Type: graphql.Int},
		"phone": &graphql.Field{Type: graphql.String},

		"bookings": &graphql.Field{Type: graphql.NewList(booking)}},

	Description: "User data"})

//User returns new user field for querying
func User(db *sql.DB) *graphql.Field {
	return &graphql.Field{
		Type: user,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return rslvUser(p, db)
		},
		Description: "user",
	}
}

//CreateUser creates new user
func CreateUser(db *sql.DB) *graphql.Field {
	return &graphql.Field{
		Type: user,
		Args: graphql.FieldConfigArgument{
			"name":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"locId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"phone": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			return rslvCreateUser(p, db)
		},
		Description: "user",
	}
}

func rslvCreateUser(p graphql.ResolveParams, db *sql.DB) (i interface{}, e error) {
	u, err := parseUserArgs(p, db, true)
	if err != nil {
		return nil, err
	}
	if err = <-u.Save(); err != nil {
		return nil, err
	}
	return u, nil
}

func rslvUser(p graphql.ResolveParams, db *sql.DB) (i interface{}, e error) {
	u, err := parseUserArgs(p, db, false)
	if err != nil {
		return nil, err
	}

	//load the basic details first
	if err = <-u.Load(); err != nil {
		return nil, err
	}

	//if booking history is requested, load it as it is costlier
	if requested(p, "bookings") {
		if err = <-u.SetBookings(); err != nil {
			return nil, err
		}
	}

	return u, err

}

func requested(p graphql.ResolveParams, fldName string) bool {
	for _, r := range p.Info.FieldASTs {
		for _, s := range r.SelectionSet.Selections {
			field := s.(*ast.Field)
			if field.Name.Value == fldName {
				return true
			}
		}

	}
	return false
}

func parseUserArgs(p graphql.ResolveParams, db *sql.DB, mutation bool) (*booker.User, error) {
	var id, locID, phone int64
	var name string

	if mutation {
		if err := validateString("name", p.Args["name"], &name); err != nil {
			return nil, err
		}
		if err := validateInt64("locId", p.Args["locId"], &locID); err != nil {
			return nil, err
		}
		if err := validateInt64("phone", p.Args["phone"], &phone); err != nil {
			return nil, err
		}

		return booker.NewUser(id, name, locID, phone, db), nil
	}

	//id not required for mutation but required for query
	if err := validateInt64("id", p.Args["id"], &id); err != nil {
		return nil, err
	}
	return booker.NewUser(id, name, locID, phone, db), nil
}

func validateInt64(arg string, val interface{}, assign *int64) error {
	if val == nil {
		return &RequiredArgumentError{arg: arg}
	}
	if v, ok := val.(int); ok {
		*assign = int64(v)
		return nil
	}
	return &InvalidArgumentError{arg: arg, expectedType: "Int64"}

}

func validateString(arg string, val interface{}, assign *string) error {
	if val == nil {
		return &RequiredArgumentError{arg: arg}
	}

	if v, ok := val.(string); ok {
		*assign = v
		return nil
	}
	return &InvalidArgumentError{arg: arg, expectedType: "String"}

}
