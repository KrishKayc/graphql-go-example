package fields_test

import (
	"bookCab/graphql/fields"
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestResolveCabLocationWithoutLocId(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	c := fields.SetLocation(db, cache)
	p := graphql.ResolveParams{Args: map[string]interface{}{"cabId": 2345}}

	_, err := c.Resolve(p)
	if err.Error() != "Argument 'locationId' is required and is missing." {
		t.Error("Error not thrown for required args")
	}
}

func TestResolveCabLocationWithoutCabId(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	c := fields.SetLocation(db, cache)
	p := graphql.ResolveParams{Args: map[string]interface{}{"locationId": 2345}}

	_, err := c.Resolve(p)
	if err.Error() != "Argument 'cabId' is required and is missing." {
		t.Error("Error not thrown for required args")
	}
}

func TestResolveCabLocationArgs(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	c := fields.SetLocation(db, cache)

	if len(c.Args) != 2 {
		t.Error("set cab location args modified")
	}
}
