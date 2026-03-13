<script>
  import { fade } from 'svelte/transition'

  let { system, history = [] } = $props()

  let view = $state('gauge') // 'gauge' | 'graph'

  const PATH = "M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"

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

  function avg(values) {
    if (!values.length) return 0
    return values.reduce((a, b) => a + b, 0) / values.length
  }

  function sparkPoints(values, W, H, minV, maxV) {
    if (values.length < 2) return { points: '', avgY: H / 2 }
    const range = maxV - minV || 1
    const points = values.map((v, i) => {
      const x = (i / (values.length - 1)) * W
      const y = H - ((v - minV) / range) * (H - 4) - 2
      return `${x.toFixed(1)},${y.toFixed(1)}`
    }).join(' ')
    const a = avg(values)
    const avgY = H - ((a - minV) / range) * (H - 4) - 2
    return { points, avgY, avgVal: a }
  }

  let cpuSeries    = $derived(history.map(h => h.cpu?.usage_percent ?? 0))
  let ramSeries    = $derived(history.map(h => h.ram?.usage_percent ?? 0))
  let diskSeries   = $derived(history.map(h => h.disk?.usage_percent ?? 0))
  let tempSeries   = $derived(history.map(h => h.temp?.cpu_temp_c ?? 0))
  let txSeries     = $derived(history.map(h => h.net?.tx_bytes_s ?? 0))
  let rxSeries     = $derived(history.map(h => h.net?.rx_bytes_s ?? 0))

  const W = 260, H = 72

  let cpuSpk  = $derived(sparkPoints(cpuSeries, W, H, 0, 100))
  let ramSpk  = $derived(sparkPoints(ramSeries, W, H, 0, 100))
  let diskSpk = $derived(sparkPoints(diskSeries, W, H, 0, 100))
  let tempSpk = $derived(sparkPoints(tempSeries, W, H, 20, 90))
  let netMax  = $derived(Math.max(...txSeries, ...rxSeries, 1))
  let txSpk   = $derived(sparkPoints(txSeries, W, H, 0, netMax))
  let rxSpk   = $derived(sparkPoints(rxSeries, W, H, 0, netMax))

  // Human-readable time window
  let windowLabel = $derived(() => {
    const s = history.length
    if (s < 60) return `${s}s`
    return `${Math.floor(s / 60)}m ${s % 60}s`
  })
</script>

<!-- Toggle -->
<div class="top-bar">
  <div class="view-toggle">
    <button class:active={view === 'gauge'} onclick={() => view = 'gauge'}>Gauge</button>
    <button class:active={view === 'graph'} onclick={() => view = 'graph'}>Graph</button>
  </div>
</div>

