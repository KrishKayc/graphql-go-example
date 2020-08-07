package fields_test

import (
	"bookCab/graphql/fields"
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestBookCabWithoutPickupPoint(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	b := fields.BookCab(db, cache)
	p := graphql.ResolveParams{Args: map[string]interface{}{"userId": 1234, "dropLocId": 2345}}
	_, err := b.Resolve(p)
	if err.Error() != "Argument 'pickUpLocId' is required and is missing." {
		t.Error("Error not thrown for required args")
	}
}

func TestBookCabWithoutDropPoint(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	b := fields.BookCab(db, cache)
	p := graphql.ResolveParams{Args: map[string]interface{}{"userId": 1234, "pickUpLocId": 2345}}
	_, err := b.Resolve(p)
	if err.Error() != "Argument 'dropLocId' is required and is missing." {
		t.Error("Error not thrown for required args")
	}
}

func TestBookCabWithoutUser(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	b := fields.BookCab(db, cache)
	p := graphql.ResolveParams{Args: map[string]interface{}{"dropLocId": 1234, "pickUpLocId": 2345}}
	_, err := b.Resolve(p)
	if err.Error() != "Argument 'userId' is required and is missing." {
		t.Error("Error not thrown for required args")
	}
}
