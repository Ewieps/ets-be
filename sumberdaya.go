package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"
)

type sumberDayaInput struct {
	Nama  *string `json:"nama"`
	Jenis *string `json:"jenis"`
	Stok  *int    `json:"stok"`
}

func validateSumberDaya(in sumberDayaInput) (SumberDaya, error) {
	if in.Nama == nil || in.Jenis == nil || in.Stok == nil {
		return SumberDaya{}, errors.New("nama, jenis, dan stok wajib ada")
	}
	nama := strings.TrimSpace(*in.Nama)
	jenis := strings.TrimSpace(*in.Jenis)

	if nama == "" {
		return SumberDaya{}, errors.New("nama tidak boleh kosong")
	}
	if utf8.RuneCountInString(nama) > 50 {
		return SumberDaya{}, errors.New("nama maksimal 50 karakter")
	}
	if jenis != "ruang" && jenis != "alat" {
		return SumberDaya{}, errors.New(`jenis harus "ruang" atau "alat"`)
	}
	if *in.Stok < 1 {
		return SumberDaya{}, errors.New("stok minimal 1")
	}
	return SumberDaya{Nama: nama, Jenis: jenis, Stok: *in.Stok}, nil
}

func (s *server) cariSumberDayaLocked(id int) (SumberDaya, bool) {
	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			return sd, true
		}
	}
	return SumberDaya{}, false
}

// POST
func (s *server) createSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in sumberDayaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "body bukan JSON yang valid")
		return
	}
	sd, err := validateSumberDaya(in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	sd.ID = s.next_SDID
	s.next_SDID++
	s.sumberDaya = append(s.sumberDaya, sd)
	WriteJSON(w, http.StatusCreated, sd)
}

// GET
func (s *server) listSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]SumberDaya, len(s.sumberDaya)) // bukan nil -> "[]"
	copy(out, s.sumberDaya)
	WriteJSON(w, http.StatusOK, out)
}

// GET (ID)
func (s *server) getSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	sd, found := s.cariSumberDayaLocked(id)
	if !found {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	WriteJSON(w, http.StatusOK, sd)
}
