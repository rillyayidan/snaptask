const dateOnlyPattern = /^(\d{4})-(\d{2})-(\d{2})$/

export function isDueDateValueValid(value) {
  const raw = typeof value === 'string' ? value.trim() : ''
  if (!raw) return false

  const dateOnly = raw.match(dateOnlyPattern)
  if (dateOnly) {
    return isValidCalendarDate(Number(dateOnly[1]), Number(dateOnly[2]), Number(dateOnly[3]))
  }

  return !Number.isNaN(new Date(raw).getTime())
}

export function toDateTimeLocalInput(value) {
  if (!value || !isDueDateValueValid(value)) return ''
  const raw = value.trim()
  if (dateOnlyPattern.test(raw)) return `${raw}T00:00`

  const date = new Date(raw)
  const pad = (part) => String(part).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

export function fromDateTimeLocalInput(value) {
  if (!value) return null

  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? null : date.toISOString()
}

function isValidCalendarDate(year, month, day) {
  const date = new Date(Date.UTC(year, month - 1, day))
  return date.getUTCFullYear() === year &&
    date.getUTCMonth() === month - 1 &&
    date.getUTCDate() === day
}
