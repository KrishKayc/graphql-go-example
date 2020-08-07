package booker

import "database/sql"

//Cab represents a cab
type Cab struct {
	ID       int64  `json:"id"`
	CabType  int8   `json:"type"`
	DriverID int64  `json:"driverId"`
	Number   string `json:"number"`
	db       *sql.DB
}

//NewCab instantiates a new cab
func NewCab(id int64, cabType int8, driverID int64, number string, db *sql.DB) *Cab {
	return &Cab{ID: id, CabType: cabType, DriverID: driverID, Number: number, db: db}
}

//Load loads the cab using ids
func (c *Cab) Load() chan error {
	return goroutine(func(out chan error) {
		r := c.db.QueryRow("select id, type, driverId, number from Cab where id = ?", c.ID)
		out <- c.loader(r)
	})

}

func (c *Cab) loader(r *sql.Row) error {
	var id, driverID int64
	var cabType int8
	var number string
	if err := r.Scan(&id, &cabType, &driverID, &number); err != nil {
		return err
	}
	c.ID = id
	c.CabType = cabType
	c.DriverID = driverID
	c.Number = number

	return nil
}
