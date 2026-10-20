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

## Apa yang harus diubah dari peminjaman untuk mengakomodasi persona

## prototypng racing-test

