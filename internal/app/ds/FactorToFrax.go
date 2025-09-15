package ds

// FactorToFrax соответствует таблице "FactorToFrax"
type FactorToFrax struct {
	ID          uint    `gorm:"primaryKey;column:id"`
	FraxID      uint    `gorm:"column:frax_id;not null"`   // Внешний ключ к FraxSearching
	FactorID    uint    `gorm:"column:factor_id;not null"` // Внешний ключ к Factors
	Description *string `gorm:"column:description;type:text"`

	// --- СВЯЗИ ---
	// Отношение "принадлежит к" для каждой из связанных таблиц.
	Frax   FraxSearching `gorm:"foreignKey:FraxID"`
	Factor Factors       `gorm:"foreignKey:FactorID"`
}
