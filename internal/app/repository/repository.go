package repository

type Order struct {
	ID          int
	Title       string
	Description string
	Terrain     string
	Speed       int
	Status      string
	Image       string
	Video       string
	Likes       []int
}

func GetOrders() []Order {
	return []Order{
		{
			ID:          1,
			Title:       "Пехота",
			Description: "Средняя скорость дневного перехода пехоты.",
			Terrain:     "Равнина",
			Speed:       4,
			Status:      "опубликован",
			Image:       "army.png",
			Video:       "army_video.mp4",
			Likes:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			ID:          2,
			Title:       "Кавалерия",
			Description: "Средняя скорость дневного перехода кавалерии.",
			Terrain:     "Равнина",
			Speed:       7,
			Status:      "опубликован",
			Image:       "army_cavalery.png",
			Video:       "army_cavalery.mp4",
			Likes:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35},
		},
		{
			ID:          3,
			Title:       "Обоз",
			Description: "Средняя скорость дневного перехода обоза.",
			Terrain:     "Равнина",
			Speed:       3,
			Status:      "опубликован",
			Image:       "army_oboz.png",
			Video:       "army_oboz.mp4",
			Likes:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45},
		},
		{
			ID:          4,
			Title:       "Артиллерия",
			Description: "Средняя скорость дневного перехода артиллерии.",
			Terrain:     "Равнина",
			Speed:       5,
			Status:      "опубликован",
			Image:       "army_artileri.png",
			Video:       "army_artilery.mp4",
			Likes:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67},
		},
		{
			ID:          5,
			Title:       "",
			Description: "",
			Terrain:     "",
			Speed:       0,
			Status:      "черновик",
			Image:       "",
			Video:       "",
			Likes:       []int{},
		},
		{
			ID:          6,
			Title:       "Удалённое подразделение",
			Description: "",
			Terrain:     "",
			Speed:       0,
			Status:      "удален",
			Image:       "",
			Video:       "",
			Likes:       []int{},
		},
	}
}

func GetOrder(id int) Order {
	orders := GetOrders()

	for _, order := range orders {
		if order.ID == id {
			return order
		}
	}

	return Order{}
}

func GetDraft() Order {
	orders := GetOrders()

	for _, order := range orders {
		if order.Status == "черновик" {
			return order
		}
	}

	return Order{}
}

func GetOrdersBySpeed(speed int) []Order {
	orders := GetOrders()
	result := []Order{}

	for _, order := range orders {
		if order.Status == "опубликован" && order.Speed == speed {
			result = append(result, order)
		}
	}

	return result
}
