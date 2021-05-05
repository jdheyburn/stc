package models

import (
	"gorm.io/gorm"
)

// FareData maps directly to rows in the fare table
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

// FareDetail contains the base of FareData, with additional columns LEFT JOINed from ticket_type table
type FareDetail struct {
	gorm.Model
	FlowID          uint
	TicketCode      string
	Fare            uint
	RestrictionCode string
	// LEFT JOINed from ticket_type
	TicketDescription string
	TicketClass       uint
	TicketType        string
	// LEFT JOINed from restriction_header
	RestrictionDesc    string
	RestrictionDescOut string
	RestrictionDescRtn string
}
