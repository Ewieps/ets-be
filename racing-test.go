package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestConcurrentSlot(t *testing.T) {
	h := NewHandler()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/sumber-daya",
		strings.NewReader(`{"nama":"Aula","jenis":"ruang","stok":1}`)))
	if rec.Code != 201 {
		t.Fatalf("setup: %d", rec.Code)
	}

	body := `{"sumber_daya_id":1,"peminjam":"Rina","tanggal":"2026-11-02","slot_mulai":3,"slot_selesai":4,"jumlah":1}`
	var wg sync.WaitGroup
	var mu sync.Mutex
	counts := map[int]int{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/peminjaman", strings.NewReader(body)))
			mu.Lock()
			counts[rr.Code]++
			mu.Unlock()
		}()
	}
	wg.Wait()

	if counts[201] != 1 || counts[409] != 49 {
		t.Fatalf("want 1x201 and 49x409, got %v", counts)
	}
}

//RUSAKK
