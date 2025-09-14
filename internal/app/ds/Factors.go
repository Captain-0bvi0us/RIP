package ds

// Factors соответствует таблице "Factors"
// Это справочник всех возможных факторов риска (например, "Курение", "Алкоголь" и т.д.).
// Factors соответствует таблице "Factors".
type Factors struct {
	ID       uint     `gorm:"primaryKey;column:id"`
	Title    string   `gorm:"column:title;size:255;not null"`
	Text     string   `gorm:"column:text;not null"`
	Image    *string  `gorm:"column:image;size:255"`
	Argument *float64 `gorm:"column:argument"`
	Status   *bool    `gorm:"column:status"`

	// --- СВЯЗИ ---
	// Отношение "один-ко-многим" к связующей таблице:
	// Один фактор может быть использован во многих сессиях.
	FraxLinks []FactorToFrax `gorm:"foreignKey:FactorID"`
}
