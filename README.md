# Golang Weekly - CLI Restaurant / Starbuck Ordering System

Aplikasi **Point of Sale (POS)** berbasis **Command Line Interface (CLI)** menggunakan bahasa **Go (Golang)**.  
Aplikasi ini memungkinkan pengguna untuk:

- Melihat daftar menu
- Menambahkan item ke keranjang
- Melakukan checkout
- Melihat riwayat transaksi
- Menghapus cache menu
- Mendapatkan data menu dari API atau cache (dengan fallback ke data default)

---

## Fitur Utama

| Fitur                | Deskripsi |
|---------------------|-----------|
| Menu Dinamis        | Menu diambil dari API lalu disimpan ke cache lokal. |
| Cache System        | Mempercepat loading menu dengan penyimpanan file `.json`. |
| Keranjang Belanja   | Tambah dan lihat item sebelum checkout. |
| Checkout            | Hitung total harga dan simpan ke riwayat. |
| Riwayat Transaksi   | Menyimpan hingga beberapa transaksi terakhir. |
| Clear Cache         | Menghapus file cache jika diperlukan. |
| Konfigurasi `.env`  | Base URL API & pengaturan cache fleksibel. |


## Instalasi & Penggunaan

### 1. Clone Repository
```bash
git clone https://github.com/yourusername/Golang-weekly.git
cd Golang-weekly
```

### 2. Buat .env

```bash
API_URL=https://example.com/products.json
CACHE_DURATION=15m
CACHE_FILE_PATH=/tmp/menu_cache.json
```
Jangan lupa .env sudah masuk ke .gitignore.

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Jalankan Program
```bash
go run main.go
```

## Hasil
**Tampilan Utama**

![Menu Utama](/public/utama.png)
<br>

**Tampilan Keranjang**

![Tampilan Keranjang](/public/keranjang.png)
<br>

**Tampilan Checkout**

![Tampilan Checkout](/public/checkout.png)
<br>

**Tampilan History**

![Tampilan History](/public/history.png)
<br>

**Tampilan Hapus Chace**

![Tampilan Hapus Chace](/public/keranjang.png)