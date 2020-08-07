package booker_test

import (
	"bookCab/booker"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserBookingsHistory(t *testing.T) {
	db, mock, cache := mock(t)
	defer db.Close()
	defer cache.Close()

	user := booker.NewUser(987367897583, "Rob Pike", 5678475, 3847587346, db)
	rows := sqlmock.NewRows([]string{"id", "cabId", "userId", "pickUpLocationId", "dropLocationId", "booked"}).AddRow(783478563487, 347538485, user.ID, 74875347389, 3847537, time.Now())
	mock.ExpectQuery("^select (.+) from Booking*").WithArgs(user.ID).WillReturnRows(rows)

	handleErr(<-user.SetBookings(), t)

	if user.Bookings[0].ID != 783478563487 {
		t.Error("user booking history not fetched properly")
	}

}
