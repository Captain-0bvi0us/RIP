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
			Title: "Алкоголизм",
			Text:  "Алкоголизм приводит к системному нарушению кальциевого обмена и снижению плотности костной ткани. Регулярное употребление алкоголя нарушает работу остеобластов - клеток, отвечающих за образование новой костной ткани, что увеличивает хрупкость костей.",
			Image: "http://localhost:9000/factors/Images/Алкоголизм.png",
		},
		{
			ID:    2,
			Title: "Курение",
			Text:  "Курение значительно ухудшает кровоснабжение костной ткани и замедляет процессы ее восстановления. Никотин и другие токсичные вещества нарушают микроциркуляцию крови, что существенно осложняет заживление переломов и повышает риск будущих повреждений.",
			Image: "http://localhost:9000/factors/Images/Курение.png",
		},
		{
			ID:    3,
			Title: "Предыдущие переломы",
			Text:  "Предыдущий перелом часто свидетельствует о наличии остеопороза или снижении общей прочности костной системы. Даже после полного заживления область бывшего перелома может оставаться уязвимой, а сама травма указывает на предрасположенность к повреждениям костей.",
			Image: "http://localhost:9000/factors/Images/Предыдущие переломы.png",
		},
		{
			ID:    4,
			Title: "Ревматоидный артрит",
			Text:  "Ревматоидный артрит создает хронический воспалительный процесс, который ускоряет потерю костной массы. Аутоиммунный характер заболевания приводит к нарушению нормальной структуры костной ткани, делая ее более хрупкой и susceptible к патологическим переломам.",
			Image: "http://localhost:9000/factors/Images/Ревматоидный артрит.png",
		},
		{
			ID:    5,
			Title: "Вторичный остеопороз",
			Text:  "Вторичный остеопороз — это снижение плотности и прочности костей, которое развивается не как самостоятельное заболевание, а как осложнение или следствие других болезней, состояний или приёма определенных лекарственных препаратов. В отличие от первичного (возрастного или постменопаузального) остеопороза, его причина может быть установлена. Ключевыми триггерами, помимо приёма глюкокортикоидов, являются эндокринные нарушения, ревматоидный артрит, заболевания почек, желудочно-кишечного тракта с нарушением всасывания, а также курение и алкоголизм.",
			Image: "http://localhost:9000/factors/Images/Вторичный остеопороз.png",
		},
		{
			ID:    6,
			Title: "Глюкокортикоиды",
			Text:  "Длительный приём глюкокортикоидов (например, преднизолона) является одной из самых частых причин вторичного остеопороза. Эти препараты напрямую угнетают деятельность остеобластов (клеток, строящих кость) и нарушают всасывание кальция в кишечнике. Кроме того, они ускоряют разрушение костной ткани и могут приводить к быстрой и значительной потере её плотности, dramatically повышая риск патологических переломов, особенно позвонков и шейки бедра.",
			Image: "http://localhost:9000/factors/Images/Глюкокортикоиды.png",
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
