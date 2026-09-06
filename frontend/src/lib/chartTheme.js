// Shared Chart.js theming helpers used by SystemStats and DockerList.

export function isDark() {
  return document.documentElement.classList.contains('dark')
}

// Tooltip color palette shared by every chart in the app. Charts with
// different sizes (full-size vs. sparkline) still set their own padding/
// cornerRadius/font on top of this -- those genuinely differ by design,
// only the colors are meant to be identical everywhere.
export function tooltipTheme(dark = isDark()) {
  return {
    backgroundColor: dark ? 'rgba(30,30,34,0.92)' : 'rgba(255,255,255,0.92)',
    titleColor: dark ? '#f0f0f5' : '#1c1c1e',
    bodyColor: dark ? '#a0a0a8' : '#6c6c70',
    borderColor: dark ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.06)',
  }
}

export function chartGridColor(dark = isDark()) {
  return dark ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.06)'
}

export function hexToRgba(hex, alpha) {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `rgba(${r},${g},${b},${alpha})`
}

// Top-to-bottom fading fill for a line chart dataset. topAlpha defaults to
// 0.3 (full-size charts); sparklines pass a lower value for a subtler fill.
export function makeGradient(ctx, hex, height, topAlpha = 0.3) {
  const grad = ctx.createLinearGradient(0, 0, 0, height)
  grad.addColorStop(0, hexToRgba(hex, topAlpha))
  grad.addColorStop(1, hexToRgba(hex, 0))
  return grad
}
