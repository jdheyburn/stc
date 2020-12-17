package models

import "time"

type FlowData struct {
	ID              string
	OriginCode      string
	DestinationCode string
	RouteCode       string
	Direction       string
	StartDate       *time.Time
	EndDate         *time.Time
}
