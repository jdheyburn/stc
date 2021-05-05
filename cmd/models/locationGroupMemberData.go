package models

import (
	"time"

	"gorm.io/gorm"
)

type LocationGroupMemberData struct {
	gorm.Model
	GroupUICCode  string
	MemberUICCode string
	MemberCRSCode string
	StartDate     *time.Time
	EndDate       *time.Time
}

func (LocationGroupMemberData) TableName() string {
	return "location_group_member"
}
