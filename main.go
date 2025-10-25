package main

import (
	"Golang-weekly/lib"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	menu := lib.NewMenu()
	cart := lib.NewCart()
	history := lib.NewHistory()

	fmt.Println("SELAMAT DATANG DI RESTAURANT")

	for {
		fmt.Println("\nPILIH MENU")
		fmt.Println("1. Pesan Menu")
		fmt.Println("2. Lihat Keranjang")
		fmt.Println("3. Checkout")
		fmt.Println("4. History")
		fmt.Println("5. Exit")
		fmt.Print("Pilih: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			menu.Show()
			fmt.Println("\n0. Kembali")
			fmt.Print("Pilih ID menu: ")

			menuID, _ := reader.ReadString('\n')
			menuID = strings.TrimSpace(menuID)

			if menuID == "0" {
				continue
			}

			item, err := menu.Get(menuID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("\n%s\n", strings.ToUpper(item.Name))
			fmt.Printf("Harga: Rp %s\n", lib.FormatCurrency(item.Price))

			fmt.Println("\n1. Tambah ke Keranjang")
			fmt.Println("2. Kembali")
			fmt.Print("Pilih: ")

			action, _ := reader.ReadString('\n')
			action = strings.TrimSpace(action)

			if action == "1" {
				fmt.Print("Masukkan jumlah: ")
				qtyStr, _ := reader.ReadString('\n')
				qtyStr = strings.TrimSpace(qtyStr)

				qty, err := strconv.Atoi(qtyStr)
				if err != nil || qty <= 0 {
					fmt.Println("Jumlah tidak valid")
					continue
				}

				err = cart.Add(menuID, qty, menu)
				if err != nil {
					fmt.Println(err)
				}
			}

		case "2":
			cart.Show()
			fmt.Print("\nTekan Enter untuk kembali...")
			reader.ReadString('\n')

		case "3":
			if cart.IsEmpty() {
				fmt.Println("Keranjang kosong. Silakan pesan menu terlebih dahulu.")
			} else {
				cart.Show()
				total := cart.GetTotal()
				fmt.Printf("\nTotal pembayaran: Rp %s\n", lib.FormatCurrency(total))

				fmt.Println("\n1. Konfirmasi Pembayaran")
				fmt.Println("2. Batal")
				fmt.Print("Pilih: ")

				confirm, _ := reader.ReadString('\n')
				confirm = strings.TrimSpace(confirm)

				if confirm == "1" {
					items := cart.GetItems()
					if len(items) > 0 {
						history.Add(items, total)
						cart.Clear()
						fmt.Println("\nOrder berhasil! Terima kasih atas pesanan Anda.")
					}
				} else {
					fmt.Println("Checkout dibatalkan.")
				}
			}

			fmt.Print("\nTekan Enter untuk kembali...")
			reader.ReadString('\n')

		case "4":
			history.Show()
			fmt.Print("\nTekan Enter untuk kembali...")
			reader.ReadString('\n')

		case "5":
			fmt.Println("Terima kasih!")
			return

		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}
