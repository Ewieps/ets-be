package main

import (
	"net/http"
	"strings"
	"time"
)

type ringkasanSD struct {
	SumberDayaID     int    `json:"sumber_daya_id"`
	Nama             string `json:"nama"`
	JumlahPeminjaman int    `json:"jumlah_peminjaman"`
}

type ringkasan struct {
	Tanggal          string        `json:"tanggal"`
	JumlahPeminjaman int           `json:"jumlah_peminjaman"`
	TotalUnitSlot    int           `json:"total_unit_slot"`
	PerSumberDaya    []ringkasanSD `json:"per_sumber_daya"`
}

// GET /ringkasan?tanggal=YYYY-MM-DD
func (s *server) getRingkasan(w http.ResponseWriter, r *http.Request) {
	tanggal := strings.TrimSpace(r.URL.Query().Get("tanggal"))
	if _, err := time.Parse("2006-01-02", tanggal); err != nil { // kosong juga gagal di sini
		writeError(w, http.StatusBadRequest, msgTanggalSalah)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Satu entri per sumber daya (urut id karena slice append-only), mulai dari 0.
	hasil := ringkasan{Tanggal: tanggal, PerSumberDaya: []ringkasanSD{}}
	idx := map[int]int{} // sumber_daya_id -> posisi di PerSumberDaya
	for _, sd := range s.sumberDaya {
		idx[sd.ID] = len(hasil.PerSumberDaya)
		hasil.PerSumberDaya = append(hasil.PerSumberDaya,
			ringkasanSD{SumberDayaID: sd.ID, Nama: sd.Nama})
	}

	for _, p := range s.peminjaman {
		if p.Tanggal != tanggal || !menahanSlot(p.Status) {
			continue
		}
		hasil.JumlahPeminjaman++
		hasil.TotalUnitSlot += p.Jumlah * (p.SlotSelesai - p.SlotMulai + 1)
		if i, ok := idx[p.SumberDayaID]; ok {
			hasil.PerSumberDaya[i].JumlahPeminjaman++
		}
	}

	WriteJSON(w, http.StatusOK, hasil)
}
