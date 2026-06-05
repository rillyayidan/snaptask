import React, { useState } from 'react'
import { createRoot } from 'react-dom/client'
import { Camera, CheckCircle2, Download, Loader2, Trash2, Upload, Wand2 } from 'lucide-react'
import { TaskList } from './components/TaskList.jsx'
import { PushButton } from './components/PushButton.jsx'
import { useShareTarget } from './hooks/useShareTarget.js'
import { useGoogleAuth } from './hooks/useGoogleAuth.js'
import './styles.css'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8081'
const MAX_IMAGE_BYTES = 10 * 1024 * 1024

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
  const [isDragging, setIsDragging] = useState(false)
  const [installPrompt, setInstallPrompt] = useState(null)
  const [isInstalled, setIsInstalled] = useState(isStandaloneDisplay())
  const extractIdRef = React.useRef(0)
  const activeExtractionRef = React.useRef(null)

  React.useEffect(() => {
    if (shared.image) {
      if (acceptImage(shared.image)) {
        extract(shared.image)
      }
    }
    if (shared.error) {
      setError(shared.error)
      setStatus('idle')
    }
  }, [shared.image, shared.error])

  React.useEffect(() => () => {
    cancelActiveExtraction()
  }, [])

  React.useEffect(() => {
    function handleBeforeInstallPrompt(event) {
      event.preventDefault()
      setInstallPrompt(event)
    }

    function handleAppInstalled() {
      setInstallPrompt(null)
      setIsInstalled(true)
    }

    const displayMode = window.matchMedia?.('(display-mode: standalone)')
    const handleDisplayModeChange = (event) => setIsInstalled(event.matches || isStandaloneDisplay())

    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
    window.addEventListener('appinstalled', handleAppInstalled)
    if (displayMode?.addEventListener) {
      displayMode.addEventListener('change', handleDisplayModeChange)
    } else {
      displayMode?.addListener?.(handleDisplayModeChange)
    }

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
      window.removeEventListener('appinstalled', handleAppInstalled)
      if (displayMode?.removeEventListener) {
        displayMode.removeEventListener('change', handleDisplayModeChange)
      } else {
        displayMode?.removeListener?.(handleDisplayModeChange)
      }
    }
  }, [])

  React.useEffect(() => {
    function handlePaste(event) {
      const pastedImage = imageFileFromClipboard(event.clipboardData)
      if (!pastedImage) return
      event.preventDefault()
      if (acceptImage(pastedImage)) {
        extract(pastedImage)
      }
    }

    window.addEventListener('paste', handlePaste)
    return () => window.removeEventListener('paste', handlePaste)
  }, [])

  const imageUrl = React.useMemo(() => (image ? URL.createObjectURL(image) : ''), [image])
  const uploadBoxClass = ['upload-box', imageUrl ? 'has-image' : '', isDragging ? 'dragging' : ''].filter(Boolean).join(' ')
  const canInstall = Boolean(installPrompt) && !isInstalled

  React.useEffect(() => () => {
    if (imageUrl) URL.revokeObjectURL(imageUrl)
  }, [imageUrl])

  async function installApp() {
    if (!installPrompt) return
    const prompt = installPrompt
    setInstallPrompt(null)
    try {
      await prompt.prompt()
      await prompt.userChoice
    } catch {
      setInstallPrompt(null)
    }
  }

  async function extract(targetImage = image) {
    if (!targetImage) return
    cancelActiveExtraction()
    const extractId = ++extractIdRef.current
    const controller = new AbortController()
    activeExtractionRef.current = controller
    setStatus('extracting')
    setError('')
    const form = new FormData()
    form.append('image', targetImage)
    try {
      const response = await fetch(`${API_BASE}/extract`, {
        method: 'POST',
        body: form,
        signal: controller.signal
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error ?? 'Extraction failed')
      if (extractId !== extractIdRef.current) return
      setItems(data.items ?? [])
      setStatus('ready')
    } catch (err) {
      if (err.name === 'AbortError') return
      if (extractId !== extractIdRef.current) return
      setError(err.message)
      setStatus('idle')
    } finally {
      if (activeExtractionRef.current === controller) {
        activeExtractionRef.current = null
      }
    }
  }

  function onFileChange(event) {
    const file = event.target.files?.[0]
    if (!file) return
    if (acceptImage(file)) {
      extract(file)
    }
    event.target.value = ''
  }

  function acceptImage(file) {
    const validationError = validateImageFile(file)
    cancelActiveExtraction()
    extractIdRef.current += 1
    setItems([])
    setStatus('idle')
    if (validationError) {
      setImage(null)
      setError(validationError)
      return false
    }
    setImage(file)
    setError('')
    return true
  }

  function onDrop(event) {
    event.preventDefault()
    setIsDragging(false)
    const droppedImage = firstImageFile(event.dataTransfer?.files)
    if (!droppedImage) {
      setError('Drop a screenshot image file.')
      return
    }
    if (acceptImage(droppedImage)) {
      extract(droppedImage)
    }
  }

  function clearCapture() {
    cancelActiveExtraction()
    extractIdRef.current += 1
    setImage(null)
    setItems([])
    setError('')
    setStatus('idle')
  }

  function cancelActiveExtraction() {
    activeExtractionRef.current?.abort()
    activeExtractionRef.current = null
  }

  return (
    <main className="app-shell">
      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">SnapTask</p>
            <h1>Screenshot in. Tasks out.</h1>
          </div>
          <div className="top-actions">
            {canInstall && (
              <button className="icon-button install-button" title="Install SnapTask" aria-label="Install SnapTask" onClick={installApp}>
                <Download size={19} />
              </button>
            )}
            <button className="auth-button" onClick={auth.signIn}>
              <CheckCircle2 size={18} />
              {auth.user ? auth.user.displayName : 'Sign in'}
            </button>
          </div>
        </header>
        {auth.error && <p className="error-banner">{auth.error}</p>}

        <div className="layout">
          <section className="capture-panel">
            <div
              className={uploadBoxClass}
              onDragEnter={(event) => {
                event.preventDefault()
                setIsDragging(true)
              }}
              onDragOver={(event) => event.preventDefault()}
              onDragLeave={() => setIsDragging(false)}
              onDrop={onDrop}
            >
              {imageUrl ? (
                <img src={imageUrl} alt="Shared screenshot preview" />
              ) : (
                <div className="empty-state">
                  <Camera size={40} />
                  <span>Share a screenshot here or upload one for desktop testing.</span>
                </div>
              )}
            </div>
            {image && (
              <div className="image-meta">
                <span>{image.name || 'Shared screenshot'}</span>
                <span>{formatBytes(image.size)}</span>
              </div>
            )}
            <div className="toolbar">
              <label className="icon-button" title="Upload screenshot">
                <Upload size={19} />
                <input type="file" accept="image/*" onChange={onFileChange} />
              </label>
              <button className="primary-button" disabled={!image || status === 'extracting'} onClick={extract}>
                {status === 'extracting' ? <Loader2 className="spin" size={18} /> : <Wand2 size={18} />}
                Extract
              </button>
              <button className="icon-button" title="Clear" disabled={!image} onClick={clearCapture}>
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

function isStandaloneDisplay() {
  const displayMode = window.matchMedia?.('(display-mode: standalone)')
  return Boolean(displayMode?.matches || window.navigator.standalone)
}

function validateImageFile(file) {
  if (file.type && !file.type.startsWith('image/')) {
    return 'Choose a screenshot image file.'
  }
  if (file.size > MAX_IMAGE_BYTES) {
    return 'Screenshot is too large. Use an image under 10 MB.'
  }
  return ''
}

function firstImageFile(files) {
  return Array.from(files ?? []).find((file) => !file.type || file.type.startsWith('image/')) ?? null
}

function imageFileFromClipboard(data) {
  const imageFromFiles = firstImageFile(data?.files)
  if (imageFromFiles) return imageFromFiles

  const imageItem = Array.from(data?.items ?? []).find((item) => item.kind === 'file' && item.type.startsWith('image/'))
  return imageItem?.getAsFile() ?? null
}

function formatBytes(value) {
  if (!value) return '0 KB'
  if (value < 1024 * 1024) return `${Math.max(1, Math.round(value / 1024))} KB`
  return `${(value / 1024 / 1024).toFixed(1)} MB`
}
