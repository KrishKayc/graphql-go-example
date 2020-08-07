package booker

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const cachePrefix = "BookCab:"

var ctx = context.Background()

//Region is a zone of locations
type Region struct {
	ID      int64
	Name    string
	MinLat  float32
	MaxLat  float32
	MinLong float32
	MaxLong float32
}

//Location represents a real time location
type Location struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`

	db    *sql.DB
	cache *redis.Client
}

//NewLocation gets a new location instance
func NewLocation(id int64, name string, latitude float32, longitude float32, db *sql.DB, c *redis.Client) *Location {
	return &Location{ID: id, Name: name, Latitude: latitude, Longitude: longitude, db: db, cache: c}
}

//Load loads the location from db
func (l *Location) Load() chan error {
	return goroutine(func(out chan error) {
		r := l.db.QueryRow("select id, name, latitude, longitude  from Location where id = ?", l.ID)
		out <- l.loader(r)
	})

}

//NearbyCabs fetches the nearby cabs in the location
func (l Location) NearbyCabs() ([]int64, error) {
	regionID := <-l.region()
	if regionID.err != nil {
		return nil, regionID.err
	}

	//get the neary cabs from the cache...
	//cache must be updated whenever a cab is active in a region i.e whenever driver starts the app in the region
	cabs, err := l.cache.LRange(ctx, cacheKey(regionID.data), 0, -1).Result()

	if err != nil {
		return nil, err
	}

	if len(cabs) == 0 {
		return nil, &NoCabsError{}
	}

	return stoi(cabs), nil
}

//AddCab adds a cab to the location's region
func (l Location) AddCab(cabID int64) error {
	regionID := <-l.region()
	if regionID.err != nil {
		return regionID.err
	}

	_, err := l.cache.RPush(ctx, cacheKey(regionID.data), cacheValue(cabID)).Result()

	return err

}

func (l Location) region() chan Int64Result {
	return goroutineInts(func(out chan Int64Result) {
		r := l.db.QueryRow("select id from Region where minLatitude < ? and maxLatitude > ? and minLongitude < ? and maxLongitude > ?", l.Latitude, l.Latitude, l.Longitude, l.Longitude)

		var id int64
		err := r.Scan(&id)

		if err != nil {
			out <- Int64Result{data: -1, err: err}
		}

		out <- Int64Result{data: id, err: nil}
	})

}

func (l *Location) loader(r *sql.Row) error {
	var id int64
	var lat, long float32
	var name string
	if err := r.Scan(&id, &name, &lat, &long); err != nil {
		return err
	}

	l.ID = id
	l.Latitude = lat
	l.Longitude = long
	l.Name = name

	return nil
}

func stoi(val []string) []int64 {
	ids := make([]int64, 0)
	for _, c := range val {
		if id, err := strconv.Atoi(c); err == nil {
			ids = append(ids, int64(id))
		}
	}

	return ids
}

func cacheKey(key int64) string {
	return cachePrefix + strconv.FormatInt(key, 10)
}

func cacheValue(value int64) string {
	return strconv.FormatInt(value, 10)
}
