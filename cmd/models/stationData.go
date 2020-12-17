package models

import "time"

type StationData struct {
	UIC         string     `db:"uic"`
	NLC         string     `db:"nlc"`
	CRS         string     `db:"crs"`
	FareGroup   string     `db:"fare_group"`
	Description string     `db:"description"`
	StartDate   *time.Time `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}
