# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Replikasi sumberdaya.go -> peminjaman.go
- Prompt: Apakah bisa mengimplementasikan peminjaman.go seperti (snippet sumberdaya.go)
- Dipakai: bisa, selama tetap mengikuti constraint awal
- Diubah: menambahkan peminjaman.go
Memastikan iimplementasi bisa serupa/tidak

## Error formatting
- Prompt: Apakah penggunaan error suah sesuai B4? (paste snippet peminjaman)
- Dipakai: sesuai, namun disarankan membuat konstanta error baru di errors.go, (men-generate list error yang sudah ada)
- Diubah: penambahan errors.go dan mengubah validasion di setiap. 
Dilakukan agar error yang serupa tidak harus duplikat & memenuhi B4.


## prototypng racing-test
- Prompt: Bagaimana membuktikan B7 berjalan.
- Dipakai: race_test.go
- Diubah: race_testt go
prototype racing test (rusak)

## Kerangka awal B1-B3 (server.go, sumberdaya.go)
- Prompt: Apa yang perlu diambil dari store.go dan api.go pertemuan-06 untuk memulai B1-B3, dan apakah butuh interface.
- Dipakai: struct server, NewHandler, writeJSON
- Diubah: NewHandler & setup mux & handler lainnya

Tidak memakai interface karena tidak ada database di ETS.

## Deteksi bentrok (B5)
- Prompt: Cara mengimplementasikan deteksi bentrok kapasitas per slot.
- Dipakai: menahanSlot(), terpakaiLocked(), bentrokLocked() di peminjaman.go, dipanggil di createPeminjaman.
- Diubah: bentrok mengembalikan 409 (sebelumnya saya pakai 400), pemanggilan diletakkan setelah cek jumlah dan sebelum ID dibuat, di dalam lock yang sama supaya aman untuk B7.

## Ubah dan hapus (B6)
- Prompt: Apakah updatePeminjaman saya sudah benar untuk B6.
- Dipakai: helper bentrokLocked dengan parameter kecuali=id.
- Diubah: implementasi bentrokLocked
Dilakukan agar PUT tidak bentrok dengan dirinya sendiri dan data lama tidak berubah jika gagal.

## Change request A: ringkasan harian
- Prompt: identifikasi apa yang berjalan
- Dipakai: getRingkasan() dan struct ringkasan di ringkasan.go, memakai menahanSlot().
- Diubah: server, tambahan error 1, dan handler
Daftar per_sumber_daya dibuat dari tabel sumber daya dulu agar yang tanpa peminjaman tetap muncul dengan nilai 0.