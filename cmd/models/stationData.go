package models

import (
	"time"

	"gorm.io/gorm"
)

type StationData struct {
	gorm.Model
	UIC         string `gorm:"column:uic"`
	NLC         string
	CRS         string
	FareGroup   string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
}

// type Tabler interface {
// 	TableName() string
// }

func (StationData) TableName() string {
	return "location"
}
