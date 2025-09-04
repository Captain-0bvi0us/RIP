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

type Factors struct { // вот наша новая структура
	ID    int    // поля структур, которые передаются в шаблон
	Title string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	Text  string
}

func (r *Repository) GetFactors() ([]Factors, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	factors := []Factors{ // массив элементов из наших структур
		{
			ID:    1,
			Title: "first order",
			Text:  "first order text",
		},
		{
			ID:    2,
			Title: "second order",
			Text:  "first order text",
		},
		{
			ID:    3,
			Title: "third order",
			Text:  "first order text",
		},
		{
			ID:    4,
			Title: "third order",
			Text:  "first order text",
		},
	}

	if len(factors) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return factors, nil
}

func (r *Repository) GetFactor(id int) (Factors, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	factors, err := r.GetFactors()
	if err != nil {
		return Factors{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, factor := range factors {
		if factor.ID == id {
			return factor, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Factors{}, fmt.Errorf("фактор не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
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

type Orders struct { // вот наша новая структура
	ID_order  int
	ID_factor int    // поля структур, которые передаются в шаблон
	Title     string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	Text      string
}

func (r *Repository) GetOrders() ([]Orders, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	orders := []Orders{ // массив элементов из наших структур
		{
			ID_order:  1,
			ID_factor: 1,
			Title:     "first order",
			Text:      "first order text",
		},
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Orders, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	orders, err := r.GetOrders()
	if err != nil {
		return Orders{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, order := range orders {
		if order.ID_order == id {
			return order, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Orders{}, fmt.Errorf("фактор не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}
