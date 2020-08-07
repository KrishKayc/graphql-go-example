package booker

import (
	"database/sql"
	"time"

	"github.com/go-redis/redis/v8"
)

//Booking represents a booking request
type Booking struct {
	ID          int64
	CabID       int64
	UserID      int64
	PickUpLocID int64
	DropLocID   int64
	Time        time.Time `json:"time"`
	Cab         *Cab      `json:"cab"` //used by graphql

	db    *sql.DB
	cache *redis.Client
}

//NewBooking return a new booking instance
func NewBooking(id int64, userID int64, pickUpLocID int64, dropLocID int64, time time.Time, db *sql.DB, c *redis.Client) *Booking {
	return &Booking{ID: id, UserID: userID, PickUpLocID: pickUpLocID, DropLocID: dropLocID, Time: time, db: db, cache: c}
}

//Save a booking
func (b *Booking) Save() chan error {
	return goroutine(func(out chan error) {
		//find a cab near to pick up point and book
		if err := b.nearbyCab(); err != nil {
			out <- err
		}

		//book that cab for the user
		res, err := b.db.Exec("insert into Booking(cabId, userId, pickUpLocationId, dropLocationId, booked) values(?,?,?,?,?)", b.CabID, b.UserID, b.PickUpLocID, b.DropLocID, b.Time)

		if err != nil {
			out <- err
		}
		if b.ID, err = res.LastInsertId(); err != nil {
			out <- err
		}

		out <- nil
	})
}

func (b *Booking) loader(r *sql.Rows) error {
	var id, cabID, userID, pickUpLocID, dropLocID int64
	var time time.Time

	if err := r.Scan(&id, &cabID, &userID, &pickUpLocID, &dropLocID, &time); err != nil {
		return err
	}
	b.ID = id
	b.CabID = cabID
	b.UserID = userID
	b.PickUpLocID = pickUpLocID
	b.DropLocID = dropLocID
	b.Time = time

	return nil
}

func (b *Booking) nearbyCab() error {
	pickUpLoc := NewLocation(b.PickUpLocID, "", -1, -1, b.db, b.cache)
	if err := <-pickUpLoc.Load(); err != nil {
		return err
	}

	nearByCabs, err := pickUpLoc.NearbyCabs()

	if err != nil {
		return err
	}

	b.CabID = nearByCabs[0]
	b.Cab = NewCab(b.CabID, -1, -1, "", b.db)
	return nil
}
