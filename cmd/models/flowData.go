package models

import (
	"time"

	"gorm.io/gorm"
)

// FareData maps directly to records in the flow table
type FlowData struct {
	gorm.Model
	FlowID          string
	OriginCode      string
	DestinationCode string
	RouteCode       string
	Direction       string
	StartDate       *time.Time
	EndDate         *time.Time
}

func (FlowData) TableName() string {
	return "flow"
}

type FlowDetail struct {
	gorm.Model
	FlowID          string
	OriginCode      string
	DestinationCode string
	Direction       string
	StartDate       *time.Time
	EndDate         *time.Time
	RouteCode       string
	// Left Join route
	RouteDesc string
}
