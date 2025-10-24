package lib

import "fmt"

type CartItem struct {
	ID    string
	Name  string
	Price int
	Qty   int
}

type Cart struct {
	items []CartItem
}

func NewCart() *Cart {
	return &Cart{
		items: make([]CartItem, 0),
	}
}

func (c *Cart) Add(id string, qty int, menu *Menu) error {
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
	return len(c.items) == 0
}

func (c *Cart) GetTotal() int {
	total := 0
	for _, item := range c.items {
		total += item.Price * item.Qty
	}
	return total
}

func (c *Cart) GetItems() []CartItem {
	result := make([]CartItem, len(c.items))
	copy(result, c.items)
	return result
}

func (c *Cart) Clear() {
	c.items = make([]CartItem, 0)
}
