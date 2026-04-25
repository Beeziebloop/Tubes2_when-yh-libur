# Tubes2_when yh libur

Program ini melakukan traversal pohon DOM untuk mencari node yang cocok dengan CSS selector.

## Penjelasan Singkat Algoritma BFS dan DFS

### BFS (Breadth-First Search)
- Diimplementasikan pada fungsi `BFS`.
- Menggunakan queue (FIFO): node yang masuk duluan akan diproses duluan.
- Traversal dilakukan per level (atas ke bawah, kiri ke kanan).
- Setiap node dicek apakah cocok dengan selector, lalu dicatat ke traversal log.
- Jika parameter `topN` terpenuhi, pencarian berhenti lebih awal.

### DFS (Depth-First Search)
- Diimplementasikan pada fungsi `DFS`.
- Menggunakan stack (LIFO): node terakhir yang dimasukkan diproses lebih dulu.
- Traversal menelusuri cabang sedalam mungkin sebelum pindah ke cabang lain.
- Children didorong ke stack secara terbalik agar urutan kunjungan tetap leftmost-first.
- Sama seperti BFS, pencarian bisa berhenti lebih awal jika `topN` sudah terpenuhi.

## Requirement Program dan Instalasi

### Requirement
- Go versi 1.25.0 (sesuai `go.mod`)
- Koneksi internet (hanya jika input berupa URL, karena program melakukan fetch HTML)

### Dependency
Dependency utama:
- `golang.org/x/net v0.53.0`

Instalasi dependency dilakukan otomatis melalui Go Modules:

```bash
go mod tidy
```

## Langkah Compile / Build / Run Program

Jalankan perintah dari root project.

### 1. Download dependency
```bash
go mod tidy
```

### 2. Menjalankan program langsung
```bash
go run ./src/backend
```

Server akan berjalan di:
- `http://localhost:8080`

### 3. Compile / build executable
```bash
go build -o app.exe ./src/backend
```

Lalu jalankan hasil build:

```bash
./app.exe
```

## Author (Identitas Pembuat)

- Mirza Tsabita Wafa'ana 13524114
- Safira Berlianti 13524128
- Varistha Devi 13524135