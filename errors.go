package main

const (
	msgBodyInvalid      = "body bukan JSON yang valid"
	msgIDBukanAngka     = "id harus berupa angka"
	msgSDTidakAda       = "sumber daya tidak ditemukan"
	msgPinjamTidakAda   = "peminjaman tidak ditemukan"
	msgFieldSDWajib     = "nama, jenis, dan stok wajib ada"
	msgFieldPinjamWajib = "semua field wajib ada"
	msgNamaKosong       = "nama tidak boleh kosong"
	msgNamaPanjang      = "nama maksimal 50 karakter"
	msgJenisSalah       = `jenis harus "ruang" atau "alat"`
	msgStokMinimal      = "stok minimal 1"
	msgPeminjamKosong   = "peminjam tidak boleh kosong"
	msgPeminjamPanjang  = "peminjam maksimal 50 karakter"
	msgTanggalSalah     = "tanggal harus YYYY-MM-DD dan valid"
	msgSlotRentang      = "slot harus antara 1 dan 8"
	msgSlotUrutan       = "slot_mulai tidak boleh lebih besar dari slot_selesai"
	msgSDIDTidakAda     = "sumber_daya_id tidak ditemukan"
	msgJumlahSalah      = "jumlah harus antara 1 dan stok sumber daya"
	msgFilterSDID       = "sumber_daya_id harus berupa angka"
	msgBentrok          = "kapasitas slot terlampaui"
	msgAlasanWajib      = "alasan penolakan wajib diisi"
	msgStatusTidakBisa  = "perubahan status tidak diizinkan dari status saat ini"
	msgAlasanPanjang    = "alasan maksimal 100 karakter"
)
