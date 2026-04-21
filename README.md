# Tubes2_PengenLibur

Project structure
- `api`
- `cmd/api/main.go` : entry point
- `api/v1/` : routes
- `deployment` : folder untuk deploy
- `doc` : laporan
- `docs` : dokumentasi kyk swagger dll
- `internal` : logic dll

notes :
- `main.go` di internal cuma untuk testing :)
- Makefile 

how to run :
jalankan `make` di

how to update the docs:
`swag init-g cmd/api/main.go` untuk updat terus documentationnya