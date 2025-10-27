package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MenuItem struct {
	ID    string
	Name  string
	Price int
}

type Menu struct {
	items map[string]*MenuItem
}

type ProductAPI struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func NewMenu() *Menu {
	menu := &Menu{
		items: make(map[string]*MenuItem),
	}

	err := menu.fetchFromAPI()
	if err != nil {
		fmt.Println("Warning: Gagal mengambil data dari API:", err)
		fmt.Println("Menggunakan data default")
		menu.initializeItems()
	}

	return menu
}

func (m *Menu) fetchFromAPI() error {
	resp, err := http.Get("https://raw.githubusercontent.com/Anggiprayoga28/Koda-B4-Golang--Weekly-Data/refs/heads/main/dataProduct.json")
	if err != nil {
		return fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var products []ProductAPI
	err = json.Unmarshal(body, &products)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	for _, product := range products {
		id := fmt.Sprintf("%d", product.ID)
		m.items[id] = &MenuItem{
			ID:    id,
			Name:  product.Name,
			Price: product.Price,
		}
	}

	fmt.Printf("Berhasil memuat %d item dari API\n", len(products))
	return nil
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
