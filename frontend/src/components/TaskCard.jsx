import { AlertTriangle, CalendarClock, ClipboardCheck, StickyNote, Trash2 } from 'lucide-react'

const icons = {
  task: ClipboardCheck,
  event: CalendarClock,
  note: StickyNote
}

function toDateTimeLocal(value) {
  if (!value) return ''
  if (/^\d{4}-\d{2}-\d{2}$/.test(value)) return `${value}T00:00`

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''

  const pad = (part) => String(part).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

function fromDateTimeLocal(value) {
  if (!value) return null
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? null : date.toISOString()
}

export function TaskCard({ item, index, onUpdate, onDelete }) {
  const Icon = icons[item.type] ?? ClipboardCheck
  const missingEventTime = item.type === 'event' && !item.due_date
  return (
    <article className="task-card">
      <div className="card-head">
        <span className={`type-pill ${item.type}`}>
          <Icon size={15} />
          {item.type}
        </span>
        <button className="small-icon" title="Delete item" onClick={onDelete}>
          <Trash2 size={16} />
        </button>
      </div>

      <label>
        Title
        <input
          aria-invalid={!item.title?.trim()}
          className={!item.title?.trim() ? 'invalid-field' : ''}
          value={item.title}
          onChange={(event) => onUpdate(index, { title: event.target.value })}
        />
      </label>
      <label>
        Detail
        <textarea value={item.detail} rows={3} onChange={(event) => onUpdate(index, { detail: event.target.value })} />
      </label>
      <div className="field-grid">
        <label>
          Type
          <select value={item.type} onChange={(event) => onUpdate(index, { type: event.target.value })}>
            <option value="task">Task</option>
            <option value="event">Event</option>
            <option value="note">Note</option>
          </select>
        </label>
        <label>
          Priority
          <select value={item.priority} onChange={(event) => onUpdate(index, { priority: event.target.value })}>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </label>
      </div>
      <label>
        Due date
        <input
          type="datetime-local"
          value={toDateTimeLocal(item.due_date)}
          onChange={(event) => onUpdate(index, { due_date: fromDateTimeLocal(event.target.value) })}
        />
      </label>
      {missingEventTime && (
        <p className="item-warning">
          <AlertTriangle size={15} />
          Add a date or time before pushing this event to Calendar.
        </p>
      )}
    </article>
  )
}
