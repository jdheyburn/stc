package models

const (
	Unknown TicketType = iota
	WeeklyStd
	WeeklyZ16TCStd
)

type TicketType int

type TicketCode string

var TicketMappings = map[TicketCode]TicketType{
	"0AQ": WeeklyStd,
	"0AS": WeeklyStd,
	"7DS": WeeklyStd,
	"7TS": WeeklyZ16TCStd,
}
