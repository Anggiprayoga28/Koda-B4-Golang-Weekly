package main

import (
	"Golang-weekly/lib"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	menu := lib.NewMenu()

	fmt.Println("\nSelamat Datang di Starbuck")

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
			} else {
				fmt.Printf("\n%s\n", strings.ToUpper(item.Name))
				fmt.Printf("Harga: Rp %s\n", lib.FormatCurrency(item.Price))
				fmt.Println("Tambahkan ke keranjang")
			}

		case "2":
			fmt.Println("Lihat Keranjang")
		case "3":
			fmt.Println("Checkout")
		case "4":
			fmt.Println("History")
		case "5":
			fmt.Println("Terima kasih")
			return
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}
