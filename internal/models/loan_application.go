package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanApplication struct {
	gorm.Model
	UUID                     *uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex"`
	IncomingOrganizationUUID *uuid.UUID    `gorm:"type:uuid;not null;index"`
	IssueOrganizationUUID    *uuid.UUID    `gorm:"type:uuid;not null;index"`
	IncomingOrganization     *Organization `gorm:"foreignKey:IncomingOrganizationUUID;references:UUID"`
	IssueOrganization        *Organization `gorm:"foreignKey:IssueOrganizationUUID;references:UUID"`
	Value                    int64         `gorm:"not null;check:value>=10000"`
	Phone                    string        `gorm:"not null; size:20"`
	Comment                  string        `gorm:"type:text"`
}
