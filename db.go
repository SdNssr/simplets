package main

import (
	"github.com/boltdb/bolt"
	"bytes"
	"time"
	"fmt"
)

type Datapoints struct {
	x []time.Time
	y []float64
}

type ServerEnv struct {
	db *bolt.DB
}

type DataPoint struct {
	Value  float64
}

func (env ServerEnv) GetSeries() ([]string, error) {
	series := make([]string, 0)

	err := env.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			series = append(series, string(name))

			return nil
		})
	})

	if err != nil {
		return nil, fmt.Errorf("accessing db: %s", err)
	}

	return series, nil
}

func (env ServerEnv) GetDataPoints(name string, duration string) (*Datapoints, error) {
	durationp, err := time.ParseDuration(duration)

	if err != nil {
		return nil, fmt.Errorf("parsing duration: %s", err)
	}

	to := []byte(time.Now().Format(time.RFC3339))
	from := []byte(time.Now().Add(-1 * durationp).Format(time.RFC3339))

	datapoints := Datapoints{}

	err = env.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(name)).Cursor()

		for k, v := c.Seek(from); k != nil && bytes.Compare(k, to) <= 0; k, v = c.Next() {
			x, err := time.Parse(time.RFC3339, string(k))

			if err != nil {
				return fmt.Errorf("parse time: %s", err)
			}

			datapoints.x = append(datapoints.x, x)
			datapoints.y = append(datapoints.y, btof(v))
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("fetch values: %s", err)
	}

	return &datapoints, nil
}

func (env ServerEnv) AddDataPoint(name string, value float64) error {
	time := time.Now().Format(time.RFC3339)

	err := env.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(name))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(time), ftob(value))

		if err != nil {
			return fmt.Errorf("put value: %s", err)
		}

		return err
	})

	if err != nil {
		return fmt.Errorf("add values: %s", err)
	}

	return nil
}
