import { AlertTriangle, CalendarClock, ClipboardCheck, StickyNote, Trash2, X } from 'lucide-react'
import { fromDateTimeLocalInput, isDueDateValueValid, toDateTimeLocalInput } from '../lib/dates.js'

const icons = {
  task: ClipboardCheck,
  event: CalendarClock,
  note: StickyNote
}

export function TaskCard({ item, index, onUpdate, onDelete }) {
  const itemType = reviewItemType(item.type)
  const priority = reviewPriority(item.priority)
  const Icon = icons[itemType] ?? ClipboardCheck
  const missingEventTime = itemType === 'event' && !item.due_date
  const invalidEventTime = itemType === 'event' && Boolean(item.due_date) && !isDueDateValueValid(item.due_date)
  return (
    <article className="task-card">
      <div className="card-head">
        <span className={`type-pill ${itemType}`}>
          <Icon size={15} />
          {itemType}
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
          <select value={itemType} onChange={(event) => onUpdate(index, { type: event.target.value })}>
            <option value="task">Task</option>
            <option value="event">Event</option>
            <option value="note">Note</option>
          </select>
        </label>
        <label>
          Priority
          <select value={priority} onChange={(event) => onUpdate(index, { priority: event.target.value })}>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </label>
      </div>
      <label>
        Due date
        <span className="input-action">
          <input
            type="datetime-local"
            className={invalidEventTime ? 'invalid-field' : ''}
            value={toDateTimeLocalInput(item.due_date)}
            onChange={(event) => onUpdate(index, { due_date: fromDateTimeLocalInput(event.target.value) })}
          />
          <button
            className="small-icon"
            type="button"
            title="Clear due date"
            disabled={!item.due_date}
            onClick={() => onUpdate(index, { due_date: null })}
          >
            <X size={16} />
          </button>
        </span>
      </label>
      {(missingEventTime || invalidEventTime) && (
        <p className="item-warning">
          <AlertTriangle size={15} />
          {invalidEventTime ? 'Fix this date before pushing the event to Calendar.' : 'Add a date or time before pushing this event to Calendar.'}
        </p>
      )}
    </article>
  )
}

function reviewItemType(value) {
  return value === 'event' || value === 'note' ? value : 'task'
}

function reviewPriority(value) {
  return value === 'high' || value === 'low' ? value : 'medium'
}
