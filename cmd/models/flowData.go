package models

import (
	"time"

	"gorm.io/gorm"
)

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
