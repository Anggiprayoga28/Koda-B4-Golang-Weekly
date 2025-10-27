package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

const (
	cacheFilePath = "/tmp/menu_cache.json"
	cacheDuration = 15 * time.Second
)

func NewMenu() *Menu {
	menu := &Menu{
		items: make(map[string]*MenuItem),
	}

	if menu.loadFromCache() {
		return menu
	}

	err := menu.fetchFromAPI()
	if err != nil {
		fmt.Println("Warning: Gagal mengambil data dari API:", err)
		fmt.Println("Menggunakan data default")
		menu.initializeItems()
	}

	return menu
}

func (m *Menu) loadFromCache() bool {
	fileInfo, err := os.Stat(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Cache tidak ditemukan, akan fetch dari API")
			return false
		}
		fmt.Println("Error membaca cache:", err)
		return false
	}

	cacheAge := time.Since(fileInfo.ModTime())
	if cacheAge >= cacheDuration {
		fmt.Printf("Cache sudah expired (%v), akan di fetch ulang\n", cacheAge.Round(time.Second))
		return false
	}

	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		fmt.Println("Error membaca file cache:", err)
		return false
	}

	var products []ProductAPI
	err = json.Unmarshal(data, &products)
	if err != nil {
		fmt.Println("Error parsing cache JSON:", err)
		return false
	}

	for _, product := range products {
		id := fmt.Sprintf("%d", product.ID)
		m.items[id] = &MenuItem{
			ID:    id,
			Name:  product.Name,
			Price: product.Price,
		}
	}

	fmt.Printf("Berhasil memuat %d item dari cache (%v)\n", len(products), cacheAge.Round(time.Second))
	return true
}

func (m *Menu) saveToCache(data []byte) error {
	err := os.WriteFile(cacheFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("gagal menyimpan cache: %w", err)
	}
	fmt.Println("Data berhasil disimpan ke cache")
	return nil
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

	err = m.saveToCache(body)
	if err != nil {
		fmt.Println("Warning:", err)
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
