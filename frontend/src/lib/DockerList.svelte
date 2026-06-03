<script>
  import { onDestroy } from 'svelte'
  import { Chart, LineController, LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip } from 'chart.js'
  Chart.register(LineController, LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip)

  let { containers = [], containerHistory = {} } = $props()

  let expanded = $state(null)
  let sparkCanvas = $state(null)
  let sparkChart = null

  function toggle(name) {
    expanded = expanded === name ? null : name
  }

  function fmtBytes(bps) {
    if (!bps || bps <= 0) return '0 B/s'
    if (bps >= 1_048_576) return (bps / 1_048_576).toFixed(1) + ' MB/s'
    if (bps >= 1024) return (bps / 1024).toFixed(0) + ' KB/s'
    return Math.round(bps) + ' B/s'
  }

  function fmtMem(mb) {
    if (!mb) return '0 MB'
    if (mb >= 1024) return (mb / 1024).toFixed(1) + ' GB'
    return Math.round(mb) + ' MB'
  }

  function shortImage(image) {
    const parts = image.split('/')
    return parts[parts.length - 1]
  }

  function clamp(v) {
    return Math.min(Math.max(v || 0, 0), 100)
  }

  function barColor(pct) {
    if (pct >= 90) return 'var(--red)'
    if (pct >= 75) return 'var(--orange)'
    return null // use default css colour
  }

  let grouped = $derived.by(() => {
    const map = {}
    for (const c of containers) {
      const g = c.group || 'Containers'
      if (!map[g]) map[g] = []
      map[g].push(c)
    }
    return Object.entries(map).sort(([a], [b]) => a.localeCompare(b))
  })

  // Init sparkline chart when a row is expanded
  $effect(() => {
    if (!sparkCanvas || !expanded) return

    sparkChart?.destroy()
    sparkChart = null

    const history = containerHistory[expanded] ?? []
    const cpuData = history.map(d => d.cpu_percent ?? 0)
    const isDark = document.documentElement.classList.contains('dark')
    const ctx = sparkCanvas.getContext('2d')
    const grad = ctx.createLinearGradient(0, 0, 0, 72)
    grad.addColorStop(0, 'rgba(0,122,255,0.22)')
    grad.addColorStop(1, 'rgba(0,122,255,0)')

    sparkChart = new Chart(sparkCanvas, {
      type: 'line',
      data: {
        labels: cpuData.map((_, i) => i),
        datasets: [{
          label: 'CPU %',
          data: cpuData,
          borderColor: '#007AFF',
          backgroundColor: grad,
          fill: true,
          tension: 0.4,
          borderWidth: 1.5,
          pointRadius: 0,
          pointHoverRadius: 3,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: false,
        plugins: {
          legend: { display: false },
          tooltip: {
            callbacks: { label: c => `CPU: ${c.parsed.y.toFixed(1)}%` },
            backgroundColor: isDark ? 'rgba(30,30,34,0.92)' : 'rgba(255,255,255,0.92)',
            titleColor: isDark ? '#f0f0f5' : '#1c1c1e',
            bodyColor: isDark ? '#a0a0a8' : '#6c6c70',
            borderColor: isDark ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.06)',
            borderWidth: 1,
            padding: 6,
            cornerRadius: 6,
          }
        },
        scales: {
          x: { display: false },
          y: {
            min: 0,
            max: 100,
            border: { display: false },
            grid: { color: isDark ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.06)' },
            ticks: {
              callback: v => v + '%',
              maxTicksLimit: 3,
              font: { size: 10 },
              color: isDark ? 'rgba(240,240,245,0.4)' : 'rgba(28,28,30,0.4)',
              padding: 4,
            }
          }
        }
      }
    })

    return () => {
      sparkChart?.destroy()
      sparkChart = null
    }
  })

  // Live-update sparkline when history changes
  $effect(() => {
    if (!sparkChart || !expanded) return
    const cpuData = (containerHistory[expanded] ?? []).map(d => d.cpu_percent ?? 0)
    sparkChart.data.labels = cpuData.map((_, i) => i)
    sparkChart.data.datasets[0].data = cpuData
    sparkChart.update('none')
  })

  onDestroy(() => sparkChart?.destroy())
</script>

<div class="docker-list">
  {#each grouped as [group, items]}
    {#if grouped.length > 1}
      <div class="group-header">
        <span class="group-name">{group}</span>
        <span class="group-count">{items.length}</span>
      </div>
    {/if}

    {#each items as c (c.name)}
      <div
        class="container-row card-surface"
        class:stopped={!c.running}
        class:is-expanded={expanded === c.name}
        onclick={() => c.running && toggle(c.name)}
        role={c.running ? 'button' : undefined}
        tabindex={c.running ? 0 : undefined}
        onkeydown={(e) => c.running && e.key === 'Enter' && toggle(c.name)}
      >
        <!-- Main row -->
        <div class="row-main">
          <span class="status-dot" class:running={c.running}></span>

          <span class="container-name">{c.name}</span>
          <span class="image-tag">{shortImage(c.image)}</span>

          {#if c.running}
            <div class="bars">
              <div class="bar-item">
                <span class="bar-label">CPU</span>
                <div class="bar-track">
                  <div
                    class="bar-fill"
                    style="width:{clamp(c.cpu_percent)}%;background:{barColor(c.cpu_percent) ?? '#007AFF'}"
                  ></div>
                </div>
                <span class="bar-val">{(c.cpu_percent || 0).toFixed(1)}%</span>
              </div>
              <div class="bar-item">
                <span class="bar-label">RAM</span>
                <div class="bar-track">
                  <div
                    class="bar-fill"
                    style="width:{clamp(c.mem_percent)}%;background:{barColor(c.mem_percent) ?? '#34c759'}"
                  ></div>
                </div>
                <span class="bar-val">{(c.mem_percent || 0).toFixed(1)}%</span>
              </div>
            </div>

            <div class="net-info">
              <span class="net-tx">↑ {fmtBytes(c.net_tx_bytes_s)}</span>
              <span class="net-rx">↓ {fmtBytes(c.net_rx_bytes_s)}</span>
            </div>
          {:else}
            <span class="stopped-text">stopped</span>
          {/if}

          <span class="uptime-text">{c.uptime}</span>

          {#if c.running}
            <span class="chevron" class:open={expanded === c.name}>›</span>
          {/if}
        </div>

        <!-- Expanded sparkline -->
        {#if expanded === c.name}
          <div class="expanded-section">
            <div class="expand-meta">
              <span class="expand-label">CPU history · {(containerHistory[c.name]?.length ?? 0)}pts</span>
              <span class="expand-mem">{fmtMem(c.mem_used_mb)} / {fmtMem(c.mem_limit_mb)} RAM</span>
            </div>
            <div class="spark-wrap">
              <canvas bind:this={sparkCanvas}></canvas>
            </div>
          </div>
        {/if}
      </div>
    {/each}
  {/each}
</div>

<style>
  .docker-list {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .group-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-top: 1rem;
    margin-bottom: 0.3rem;
  }

  .group-header:first-child { margin-top: 0; }

  .group-name {
    font-size: 0.68rem;
    font-weight: 700;
    letter-spacing: 0.07em;
    text-transform: uppercase;
    color: var(--text-3);
  }

  .group-count {
    font-size: 0.62rem;
    font-weight: 700;
    padding: 0.1rem 0.4rem;
    border-radius: 999px;
    background: var(--pill);
    color: var(--text-2);
  }

  .container-row {
    border-radius: 12px;
    padding: 0.65rem 1rem;
    cursor: default;
    transition: background 0.15s ease;
  }

  .container-row[role="button"] { cursor: pointer; }

  .container-row.stopped {
    opacity: 0.55;
  }

  /* Main row layout */
  .row-main {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    min-height: 28px;
  }

  .status-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: var(--red);
    flex-shrink: 0;
  }

  .status-dot.running { background: var(--green); }

  .container-name {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text);
    white-space: nowrap;
    min-width: 100px;
    max-width: 160px;
    overflow: hidden;
    text-overflow: ellipsis;
    flex-shrink: 0;
  }

  .image-tag {
    font-size: 0.72rem;
    color: var(--text-3);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 140px;
    flex-shrink: 1;
  }

  /* Resource bars */
  .bars {
    display: flex;
    gap: 0.75rem;
    flex: 1;
    min-width: 0;
  }

  .bar-item {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    min-width: 120px;
  }

  .bar-label {
    font-size: 0.68rem;
    font-weight: 700;
    color: var(--text-3);
    text-transform: uppercase;
    letter-spacing: 0.03em;
    width: 26px;
    flex-shrink: 0;
  }

  .bar-track {
    flex: 1;
    height: 5px;
    background: var(--ring-track);
    border-radius: 999px;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    border-radius: 999px;
    transition: width 0.6s cubic-bezier(0.34, 1.2, 0.64, 1);
  }

  .bar-val {
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--text-2);
    width: 36px;
    text-align: right;
    flex-shrink: 0;
    font-variant-numeric: tabular-nums;
  }

  /* Network */
  .net-info {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
    flex-shrink: 0;
  }

  .net-tx, .net-rx {
    font-size: 0.7rem;
    font-weight: 600;
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
  }

  .net-tx { color: var(--green); }
  .net-rx { color: var(--accent); }

  .stopped-text {
    font-size: 0.72rem;
    color: var(--text-3);
    flex: 1;
  }

  .uptime-text {
    font-size: 0.7rem;
    color: var(--text-3);
    white-space: nowrap;
    flex-shrink: 0;
    margin-left: auto;
  }

  .chevron {
    font-size: 1.1rem;
    color: var(--text-3);
    flex-shrink: 0;
    transition: transform 0.2s ease;
    line-height: 1;
  }

  .chevron.open { transform: rotate(90deg); }

  /* Expanded sparkline */
  .expanded-section {
    margin-top: 0.75rem;
    padding-top: 0.75rem;
    border-top: 1px solid var(--border-subtle);
  }

  .expand-meta {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    margin-bottom: 0.5rem;
  }

  .expand-label {
    font-size: 0.68rem;
    font-weight: 700;
    color: var(--text-3);
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .expand-mem {
    font-size: 0.7rem;
    color: var(--text-3);
    font-variant-numeric: tabular-nums;
  }

  .spark-wrap {
    height: 72px;
    position: relative;
  }

  .spark-wrap canvas {
    width: 100% !important;
    height: 100% !important;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .image-tag { display: none; }
    .bar-item { min-width: 90px; }
    .net-info { display: none; }
  }

  @media (max-width: 480px) {
    .bars { gap: 0.4rem; }
    .bar-item { min-width: 70px; }
  }
</style>
