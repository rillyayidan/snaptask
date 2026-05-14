import { useMemo, useState } from 'react'
import { initializeApp } from 'firebase/app'
import { getAuth, GoogleAuthProvider, signInWithPopup } from 'firebase/auth'

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
  appId: import.meta.env.VITE_FIREBASE_APP_ID
}

export function useGoogleAuth() {
  const [user, setUser] = useState(null)
  const [accessToken, setAccessToken] = useState('')
  const [error, setError] = useState('')

  const auth = useMemo(() => {
    if (!firebaseConfig.apiKey) return null
    return getAuth(initializeApp(firebaseConfig))
  }, [])

  async function signIn() {
    setError('')
    if (!auth) {
      setError('Firebase config is missing. Fill frontend env values before Google sign-in.')
      return
    }
    try {
      const provider = new GoogleAuthProvider()
      provider.addScope('https://www.googleapis.com/auth/tasks')
      provider.addScope('https://www.googleapis.com/auth/calendar.events')
      const result = await signInWithPopup(auth, provider)
      const credential = GoogleAuthProvider.credentialFromResult(result)
      setUser(result.user)
      setAccessToken(credential?.accessToken ?? '')
    } catch (err) {
      setError(err.message)
    }
  }

  return { user, accessToken, error, configured: Boolean(auth), signIn }
}
