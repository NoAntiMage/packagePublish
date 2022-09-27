package model

import "gorm.io/gorm"

//NEXT VERSION when publish package
type PublishPlanBatch struct {
	gorm.Model
	BatchName       string
	AreaInfoID      int
	AreaInfo        AreaInfo
	PublishPlanLogs []*PublishPlanLog
}
