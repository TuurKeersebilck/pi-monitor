// Shared formatting/math helpers used by SystemStats and DockerList.

export function fmtBytes(bytesPerSec) {
  if (!bytesPerSec || bytesPerSec <= 0) return '0 B/s'
  if (bytesPerSec >= 1_048_576) return (bytesPerSec / 1_048_576).toFixed(1) + ' MB/s'
  if (bytesPerSec >= 1024) return (bytesPerSec / 1024).toFixed(0) + ' KB/s'
  return Math.round(bytesPerSec) + ' B/s'
}

// Same thresholds as fmtBytes but without the "/s" suffix, for chart axis
// ticks that show a plain byte count rather than a rate.
export function fmtBytesUnit(bytes) {
  if (!bytes || bytes <= 0) return '0 B'
  if (bytes >= 1_048_576) return (bytes / 1_048_576).toFixed(1) + ' MB'
  if (bytes >= 1024) return (bytes / 1024).toFixed(0) + ' KB'
  return Math.round(bytes) + ' B'
}

export function fmtMem(mb) {
  if (!mb) return '0 MB'
  if (mb >= 1024) return (mb / 1024).toFixed(1) + ' GB'
  return Math.round(mb) + ' MB'
}

export function avg(numbers) {
  if (!numbers?.length) return 0
  return numbers.reduce((a, b) => a + b, 0) / numbers.length
}
