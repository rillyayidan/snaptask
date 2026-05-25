import { CalendarClock, ClipboardCheck, Funnel, MinusCircle, Plus, StickyNote } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import { TaskCard } from './TaskCard.jsx'

const emptyItem = {
  type: 'task',
  title: '',
  detail: '',
  due_date: null,
  priority: 'medium'
}

export function TaskList({ items, onChange }) {
  const [filter, setFilter] = useState('all')
  const counts = items.reduce((summary, item) => {
    const type = item.type === 'event' || item.type === 'note' ? item.type : 'task'
    summary[type] += 1
    return summary
  }, { task: 0, event: 0, note: 0 })
  const activeFilterCount = filter === 'all' ? items.length : counts[filter]

  useEffect(() => {
    if (filter !== 'all' && activeFilterCount === 0) {
      setFilter('all')
    }
  }, [activeFilterCount, filter])

  const filteredItems = useMemo(() => {
    return items
      .map((item, index) => ({ item, index }))
      .filter(({ item }) => filter === 'all' || item.type === filter)
  }, [filter, items])

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
      {items.length > 0 && (
        <>
          <div className="review-metrics" aria-label="Extraction summary">
            <span><ClipboardCheck size={15} /> {counts.task} task{counts.task === 1 ? '' : 's'}</span>
            <span><CalendarClock size={15} /> {counts.event} event{counts.event === 1 ? '' : 's'}</span>
            <span><StickyNote size={15} /> {counts.note} note{counts.note === 1 ? '' : 's'}</span>
          </div>
          <div className="review-controls">
            <div className="filter-group" role="tablist" aria-label="Filter extracted items">
              {filterOptions(counts).map((option) => (
                <button
                  key={option.value}
                  className={`filter-chip ${filter === option.value ? 'active' : ''}`}
                  type="button"
                  disabled={option.value !== 'all' && option.count === 0}
                  onClick={() => setFilter(option.value)}
                >
                  {option.icon}
                  {option.label}
                </button>
              ))}
            </div>
            {counts.note > 0 && (
              <button
                className="text-button"
                type="button"
                onClick={() => onChange(items.filter((item) => item.type !== 'note'))}
              >
                <MinusCircle size={16} />
                Remove notes
              </button>
            )}
          </div>
        </>
      )}

      {items.length === 0 ? (
        <div className="review-empty">Extracted tasks, events, and notes will appear here before you push them to Google.</div>
      ) : (
        <div className="task-list">
          {filteredItems.map(({ item, index }) => (
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

function filterOptions(counts) {
  return [
    { value: 'all', label: `All ${counts.task + counts.event + counts.note}`, count: counts.task + counts.event + counts.note, icon: <Funnel size={14} /> },
    { value: 'task', label: `Tasks ${counts.task}`, count: counts.task, icon: <ClipboardCheck size={14} /> },
    { value: 'event', label: `Events ${counts.event}`, count: counts.event, icon: <CalendarClock size={14} /> },
    { value: 'note', label: `Notes ${counts.note}`, count: counts.note, icon: <StickyNote size={14} /> }
  ]
}