<!-- Gauge view -->
{#if view === 'gauge'}
  <div class="gauges" transition:fade={{ duration: 180 }}>
    <!-- CPU -->
    <div class="gauge-card">
      <svg class="ring" viewBox="0 0 36 36">
        <defs>
          <filter id="glow-cpu" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur stdDeviation="1.2" result="blur"/>
            <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
          </filter>
        </defs>
        <path class="bg" d={PATH}/>
        <path class="fill" style="stroke:#007AFF" stroke-dasharray="{system.cpu.usage_percent},100" d={PATH} filter="url(#glow-cpu)"/>
        <text x="18" y="16.5" class="val">{system.cpu.usage_percent}%</text>
      </svg>
      <span class="label">CPU</span>
    </div>

    <!-- RAM -->
    <div class="gauge-card">
      <svg class="ring" viewBox="0 0 36 36">
        <defs>
          <filter id="glow-ram" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur stdDeviation="1.2" result="blur"/>
            <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
          </filter>
        </defs>
        <path class="bg" d={PATH}/>
        <path class="fill" style="stroke:#34c759" stroke-dasharray="{system.ram.usage_percent},100" d={PATH} filter="url(#glow-ram)"/>
        <text x="18" y="16.5" class="val">{system.ram.usage_percent}%</text>
      </svg>
      <span class="label">RAM</span>
      <span class="sub">{system.ram.used_mb} / {system.ram.total_mb} MB</span>
    </div>

    <!-- Disk -->
    <div class="gauge-card">
      <svg class="ring" viewBox="0 0 36 36">
        <defs>
          <filter id="glow-disk" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur stdDeviation="1.2" result="blur"/>
            <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
          </filter>
        </defs>
        <path class="bg" d={PATH}/>
        <path class="fill" style="stroke:#ff9500" stroke-dasharray="{system.disk.usage_percent},100" d={PATH} filter="url(#glow-disk)"/>
        <text x="18" y="16.5" class="val">{system.disk.usage_percent}%</text>
      </svg>
      <span class="label">Disk</span>
      <span class="sub">{system.disk.used_gb} / {system.disk.total_gb} GB</span>
    </div>

    <!-- Temp -->
    {#if system.temp?.cpu_temp_c}
      <div class="gauge-card">
        <svg class="ring" viewBox="0 0 36 36">
          <defs>
            <filter id="glow-temp" x="-50%" y="-50%" width="200%" height="200%">
              <feGaussianBlur stdDeviation="1.2" result="blur"/>
              <feMerge><feMergeNode in="blur"/><feMergeNode in="SourceGraphic"/></feMerge>
            </filter>
          </defs>
          <path class="bg" d={PATH}/>
          <path class="fill" style="stroke:{tempColor(system.temp.cpu_temp_c)}" stroke-dasharray="{system.temp.cpu_temp_c},100" d={PATH} filter="url(#glow-temp)"/>
          <text x="18" y="16.5" class="val">{system.temp.cpu_temp_c}°</text>
        </svg>
        <span class="label">Temp</span>
        <span class="sub">{system.temp.cpu_temp_c}°C</span>
      </div>
    {/if}

    <!-- Network -->
    {#if system.net}
      <div class="gauge-card net-card">
        <div class="net-rows">
          <div class="net-row">
            <svg viewBox="0 0 24 24" fill="none" stroke="#34c759" stroke-width="2.5" width="15" height="15">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
              <polyline points="17 6 23 6 23 12"/>
            </svg>
            <span class="net-val">{fmtBytes(system.net.tx_bytes_s)}</span>
          </div>
          <div class="net-row">
            <svg viewBox="0 0 24 24" fill="none" stroke="#007AFF" stroke-width="2.5" width="15" height="15">
              <polyline points="23 18 13.5 8.5 8.5 13.5 1 6"/>
              <polyline points="17 18 23 18 23 12"/>
            </svg>
            <span class="net-val">{fmtBytes(system.net.rx_bytes_s)}</span>
          </div>
        </div>
        <span class="label">Network</span>
        <span class="sub">↑ up · ↓ down</span>
      </div>
    {/if}
  </div>
{/if}

<!-- Graph view -->
{#if view === 'graph'}
  <div class="graphs" transition:fade={{ duration: 180 }}>
    <!-- CPU graph -->
    <div class="graph-card">
      <div class="graph-header">
        <span class="graph-label">CPU</span>
        <span class="graph-current" style="color:#007AFF">{system.cpu.usage_percent}%</span>
      </div>
      <svg class="sparkline" viewBox="0 0 {W} {H}" preserveAspectRatio="none">
        <defs>
          <linearGradient id="g-cpu" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#007AFF" stop-opacity="0.3"/>
            <stop offset="100%" stop-color="#007AFF" stop-opacity="0"/>
          </linearGradient>
        </defs>
        {#if cpuSpk.points}
          <polygon points="{cpuSpk.points} {W},{H} 0,{H}" fill="url(#g-cpu)"/>
          <polyline points={cpuSpk.points} fill="none" stroke="#007AFF" stroke-width="1.8" stroke-linejoin="round" stroke-linecap="round"/>
          <line x1="0" x2={W} y1={cpuSpk.avgY} y2={cpuSpk.avgY} stroke="#007AFF" stroke-width="0.8" stroke-dasharray="4,3" opacity="0.45"/>
        {/if}
      </svg>
      <div class="graph-footer">
        <span class="avg-label">avg {cpuSpk.avgVal?.toFixed(1) ?? '—'}%</span>
        <span class="graph-sub">last {windowLabel()}</span>
      </div>
    </div>

    <!-- RAM graph -->
    <div class="graph-card">
      <div class="graph-header">
        <span class="graph-label">RAM</span>
        <span class="graph-current" style="color:#34c759">{system.ram.usage_percent}%</span>
      </div>
      <svg class="sparkline" viewBox="0 0 {W} {H}" preserveAspectRatio="none">
        <defs>
          <linearGradient id="g-ram" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#34c759" stop-opacity="0.3"/>
            <stop offset="100%" stop-color="#34c759" stop-opacity="0"/>
          </linearGradient>
        </defs>
        {#if ramSpk.points}
          <polygon points="{ramSpk.points} {W},{H} 0,{H}" fill="url(#g-ram)"/>
          <polyline points={ramSpk.points} fill="none" stroke="#34c759" stroke-width="1.8" stroke-linejoin="round" stroke-linecap="round"/>
          <line x1="0" x2={W} y1={ramSpk.avgY} y2={ramSpk.avgY} stroke="#34c759" stroke-width="0.8" stroke-dasharray="4,3" opacity="0.45"/>
        {/if}
      </svg>
      <div class="graph-footer">
        <span class="avg-label">avg {ramSpk.avgVal?.toFixed(1) ?? '—'}%</span>
        <span class="graph-sub">{system.ram.used_mb} / {system.ram.total_mb} MB · {windowLabel()}</span>
      </div>
    </div>

    <!-- Disk graph -->
    <div class="graph-card">
      <div class="graph-header">
        <span class="graph-label">Disk</span>
        <span class="graph-current" style="color:#ff9500">{system.disk.usage_percent}%</span>
      </div>
      <svg class="sparkline" viewBox="0 0 {W} {H}" preserveAspectRatio="none">
        <defs>
          <linearGradient id="g-disk" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#ff9500" stop-opacity="0.3"/>
            <stop offset="100%" stop-color="#ff9500" stop-opacity="0"/>
          </linearGradient>
        </defs>
        {#if diskSpk.points}
          <polygon points="{diskSpk.points} {W},{H} 0,{H}" fill="url(#g-disk)"/>
          <polyline points={diskSpk.points} fill="none" stroke="#ff9500" stroke-width="1.8" stroke-linejoin="round" stroke-linecap="round"/>
          <line x1="0" x2={W} y1={diskSpk.avgY} y2={diskSpk.avgY} stroke="#ff9500" stroke-width="0.8" stroke-dasharray="4,3" opacity="0.45"/>
        {/if}
      </svg>
      <div class="graph-footer">
        <span class="avg-label">avg {diskSpk.avgVal?.toFixed(1) ?? '—'}%</span>
        <span class="graph-sub">{system.disk.used_gb} / {system.disk.total_gb} GB · {windowLabel()}</span>
      </div>
    </div>

    <!-- Temp graph -->
    {#if system.temp?.cpu_temp_c}
      {@const color = tempColor(system.temp.cpu_temp_c)}
      <div class="graph-card">
        <div class="graph-header">
          <span class="graph-label">Temp</span>
          <span class="graph-current" style="color:{color}">{system.temp.cpu_temp_c}°C</span>
        </div>
        <svg class="sparkline" viewBox="0 0 {W} {H}" preserveAspectRatio="none">
          <defs>
            <linearGradient id="g-temp" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="{color}" stop-opacity="0.3"/>
              <stop offset="100%" stop-color="{color}" stop-opacity="0"/>
            </linearGradient>
          </defs>
          {#if tempSpk.points}
            <polygon points="{tempSpk.points} {W},{H} 0,{H}" fill="url(#g-temp)"/>
            <polyline points={tempSpk.points} fill="none" stroke={color} stroke-width="1.8" stroke-linejoin="round" stroke-linecap="round"/>
            <line x1="0" x2={W} y1={tempSpk.avgY} y2={tempSpk.avgY} stroke={color} stroke-width="0.8" stroke-dasharray="4,3" opacity="0.45"/>
          {/if}
        </svg>
        <div class="graph-footer">
          <span class="avg-label">avg {tempSpk.avgVal?.toFixed(1) ?? '—'}°C</span>
          <span class="graph-sub">cpu temp · {windowLabel()}</span>
        </div>
      </div>
    {/if}

    <!-- Network graph -->
    {#if system.net}
      <div class="graph-card graph-card--wide">
        <div class="graph-header">
          <span class="graph-label">Network</span>
          <div class="net-legend">
            <span style="color:#34c759">↑ {fmtBytes(system.net.tx_bytes_s)}</span>
            <span style="color:#007AFF">↓ {fmtBytes(system.net.rx_bytes_s)}</span>
          </div>
        </div>
        <svg class="sparkline" viewBox="0 0 {W} {H}" preserveAspectRatio="none">
          <defs>
            <linearGradient id="g-tx" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#34c759" stop-opacity="0.22"/>
              <stop offset="100%" stop-color="#34c759" stop-opacity="0"/>
            </linearGradient>
            <linearGradient id="g-rx" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#007AFF" stop-opacity="0.22"/>
              <stop offset="100%" stop-color="#007AFF" stop-opacity="0"/>
            </linearGradient>
          </defs>
          {#if txSpk.points}
            <polygon points="{txSpk.points} {W},{H} 0,{H}" fill="url(#g-tx)"/>
            <polyline points={txSpk.points} fill="none" stroke="#34c759" stroke-width="1.8" stroke-linejoin="round" stroke-linecap="round"/>
          {/if}
          {#if rxSpk.points}
            <polygon points="{rxSpk.points} {W},{H} 0,{H}" fill="url(#g-rx)"/>
            <polyline points={rxSpk.points} fill="none" stroke="#007AFF" stroke-width="1.8" stroke-linejoin="round" stroke-linecap="round"/>
          {/if}
        </svg>
        <div class="graph-footer">
          <span class="avg-label" style="color:#34c759">↑ avg {fmtBytes(txSpk.avgVal ?? 0)}</span>
          <span class="graph-sub">{windowLabel()}</span>
          <span class="avg-label" style="color:#007AFF">↓ avg {fmtBytes(rxSpk.avgVal ?? 0)}</span>
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
    background: var(--glass-bg);
    backdrop-filter: blur(12px) saturate(1.8);
    -webkit-backdrop-filter: blur(12px) saturate(1.8);
    border: 1px solid var(--glass-border);
    border-radius: 10px;
    padding: 3px;
    gap: 2px;
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
    background: rgba(255, 255, 255, 0.18);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    color: var(--text);
    box-shadow: 0 1px 5px rgba(0, 0, 0, 0.14), inset 0 0.5px 0 rgba(255,255,255,0.25);
    border: 1px solid rgba(255, 255, 255, 0.18);
  }

  /* ── Shared glass card mixin ── */
  .gauge-card,
  .graph-card {
    background: var(--glass-bg);
    backdrop-filter: blur(20px) saturate(1.8);
    -webkit-backdrop-filter: blur(20px) saturate(1.8);
    border: 1px solid var(--glass-border);
    border-radius: 20px;
    box-shadow:
      0 4px 24px var(--glass-shadow),
      0 1px 3px var(--glass-shadow),
      inset 0 1px 0 var(--glass-shine);
    /* Prevent backdrop-filter bleed artifacts */
    transform: translateZ(0);
    will-change: transform;
    isolation: isolate;
    transition: background 0.2s ease, box-shadow 0.2s ease;
  }

  .gauge-card:hover,
  .graph-card:hover {
    background: var(--glass-bg-hover);
    box-shadow:
      0 6px 32px var(--glass-shadow),
      0 1px 4px var(--glass-shadow),
      inset 0 1px 0 var(--glass-shine);
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

  .ring { width: 120px; height: 120px; }

  .bg {
    fill: none;
    stroke: rgba(120, 120, 128, 0.2);
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

  /* ── Graphs ── */
  .graphs {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 0.875rem;
  }

  .graph-card--wide {
    grid-column: span 2;
  }

  .graph-card {
    padding: 1rem 1.1rem 0.85rem;
    display: flex;
    flex-direction: column;
    gap: 0.55rem;
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

  .sparkline {
    width: 100%;
    height: 72px;
    display: block;
    border-radius: 8px;
    overflow: visible;
  }

  .graph-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .avg-label {
    font-size: 0.68rem;
    font-weight: 600;
    color: var(--text);
    opacity: 0.55;
    letter-spacing: -0.01em;
  }

  .graph-sub {
    font-size: 0.65rem;
    color: var(--text);
    opacity: 0.4;
  }

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
