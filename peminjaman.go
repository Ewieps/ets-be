package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type peminjamanInput struct {
	SumberDayaID *int    `json:"sumber_daya_id"`
	Peminjam     *string `json:"peminjam"`
	Tanggal      *string `json:"tanggal"`
	SlotMulai    *int    `json:"slot_mulai"`
	SlotSelesai  *int    `json:"slot_selesai"`
	Jumlah       *int    `json:"jumlah"`
}

func validatePeminjaman(in peminjamanInput) (Peminjaman, error) {
	if in.SumberDayaID == nil || in.Peminjam == nil || in.Tanggal == nil ||
		in.SlotMulai == nil || in.SlotSelesai == nil || in.Jumlah == nil {
		return Peminjaman{}, errors.New("semua field wajib ada")
	}
	peminjam := strings.TrimSpace(*in.Peminjam)
	tanggal := strings.TrimSpace(*in.Tanggal)

	if peminjam == "" {
		return Peminjaman{}, errors.New("peminjam tidak boleh kosong")
	}
	if utf8.RuneCountInString(peminjam) > 50 {
		return Peminjaman{}, errors.New("peminjam maksimal 50 karakter")
	}
	if _, err := time.Parse("2006-01-02", tanggal); err != nil {
		return Peminjaman{}, errors.New("tanggal harus YYYY-MM-DD dan valid")
	}
	m, sel := *in.SlotMulai, *in.SlotSelesai
	if m < 1 || m > 8 || sel < 1 || sel > 8 {
		return Peminjaman{}, errors.New("slot harus antara 1 dan 8")
	}
	if m > sel {
		return Peminjaman{}, errors.New("slot_mulai tidak boleh lebih besar dari slot_selesai")
	}
	return Peminjaman{
		SumberDayaID: *in.SumberDayaID,
		Peminjam:     peminjam,
		Tanggal:      tanggal,
		SlotMulai:    m,
		SlotSelesai:  sel,
		Jumlah:       *in.Jumlah,
	}, nil
}

// POST
func (s *server) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "body bukan JSON yang valid")
		return
	}
	p, err := validatePeminjaman(in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	sd, found := s.cariSumberDayaLocked(p.SumberDayaID)
	if !found {
		writeError(w, http.StatusBadRequest, "sumber_daya_id tidak ditemukan")
		return
	}
	if p.Jumlah < 1 || p.Jumlah > sd.Stok {
		writeError(w, http.StatusBadRequest, "jumlah harus antara 1 dan stok sumber daya")
		return
	}

	p.ID = s.next_PID
	s.next_PID++
	s.peminjaman = append(s.peminjaman, p)
	WriteJSON(w, http.StatusCreated, p)
}

// GET
func (s *server) listPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var sdFilter int
	hasSD := false
	if v := q.Get("sumber_daya_id"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka")
			return
		}
		sdFilter, hasSD = n, true
	}
	tanggal := q.Get("tanggal")
	status := q.Get("status")

	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Peminjaman{}
	for _, p := range s.peminjaman {
		if hasSD && p.SumberDayaID != sdFilter {
			continue
		}
		if tanggal != "" && p.Tanggal != tanggal {
			continue
		}
		if status != "" && p.Status != status {
			continue
		}
		out = append(out, p)
	}
	WriteJSON(w, http.StatusOK, out)
}

// GET (ID)
func (s *server) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.peminjaman {
		if p.ID == id {
			WriteJSON(w, http.StatusOK, p)
			return
		}
	}
	writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
}
