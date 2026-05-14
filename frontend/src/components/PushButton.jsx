import { Loader2, Send } from 'lucide-react'
import { useState } from 'react'

export function PushButton({ apiBase, accessToken, items, disabled, disabledReason }) {
  const [state, setState] = useState('idle')
  const [message, setMessage] = useState('')

  async function pushAll() {
    setState('pushing')
    setMessage('')
    try {
      const response = await fetch(`${apiBase}/push`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ access_token: accessToken, items })
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error ?? 'Push failed')
      const failed = (data.results ?? []).filter((result) => result.status !== 'ok')
      setMessage(failed.length ? `${failed.length} item failed. Check Google permissions.` : 'Pushed to Google Tasks and Calendar.')
      setState('done')
    } catch (err) {
      setMessage(err.message)
      setState('idle')
    }
  }

  return (
    <div className="push-bar">
      <button className="primary-button wide" disabled={disabled || state === 'pushing'} onClick={pushAll}>
        {state === 'pushing' ? <Loader2 className="spin" size={18} /> : <Send size={18} />}
        Push all
      </button>
      {disabled && disabledReason && items.length > 0 && <p className="hint-text">{disabledReason}</p>}
      {message && <p className={state === 'done' ? 'success-text' : 'error-text'}>{message}</p>}
    </div>
  )
}
