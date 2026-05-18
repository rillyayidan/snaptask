import React, { useState } from 'react'
import { createRoot } from 'react-dom/client'
import { Camera, CheckCircle2, Loader2, Trash2, Upload, Wand2 } from 'lucide-react'
import { TaskList } from './components/TaskList.jsx'
import { PushButton } from './components/PushButton.jsx'
import { useShareTarget } from './hooks/useShareTarget.js'
import { useGoogleAuth } from './hooks/useGoogleAuth.js'
import './styles.css'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8081'

if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js')
}

function App() {
  const shared = useShareTarget()
  const auth = useGoogleAuth()
  const [image, setImage] = useState(shared.image)
  const [items, setItems] = useState([])
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState('')

  React.useEffect(() => {
    if (shared.image) {
      setImage(shared.image)
      setItems([])
      setError('')
      setStatus('idle')
    }
    if (shared.error) {
      setError(shared.error)
      setStatus('idle')
    }
  }, [shared.image, shared.error])

  const imageUrl = React.useMemo(() => (image ? URL.createObjectURL(image) : ''), [image])

  React.useEffect(() => () => {
    if (imageUrl) URL.revokeObjectURL(imageUrl)
  }, [imageUrl])

  async function extract() {
    if (!image) return
    setStatus('extracting')
    setError('')
    const form = new FormData()
    form.append('image', image)
    try {
      const response = await fetch(`${API_BASE}/extract`, {
        method: 'POST',
        body: form
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error ?? 'Extraction failed')
      setItems(data.items ?? [])
      setStatus('ready')
    } catch (err) {
      setError(err.message)
      setStatus('idle')
    }
  }

  function onFileChange(event) {
    const file = event.target.files?.[0]
    if (!file) return
    setImage(file)
    setItems([])
    setError('')
    setStatus('idle')
  }

  return (
    <main className="app-shell">
      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">SnapTask</p>
            <h1>Screenshot in. Tasks out.</h1>
          </div>
          <button className="auth-button" onClick={auth.signIn}>
            <CheckCircle2 size={18} />
            {auth.user ? auth.user.displayName : 'Sign in'}
          </button>
        </header>
        {auth.error && <p className="error-banner">{auth.error}</p>}

        <div className="layout">
          <section className="capture-panel">
            <div className="upload-box">
              {imageUrl ? (
                <img src={imageUrl} alt="Shared screenshot preview" />
              ) : (
                <div className="empty-state">
                  <Camera size={40} />
                  <span>Share a screenshot here or upload one for desktop testing.</span>
                </div>
              )}
            </div>
            <div className="toolbar">
              <label className="icon-button" title="Upload screenshot">
                <Upload size={19} />
                <input type="file" accept="image/*" onChange={onFileChange} />
              </label>
              <button className="primary-button" disabled={!image || status === 'extracting'} onClick={extract}>
                {status === 'extracting' ? <Loader2 className="spin" size={18} /> : <Wand2 size={18} />}
                Extract
              </button>
              <button className="icon-button" title="Clear" disabled={!image} onClick={() => {
                setImage(null)
                setItems([])
                setError('')
                setStatus('idle')
              }}>
                <Trash2 size={19} />
              </button>
            </div>
            {error && <p className="error-text">{error}</p>}
          </section>

          <section className="review-panel">
            <TaskList items={items} onChange={setItems} />
            <PushButton
              apiBase={API_BASE}
              accessToken={auth.accessToken}
              items={items}
              disabled={!auth.accessToken || items.length === 0}
              disabledReason={!auth.configured ? 'Firebase is not configured yet.' : !auth.accessToken ? 'Sign in with Google before pushing.' : ''}
            />
          </section>
        </div>
      </section>
    </main>
  )
}

createRoot(document.getElementById('root')).render(<App />)
