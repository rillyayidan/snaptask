# SnapTask Session Summary

Date: 2026-05-15

## What Was Built

- Go/Fiber backend with:
  - `GET /health`
  - `GET /healthz`
  - `POST /extract`
  - `POST /push`
- React/Vite frontend PWA with:
  - screenshot preview
  - share-target support
  - task review UI
  - Google sign-in flow
  - push button for Google Tasks and Calendar
- Firebase Hosting config for the frontend
- Cloud Run deploy config for backend and frontend

## Deployed Services

- Backend: `https://snaptask-api-872640630070.asia-southeast2.run.app`
- Backend health: `https://snaptask-api-872640630070.asia-southeast2.run.app/health`
- Frontend: `https://snaptask-web-872640630070.asia-southeast2.run.app`

## Local Ports

- Backend dev default: `http://localhost:8081`
- Frontend dev: `http://localhost:5173`

## Repo Files Added

- `backend/`
- `frontend/`
- `cloudbuild.yaml`
- `frontend-cloudbuild.yaml`
- `frontend/firebase.json`
- `frontend/.firebaserc`
- `README.md`

## Remaining Credentials / Config

You already have:

- `GEMINI_API_KEY`

Still needed:

- `VITE_FIREBASE_API_KEY`
- `VITE_FIREBASE_AUTH_DOMAIN`
- `VITE_FIREBASE_PROJECT_ID`
- `VITE_FIREBASE_STORAGE_BUCKET`
- `VITE_FIREBASE_MESSAGING_SENDER_ID`
- `VITE_FIREBASE_APP_ID`

Firebase/Auth setup still needed:

- Enable Google sign-in provider in Firebase Authentication
- Add authorized domains for the deployed app
- Add your test/demo Google account as an OAuth test user
- Add scopes:
  - `https://www.googleapis.com/auth/tasks`
  - `https://www.googleapis.com/auth/calendar.events`

Cloud Run secret wiring still needed:

- Put `GEMINI_API_KEY` into Secret Manager if not already done
- Attach the secret to `snaptask-api`

## Commands Used Most Recently

```powershell
cd C:\Users\zidan\Documents\snaptask\backend
go run ./cmd
```

```powershell
cd C:\Users\zidan\Documents\snaptask\frontend
npm.cmd run dev -- --host 0.0.0.0
```

```powershell
cd C:\Users\zidan\Documents\snaptask
gcloud.cmd builds submit --config=cloudbuild.yaml --substitutions=_ALLOWED_ORIGINS=*
gcloud.cmd builds submit --config=frontend-cloudbuild.yaml --substitutions=_API_BASE_URL=https://snaptask-api-872640630070.asia-southeast2.run.app
```
