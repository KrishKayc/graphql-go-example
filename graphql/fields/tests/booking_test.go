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

func TestBookCabField(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	b := fields.BookCab(db, cache)
	if b.Type.Name() != "booking" {
		t.Error("Wrong type for booking field")
	}
}

func TestBookCabArgs(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	b := fields.BookCab(db, cache)

	if len(b.Args) != 3 {
		t.Error("arguments modified for booking cabs")
	}
	if b.Args["userId"] == nil || b.Args["pickUpLocId"] == nil || b.Args["dropLocId"] == nil {
		t.Error("argument names modified for booking cabs")
	}
}
