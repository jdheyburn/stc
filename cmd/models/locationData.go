package models

import (
	"time"

	"gorm.io/gorm"
)

// LocationData represents the records in location table
type LocationData struct {
	gorm.Model
	UIC         string `gorm:"column:uic"`
	NLC         string
	CRS         string
	FareGroup   string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
}

func (LocationData) TableName() string {
	return "location"
}

// LocationWithGroupData is the record the subquery method retrieves back
type LocationWithGroupData struct {
	gorm.Model
	UIC              string `gorm:"column:uic"`
	NLC              string
	CRS              string
	FareGroup        string
	Description      string
	StartDate        *time.Time
	EndDate          *time.Time
	GroupUicCode     string
	GroupDescription string
}

// LocationGroup contains fields pulled from the subquery method
type LocationGroup struct {
	UIC         string
	Description string
}

// LocationWithGroups is a flattened struct from LocationWithGroupData
type LocationWithGroups struct {
	UIC         string
	NLC         string
	CRS         string
	FareGroup   string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
	Groups      []*LocationGroup
}

func (l LocationWithGroups) IsGroupedStation() bool {
	return len(l.Groups) > 0
}
