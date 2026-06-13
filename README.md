# ZvA GAME — Web Edition (Tugas 12 Cloud Computing)

Game turn-based **Zhongli vs Azhdaha** (versi web dari game terminal Go-ku).
Server Go menjalankan logika battle asli, browser menampilkan UI bertema Liyue.

Dibuat oleh: **RASHAQA NASHWAN MOYA — 103012300058 — IF_47_11**

## Struktur
```
zva-game/
├── main.go       # server Go + logika battle (port dari ZvA.go)
├── index.html    # UI battle (di-embed ke binary lewat go:embed)
├── go.mod
└── README.md
```

## Cara jalanin lokal (buat ngetes dulu)
```bash
go run .
# buka http://localhost:5000
```

## Cara Deploy ke Render (PaaS)

### 1. Naikin ke GitHub
- Bikin repo baru di github.com (misal `zva-game`) -> Public
- Upload semua file (main.go, index.html, go.mod, README.md)

### 2. Bikin Web Service di Render
- render.com -> **New +** -> **Web Service** -> pilih repo `zva-game`
- Isi konfigurasi:

| Field | Isi |
|---|---|
| Language / Runtime | **Go** |
| Build Command | `go build -o app .` |
| Start Command | `./app` |
| Instance Type | **Free** |

### 3. Deploy
- Klik **Create Web Service**, tunggu build (~1-2 menit)
- Dapet URL: `https://zva-game.onrender.com` -> mainkan & lampirkan di laporan

## Catatan buat laporan (biar keliatan paham PaaS)
- **PaaS = Platform as a Service**: kita cuma push kode Go, Render yang urus
  compile, OS, runtime, networking, dan scaling. Kita nggak sewa server mentah.
- **Kenapa game terminal harus diubah?** PaaS menyajikan **HTTP (web)**, bukan
  terminal. Jadi input `fmt.Scan` diganti jadi tombol di browser yang manggil
  endpoint `/api/turn`.
- **Server stateless:** state game (HP, energy) disimpan di browser dan dikirim
  tiap giliran. Ini bikin app tahan terhadap "spin-down" free tier Render —
  walau server tidur lalu bangun, game tetap lanjut.
- **`os.Getenv("PORT")`:** port ditentukan PaaS, bukan di-hardcode.
- **Trade-off free tier:** service tidur kalau idle, request pertama ~1 menit
  buat bangun (spin-up). Wajar buat demo.
