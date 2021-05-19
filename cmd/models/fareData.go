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
	FlowID          uint `header:"flowId"`
	TicketCode      string `header:"ticket_code"`
	Fare            uint `header:"fare"`
	RestrictionCode string `header:"restriction_code"`
	// LEFT JOINed from ticket_type
	TicketDescription string `header:"tkt_desc"`
	TicketClass       uint `header:"tkt_class"`
	TicketType        string `header:"tkt_type"`
	// LEFT JOINed from restriction_header
	RestrictionDesc    string `header:"restriction_desc"`
	RestrictionDescOut string `header:"restriction_desc_out"`
	RestrictionDescRtn string `header:"restriction_desc_rtn"`
}
