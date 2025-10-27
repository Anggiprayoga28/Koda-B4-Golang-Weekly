package lib

import (
	"fmt"
	"sync"
	"time"
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
}

type History struct {
	orders     []Order
	maxHistory int
	mu         sync.Mutex
}

func NewHistory() *History {
	return &History{
		orders:     make([]Order, 0),
		maxHistory: 5,
	}
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
