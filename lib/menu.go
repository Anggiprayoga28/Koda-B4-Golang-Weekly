package lib

import "fmt"

type MenuItem struct {
	ID    string
	Name  string
	Price int
}

type Menu struct {
	items map[string]*MenuItem
}

func NewMenu() *Menu {
	menu := &Menu{
		items: make(map[string]*MenuItem),
	}
	menu.initializeItems()
	return menu
}

func (m *Menu) initializeItems() {
	menuData := []MenuItem{
		{"1", "Caffè Americano (Tall)", 39000},
		{"2", "Caffè Latte (Tall)", 45000},
		{"3", "Cappuccino (Tall)", 45000},
		{"4", "Caramel Macchiato (Tall)", 55000},
		{"5", "Espresso (Double Shot)", 35000},
		{"6", "Mocha Frappuccino", 58000},
		{"7", "Java Chip Frappuccino", 60000},
		{"8", "Green Tea Latte", 55000},
		{"9", "Signature Chocolate", 52000},
		{"10", "Vanilla Sweet Cream Cold Brew", 56000},
	}

	for _, item := range menuData {
		m.items[item.ID] = &MenuItem{
			ID:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		}
	}
}

func (m *Menu) Show() {
	fmt.Println("\nDAFTAR MENU")
	for i := 1; i <= len(m.items); i++ {
		id := fmt.Sprintf("%d", i)
		if item, exists := m.items[id]; exists {
			fmt.Printf("%s. %s - Rp %s\n", id, item.Name, FormatCurrency(item.Price))
		}
	}
}

func (m *Menu) Get(id string) (*MenuItem, error) {
	if item, exists := m.items[id]; exists {
		return item, nil
	}
	return nil, fmt.Errorf("menu tidak ditemukan")
}
