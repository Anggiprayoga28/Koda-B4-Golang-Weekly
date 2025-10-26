package main

import (
	"Golang-weekly/lib"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	reader  *bufio.Reader
	menu    *lib.Menu
	cart    *lib.Cart
	history *lib.History
}

func NewApplication() *Application {
	return &Application{
		reader:  bufio.NewReader(os.Stdin),
		menu:    lib.NewMenu(),
		cart:    lib.NewCart(),
		history: lib.NewHistory(),
	}
}

func (app *Application) ask(question string) string {
	fmt.Print(question)
	input, _ := app.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (app *Application) showMainMenu() {
	fmt.Println("\nPILIH MENU")
	fmt.Println("1. Pesan Menu")
	fmt.Println("2. Lihat Keranjang")
	fmt.Println("3. Checkout")
	fmt.Println("4. History")
	fmt.Println("5. Exit")
}

func (app *Application) handleOrder() {
	app.menu.Show()
	fmt.Println("\n0. Kembali")
	menuID := app.ask("Pilih ID menu: ")

	if menuID == "0" {
		return
	}

	item, err := app.menu.Get(menuID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n%s\n", strings.ToUpper(item.Name))
	fmt.Printf("Harga: Rp %s\n", lib.FormatCurrency(item.Price))

	fmt.Println("\n1. Tambah ke Keranjang")
	fmt.Println("2. Kembali")
	action := app.ask("Pilih: ")

	if action == "1" {
		qtyStr := app.ask("Masukkan jumlah: ")
		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty <= 0 {
			fmt.Println("Jumlah tidak valid")
			return
		}

		err = app.cart.Add(menuID, qty, app.menu)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (app *Application) handleCheckout() {
	if app.cart.IsEmpty() {
		fmt.Println("Keranjang kosong. Silakan pesan menu terlebih dahulu.")
	} else {
		app.cart.Show()
		total := app.cart.GetTotal()
		fmt.Printf("\nTotal pembayaran: Rp %s\n", lib.FormatCurrency(total))

		fmt.Println("\n1. Konfirmasi Pembayaran")
		fmt.Println("2. Batal")
		confirm := app.ask("Pilih: ")

		if confirm == "1" {
			items := app.cart.GetItems()
			if len(items) > 0 {
				app.history.Add(items, total)
				app.cart.Clear()
				fmt.Println("\nOrder berhasil! Terima kasih atas pesanan Anda.")
			}
		} else {
			fmt.Println("Checkout dibatalkan.")
		}
	}

	app.ask("\nTekan Enter untuk kembali...")
}

func (app *Application) Run() {
	fmt.Println("SELAMAT DATANG DI RESTAURANT")

	for {
		app.showMainMenu()
		choice := app.ask("Pilih: ")

		switch choice {
		case "1":
			app.handleOrder()
		case "2":
			app.cart.Show()
			app.ask("\nTekan Enter untuk kembali...")
		case "3":
			app.handleCheckout()
		case "4":
			app.history.Show()
			app.ask("\nTekan Enter untuk kembali...")
		case "5":
			fmt.Println("Terima kasih!")
			return
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}

func main() {
	app := NewApplication()
	app.Run()
}
