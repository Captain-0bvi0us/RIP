package repository

import (
	"RIP/internal/app/ds"
	"fmt"
)

func (r *Repository) GetAllFactors() ([]ds.Factors, error) {
	var factors []ds.Factors

	err := r.db.Find(&factors).Error
	if err != nil {
		return nil, err
	}

	if len(factors) == 0 {
		return nil, fmt.Errorf("factors not found")
	}
	return factors, nil
}

func (r *Repository) SearchFactorsByName(title string) ([]ds.Factors, error) {
	var factors []ds.Factors
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&factors).Error // добавили условие
	if err != nil {
		return nil, err
	}
	return factors, nil
}

func (r *Repository) GetFactorByID(id int) (*ds.Factors, error) {
	var factor ds.Factors
	err := r.db.First(&factor, id).Error
	if err != nil {
		return nil, err
	}
	return &factor, nil
}
