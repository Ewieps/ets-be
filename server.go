package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

// Objek peminjaman:
// {
//   "id": 1,
//   "sumber_daya_id": 2,
//   "peminjam": "Rina",
//   "tanggal": "2026-11-02",
//   "slot_mulai": 3,
//   "slot_selesai": 4,
//   "jumlah": 1,
//   "status": "aktif"
// }

//Objek sumber daya:
//{"id": 1, "nama": "Aula A", "jenis": "ruang", "stok": 1}

type server struct {
	mu         sync.Mutex
	sumberDaya []SumberDaya
	peminjaman []Peminjaman
	next_SDID  int
	next_PID   int
}

type errorBody struct {
	Error string `json:"error"`
}

func NewHandler() http.Handler {
	s := &server{next_SDID: 1, next_PID: 1}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sumber-daya", s.createSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.listSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.getSumberDaya)

	mux.HandleFunc("POST /peminjaman", s.createPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.listPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.getPeminjaman)

	return mux
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// TODO: func errror
func writeError(w http.ResponseWriter, status int, errmsg string) {
	WriteJSON(w, status, errorBody{Error: errmsg})
}

// TODO: parseID
func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return 0, false
	}
	return id, true
}

// func (s *server) listSumberDaya(w http.ResponseWriter, r *http.Request) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	WriteJSON(w, http.StatusOK, s.sumberDaya)
// }
