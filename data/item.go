package data

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Item struct {
	ID          uint          `gorm:"primarykey"`
	Appid       uint32        `gorm:"uniqueIndex:idx_item_appid_hashname"`
	HashName    string        `gorm:"uniqueIndex:idx_item_appid_hashname"`
	MarketPrice []MarketPrice `gorm:"foreignKey:ItemRefer"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (db *DB) UpsertItem(item *Item) {
	db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "appid"},
				{Name: "hash_name"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"updated_at",
			}),
		},
	).Omit(
		"market_price",
	).Create(
		item,
	)

	for _, price := range item.MarketPrice {
		price.ItemRefer = item.ID
		db.UpsertMarketPrice(&price)
	}
}
