# PROSES

## 1. Persona dan masalahnya
Pengguna persona adalah admin himpunan yang mengatur persetujuan penggunaan ruang himpunan. Masalahnya, ruang tidak boleh langsung terpakai hanya karena ada yang meminta: setiap permintaan harus disetujui dulu, dan penolakan harus disertai alasan. Fitur khasnya (alur diajukan, disetujui/ditolak, selesai) membuat proses itu punya urutan yang jelas dan riwayatnya bisa dilacak.

## 2. Rancangan struktur data
Semua data disimpan dalam satu struct `server` yang dibuat di dalam `NewHandler()`, sehingga tiap pemanggilan menghasilkan penyimpanan baru yang kosong. Handler berupa method dari `*server` (pointer) supaya semua handler berbagi data yang sama.

type SumberDaya struct { ID, Nama, Jenis, Stok }
type Peminjaman struct { ID, SumberDayaID, Peminjam, Tanggal, SlotMulai, SlotSelesai, Jumlah, Status, Alasan }

type server struct {
    mu         sync.Mutex
    sumberDaya []SumberDaya
    peminjaman []Peminjaman
    nextSDID   int
    nextPID    int
}

Saya memakai slice, bukan map, karena data kecil dan daftar harus terurut menurut id naik: dengan append id yang selalu naik, urutan itu sudah terjadi dengan sendirinya. Id tidak dipakai ulang setelah DELETE.

Perubahan di tengah jalan:
- Menambah slice `peminjaman` dan `Alasan` di struct Peminjaman (dari kartu persona C).
- Nama field `next_SDID`/`next_PID` diubah menjadi `nextSDID`/`nextPID` agar konsisten.

## 3. Satu keputusan desain

Pemakaian slot tidak disimpan di struktur terpisah, tetapi dihitung dari daftar peminjaman setiap kali dibutuhkan (`terpakaiLocked`), dan cek bentrok dilakukan per slot, bukan per rentang.
- Alternatif: menyimpan penghitung per sumber daya, tanggal, dan slot.
- Alasan: dengan penghitung terpisah, saya harus menjaganya tetap sinkron di POST, PUT, DELETE, dan tiap perubahan status; satu update yang terlewat membuat datanya salah. Dengan menghitung dari peminjaman, sumber datanya satu. Cek per slot diperlukan agar rentang 3-5 dan 4-4 terdeteksi bentrok, sedangkan 1-2 dan 3-4 tidak.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: `go test` tidak menemukan test untuk race (racing-test).
- Cara menemukan: menjalankan `go test ./...`.
- Penyebab dan solusi: tidak diketahui

### Bug 2
- Gejala: PUT /peminjaman/{id} pada id yang tidak ada dengan body rusak menghasilkan 400, bukan 404; bentrok kapasitas menghasilkan 400, bukan 409.
- Cara menemukan: membandingkan urutan pemeriksaan dan kode status di kode saya dengan kontrak ETS.
- Penyebab dan solusi: pengecekan keberadaan data ada di akhir handler, dan bentrok memakai `StatusBadRequest`. Saya pindahkan pencarian id sebelum pemeriksaan body (urutan 400, 404, 400, 409) dan mengganti kode bentrok menjadi 409.


## 5. Change request
Diminta endpoint `GET /ringkasan?tanggal=YYYY-MM-DD` berisi jumlah peminjaman yang menahan slot, total unit-slot, dan jumlah per sumber daya. Saya menambah `ringkasan.go` dan satu route di `NewHandler()`; endpoint lain tidak diubah.
Keputusan desain: daftar `per_sumber_daya` dibuat dari tabel sumber daya lebih dulu (dengan nilai 0), baru dihitung dari peminjaman. Jika dikelompokkan dari peminjaman, sumber daya tanpa peminjaman tidak akan muncul, padahal kartu mensyaratkannya. Saya juga memakai ulang `menahanSlot` agar ringkasan dan cek bentrok sepakat soal status yang menahan slot.