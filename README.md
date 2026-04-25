# Tubes2_PengenLibur_BE

## Deskripsi Program
Tubes2_PengenLibur_BE adalah aplikasi backend yang menangani proses traversal tree dan pencarian relasi node.

Aplikasi ini terhubung dengan frontend berikut:

- [Tubes2_PengenLibur_FE](https://github.com/ethj0r/Tubes2_PengenLibur_FE)

## Fitur-Fitur

| No. | Fitur | Deskripsi |
| --- | --- | --- |
| 1 | BFS | Menelusuri tree secara bertahap dari level atas ke level bawah |
| 2 | DFS | Menelusuri tree sedalam mungkin sebelum pindah ke cabang lain |
| 3 | Selector Search | Mencari elemen berdasarkan selector seperti tag, class, id, dan kombinasi selector |
| 4 | Traversal Log | Menampilkan urutan node yang dikunjungi dan node yang cocok |
| 5 | LCA | Mencari Lowest Common Ancestor dari dua node yang dipilih |
| 6 | Live Visualization | Menampilkan proses traversal secara visual di layar |

## Syarat Menjalankan Program

Sebelum menjalankan project, pastikan program berikut sudah terpasang:

| No. | Required Program | Keterangan |
| --- | --- | --- |
| 1 | [Go](https://go.dev/dl/) (1.25+) | Untuk menjalankan aplikasi backend dan command pada Makefile |
| 2 | [Git](https://git-scm.com/downloads) | Untuk clone repository |
| 3 | [Make](https://www.gnu.org/software/make/) (opsional) | Untuk shortcut perintah `run`, `build`, `test`, dan `tidy` |
| 4 | [Docker](https://www.docker.com/) (opsional) | Untuk menjalankan backend dalam container |

## Cara Menjalankan Program

Ada dua cara untuk menjalankan program ini.

### Opsi Pertama (Local)

1. Clone repository ini.

```bash
git clone https://github.com/ethj0r/Tubes2_PengenLibur_BE.git
```

2. Masuk ke folder project.

3. Download dependency Go.

```bash
go mod tidy
```

4. Jalankan backend.

Dengan Makefile:

```bash
make run
```

Atau langsung dengan Go:

```bash
go run ./cmd/api
```

5. Akses backend.

```text
http://localhost:8080
```

Catatan konfigurasi:

- Port default: `8080` (bisa diubah via environment variable `PORT`)
- Origin frontend default yang diizinkan CORS: `http://localhost:3000` dan `http://localhost:3001`
- Bisa override origin CORS dengan environment variable `CORS_ALLOWED_ORIGIN`

### Opsi Kedua (Docker)

Build image backend:

```bash
docker build -t pengenlibur-be .
```

Jalankan container:

```bash
docker run --rm -p 8080:8080 --name pengenlibur-be pengenlibur-be
```
Contoh dengan custom port:

```bash
docker run --rm -p 9000:9000 -e PORT=9000 --name pengenlibur-be pengenlibur-be
```
## Kontributor

| NIM | Nama |
| --- | --- |
| 13524026 | Made Branenda Jordhy |
| 13524030 | Irvin Tandiarrang Sumual |
| 13524104 | Valentino Daniel Kusumo |
