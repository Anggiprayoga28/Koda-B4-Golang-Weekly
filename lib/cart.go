package lib

import (
	"fmt"
	"sync"
)

type CartItem struct {
	ID    string
	Name  string
	Price int
	Qty   int
}

type CartInterface interface {
	Add(id string, qty int, menu MenuInterface) error
	Show()
	IsEmpty() bool
	GetTotal() int
	GetItems() []CartItem
	Clear()
}

type Cart struct {
	items []CartItem
	mu    sync.Mutex
}

func NewCart() *Cart {
	return &Cart{
		items: make([]CartItem, 0),
	}
}

func (c *Cart) Add(id string, qty int, menu MenuInterface) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if qty <= 0 {
		return fmt.Errorf("jumlah harus lebih dari 0")
	}

	menuItem, err := menu.Get(id)
	if err != nil {
		return err
	}

	for i := range c.items {
		if c.items[i].ID == id {
			c.items[i].Qty += qty
			fmt.Printf("%s ditambahkan. Total: x%d\n", menuItem.Name, c.items[i].Qty)
			return nil
		}
	}

	c.items = append(c.items, CartItem{
		ID:    id,
		Name:  menuItem.Name,
		Price: menuItem.Price,
		Qty:   qty,
	})

	fmt.Printf("%s x%d ditambahkan ke keranjang\n", menuItem.Name, qty)
	return nil
}

func (c *Cart) Show() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) == 0 {
		fmt.Println("\nKeranjang kosong")
		return
	}

	fmt.Println("\nISI KERANJANG")
	total := 0

	for i, item := range c.items {
		subtotal := item.Price * item.Qty
		fmt.Printf("%d. %s\n", i+1, item.Name)
		fmt.Printf("   %d x Rp %s = Rp %s\n", item.Qty, FormatCurrency(item.Price), FormatCurrency(subtotal))
		total += subtotal
	}

	fmt.Printf("TOTAL: Rp %s\n", FormatCurrency(total))
}

func (c *Cart) IsEmpty() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items) == 0
}

func (c *Cart) GetTotal() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) == 0 {
		return 0
	}

	resultChan := make(chan int, len(c.items))
	var wg sync.WaitGroup

	for _, item := range c.items {
		wg.Add(1)
		go func(i CartItem) {
			defer wg.Done()
			resultChan <- i.Price * i.Qty
		}(item)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	total := 0
	for subtotal := range resultChan {
		total += subtotal
	}

	return total
}

func (c *Cart) GetItems() []CartItem {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make([]CartItem, len(c.items))
	copy(result, c.items)
	return result
}

func (c *Cart) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make([]CartItem, 0)
}

var _ CartInterface = (*Cart)(nil)
