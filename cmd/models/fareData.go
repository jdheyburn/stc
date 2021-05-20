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
	FlowID          uint   `header:"flowId"`
	TicketCode      string `header:"ticket_code"`
	Fare            uint   `header:"fare"`
	RestrictionCode string `header:"restriction_code"`
	// LEFT JOINed from ticket_type
	TicketDescription string `header:"tkt_desc"`
	TicketClass       uint   `header:"tkt_class"`
	TicketType        string `header:"tkt_type"`
	// LEFT JOINed from restriction_header
	RestrictionDesc    string `header:"restriction_desc"`
	RestrictionDescOut string `header:"restriction_desc_out"`
	RestrictionDescRtn string `header:"restriction_desc_rtn"`
}

// FareDetailExtreme has all the fields I think we'll need to hit the db just once
type FareDetailExtreme struct {
	gorm.Model
	FlowID          uint   `header:"flow_id"`
	OriginCode      string `header:"origin_code"`
	OriginName      string `header:"origin_name"`
	DestinationCode string `header:"destination_code"`
	DestinationName string `header:"destination_name"`
	RouteCode       string `header:"route_code"`
	RouteDesc       string `header:"route_desc"`
	RouteAaaDesc    string `header:"route_aaa_desc"`
	StatusCode      string `header:"status_code"`
	UsageCode       string `header:"usage_code"`
	TOC             string `header:"toc"`
	FareID          string `header:"fare_id"`
	TicketCode      string `header:"ticket_code"`
	TicketDesc      string `header:"tkt_desc"`
	TicketClass     uint   `header:"tkt_class"`
	TicketType      string `header:"tkt_type"`
	AdultFare       uint   `header:"adult_fare"`
	ChildFare       uint  `header:"child_fare"`
	RestrictionCode string `header:"restriction_code"`
	RestrictionDesc string `header:"restriction_desc"`
}
