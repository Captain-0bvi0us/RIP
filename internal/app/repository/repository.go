package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// factors

type Factors struct {
	ID    int
	Title string
	Text  string
	Image string
}

func (r *Repository) GetFactors() ([]Factors, error) {
	factors := []Factors{
		{
			ID:    1,
			Title: "first order",
			Text:  "first order text",
			Image: "http://localhost:9000/factors/Images/Алкоголизм.png",
		},
		{
			ID:    2,
			Title: "second order",
			Text:  "first order text",
			Image: "http://localhost:9000/factors/Images/Алкоголизм.png",
		},
		{
			ID:    3,
			Title: "third order",
			Text:  "first order text",
			Image: "http://localhost:9000/factors/Images/Алкоголизм.png",
		},
		{
			ID:    4,
			Title: "third order",
			Text:  "first order text",
			Image: "http://localhost:9000/factors/Images/Алкоголизм.png",
		},
	}

	if len(factors) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return factors, nil
}

func (r *Repository) GetFactor(id int) (Factors, error) {
	factors, err := r.GetFactors()
	if err != nil {
		return Factors{}, err
	}

	for _, factor := range factors {
		if factor.ID == id {
			return factor, nil
		}
	}

	return Factors{}, fmt.Errorf("фактор не найден")
}

func (r *Repository) GetFactorsByTitle(title string) ([]Factors, error) {
	factors, err := r.GetFactors()
	if err != nil {
		return []Factors{}, err
	}

	var result []Factors
	for _, factor := range factors {
		if strings.Contains(strings.ToLower(factor.Title), strings.ToLower(title)) {
			result = append(result, factor)
		}
	}

	return result, nil
}

// orders

type Orders struct {
	ID_order int
	Factors  []Factors
}

func (r *Repository) GetOrders() ([]Orders, error) {
	factors, err := r.GetFactors()
	if err != nil {
		return nil, err
	}

	orders := []Orders{
		{
			ID_order: 1,
			Factors: []Factors{
				factors[0],
				factors[1],
				factors[3],
			},
		},
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Orders, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Orders{}, err
	}

	for _, order := range orders {
		if order.ID_order == id {
			return order, nil
		}
	}

	return Orders{}, fmt.Errorf("фактор не найден")
}
