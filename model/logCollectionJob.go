package model

import (
	"time"

	"gorm.io/gorm"
)

type LogCollectionJob struct {
	gorm.Model
	JobName   string
	LogDate   time.Time `gorm:"type: date"`
	JobStatus string
}
