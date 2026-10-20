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
	Status *string `json:"status"`
	Alasan *string `json:"alasan"`

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
		return Peminjaman{}, errors.New(msgFieldPinjamWajib)
	}
	peminjam := strings.TrimSpace(*in.Peminjam)
	tanggal := strings.TrimSpace(*in.Tanggal)

	if peminjam == "" {
		return Peminjaman{}, errors.New(msgPeminjamKosong)
	}
	if utf8.RuneCountInString(peminjam) > 50 {
		return Peminjaman{}, errors.New(msgPeminjamPanjang)
	}
	if _, err := time.Parse("2006-01-02", tanggal); err != nil {
		return Peminjaman{}, errors.New(msgTanggalSalah)
	}
	m, sel := *in.SlotMulai, *in.SlotSelesai
	if m < 1 || m > 8 || sel < 1 || sel > 8 {
		return Peminjaman{}, errors.New(msgSlotRentang)
	}
	if m > sel {
		return Peminjaman{}, errors.New(msgSlotUrutan)
	}
	return Peminjaman{
		SumberDayaID: *in.SumberDayaID, Peminjam: peminjam, Tanggal: tanggal,
		SlotMulai: m, SlotSelesai: sel, Jumlah: *in.Jumlah,
	}, nil
}

// ---------- B5: bentrok ----------

func menahanSlot(status string) bool {
	return status == "diajukan" || status == "disetujui"
}

// Pemanggil harus memegang s.mu. kecuali = id yang diabaikan (0 = tidak ada).
func (s *server) terpakaiLocked(sdID int, tanggal string, slot, kecuali int) int {
	total := 0
	for _, p := range s.peminjaman {
		if p.ID == kecuali {
			continue
		}
		if p.SumberDayaID != sdID || p.Tanggal != tanggal {
			continue
		}
		if !menahanSlot(p.Status) {
			continue
		}
		if slot >= p.SlotMulai && slot <= p.SlotSelesai {
			total += p.Jumlah
		}
	}
	return total
}

//helper

func (s *server) bentrokLocked(p Peminjaman, stok, kecuali int) bool {
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		if s.terpakaiLocked(p.SumberDayaID, p.Tanggal, slot, kecuali)+p.Jumlah > stok {
			return true
		}
	}
	return false
}

func (s *server) cariPeminjamanIdxLocked(id int) int {
	for i, p := range s.peminjaman {
		if p.ID == id {
			return i
		}
	}
	return -1
}

// POST /peminjaman
func (s *server) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, msgBodyInvalid)
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
		writeError(w, http.StatusBadRequest, msgSDIDTidakAda)
		return
	}
	if p.Jumlah < 1 || p.Jumlah > sd.Stok {
		writeError(w, http.StatusBadRequest, msgJumlahSalah)
		return
	}
	if s.bentrokLocked(p, sd.Stok, 0) {
		writeError(w, http.StatusConflict, msgBentrok)
		return
	}

	p.ID = s.next_PID
	s.next_PID++
	//p.Status = statusDiajukan
	s.peminjaman = append(s.peminjaman, p)
	WriteJSON(w, http.StatusCreated, p)
}

// GET /peminjaman
func (s *server) listPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var sdFilter int
	hasSD := false
	if v := q.Get("sumber_daya_id"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, msgFilterSDID)
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

// GET /peminjaman/{id}
func (s *server) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := s.cariPeminjamanIdxLocked(id)
	if idx < 0 {
		writeError(w, http.StatusNotFound, msgPinjamTidakAda)
		return
	}
	WriteJSON(w, http.StatusOK, s.peminjaman[idx])
}

// PUT
func (s *server) updatePeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	// Baca dan validasi body tanpa lock, tetapi balas setelah cek 404.
	var p Peminjaman
	bodyErr := ""
	var in peminjamanInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		bodyErr = msgBodyInvalid
	} else if v, err := validatePeminjaman(in); err != nil {
		bodyErr = err.Error()
	} else {
		p = v
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.cariPeminjamanIdxLocked(id)
	if idx < 0 {
		writeError(w, http.StatusNotFound, msgPinjamTidakAda)
		return
	}
	if bodyErr != "" {
		writeError(w, http.StatusBadRequest, bodyErr)
		return
	}

	sd, found := s.cariSumberDayaLocked(p.SumberDayaID)
	if !found {
		writeError(w, http.StatusBadRequest, msgSDIDTidakAda)
		return
	}
	if p.Jumlah < 1 || p.Jumlah > sd.Stok {
		writeError(w, http.StatusBadRequest, msgJumlahSalah)
		return
	}

	p.ID = id
	p.Status = s.peminjaman[idx].Status
	p.Alasan = s.peminjaman[idx].Alasan

	if s.bentrokLocked(p, sd.Stok, id) {
		writeError(w, http.StatusConflict, msgBentrok)
		return
	}

	s.peminjaman[idx] = p
	WriteJSON(w, http.StatusOK, p)
}

// DELETE
func (s *server) deletePeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.cariPeminjamanIdxLocked(id)
	if idx < 0 {
		writeError(w, http.StatusNotFound, msgPinjamTidakAda)
		return
	}
	s.peminjaman = append(s.peminjaman[:idx], s.peminjaman[idx+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

// PERSONA

func (s *server) ubahStatus(w http.ResponseWriter, r *http.Request, dari []string, ke string, perluAlasan bool) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	alasan, bodyErr := "", ""
	if perluAlasan {
		var in struct {
			Alasan *string `json:"alasan"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			bodyErr = msgBodyInvalid
		} else if in.Alasan == nil || strings.TrimSpace(*in.Alasan) == "" {
			bodyErr = msgAlasanWajib
		} else if utf8.RuneCountInString(strings.TrimSpace(*in.Alasan)) > 100 {
			bodyErr = msgAlasanPanjang
		} else {
			alasan = strings.TrimSpace(*in.Alasan)
		}

	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.cariPeminjamanIdxLocked(id)
	if idx < 0 {
		writeError(w, http.StatusNotFound, msgPinjamTidakAda)
		return
	}
	if bodyErr != "" {
		writeError(w, http.StatusBadRequest, bodyErr)
		return
	}

	boleh := false
	for _, d := range dari {
		if s.peminjaman[idx].Status == d {
			boleh = true
		}
	}
	if !boleh {
		writeError(w, http.StatusConflict, msgStatusTidakBisa)
		return
	}

	s.peminjaman[idx].Status = ke
	if perluAlasan {
		s.peminjaman[idx].Alasan = alasan
	}
	WriteJSON(w, http.StatusOK, s.peminjaman[idx])
}

func (s *server) setujuiPeminjaman(w http.ResponseWriter, r *http.Request) {
	s.ubahStatus(w, r, []string{"diajukan"}, "disetujui", false)
}
func (s *server) tolakPeminjaman(w http.ResponseWriter, r *http.Request) {
	s.ubahStatus(w, r, []string{"diajukan"}, "ditolak", true)
}
func (s *server) selesaiPeminjaman(w http.ResponseWriter, r *http.Request) {
	s.ubahStatus(w, r, []string{"disetujui"}, "selesai", false)
}
