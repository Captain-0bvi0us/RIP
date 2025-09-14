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
		// GORM сам возвращает специальную ошибку gorm.ErrRecordNotFound
		return nil, err
	}
	return &factor, nil
}

// func (r *Repository) GetFactor(id int) ([]ds.Factors, error) {
// 	factors, err := r.GetAllFactors()
// 	if err != nil {
// 		return Factor{}, err
// 	}

// 	for _, factor := range factors {
// 		if factor.FactorID == id {
// 			return factor, nil
// 		}
// 	}

// 	return Factor{}, fmt.Errorf("фактор не найден")
// }

// func (r *Repository) GetFactorsByTitle(title string) ([]Factor, error) {
// 	factors, err := r.GetFactors()
// 	if err != nil {
// 		return []Factor{}, err
// 	}

// 	var result []Factor
// 	for _, factor := range factors {
// 		if strings.Contains(strings.ToLower(factor.FactorTitle), strings.ToLower(title)) {
// 			result = append(result, factor)
// 		}
// 	}

// 	return result, nil
// }

// // frax

// type FraxPage struct {
// 	Age          int
// 	Gender       int
// 	Weight       int
// 	Height       int
// 	Factors      []FactorsToFrax
// 	FirstResult  string
// 	SecondResult string
// }

// type FactorsToFrax struct {
// 	Factor      Factor
// 	Description string
// }

// var fraxPages = map[int]FraxPage{
// 	1: {
// 		Age:    56,
// 		Gender: 1,
// 		Weight: 97,
// 		Height: 174,
// 		Factors: []FactorsToFrax{
// 			{Factor: factors[0], Description: "Зависимость продолжается на протяжении 5 месяцев."},
// 			{Factor: factors[1], Description: "Привычка наблюдается на протяжении 6 лет."},
// 			{Factor: factors[2], Description: "Был перелом бедренной кости 7 лет назад."},
// 		},
// 		FirstResult:  "33%",
// 		SecondResult: "24%",
// 	},
// }

// func (r *Repository) GetFraxPage(id int) (FraxPage, error) {
// 	return fraxPages[id], nil
// }
