package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"strconv"
	"time"
)

func (r *Repository) GetDraftFrax(userID uint) (*ds.FraxSearching, error) {
	var frax ds.FraxSearching
	err := r.db.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&frax).Error
	if err != nil {
		return nil, err
	}
	return &frax, nil
}

func (r *Repository) CreateFrax(frax *ds.FraxSearching) error {
	return r.db.Create(frax).Error
}

func (r *Repository) GetFraxWithFactors(fraxID uint) (*ds.FraxSearching, error) {
	var frax ds.FraxSearching
	err := r.db.Preload("FactorsLink.Factor").Preload("Creator").Preload("Moderator").First(&frax, fraxID).Error
	if err != nil {
		return nil, err
	}

	if frax.Status == ds.StatusDeleted {
		return nil, errors.New("frax page not found or has been deleted")
	}

	return &frax, nil
}

func (r *Repository) LogicallyDeleteFrax(fraxID uint) error {
	result := r.db.Exec("UPDATE frax_searchings SET status = ? WHERE id = ?", ds.StatusDeleted, fraxID)
	return result.Error
}

// Методы для работы с FRAX

func (r *Repository) FraxListFiltered(status, from, to string) ([]ds.FraxDTO, error) {
	var fraxList []ds.FraxSearching
	query := r.db.Preload("Creator").Preload("Moderator")

	// Исключаем удаленные и черновики
	query = query.Where("status != ? AND status != ?", ds.StatusDeleted, ds.StatusDraft)

	// Фильтрация по статусу
	if status != "" {
		if statusInt, err := strconv.Atoi(status); err == nil {
			query = query.Where("status = ?", statusInt)
		}
	}

	// Фильтрация по дате формирования
	if from != "" {
		if fromTime, err := time.Parse("2006-01-02", from); err == nil {
			query = query.Where("forming_date >= ?", fromTime)
		}
	}
	if to != "" {
		if toTime, err := time.Parse("2006-01-02", to); err == nil {
			query = query.Where("forming_date <= ?", toTime)
		}
	}

	if err := query.Find(&fraxList).Error; err != nil {
		return nil, err
	}

	var result []ds.FraxDTO
	for _, frax := range fraxList {
		dto := ds.FraxDTO{
			ID:             frax.ID,
			Status:         frax.Status,
			CreationDate:   frax.CreationDate,
			CreatorLogin:   frax.Creator.Username,
			ModeratorLogin: nil,
			FormingDate:    frax.FormingDate,
			ComplitionDate: frax.ComplitionDate,
			Age:            frax.Age,
			Gender:         frax.Gender,
			Weight:         frax.Weight,
			Height:         frax.Height,
			POF:            frax.POF,
			PHF:            frax.PHF,
		}

		if frax.ModeratorID != nil {
			dto.ModeratorLogin = &frax.Moderator.Username
		}

		result = append(result, dto)
	}

	return result, nil
}

func (r *Repository) UpdateFraxUserFields(id uint, req ds.FraxUpdateRequest) error {
	updates := make(map[string]interface{})

	if req.Age != nil {
		updates["age"] = *req.Age
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Height != nil {
		updates["height"] = *req.Height
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&ds.FraxSearching{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) FormFrax(id uint, creatorID uint) error {
	var frax ds.FraxSearching
	if err := r.db.First(&frax, id).Error; err != nil {
		return err
	}

	// Проверяем права доступа
	if frax.CreatorID != creatorID {
		return errors.New("only creator can form frax")
	}

	// Проверяем статус
	if frax.Status != ds.StatusDraft {
		return errors.New("only draft frax can be formed")
	}

	// Проверяем обязательные поля
	if frax.Age == nil || frax.Gender == nil || frax.Weight == nil || frax.Height == nil {
		return errors.New("age, gender, weight and height are required")
	}

	// Обновляем статус и дату формирования
	now := time.Now()
	return r.db.Model(&frax).Updates(map[string]interface{}{
		"status":       ds.StatusFormed,
		"forming_date": now,
	}).Error
}

func (r *Repository) ResolveFrax(id uint, moderatorID uint, action string) error {
	var frax ds.FraxSearching
	if err := r.db.Preload("FactorsLink.Factor").First(&frax, id).Error; err != nil {
		return err
	}

	// Проверяем статус
	if frax.Status != ds.StatusFormed {
		return errors.New("only formed frax can be resolved")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"moderator_id":    moderatorID,
		"complition_date": now,
	}

	if action == "complete" {
		updates["status"] = ds.StatusCompleted
		// Рассчитываем POF и PHF
		pof, phf := r.calculateFRAX(frax)
		updates["POF"] = pof
		updates["PHF"] = phf
	} else if action == "reject" {
		updates["status"] = ds.StatusRejected
	} else {
		return errors.New("invalid action, must be 'complete' or 'reject'")
	}

	return r.db.Model(&frax).Updates(updates).Error
}

func (r *Repository) RemoveFactorFromFrax(fraxID, factorID uint) error {
	return r.db.Where("frax_id = ? AND factor_id = ?", fraxID, factorID).Delete(&ds.FactorToFrax{}).Error
}

func (r *Repository) UpdateMM(fraxID, factorID uint, updateData ds.FactorToFrax) error {
	var link ds.FactorToFrax
	if err := r.db.Where("frax_id = ? AND factor_id = ?", fraxID, factorID).First(&link).Error; err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if updateData.Description != nil {
		updates["description"] = *updateData.Description
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&link).Updates(updates).Error
}

// Упрощенная формула расчета FRAX (идеализированная)
func (r *Repository) calculateFRAX(frax ds.FraxSearching) (float64, float64) {
	// Базовые коэффициенты (примерные)
	age := float64(*frax.Age)
	gender := 0.0
	if *frax.Gender {
		gender = 1.0
	}

	// BMI
	bmi := float64(*frax.Weight) / ((float64(*frax.Height) / 100) * (float64(*frax.Height) / 100))

	// Сумма факторов риска
	factorSum := 0.0
	for _, link := range frax.FactorsLink {
		if link.Factor.Argument != nil {
			factorSum += *link.Factor.Argument
		}
	}

	// Упрощенная формула для POF (основные остеопатические переломы)
	pof := 0.1*age + 0.2*gender + 0.05*bmi + 0.3*factorSum + 0.5

	// Упрощенная формула для PHF (перелом шейки бедра)
	phf := 0.05*age + 0.15*gender + 0.03*bmi + 0.2*factorSum + 0.2

	// Ограничиваем значения от 0 до 100
	if pof < 0 {
		pof = 0
	}
	if pof > 100 {
		pof = 100
	}
	if phf < 0 {
		phf = 0
	}
	if phf > 100 {
		phf = 100
	}

	return pof, phf
}
