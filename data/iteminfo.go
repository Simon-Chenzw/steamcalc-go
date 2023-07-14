package data

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemInfo struct {
	ItemID     int  `gorm:"primarykey"`
	Item       Item `gorm:"foreignKey:ItemID"`
	ItemNameid int
	LocaleName string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (db *DB) UpsertItemInfo(info *ItemInfo) {
	db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "item_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"updated_at",
				"item_nameid",
				"locale_name",
			}),
		},
	).Create(
		info,
	)
}
