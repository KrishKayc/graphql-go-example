package booker_test

import (
	"bookCab/booker"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLoadCab(t *testing.T) {
	db, mock, cache := mock(t)
	defer db.Close()
	defer cache.Close()

	cab := booker.NewCab(9888875353, 7, 76347583784, "KA11B0545", db)
	rows := sqlmock.NewRows([]string{"id", "type", "driverId", "number"}).AddRow(cab.ID, cab.CabType, cab.DriverID, cab.Number)
	mock.ExpectQuery("^select (.+) from Cab*").WithArgs(cab.ID).WillReturnRows(rows)
	handleErr(<-cab.Load(), t)
}
