import { Plus } from 'lucide-react'
import { TaskCard } from './TaskCard.jsx'

const emptyItem = {
  type: 'task',
  title: '',
  detail: '',
  due_date: null,
  priority: 'medium'
}

export function TaskList({ items, onChange }) {
  function update(index, patch) {
    onChange(items.map((item, itemIndex) => (itemIndex === index ? { ...item, ...patch } : item)))
  }

  return (
    <div className="review-stack">
      <div className="section-head">
        <div>
          <p className="eyebrow">Review</p>
          <h2>{items.length ? `${items.length} extracted item${items.length === 1 ? '' : 's'}` : 'No extracted items yet'}</h2>
        </div>
        <button className="icon-button" title="Add item" onClick={() => onChange([...items, emptyItem])}>
          <Plus size={18} />
        </button>
      </div>

      {items.length === 0 ? (
        <div className="review-empty">Extracted tasks, events, and notes will appear here before you push them to Google.</div>
      ) : (
        <div className="task-list">
          {items.map((item, index) => (
            <TaskCard
              key={`${item.title}-${index}`}
              item={item}
              index={index}
              onUpdate={update}
              onDelete={() => onChange(items.filter((_, itemIndex) => itemIndex !== index))}
            />
          ))}
        </div>
      )}
    </div>
  )
}
