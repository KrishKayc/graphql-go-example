package booker

import (
	"database/sql"
)

//User represents a user
type User struct {
	ID         int64
	Name       string
	LocationID int64 `json:"locId"`
	PhoneNo    int64 `json:"phone"`

	Bookings []Booking `json:"bookings"`
	db       *sql.DB
}

//NewUser returns a new user
func NewUser(id int64, name string, locID int64, phoneNo int64, db *sql.DB) *User {
	return &User{ID: id, Name: name, LocationID: locID, PhoneNo: phoneNo, db: db}
}

//Load gives a new user for the id
func (u *User) Load() chan error {
	return goroutine(func(out chan error) {
		row := u.db.QueryRow("select id, name, locationId, phoneNumber from User where id=?", u.ID)
		out <- u.loader(row)
	})

}

//Save writes a new user to DB
func (u *User) Save() chan error {
	return goroutine(func(out chan error) {
		res, err := u.db.Exec("insert into User(name, locationId, phoneNumber) values(?,?,?)", u.Name, u.LocationID, u.PhoneNo)

		if err != nil {
			out <- err
		}
		id, err := res.LastInsertId()

		if err != nil {
			out <- err
		}
		u.ID = id

		out <- nil
	})
}

//SetBookings sets history of the user
func (u *User) SetBookings() chan error {
	return goroutine(func(out chan error) {
		r, err := u.db.Query("select id, cabId, userId, pickUpLocationId, dropLocationId, booked from Booking where userId=?", u.ID)
		defer r.Close()

		if err != nil {
			out <- err
		}

		u.Bookings = make([]Booking, 0)
		for r.Next() {
			var b Booking
			if err = b.loader(r); err != nil {
				out <- err
			}
			u.Bookings = append(u.Bookings, b)
		}
		out <- nil
	})

}

func (u *User) loader(r *sql.Row) error {
	var id, phoneNo, locID int64
	var name string
	if err := r.Scan(&id, &name, &locID, &phoneNo); err != nil {
		return err
	}
	u.ID = id
	u.Name = name
	u.LocationID = locID
	u.PhoneNo = phoneNo

	return nil
}
