package ds

import "time"

// FraxSearching соответствует таблице "FraxSearching".
type FraxSearching struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	Status         int        `gorm:"column:status;not null"`
	CreationDate   time.Time  `gorm:"column:creation_date;not null"`
	CreatorID      uint       `gorm:"column:creator_id;not null"`
	ModeratorID    *uint      `gorm:"column:moderator_id"`
	FormingDate    *time.Time `gorm:"column:forming_date"`
	ComplitionDate *time.Time `gorm:"column:complition_date"`
	Age            *int       `gorm:"column:age"`
	Gender         *bool      `gorm:"column:gender"`
	Weight         *int       `gorm:"column:weight"`
	Height         *int       `gorm:"column:height"`
	POF            *float64   `gorm:"column:POF"`
	PHF            *float64   `gorm:"column:PHF"`

	Creator     Users          `gorm:"foreignKey:CreatorID"`
	Moderator   *Users         `gorm:"foreignKey:ModeratorID"`
	FactorsLink []FactorToFrax `gorm:"foreignKey:FraxID"`
}
