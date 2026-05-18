# SnapTask

SnapTask turns a shared screenshot into reviewed tasks and events, then pushes them to Google Tasks and Google Calendar.

## Stack

- Frontend: React PWA + Vite
- Backend: Go + Fiber
- AI: Gemini Vision via `GEMINI_API_KEY`
- Auth: Firebase Google sign-in with Tasks and Calendar scopes
- Deploy: Cloud Run

## Local Setup

```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

Fill the Gemini and Firebase values, then run:

```bash
cd backend
go run ./cmd
```

```bash
cd frontend
npm.cmd install
npm.cmd run dev
```

The frontend env file needs `VITE_API_BASE_URL` plus the Firebase web app config values. Google sign-in must request Tasks and Calendar scopes, so keep demo accounts listed as test users until the OAuth consent screen is verified for production.

For real mobile share-target testing, deploy the frontend to a valid HTTPS origin and install the PWA. Local desktop upload works without install.

## Backend API

- Local backend default: `http://localhost:8081`
- `GET /health`
- `GET /healthz`
- `POST /extract` with multipart field `image`
- `POST /push` with `{ "access_token": "...", "items": [...] }`

`/push` sends `task` items to Google Tasks and `event` items to Google Calendar. `note` items are returned as `skipped` because they are review-only.

### API Examples

Extract tasks from a screenshot:

```bash
curl -X POST http://localhost:8081/extract \
  -F "image=@./sample-screenshot.png"
```

Push reviewed items to Google:

```bash
curl -X POST http://localhost:8081/push \
  -H "Content-Type: application/json" \
  -d '{
    "access_token": "GOOGLE_OAUTH_ACCESS_TOKEN",
    "items": [
      {
        "type": "task",
        "title": "Send project brief",
        "detail": "Original screenshot text",
        "due_date": "2026-05-22T09:00:00+07:00",
        "priority": "medium"
      }
    ]
  }'
```

## Cloud Run Deployment

Create an Artifact Registry Docker repo once:

```bash
gcloud artifacts repositories create cloud-run --repository-format=docker --location=asia-southeast2
```

Store the Gemini key in Secret Manager:

```bash
printf "YOUR_GEMINI_API_KEY" | gcloud secrets create GEMINI_API_KEY --data-file=-
```

Grant Cloud Run access to the secret if needed:

```bash
gcloud secrets add-iam-policy-binding GEMINI_API_KEY --member="serviceAccount:PROJECT_NUMBER-compute@developer.gserviceaccount.com" --role="roles/secretmanager.secretAccessor"
```

Attach it to the deployed API after the secret exists:

```bash
gcloud run services update snaptask-api --region=asia-southeast2 --set-secrets=GEMINI_API_KEY=GEMINI_API_KEY:latest
```

Submit the backend build and deploy:

```bash
gcloud builds submit --config=cloudbuild.yaml --substitutions=_ALLOWED_ORIGINS=https://YOUR_FRONTEND_ORIGIN
```

Deploy the frontend PWA to Cloud Run after the API URL is known:

```bash
gcloud builds submit frontend --tag asia-southeast2-docker.pkg.dev/PROJECT_ID/cloud-run/snaptask-web
gcloud run deploy snaptask-web --image=asia-southeast2-docker.pkg.dev/PROJECT_ID/cloud-run/snaptask-web --region=asia-southeast2 --allow-unauthenticated
```

For the final frontend image, pass build args for `VITE_API_BASE_URL` and Firebase config. One practical flow for the competition is:

```bash
gcloud builds submit frontend --config=frontend-cloudbuild.yaml --substitutions=_API_BASE_URL=https://YOUR_API_URL,_FIREBASE_API_KEY=...,_FIREBASE_AUTH_DOMAIN=...,_FIREBASE_PROJECT_ID=...,_FIREBASE_STORAGE_BUCKET=...,_FIREBASE_MESSAGING_SENDER_ID=...,_FIREBASE_APP_ID=...
```

## Firebase Hosting

The frontend also has Firebase Hosting wired in through [frontend/firebase.json](frontend/firebase.json) and [frontend/.firebaserc](frontend/.firebaserc). From the `frontend/` directory:

```bash
npm install
npm run build
firebase login
firebase deploy --only hosting
```

The hosting site is configured as an SPA, so all routes rewrite to `index.html`, and `sw.js` plus `manifest.json` are served with `no-cache` headers.
