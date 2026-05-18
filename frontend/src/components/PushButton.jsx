import { Loader2, Send } from 'lucide-react'
import { useEffect, useState } from 'react'

export function PushButton({ apiBase, accessToken, items, disabled, disabledReason }) {
  const [state, setState] = useState('idle')
  const [message, setMessage] = useState('')
  const pushableItems = items.filter((item) => item.title?.trim())
  const missingTitleReason = items.length > 0 && pushableItems.length === 0 ? 'Add a title before pushing.' : ''
  const blockedReason = disabledReason || missingTitleReason
  const pushDisabled = disabled || pushableItems.length === 0

  useEffect(() => {
    setState('idle')
    setMessage('')
  }, [accessToken, items])

  async function pushAll() {
    if (pushableItems.length === 0) {
      setMessage('Add at least one titled item before pushing.')
      return
    }

    setState('pushing')
    setMessage('')
    try {
      const response = await fetch(`${apiBase}/push`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ access_token: accessToken, items: pushableItems })
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error ?? 'Push failed')
      const results = data.results ?? []
      const failed = results.filter((result) => result.status === 'error')
      const skipped = results.filter((result) => result.status === 'skipped')
      if (failed.length) {
        setMessage(`${failed.length} item failed. Check Google permissions.`)
      } else if (skipped.length) {
        setMessage(`Pushed to Google. ${skipped.length} note ${skipped.length === 1 ? 'was' : 'were'} kept review-only.`)
      } else {
        setMessage('Pushed to Google Tasks and Calendar.')
      }
      setState('done')
    } catch (err) {
      setMessage(err.message)
      setState('idle')
    }
  }

  return (
    <div className="push-bar">
      <button className="primary-button wide" disabled={pushDisabled || state === 'pushing'} onClick={pushAll}>
        {state === 'pushing' ? <Loader2 className="spin" size={18} /> : <Send size={18} />}
        Push all
      </button>
      {blockedReason && items.length > 0 && <p className="hint-text">{blockedReason}</p>}
      {message && <p className={state === 'done' ? 'success-text' : 'error-text'}>{message}</p>}
    </div>
  )
}
