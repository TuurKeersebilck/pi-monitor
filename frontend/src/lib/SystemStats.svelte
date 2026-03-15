<script>
  import { tweened } from 'svelte/motion'
  import { cubicOut } from 'svelte/easing'
  import { fade } from 'svelte/transition'
  import { onMount, onDestroy } from 'svelte'
  import { Chart, LineController, LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip, Legend } from 'chart.js'

  Chart.register(LineController, LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip, Legend)

  let { system, history = [] } = $props()

  let view = $state('gauge')

  const PATH = "M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"

  // Tweened gauge values for smooth animation
  const cpuTween = tweened(0, { duration: 600, easing: cubicOut })
  const ramTween = tweened(0, { duration: 600, easing: cubicOut })
  const diskTween = tweened(0, { duration: 600, easing: cubicOut })
  const tempTween = tweened(0, { duration: 600, easing: cubicOut })

  $effect(() => {
    if (system?.cpu) cpuTween.set(system.cpu.usage_percent)
    if (system?.ram) ramTween.set(system.ram.usage_percent)
    if (system?.disk) diskTween.set(system.disk.usage_percent)
    if (system?.temp) tempTween.set(system.temp.cpu_temp_c)
  })

  function tempColor(t) {
    if (t >= 75) return '#ff3b30'
    if (t >= 60) return '#ff9500'
    return '#34c759'
  }

  function fmtBytes(bps) {
    if (bps >= 1_048_576) return (bps / 1_048_576).toFixed(1) + ' MB/s'
    if (bps >= 1024) return (bps / 1024).toFixed(0) + ' KB/s'
    return bps.toFixed(0) + ' B/s'
  }

  // ── Chart.js setup ──
  let cpuCanvas = $state()
  let ramCanvas = $state()
  let diskCanvas = $state()
  let tempCanvas = $state()
  let netCanvas = $state()
  let cpuChart, ramChart, diskChart, tempChart, netChart

  const MAX_POINTS = 120

  function isDark() {
    return document.documentElement.classList.contains('dark')
  }

  function chartTextColor() {
    return isDark() ? 'rgba(240,240,245,0.55)' : 'rgba(28,28,30,0.55)'
  }

  function chartGridColor() {
    return isDark() ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.06)'
  }

  function createGradient(ctx, color, height) {
    const gradient = ctx.createLinearGradient(0, 0, 0, height)
    gradient.addColorStop(0, color.replace(')', ',0.3)').replace('rgb', 'rgba'))
    gradient.addColorStop(1, color.replace(')', ',0)').replace('rgb', 'rgba'))
    return gradient
  }

  function hexToRgba(hex, alpha) {
    const r = parseInt(hex.slice(1, 3), 16)
    const g = parseInt(hex.slice(3, 5), 16)
    const b = parseInt(hex.slice(5, 7), 16)
    return `rgba(${r},${g},${b},${alpha})`
  }

  function makeGradient(ctx, hex, h) {
    const grad = ctx.createLinearGradient(0, 0, 0, h)
    grad.addColorStop(0, hexToRgba(hex, 0.3))
    grad.addColorStop(1, hexToRgba(hex, 0))
    return grad
  }

  function baseOptions(min = 0, max = 100) {
    return {
      responsive: true,
      maintainAspectRatio: false,
      animation: { duration: 400, easing: 'easeOutCubic' },
      interaction: { mode: 'index', intersect: false },
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: isDark() ? 'rgba(30,30,34,0.92)' : 'rgba(255,255,255,0.92)',
          titleColor: isDark() ? '#f0f0f5' : '#1c1c1e',
          bodyColor: isDark() ? '#a0a0a8' : '#6c6c70',
          borderColor: isDark() ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.06)',
          borderWidth: 1,
          padding: 8,
          cornerRadius: 8,
          titleFont: { size: 11, weight: 600 },
          bodyFont: { size: 11 },
          displayColors: true,
          boxWidth: 8,
          boxHeight: 8,
          boxPadding: 4,
        },
      },
      scales: {
        x: {
          display: false,
        },
        y: {
          min,
          max,
          display: true,
          grid: { color: chartGridColor(), drawBorder: false },
          border: { display: false },
          ticks: {
            color: chartTextColor(),
            font: { size: 10, weight: 500 },
            maxTicksLimit: 4,
            padding: 4,
          },
        },
      },
      elements: {
        point: { radius: 0, hoverRadius: 4, hoverBorderWidth: 2 },
        line: { tension: 0.4, borderWidth: 2 },
      },
    }
  }

  function makeDataset(label, color, canvasEl) {
    const ctx = canvasEl.getContext('2d')
    const h = canvasEl.parentElement?.clientHeight || 72
    return {
      label,
      data: [],
      borderColor: color,
      backgroundColor: makeGradient(ctx, color, h),
      fill: true,
      pointBackgroundColor: color,
    }
  }

  function createChart(canvas, label, color, opts = {}) {
    if (!canvas) return null
    const options = baseOptions(opts.min ?? 0, opts.max ?? 100)
    if (opts.suffix) {
      options.scales.y.ticks.callback = (v) => v + opts.suffix
      options.plugins.tooltip.callbacks = {
        label: (ctx) => `${ctx.dataset.label}: ${ctx.parsed.y}${opts.suffix}`
      }
    }
    return new Chart(canvas, {
      type: 'line',
      data: {
        labels: [],
        datasets: [makeDataset(label, color, canvas)],
      },
      options,
    })
  }

  function createNetChart(canvas) {
    if (!canvas) return null
    const ctx = canvas.getContext('2d')
    const h = canvas.parentElement?.clientHeight || 72
    const options = baseOptions(0, 100)
    options.scales.y.ticks.callback = (v) => {
      if (v >= 1048576) return (v / 1048576).toFixed(1) + ' MB'
      if (v >= 1024) return (v / 1024).toFixed(0) + ' KB'
      return v + ' B'
    }
    options.plugins.tooltip.callbacks = {
      label: (ctx) => `${ctx.dataset.label}: ${fmtBytes(ctx.parsed.y)}`
    }
    return new Chart(canvas, {
      type: 'line',
      data: {
        labels: [],
        datasets: [
          {
            label: '↑ TX',
            data: [],
            borderColor: '#34c759',
            backgroundColor: makeGradient(ctx, '#34c759', h),
            fill: true,
            pointBackgroundColor: '#34c759',
          },
          {
            label: '↓ RX',
            data: [],
            borderColor: '#007AFF',
            backgroundColor: makeGradient(ctx, '#007AFF', h),
            fill: true,
            pointBackgroundColor: '#007AFF',
          },
        ],
      },
      options,
    })
  }

  function initCharts() {
    cpuChart = createChart(cpuCanvas, 'CPU', '#007AFF', { suffix: '%' })
    ramChart = createChart(ramCanvas, 'RAM', '#34c759', { suffix: '%' })
    diskChart = createChart(diskCanvas, 'Disk', '#ff9500', { suffix: '%' })
    tempChart = createChart(tempCanvas, 'Temp', '#34c759', { min: 20, max: 90, suffix: '°C' })
    netChart = createNetChart(netCanvas)
  }

  function updateChart(chart, values, maxVal) {
    if (!chart) return
    const labels = values.map((_, i) => i)
    chart.data.labels = labels
    chart.data.datasets[0].data = values
    if (maxVal !== undefined) {
      chart.options.scales.y.max = maxVal
    }
    chart.update('none') // 'none' = skip animation for streaming data
  }

  function updateNetChart(chart, tx, rx) {
    if (!chart) return
    const labels = tx.map((_, i) => i)
    const maxVal = Math.max(...tx, ...rx, 1024) * 1.2
    chart.data.labels = labels
    chart.data.datasets[0].data = tx
    chart.data.datasets[1].data = rx
    chart.options.scales.y.max = maxVal
    chart.update('none')
  }

  let chartsReady = false

  $effect(() => {
    if (view === 'graph' && !chartsReady) {
      // Wait a tick for canvas elements to mount
      requestAnimationFrame(() => {
        initCharts()
        chartsReady = true
        // Populate with existing history
        if (history.length > 0) {
          updateChart(cpuChart, history.map(h => h.cpu?.usage_percent ?? 0))
          updateChart(ramChart, history.map(h => h.ram?.usage_percent ?? 0))
          updateChart(diskChart, history.map(h => h.disk?.usage_percent ?? 0))
          updateChart(tempChart, history.map(h => h.temp?.cpu_temp_c ?? 0))
          const tx = history.map(h => h.net?.tx_bytes_s ?? 0)
          const rx = history.map(h => h.net?.rx_bytes_s ?? 0)
          updateNetChart(netChart, tx, rx)
        }
      })
    }
  })

  // Update charts when history changes
  $effect(() => {
    if (!chartsReady || view !== 'graph') return
    const h = history // track dependency
    updateChart(cpuChart, h.map(d => d.cpu?.usage_percent ?? 0))
    updateChart(ramChart, h.map(d => d.ram?.usage_percent ?? 0))
    updateChart(diskChart, h.map(d => d.disk?.usage_percent ?? 0))
    updateChart(tempChart, h.map(d => d.temp?.cpu_temp_c ?? 0))
    const tx = h.map(d => d.net?.tx_bytes_s ?? 0)
    const rx = h.map(d => d.net?.rx_bytes_s ?? 0)
    updateNetChart(netChart, tx, rx)
  })

  let windowLabel = $derived(() => {
    const s = history.length
    if (s < 60) return `${s}s`
    return `${Math.floor(s / 60)}m ${s % 60}s`
  })

  function avg(arr) {
    if (!arr.length) return 0
    return arr.reduce((a, b) => a + b, 0) / arr.length
  }

  onDestroy(() => {
    cpuChart?.destroy()
    ramChart?.destroy()
    diskChart?.destroy()
    tempChart?.destroy()
    netChart?.destroy()
  })
