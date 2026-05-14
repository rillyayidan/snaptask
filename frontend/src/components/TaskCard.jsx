import { CalendarClock, ClipboardCheck, StickyNote, Trash2 } from 'lucide-react'

const icons = {
  task: ClipboardCheck,
  event: CalendarClock,
  note: StickyNote
}

export function TaskCard({ item, index, onUpdate, onDelete }) {
  const Icon = icons[item.type] ?? ClipboardCheck
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
        <input value={item.title} onChange={(event) => onUpdate(index, { title: event.target.value })} />
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
        <input value={item.due_date ?? ''} placeholder="2026-05-15T15:00:00+07:00" onChange={(event) => onUpdate(index, { due_date: event.target.value || null })} />
      </label>
    </article>
  )
}
