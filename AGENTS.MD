Kompleksitas: 🟡 Medium — input-nya unik (screenshot dari share sheet), tapi arsitekturnya lebih lean dari NowHere. Core-nya: image → Gemini Vision OCR + extraction → structured tasks → Google Tasks/Calendar API. Let's blueprint it.SnapTask lebih lean dari NowHere — medium complexity. Core mechanic-nya simpel tapi powerful: satu gesture (share screenshot) → Gemini Vision baca → tasks terstruktur muncul. Dua diagram cukup: system overview dulu, lalu request flow.Tiga layer — jauh lebih lean dari NowHere. Tidak ada Firestore preference loop, tidak ada context aggregation. Sekarang request flow-nya — dari screenshot di-share sampai task muncul di Google Tasks.Perhatikan percabangan di step Gemini — output-nya sudah diklasifikasikan sebelum sampai ke user review. Task → Google Tasks, Event → Google Calendar, langsung paralel. Sekarang blueprint lengkapnya.

---

## Tech stack

| Layer | Technology | Why |
|---|---|---|
| Frontend | React PWA + Vite | Web Share Target API butuh PWA; React untuk task review UI yang interaktif |
| Styling | Tailwind CSS | Mobile-first cepat, utility classes |
| Backend | Go + Fiber | Singkat, fast image processing, concurrency untuk parallel Google API push |
| AI | Gemini 2.0 Flash Vision | Multimodal — bisa baca teks DAN visual context dari screenshot sekaligus |
| Auth | Firebase Auth | Google OAuth bawaan — satu klik, langsung dapat akses Tasks + Calendar |
| Deploy | Cloud Run | Scales to zero, container-based, cocok untuk hackathon |
| External | Google Tasks API | Tempat akhir semua task — sudah built-in di semua Android |
| External | Google Calendar API | Untuk item yang punya tanggal/waktu spesifik |

> Tidak ada database yang dibutuhkan untuk MVP — SnapTask adalah stateless by design. Screenshot masuk, tasks keluar. Tidak ada yang perlu disimpan server-side.

---

## Prompt engineering note

Ini yang paling kritis di SnapTask. Gemini harus output JSON yang bisa langsung di-parse, bukan prosa.

```
You are a task extractor. Analyze this screenshot and extract ALL action items, commitments, and time-sensitive information.

Respond ONLY with a JSON array. No explanation, no markdown, no preamble.

Format:
[
  {
    "type": "task" | "event" | "note",
    "title": "short action title",
    "detail": "original text from screenshot",
    "due_date": "ISO 8601 or null",
    "priority": "high" | "medium" | "low"
  }
]

Extract implicit commitments too. "Eh nanti kirim filenya ya" = task: "Kirim file ke [sender]".
```

---

## Project structure

```
snaptask/
├── frontend/
│   ├── public/
│   │   └── manifest.json         # PWA manifest + share_target config
│   ├── src/
│   │   ├── components/
│   │   │   ├── TaskCard.jsx       # Review UI per task
│   │   │   ├── TaskList.jsx       # List semua hasil ekstraksi
│   │   │   └── PushButton.jsx     # Confirm & push ke Google
│   │   ├── hooks/
│   │   │   ├── useShareTarget.js  # Handle incoming share
│   │   │   └── useGoogleAuth.js   # Firebase + OAuth scope
│   │   └── main.jsx
│   └── vite.config.js
│
├── backend/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/
│   │   │   └── extract.go         # POST /extract — core endpoint
│   │   ├── service/
│   │   │   ├── gemini.go          # Vision API call + prompt
│   │   │   ├── classifier.go      # Task vs event vs note
│   │   │   └── google.go          # Tasks + Calendar push
│   │   └── model/
│   │       └── task.go
│   ├── .env.example
│   └── Dockerfile
│
└── README.md
```

---

## Development roadmap

**Phase 1 — Share works** *(Hari 1)*
- Setup PWA + `manifest.json` dengan `share_target` config
- Screenshot dari WhatsApp bisa masuk ke app
- Tampilkan preview image di UI (buktikan share mechanic jalan)
- Milestone demo: "Share dari WA, muncul di app"

**Phase 2 — Gemini reads** *(Hari 2–3)*
- Backend endpoint `POST /extract` jalan
- Gemini Vision baca screenshot, output JSON tasks
- Task list muncul di UI dengan edit/delete
- Milestone demo: "Screenshot chat → list to-do muncul otomatis"

**Phase 3 — Google push** *(Hari 4)*
- Firebase Auth + Google OAuth scope (Tasks + Calendar)
- Task → Google Tasks, Event → Google Calendar, paralel
- Konfirmasi sukses di UI
- Milestone demo: "Satu share, tasks langsung di Google Tasks"

**Phase 4 — Polish** *(Hari 5)*
- Error handling: gambar blur, screenshot kosong, API timeout
- Loading state yang smooth (ini yang keliatan di demo)
- Deploy ke Cloud Run
- Milestone demo: Full flow dari share → Google Tasks dalam <5 detik

---

## Risks & considerations

- Web Share Target API hanya jalan di PWA yang sudah di-install — pastikan onboarding ada instruksi "Add to Home Screen", ini sering dilewat
- Gemini Vision bisa salah baca teks yang blur atau font kecil — tambahkan fallback: "Tidak bisa baca screenshot ini, coba yang lebih jelas"
- Google OAuth scope untuk Tasks + Calendar butuh approval kalau app dipublish — untuk demo/kompetisi, pakai test users di Google Cloud Console, tidak perlu production approval
- Scope risk: jangan buat "smart grouping" atau "priority AI" dulu — MVP adalah extract + push, itu saja sudah cukup wow
- Gemini kadang over-extract (ambil terlalu banyak hal) — prompt perlu di-tune dengan beberapa test screenshot nyata sebelum demo

---

## Demo script (untuk pitch)

Urutan yang paling impactful:

1. Buka WhatsApp, tunjukkan chat yang ada commitment ("nanti kirim file ya", "meeting Jumat jam 3")
2. Screenshot → Share → pilih SnapTask
3. Tunjukkan task list yang muncul otomatis — AI sudah baca implicit commitment
4. Satu tap "Push semua" → buka Google Tasks di HP, tasks sudah ada
5. Buka Google Calendar — event Jumat jam 3 sudah masuk

Total demo: kurang dari 60 detik. Tidak perlu penjelasan.

---

## Quick start

```bash
# 1. Clone dan setup
git clone https://github.com/rillyayidan/snaptask && cd snaptask
cp backend/.env.example backend/.env
# isi: GEMINI_API_KEY, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET

# 2. Jalankan backend
cd backend && go run ./cmd/main.go
# → localhost:8080

# 3. Jalankan frontend (butuh HTTPS untuk share target)
cd ../frontend && npm install && npm run dev -- --https
# → https://localhost:5173
```

> Untuk test share target di mobile: deploy ke Cloud Run dulu (atau pakai ngrok untuk dev), karena share target butuh domain yang valid, bukan localhost.

---

Mau langsung `/project-builder` sekarang? Aku bisa generate boilerplate — `manifest.json`, `extract.go`, Gemini prompt, dan `TaskCard.jsx` — siap coding hari ini.