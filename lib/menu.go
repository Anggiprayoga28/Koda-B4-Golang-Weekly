package lib

import (
	"context"
	"fmt"
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

func NewMenu() *Menu {
	menu := &Menu{
		items: make(map[string]*MenuItem),
	}

	databaseURL := GetDatabaseURL()
	if databaseURL == "" {
		fmt.Println("Error: DATABASE_URL tidak ditemukan")
		return menu
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Println("Error: Tidak bisa connect ke database:", err)
		return menu
	}

	menu.conn = conn
	menu.loadFromDatabase()
	return menu
}

func (m *Menu) Close() {
	if m.conn != nil {
		m.conn.Close(context.Background())
	}
}

func (m *Menu) loadFromDatabase() error {
	if m.conn == nil {
		return fmt.Errorf("koneksi database tidak tersedia")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	query := "SELECT id, name, price FROM menu_items ORDER BY id"
	rows, err := m.conn.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error query database:", err)
		return err
	}
	defer rows.Close()

	m.items = make(map[string]*MenuItem)
	count := 0

	for rows.Next() {
		var id int
		var name string
		var price int

		err := rows.Scan(&id, &name, &price)
		if err != nil {
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
		fmt.Println("Tidak ada data menu di database!")
		fmt.Println("Silakan insert data menu terlebih dahulu.")
		return fmt.Errorf("tidak ada data menu")
	}

	return nil
}

func (m *Menu) Show() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.items) == 0 {
		fmt.Println("\nTidak ada menu tersedia")
		return
	}

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
