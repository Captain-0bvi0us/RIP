package repository

import (
	"RIP/internal/app/ds"
	"errors"
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

func (r *Repository) AddFactorToFrax(fraxID, factorID uint) error {
	var count int64

	r.db.Model(&ds.FactorToFrax{}).Where("frax_id = ? AND factor_id = ?", fraxID, factorID).Count(&count)
	if count > 0 {
		return errors.New("factor already in frax")
	}

	link := ds.FactorToFrax{
		FraxID:   fraxID,
		FactorID: factorID,
	}
	return r.db.Create(&link).Error
}

func (r *Repository) GetFraxWithFactors(fraxID uint) (*ds.FraxSearching, error) {
	var frax ds.FraxSearching

	err := r.db.Preload("FactorsLink.Factor").First(&frax, fraxID).Error
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
