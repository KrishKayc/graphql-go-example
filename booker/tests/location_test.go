package booker_test

import (
	"bookCab/booker"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

func TestLoadLocation(t *testing.T) {

	db, mock, cache := mock(t)
	defer db.Close()
	defer cache.Close()

	l := booker.NewLocation(6287609876543287654, "chennai", 34.23, 43.23, db, cache)
	rows := sqlmock.NewRows([]string{"id", "name", "latitude", "longitude"}).AddRow(l.ID, l.Name, l.Latitude, l.Longitude)

	mock.ExpectQuery("^select (.+) from Location*").WithArgs(l.ID).WillReturnRows(rows)

	handleErr(<-l.Load(), t)

	if l.Name != "chennai" {
		t.Error("Failed loading location")
	}

}

func TestNoNearbyCabsError(t *testing.T) {
	db, mock, cache := mock(t)
	defer db.Close()
	defer cache.Close()

	l := booker.NewLocation(342, "", 78.34, 12.45, db, cache)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(45678)
	mock.ExpectQuery("^select id from Region*").WithArgs(l.Latitude, l.Latitude, l.Longitude, l.Longitude).WillReturnRows(rows)

	_, err := l.NearbyCabs()

	if err.Error() != "No nearby cabs currently. Please refresh" {
		t.Error("No cabs error message not shown")
	}

}

func TestAvailableNearbyCabs(t *testing.T) {
	db, mock, cache := mock(t)
	defer db.Close()
	defer cache.Close()

	//Add cab in the location
	var cabID int64 = 4287609876543287654
	l := booker.NewLocation(342, "", 78.34, 12.45, db, cache)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(45678)
	mock.ExpectQuery("^select id from Region*").WithArgs(l.Latitude, l.Latitude, l.Longitude, l.Longitude).WillReturnRows(rows)

	err := l.AddCab(cabID)
	handleErr(err, t)

	//Then check for cabs in the location
	rows = sqlmock.NewRows([]string{"id"}).AddRow(45678)
	mock.ExpectQuery("^select id from Region*").WithArgs(l.Latitude, l.Latitude, l.Longitude, l.Longitude).WillReturnRows(rows)
	cabs, err := l.NearbyCabs()
	handleErr(err, t)

	if cabs[0] != cabID {
		t.Error("Nearby cabs not found")
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
