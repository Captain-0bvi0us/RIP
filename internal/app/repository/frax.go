package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// GET /api/frax/cart - иконка корзины
func (r *Repository) GetDraftFrax(userID uint) (*ds.FraxSearching, error) {
	var frax ds.FraxSearching
	err := r.db.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&frax).Error
	if err != nil {
		return nil, err
	}
	return &frax, nil
}

// GET /api/frax/cart - иконка корзины
// GET /api/frax/:id - одна заявка с услугами
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

// GET /api/frax - список заявок с фильтрацией
func (r *Repository) FraxListFiltered(userID uint, isModerator bool, status, from, to string) ([]ds.FraxDTO, error) {
	var fraxList []ds.FraxSearching
	query := r.db.Preload("Creator").Preload("Moderator")

	query = query.Where("status != ? AND status != ?", ds.StatusDeleted, ds.StatusDraft)

	if !isModerator {
		query = query.Where("creator_id = ?", userID)
	}

	if status != "" {
		if statusInt, err := strconv.Atoi(status); err == nil {
			query = query.Where("status = ?", statusInt)
		}
	}

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
			CreatorID:      frax.Creator.ID,
			ModeratorID:    nil,
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
			dto.ModeratorID = &frax.Moderator.ID
		}
		result = append(result, dto)
	}
	return result, nil
}

// PUT /api/frax/:id - изменение полей заявки
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

// PUT /api/frax/:id/form - сформировать заявку
func (r *Repository) FormFrax(id uint, creatorID uint) error {
	var frax ds.FraxSearching
	if err := r.db.First(&frax, id).Error; err != nil {
		return err
	}

	if frax.CreatorID != creatorID {
		return errors.New("only creator can form frax")
	}

	if frax.Status != ds.StatusDraft {
		return errors.New("only draft frax can be formed")
	}

	if frax.Age == nil || frax.Gender == nil || frax.Weight == nil || frax.Height == nil {
		return errors.New("age, gender, weight and height are required")
	}

	now := time.Now()
	return r.db.Model(&frax).Updates(map[string]interface{}{
		"status":       ds.StatusFormed,
		"forming_date": now,
	}).Error
}

// DELETE /api/frax/:id - удаление заявки
func (r *Repository) LogicallyDeleteFrax(fraxID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var frax ds.FraxSearching

		if err := tx.Preload("FactorsLink").First(&frax, fraxID).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":       ds.StatusDeleted,
			"forming_date": time.Now(),
		}

		if err := tx.Model(&ds.FraxSearching{}).Where("id = ?", fraxID).Updates(updates).Error; err != nil {
			return err
		}

		var factorIDs []uint
		for _, link := range frax.FactorsLink {
			factorIDs = append(factorIDs, link.FactorID)
		}

		if len(factorIDs) > 0 {
			if err := tx.Model(&ds.Factors{}).Where("id IN ?", factorIDs).Update("status", false).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DELETE /api/frax/:id/factors/:factor_id - удаление фактора из заявки
func (r *Repository) RemoveFactorFromFrax(fraxID, factorID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		result := tx.Where("frax_id = ? AND factor_id = ?", fraxID, factorID).Delete(&ds.FactorToFrax{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("factor not found in this frax")
		}

		if err := tx.Model(&ds.Factors{}).Where("id = ?", factorID).Update("status", false).Error; err != nil {
			return err
		}

		var remainingCount int64
		if err := tx.Model(&ds.FactorToFrax{}).Where("frax_id = ?", fraxID).Count(&remainingCount).Error; err != nil {
			return err
		}

		if remainingCount == 0 {
			updates := map[string]interface{}{
				"status":       ds.StatusDeleted,
				"forming_date": time.Now(),
			}
			if err := tx.Model(&ds.FraxSearching{}).Where("id = ?", fraxID).Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// PUT /api/frax/:id/factors/:factor_id - изменение м-м связи
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

func (r *Repository) UpdateFraxResults(id uint, pof, phf float64) error {
	return r.db.Model(&ds.FraxSearching{}).Where("id = ?", id).Updates(map[string]interface{}{
		"POF": pof,
		"PHF": phf,
	}).Error
}

func (r *Repository) ResolveFrax(id uint, moderatorID uint, action string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var frax ds.FraxSearching
		if err := tx.First(&frax, id).Error; err != nil {
			return err
		}

		if frax.Status != ds.StatusFormed {
			return errors.New("only formed frax can be resolved")
		}

		updates := map[string]interface{}{
			"moderator_id":    moderatorID,
			"complition_date": time.Now(),
		}

		switch action {
		case "complete":
			updates["status"] = ds.StatusCompleted
		case "reject":
			updates["status"] = ds.StatusRejected
		default:
			return errors.New("invalid action")
		}

		if err := tx.Model(&frax).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})
}
