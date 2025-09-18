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

type Factor struct {
	FactorID    int
	FactorTitle string
	FactorText  string
	FactorImage string
}

var factors = []Factor{
	{
		FactorID:    1,
		FactorTitle: "Алкоголизм",
		FactorText:  "Алкоголизм приводит к системному нарушению кальциевого обмена и снижению плотности костной ткани. Регулярное употребление алкоголя нарушает работу остеобластов - клеток, отвечающих за образование новой костной ткани, что увеличивает хрупкость костей.",
		FactorImage: "http://localhost:9000/factors/Images/Алкоголизм.png",
	},
	{
		FactorID:    2,
		FactorTitle: "Курение",
		FactorText:  "Курение значительно ухудшает кровоснабжение костной ткани и замедляет процессы ее восстановления. Никотин и другие токсичные вещества нарушают микроциркуляцию крови, что существенно осложняет заживление переломов и повышает риск будущих повреждений.",
		FactorImage: "http://localhost:9000/factors/Images/Курение.png",
	},
	{
		FactorID:    3,
		FactorTitle: "Предыдущие переломы",
		FactorText:  "Предыдущий перелом часто свидетельствует о наличии остеопороза или снижении общей прочности костной системы. Даже после полного заживления область бывшего перелома может оставаться уязвимой, а сама травма указывает на предрасположенность к повреждениям костей.",
		FactorImage: "http://localhost:9000/factors/Images/Предыдущие переломы.png",
	},
	{
		FactorID:    4,
		FactorTitle: "Ревматоидный артрит",
		FactorText:  "Ревматоидный артрит создает хронический воспалительный процесс, который ускоряет потерю костной массы. Аутоиммунный характер заболевания приводит к нарушению нормальной структуры костной ткани, делая ее более хрупкой и susceptible к патологическим переломам.",
		FactorImage: "http://localhost:9000/factors/Images/Ревматоидный артрит.png",
	},
	{
		FactorID:    5,
		FactorTitle: "Вторичный остеопороз",
		FactorText:  "Вторичный остеопороз — это снижение плотности и прочности костей, которое развивается не как самостоятельное заболевание, а как осложнение или следствие других болезней, состояний или приёма определенных лекарственных препаратов. В отличие от первичного (возрастного или постменопаузального) остеопороза, его причина может быть установлена. Ключевыми триггерами, помимо приёма глюкокортикоидов, являются эндокринные нарушения, ревматоидный артрит, заболевания почек, желудочно-кишечного тракта с нарушением всасывания, а также курение и алкоголизм.",
		FactorImage: "http://localhost:9000/factors/Images/Вторичный остеопороз.png",
	},
	{
		FactorID:    6,
		FactorTitle: "Глюкокортикоиды",
		FactorText:  "Длительный приём глюкокортикоидов (например, преднизолона) является одной из самых частых причин вторичного остеопороза. Эти препараты напрямую угнетают деятельность остеобластов (клеток, строящих кость) и нарушают всасывание кальция в кишечнике. Кроме того, они ускоряют разрушение костной ткани и могут приводить к быстрой и значительной потере её плотности, dramatically повышая риск патологических переломов, особенно позвонков и шейки бедра.",
		FactorImage: "http://localhost:9000/factors/Images/Глюкокортикоиды.png",
	},
}

func (r *Repository) GetFactors() ([]Factor, error) {
	if len(factors) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return factors, nil
}

func (r *Repository) GetFactor(id int) (Factor, error) {
	factors, err := r.GetFactors()
	if err != nil {
		return Factor{}, err
	}

	for _, factor := range factors {
		if factor.FactorID == id {
			return factor, nil
		}
	}

	return Factor{}, fmt.Errorf("фактор не найден")
}

func (r *Repository) GetFactorsByTitle(title string) ([]Factor, error) {
	factors, err := r.GetFactors()
	if err != nil {
		return []Factor{}, err
	}

	var result []Factor
	for _, factor := range factors {
		if strings.Contains(strings.ToLower(factor.FactorTitle), strings.ToLower(title)) {
			result = append(result, factor)
		}
	}

	return result, nil
}

// frax

type FraxPage struct {
	Age     int
	Gender  int
	Weight  int
	Height  int
	Factors []FactorsToFrax
	POF     string
	PHF     string
}

type FactorsToFrax struct {
	Factor      Factor
	Description string
}

var fraxPages = map[int]FraxPage{
	1: {
		Age:    56,
		Gender: 1,
		Weight: 97,
		Height: 174,
		Factors: []FactorsToFrax{
			{Factor: factors[0], Description: "Зависимость продолжается на протяжении 5 месяцев."},
			{Factor: factors[1], Description: "Привычка наблюдается на протяжении 6 лет."},
			{Factor: factors[2], Description: "Был перелом бедренной кости 7 лет назад."},
		},
		POF: "33%",
		PHF: "24%",
	},
}

func (r *Repository) GetFrax(id int) (FraxPage, error) {
	return fraxPages[id], nil
}
