import { DateTime } from 'luxon'

export const isoToRelative = (iso: string): string => {
  const dt = DateTime.fromISO(iso, { locale: 'en' })
  return dt.toRelative() || dt.toLocaleString(DateTime.DATETIME_MED)
}
