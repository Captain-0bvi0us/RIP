package ds

import "time"

// FraxSearching соответствует таблице "FraxSearching".
type FraxSearching struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	Status         int        `gorm:"column:status;not null"`
	CreationDate   time.Time  `gorm:"column:creation_date;not null"`
	CreatorID      uint       `gorm:"column:creator_id;not null"` // Внешний ключ
	Moderator      *bool      `gorm:"column:moderator"`
	FormingDate    *time.Time `gorm:"column:forming_date"`
	ComplitionDate *time.Time `gorm:"column:complition_date"`
	Age            *int       `gorm:"column:age"`
	Gender         *bool      `gorm:"column:gender"`
	Weight         *int       `gorm:"column:weight"`
	Height         *int       `gorm:"column:height"`
	POF            *float64   `gorm:"column:POF"`
	PHF            *float64   `gorm:"column:PHF"`

	// --- СВЯЗИ ---
	// Отношение "принадлежит к": каждая сессия принадлежит одному пользователю.
	Creator Users `gorm:"foreignKey:CreatorID"`
	// Отношение "один-ко-многим" к связующей таблице:
	// У одной сессии может быть много записей-факторов.
	FactorsLink []FactorToFrax `gorm:"foreignKey:FraxID"`
}
