package fields_test

import (
	"bookCab/booker"
	"bookCab/graphql/fields"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/graphql-go/graphql"
)

func TestResolveUser(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	f := fields.User(db)
	p := graphql.ResolveParams{Args: map[string]interface{}{"id": 1}}

	rows := sqlmock.NewRows([]string{"id", "name", "locationId", "phoneNumber"}).AddRow(1, "Rob Pike", 34567, 543456)
	mock.ExpectQuery("^select (.+) from User*").WithArgs(1).WillReturnRows(rows)

	u, err := f.Resolve(p)
	handleErr(err, t)

	user := u.(*booker.User)
	if user.Name != "Rob Pike" {
		t.Error("user resolve function not working")
	}

	if len(user.Bookings) > 0 {
		t.Error("bookings must not be loaded by resolved until requested")
	}

}

func TestResolveUserWithoutId(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	f := fields.User(db)
	p := graphql.ResolveParams{Args: map[string]interface{}{}}

	_, err := f.Resolve(p)
	if err.Error() != "Argument 'id' is required and is missing." {
		t.Error("Error not thrown for required args")
	}

}

func TestResolveCreateUserWithoutName(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	f := fields.CreateUser(db)
	p := graphql.ResolveParams{Args: map[string]interface{}{}}

	_, err := f.Resolve(p)
	if err.Error() != "Argument 'name' is required and is missing." {
		t.Error("Error not thrown for required args")
	}
}

func TestResolveCreateUserInvalidPhone(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	f := fields.CreateUser(db)
	p := graphql.ResolveParams{Args: map[string]interface{}{"name": "test", "locId": 1234, "phone": "abcd"}}

	_, err := f.Resolve(p)
	if err.Error() != "Invalid Type for argument 'phone'. Expect type 'Int64'" {
		t.Error("Error not thrown for invalid arg types")
	}
}

func TestCreateUserArgs(t *testing.T) {
	db, mock, cache := mock(t)
	fmt.Println(mock)
	defer db.Close()
	defer cache.Close()

	f := fields.CreateUser(db)
	if len(f.Args) != 3 {
		t.Error("create user args modified")
	}
}
func mock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *redis.Client) {
	db, mock, err := sqlmock.New()
	handleErr(err, t)

	s, err := miniredis.Run()
	handleErr(err, t)

	cache := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	return db, mock, cache
}

func handleErr(err error, t *testing.T) {
	if err != nil {
		t.Error(err.Error())
	}
}
