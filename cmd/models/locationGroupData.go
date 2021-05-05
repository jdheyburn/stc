package models

import (
	"time"

	"gorm.io/gorm"
)

type LocationGroupData struct {
	gorm.Model
	GroupUICCode string
	Description  string
	ERSCountry   string
	ERSCode      string
	StartDate    *time.Time
	EndDate      *time.Time
	QuoteDate    *time.Time
}

func (LocationGroupData) TableName() string {
	return "location_group"
}
