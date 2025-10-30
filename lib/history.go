package lib

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

type Order struct {
	Items []CartItem
	Total int
	Date  string
	Time  string
}

type HistoryInterface interface {
	Add(items []CartItem, total int) bool
	Show()
	Close()
}

type History struct {
	orders     []Order
	maxHistory int
	mu         sync.Mutex
	conn       *pgx.Conn
}

func NewHistory() *History {
	history := &History{
		orders:     make([]Order, 0),
		maxHistory: 5,
	}

	databaseURL := GetDatabaseURL()
	if databaseURL != "" {
		conn, err := pgx.Connect(context.Background(), databaseURL)
		if err != nil {
			fmt.Println("Warning: History tidak bisa connect ke database:", err)
		} else {
			history.conn = conn
			history.loadFromDatabase()
		}
	}

	return history
}

func (h *History) Close() {
	if h.conn != nil {
		h.conn.Close(context.Background())
	}
}

func (h *History) loadFromDatabase() {
	if h.conn == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	query := `SELECT id, total, order_date, order_time 
			  FROM orders 
			  ORDER BY order_date DESC, order_time DESC 
			  LIMIT $1`

	rows, err := h.conn.Query(context.Background(), query, h.maxHistory)
	if err != nil {
		fmt.Println("Error loading history from database:", err)
		return
	}
	defer rows.Close()

	var orderIDs []int
	tempOrders := make([]Order, 0)

	for rows.Next() {
		var id int
		var total int
		var orderDate, orderTime string

		err := rows.Scan(&id, &total, &orderDate, &orderTime)
		if err != nil {
			continue
		}

		orderIDs = append(orderIDs, id)
		tempOrders = append(tempOrders, Order{
			Items: make([]CartItem, 0),
			Total: total,
			Date:  orderDate,
			Time:  orderTime,
		})
	}

	for i, orderID := range orderIDs {
		itemQuery := `SELECT menu_id, name, price, quantity 
					  FROM order_items 
					  WHERE order_id = $1`

		itemRows, err := h.conn.Query(context.Background(), itemQuery, orderID)
		if err != nil {
			continue
		}

		for itemRows.Next() {
			var menuID, name string
			var price, qty int

			err := itemRows.Scan(&menuID, &name, &price, &qty)
			if err != nil {
				continue
			}

			tempOrders[i].Items = append(tempOrders[i].Items, CartItem{
				ID:    menuID,
				Name:  name,
				Price: price,
				Qty:   qty,
			})
		}
		itemRows.Close()
	}

	h.orders = tempOrders
	if len(h.orders) > 0 {
		fmt.Printf("Berhasil memuat %d history dari database\n", len(h.orders))
	}
}

func (h *History) saveToDatabase(order Order) error {
	if h.conn == nil {
		return nil
	}

	var orderID int
	insertOrderQuery := `INSERT INTO orders (total, order_date, order_time) 
						 VALUES ($1, $2, $3) 
						 RETURNING id`

	err := h.conn.QueryRow(context.Background(), insertOrderQuery, order.Total, order.Date, order.Time).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("gagal menyimpan order: %w", err)
	}

	for _, item := range order.Items {
		insertItemQuery := `INSERT INTO order_items (order_id, menu_id, name, price, quantity) 
							VALUES ($1, $2, $3, $4, $5)`

		_, err := h.conn.Exec(context.Background(), insertItemQuery, orderID, item.ID, item.Name, item.Price, item.Qty)
		if err != nil {
			fmt.Printf("Warning: gagal menyimpan item %s: %v\n", item.Name, err)
		}
	}

	return nil
}

func (h *History) Add(items []CartItem, total int) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(items) == 0 || total <= 0 {
		return false
	}

	itemsCopy := make([]CartItem, len(items))
	copy(itemsCopy, items)

	now := time.Now()

	order := Order{
		Items: itemsCopy,
		Total: total,
		Date:  now.Format("02/01/2006"),
		Time:  now.Format("15:04:05"),
	}

	err := h.saveToDatabase(order)
	if err != nil {
		fmt.Println("Warning: gagal menyimpan history ke database:", err)
	}

	h.orders = append([]Order{order}, h.orders...)

	if len(h.orders) > h.maxHistory {
		h.orders = h.orders[:h.maxHistory]
	}

	return true
}

func (h *History) Show() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.orders) == 0 {
		fmt.Println("\nTidak ada riwayat pesanan")
		return
	}

	fmt.Println("\nHISTORY PESANAN")

	for i, order := range h.orders {
		fmt.Printf("\n%d. Tanggal: %s | Waktu: %s\n", i+1, order.Date, order.Time)
		fmt.Printf("   Total: Rp %s\n", FormatCurrency(order.Total))

		if len(order.Items) > 0 {
			fmt.Println("   Item:")
			for _, item := range order.Items {
				fmt.Printf("   - %s x%d\n", item.Name, item.Qty)
			}
		}
	}
}

var _ HistoryInterface = (*History)(nil)