</script>

<!-- Shared SVG filter for ring glow (rendered once, invisible) -->
<svg style="position:absolute;width:0;height:0" aria-hidden="true">
  <defs>
    <filter id="ring-glow" x="-50%" y="-50%" width="200%" height="200%">
      <feGaussianBlur stdDeviation="1.2" result="blur"/>
      <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
    </filter>
  </defs>
</svg>

<!-- Toggle -->
<div class="top-bar">
  <div class="view-toggle card-surface">
    <button class:active={view === 'gauge'} onclick={() => view = 'gauge'}>Gauge</button>
    <button class:active={view === 'graph'} onclick={() => view = 'graph'}>Graph</button>
  </div>
</div>

<!-- Conditional rendering to prevent stretching -->
{#if view === 'gauge'}
  <div class="gauges" in:fade={{ duration: 200 }}>
    <div class="gauge-card card-surface">
      <svg class="ring" viewBox="0 0 36 36">
        <path class="bg" d={PATH}/>
        <path class="fill" style="stroke:#007AFF" stroke-dasharray="{$cpuTween},100" d={PATH} filter="url(#ring-glow)"/>
        <text x="18" y="16.5" class="val">{Math.round($cpuTween)}%</text>
      </svg>
      <span class="label">CPU</span>
    </div>

    <div class="gauge-card card-surface">
      <svg class="ring" viewBox="0 0 36 36">
        <path class="bg" d={PATH}/>
        <path class="fill" style="stroke:#34c759" stroke-dasharray="{$ramTween},100" d={PATH} filter="url(#ring-glow)"/>
        <text x="18" y="16.5" class="val">{Math.round($ramTween)}%</text>
      </svg>
      <span class="label">RAM</span>
      <span class="sub">{system.ram.used_mb} / {system.ram.total_mb} MB</span>
    </div>

    <div class="gauge-card card-surface">
      <svg class="ring" viewBox="0 0 36 36">
        <path class="bg" d={PATH}/>
        <path class="fill" style="stroke:#ff9500" stroke-dasharray="{$diskTween},100" d={PATH} filter="url(#ring-glow)"/>
        <text x="18" y="16.5" class="val">{Math.round($diskTween)}%</text>
      </svg>
      <span class="label">Disk</span>
      <span class="sub">{system.disk.used_gb} / {system.disk.total_gb} GB</span>
    </div>

    {#if system.temp?.cpu_temp_c}
      <div class="gauge-card card-surface">
        <svg class="ring" viewBox="0 0 36 36">
          <path class="bg" d={PATH}/>
          <path class="fill" style="stroke:{tempColor($tempTween)}" stroke-dasharray="{$tempTween},100" d={PATH} filter="url(#ring-glow)"/>
          <text x="18" y="16.5" class="val">{Math.round($tempTween)}°</text>
        </svg>
        <span class="label">Temp</span>
        <span class="sub">{system.temp.cpu_temp_c}°C</span>
      </div>
    {/if}

    {#if system.net}
      <div class="gauge-card card-surface net-card">
        <div class="net-rows">
          <div class="net-row">
            <svg viewBox="0 0 24 24" fill="none" stroke="#34c759" stroke-width="2.5" width="15" height="15">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/>
            </svg>
            <span class="net-val">{fmtBytes(system.net.tx_bytes_s)}</span>
          </div>
          <div class="net-row">
            <svg viewBox="0 0 24 24" fill="none" stroke="#007AFF" stroke-width="2.5" width="15" height="15">
              <polyline points="23 18 13.5 8.5 8.5 13.5 1 6"/><polyline points="17 18 23 18 23 12"/>
            </svg>
            <span class="net-val">{fmtBytes(system.net.rx_bytes_s)}</span>
          </div>
        </div>
        <span class="label">Network</span>
        <span class="sub">↑ up · ↓ down</span>
      </div>
    {/if}
  </div>
{:else}
  <div class="graphs" in:fade={{ duration: 200 }}>
    <div class="graph-card card-surface">
      <div class="graph-header">
        <span class="graph-label">CPU</span>
        <span class="graph-current" style="color:#007AFF">{system.cpu.usage_percent}%</span>
      </div>
      <div class="chart-wrap"><canvas bind:this={cpuCanvas}></canvas></div>
      <div class="graph-footer">
        <span class="graph-meta">avg {avg(history.map(h => h.cpu?.usage_percent ?? 0)).toFixed(1)}%</span>
        <span class="graph-meta dim">last {windowLabel()}</span>
      </div>
    </div>

    <div class="graph-card card-surface">
      <div class="graph-header">
        <span class="graph-label">RAM</span>
        <span class="graph-current" style="color:#34c759">{system.ram.usage_percent}%</span>
      </div>
      <div class="chart-wrap"><canvas bind:this={ramCanvas}></canvas></div>
      <div class="graph-footer">
        <span class="graph-meta">avg {avg(history.map(h => h.ram?.usage_percent ?? 0)).toFixed(1)}%</span>
        <span class="graph-meta dim">{system.ram.used_mb} / {system.ram.total_mb} MB · {windowLabel()}</span>
      </div>
    </div>

    <div class="graph-card card-surface">
      <div class="graph-header">
        <span class="graph-label">Disk</span>
        <span class="graph-current" style="color:#ff9500">{system.disk.usage_percent}%</span>
      </div>
      <div class="chart-wrap"><canvas bind:this={diskCanvas}></canvas></div>
      <div class="graph-footer">
        <span class="graph-meta">avg {avg(history.map(h => h.disk?.usage_percent ?? 0)).toFixed(1)}%</span>
        <span class="graph-meta dim">{system.disk.used_gb} / {system.disk.total_gb} GB · {windowLabel()}</span>
      </div>
    </div>

    {#if system.temp?.cpu_temp_c}
      <div class="graph-card card-surface">
        <div class="graph-header">
          <span class="graph-label">Temp</span>
          <span class="graph-current" style="color:{tempColor(system.temp.cpu_temp_c)}">{system.temp.cpu_temp_c}°C</span>
        </div>
        <div class="chart-wrap"><canvas bind:this={tempCanvas}></canvas></div>
        <div class="graph-footer">
          <span class="graph-meta">avg {avg(history.map(h => h.temp?.cpu_temp_c ?? 0)).toFixed(1)}°C</span>
          <span class="graph-meta dim">cpu temp · {windowLabel()}</span>
        </div>
      </div>
    {/if}

    {#if system.net}
      <div class="graph-card card-surface graph-card--wide">
        <div class="graph-header">
          <span class="graph-label">Network</span>
          <div class="net-legend">
            <span style="color:#34c759">↑ {fmtBytes(system.net.tx_bytes_s)}</span>
            <span style="color:#007AFF">↓ {fmtBytes(system.net.rx_bytes_s)}</span>
          </div>
        </div>
        <div class="chart-wrap"><canvas bind:this={netCanvas}></canvas></div>
        <div class="graph-footer">
          <span class="graph-meta" style="color:#34c759">↑ avg {fmtBytes(avg(history.map(h => h.net?.tx_bytes_s ?? 0)))}</span>
          <span class="graph-meta dim">{windowLabel()}</span>
          <span class="graph-meta" style="color:#007AFF">↓ avg {fmtBytes(avg(history.map(h => h.net?.rx_bytes_s ?? 0)))}</span>
        </div>
      </div>
    {/if}
  </div>
{/if}

<style>
  /* ── Toggle ── */
  .top-bar {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 0.875rem;
  }

  .view-toggle {
    display: flex;
    padding: 3px;
    gap: 2px;
    border-radius: 10px;
  }

  .view-toggle button {
    padding: 0.28rem 0.85rem;
    border-radius: 7px;
    border: none;
    background: transparent;
    color: var(--text-2);
    font-size: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    letter-spacing: -0.01em;
    transition: background 0.18s ease, color 0.18s ease, box-shadow 0.18s ease;
  }

  .view-toggle button.active {
    background: var(--pill-hover);
    color: var(--text);
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  }

  /* ── Gauges ── */
  .gauges {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(145px, 1fr));
    gap: 0.875rem;
  }

  .gauge-card {
    padding: 1.5rem 1rem 1.25rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
  }

  .ring { width: 120px; height: 120px; overflow: visible; }

  .bg {
    fill: none;
    stroke: var(--ring-track);
    stroke-width: 2.8;
  }

  .fill {
    fill: none;
    stroke-width: 2.8;
    stroke-linecap: round;
    transition: stroke-dasharray 0.5s cubic-bezier(0.34, 1.2, 0.64, 1);
  }

  .val {
    fill: var(--text);
    font-size: 7.5px;
    font-weight: 700;
    text-anchor: middle;
    dominant-baseline: middle;
    letter-spacing: -0.02em;
  }

  .label {
    font-size: 0.78rem;
    color: var(--text);
    font-weight: 600;
    text-align: center;
    letter-spacing: -0.01em;
    opacity: 0.75;
  }

  .sub {
    font-size: 0.67rem;
    color: var(--text);
    text-align: center;
    opacity: 0.45;
  }

  .net-card { justify-content: center; gap: 0.875rem; }

  .net-rows {
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
    width: 100%;
    padding: 0 0.5rem;
  }

  .net-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .net-val {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--text);
    letter-spacing: -0.02em;
    white-space: nowrap;
  }

  /* ── Graphs (Chart.js) ── */
  .graphs {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 0.875rem;
  }

  .graph-card--wide { grid-column: span 2; }

  .graph-card {
    padding: 1rem 1.1rem 0.85rem;
    display: flex;
    flex-direction: column;
    gap: 0.55rem;
  }

  .chart-wrap {
    position: relative;
    width: 100%;
    height: 80px;
  }

  .chart-wrap canvas {
    width: 100% !important;
    height: 100% !important;
  }

  .graph-header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
  }

  .graph-label {
    font-size: 0.72rem;
    font-weight: 700;
    color: var(--text);
    opacity: 0.6;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .graph-current {
    font-size: 1.1rem;
    font-weight: 700;
    letter-spacing: -0.03em;
  }

  .graph-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .graph-meta {
    font-size: 0.68rem;
    font-weight: 600;
    color: var(--text);
    opacity: 0.55;
    letter-spacing: -0.01em;
  }

  .graph-meta.dim { opacity: 0.35; }

  .net-legend {
    display: flex;
    gap: 0.75rem;
    font-size: 0.8rem;
    font-weight: 600;
    letter-spacing: -0.02em;
  }

  @media (max-width: 768px) {
    .gauges { grid-template-columns: repeat(2, 1fr); }
    .graphs { grid-template-columns: 1fr; }
    .graph-card--wide { grid-column: span 1; }
  }
</style>
