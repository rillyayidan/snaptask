import { AlertTriangle, CheckCircle2, CircleSlash, Loader2, Send } from 'lucide-react'
import { useEffect, useState } from 'react'

export function PushButton({ apiBase, accessToken, items, disabled, disabledReason }) {
  const [state, setState] = useState('idle')
  const [message, setMessage] = useState('')
  const [results, setResults] = useState([])
  const pushableItems = items.filter((item) => item.title?.trim())
  const pushSummary = summarizePushableItems(pushableItems)
  const missingTitleReason = items.length > 0 && pushableItems.length === 0 ? 'Add a title before pushing.' : ''
  const blockedReason = disabledReason || missingTitleReason
  const pushDisabled = disabled || pushableItems.length === 0
  const hasResultErrors = results.some((result) => result.status === 'error')
  const messageClass = state === 'done' && !hasResultErrors ? 'success-text' : 'error-text'

  useEffect(() => {
    setState('idle')
    setMessage('')
    setResults([])
  }, [accessToken, items])

  async function pushAll() {
    if (pushableItems.length === 0) {
      setMessage('Add at least one titled item before pushing.')
      setResults([])
      return
    }

    setState('pushing')
    setMessage('')
    setResults([])
    try {
      const response = await fetch(`${apiBase}/push`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ access_token: accessToken, items: pushableItems })
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error ?? 'Push failed')
      const results = data.results ?? []
      setResults(results)
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
      setResults([])
      setState('idle')
    }
  }

  return (
    <div className="push-bar">
      <button className="primary-button wide" disabled={pushDisabled || state === 'pushing'} onClick={pushAll}>
        {state === 'pushing' ? <Loader2 className="spin" size={18} /> : <Send size={18} />}
        Push all
      </button>
      {!blockedReason && pushSummary && <p className="hint-text">{pushSummary}</p>}
      {blockedReason && items.length > 0 && <p className="hint-text">{blockedReason}</p>}
      {message && <p className={messageClass}>{message}</p>}
      {results.length > 0 && (
        <div className="push-results">
          {results.map((result, index) => (
            <PushResultRow key={`${result.title}-${index}`} result={result} />
          ))}
        </div>
      )}
    </div>
  )
}

function summarizePushableItems(items) {
  if (items.length === 0) return ''

  const tasks = items.filter((item) => item.type !== 'event' && item.type !== 'note').length
  const events = items.filter((item) => item.type === 'event').length
  const notes = items.filter((item) => item.type === 'note').length
  const pushed = tasks + events
  const parts = []

  if (pushed > 0) {
    parts.push(`${pushed} ${pushed === 1 ? 'item' : 'items'} will be pushed`)
  }
  if (notes > 0) {
    parts.push(`${notes} ${notes === 1 ? 'note stays' : 'notes stay'} review-only`)
  }
  return parts.join('; ') + '.'
}

function PushResultRow({ result }) {
  const Icon = result.status === 'ok' ? CheckCircle2 : result.status === 'skipped' ? CircleSlash : AlertTriangle
  const detail = result.status === 'ok' ? result.type : result.error

  return (
    <div className={`push-result ${result.status}`}>
      <Icon size={16} />
      <span className="push-result-title">{result.title}</span>
      <span className="push-result-detail">{detail}</span>
    </div>
  )
}
