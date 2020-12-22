package models

import (
	"gorm.io/gorm"
)

type FareData struct {
	gorm.Model
	FlowID          uint
	TicketCode      string
	Fare            uint
	RestrictionCode string
}

func (FareData) TableName() string {
	return "fare"
}

type FareDetail struct {
	gorm.Model
	FlowID          uint
	TicketCode      string
	Fare            uint
	RestrictionCode string
	Description     string
	TicketClass     uint
}
