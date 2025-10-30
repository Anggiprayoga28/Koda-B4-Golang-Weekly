package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5"
)

type MenuItem struct {
	ID    string
	Name  string
	Price int
}

type MenuInterface interface {
	Show()
	Get(id string) (*MenuItem, error)
	Close()
}

type Menu struct {
	items map[string]*MenuItem
	mu    sync.RWMutex
	conn  *pgx.Conn
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

	databaseURL := GetDatabaseURL()
	if databaseURL != "" {
		conn, err := pgx.Connect(context.Background(), databaseURL)
		if err != nil {
			fmt.Println("Warning: Tidak bisa connect ke database:", err)
		} else {
			menu.conn = conn
			fmt.Println("Berhasil terhubung ke database")
		}
	}

	if menu.loadFromDatabase() {
		return menu
	}

	// if menu.loadFromCache() {
	// 	menu.saveToDatabase()
	// 	return menu
	// }

	err := menu.fetchFromAPI()
	if err != nil {
		fmt.Println("Warning: Gagal mengambil data dari API:", err)
		fmt.Println("Menggunakan data default")
		menu.initializeItems()
	} else {
		menu.saveToDatabase()
	}

	return menu
}

func (m *Menu) Close() {
	if m.conn != nil {
		m.conn.Close(context.Background())
	}
}

func (m *Menu) loadFromDatabase() bool {
	if m.conn == nil {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	query := "SELECT id, name, price FROM menu_items ORDER BY id"
	rows, err := m.conn.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error query database:", err)
		return false
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var price int

		err := rows.Scan(&id, &name, &price)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		idStr := fmt.Sprintf("%d", id)
		m.items[idStr] = &MenuItem{
			ID:    idStr,
			Name:  name,
			Price: price,
		}
		count++
	}

	if count == 0 {
		fmt.Println("Database kosong, akan mencoba sumber lain")
		return false
	}

	fmt.Printf("Berhasil memuat %d item dari database\n", count)
	return true
}

func (m *Menu) saveToDatabase() error {
	if m.conn == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	_, err := m.conn.Exec(context.Background(), "DELETE FROM menu_items")
	if err != nil {
		fmt.Println("Warning: gagal menghapus data lama dari database:", err)
		return err
	}

	savedCount := 0
	for _, item := range m.items {
		query := `INSERT INTO menu_items (id, name, price) 
				  VALUES ($1, $2, $3) 
				  ON CONFLICT (id) DO UPDATE 
				  SET name = EXCLUDED.name, price = EXCLUDED.price`

		_, err := m.conn.Exec(context.Background(), query, item.ID, item.Name, item.Price)
		if err != nil {
			fmt.Printf("Warning: gagal menyimpan item %s ke database: %v\n", item.ID, err)
		} else {
			savedCount++
		}
	}

	if savedCount > 0 {
		fmt.Printf("Berhasil menyimpan %d item ke database\n", savedCount)
	}
	return nil
}

/*
func (m *Menu) loadFromCache() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	cacheFilePath := GetCacheFilePath()
	cacheDuration := GetCacheDuration()

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

	type result struct {
		id   string
		item *MenuItem
	}

	resultChan := make(chan result, len(products))
	var wg sync.WaitGroup

	for _, product := range products {
		wg.Add(1)
		go func(p ProductAPI) {
			defer wg.Done()
			id := fmt.Sprintf("%d", p.ID)
			item := &MenuItem{
				ID:    id,
				Name:  p.Name,
				Price: p.Price,
			}
			resultChan <- result{id: id, item: item}
		}(product)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		m.items[res.id] = res.item
	}

	fmt.Printf("Berhasil memuat %d item dari cache (%v)\n", len(products), cacheAge.Round(time.Second))
	return true
}
*/

/*
func (m *Menu) saveToCache(data []byte) error {
	cacheFilePath := GetCacheFilePath()
	err := os.WriteFile(cacheFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("gagal menyimpan cache: %w", err)
	}
	fmt.Println("Data berhasil disimpan ke cache")
	return nil
}
*/

func (m *Menu) fetchFromAPI() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	apiURL := GetAPIURL()

	resp, err := http.Get(apiURL)
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

	// err = m.saveToCache(body)
	// if err != nil {
	// 	fmt.Println("Warning:", err)
	// }

	fmt.Printf("Berhasil memuat %d item dari API\n", len(products))
	return nil
}

func (m *Menu) initializeItems() {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Println("\nDAFTAR MENU")
	for i := 1; i <= len(m.items); i++ {
		id := fmt.Sprintf("%d", i)
		if item, exists := m.items[id]; exists {
			fmt.Printf("%s. %s - Rp %s\n", id, item.Name, FormatCurrency(item.Price))
		}
	}
}

func (m *Menu) Get(id string) (*MenuItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, exists := m.items[id]; exists {
		return item, nil
	}
	return nil, fmt.Errorf("menu tidak ditemukan")
}

var _ MenuInterface = (*Menu)(nil)
