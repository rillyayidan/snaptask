# Contributing

Panduan ini dipakai untuk latihan kontribusi kecil di SnapTask memakai branch, commit, push, dan merge.

## Branch Workflow

1. Pastikan branch utama terbaru.

```bash
git switch main
git pull origin main
```

2. Buat branch feature.

```bash
git switch -c feature/nama-feature
```

Jika branch sudah dibuat, cukup pindah ke branch itu.

```bash
git switch feature/nama-feature
```

3. Kerjakan perubahan, lalu cek file yang berubah.

```bash
git status
git diff
```

4. Jalankan verifikasi yang relevan sebelum commit.

```bash
cd backend
go test ./...
```

```bash
cd frontend
npm.cmd run build
```

5. Commit perubahan.

```bash
git add README.md CONTRIBUTING.md
git commit -m "Add contribution workflow guide"
```

6. Push branch ke origin.

```bash
git push -u origin feature/nama-feature
```

Flag `-u` menghubungkan branch lokal dengan branch remote, jadi push berikutnya cukup memakai `git push`.

7. Merge lewat pull request.

- Buka repository di GitHub.
- Buat pull request dari `feature/nama-feature` ke `main`.
- Pastikan hasil test/build aman.
- Merge pull request.
- Setelah merge, update lokal.

```bash
git switch main
git pull origin main
```

8. Bersihkan branch feature jika sudah tidak dipakai.

```bash
git branch -d feature/nama-feature
git push origin --delete feature/nama-feature
```

## Kapan Push?

Push branch ke origin setelah ada commit yang ingin disimpan, dibagikan, atau dibuat pull request. Push branch kosong boleh saja untuk backup atau kolaborasi awal, tapi biasanya urutannya adalah edit, test, commit, lalu push.
