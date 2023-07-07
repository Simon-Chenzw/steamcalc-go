package data

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MarketPrice struct {
	ItemRefer uint   `gorm:"uniqueIndex:idx_marketprice_itemrefer_source"`
	Source    string `gorm:"uniqueIndex:idx_marketprice_itemrefer_source"`
	Price     float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (db *DB) UpsertMarketPrice(price *MarketPrice) {
	db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "item_refer"},
				{Name: "source"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"price",
				"updated_at",
			}),
		},
	).Create(
		price,
	)
}
