<script>
  let { system } = $props()

  const PATH = "M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"

  function tempColor(t) {
    if (t >= 75) return '#f85149'
    if (t >= 60) return '#d29922'
    return '#2dd4bf'
  }

  function fmtBytes(bps) {
    if (bps >= 1_048_576) return (bps / 1_048_576).toFixed(1) + ' MB/s'
    if (bps >= 1024) return (bps / 1024).toFixed(0) + ' KB/s'
    return bps.toFixed(0) + ' B/s'
  }
</script>

<div class="gauges">
  <!-- CPU -->
  <div class="gauge-card" style="border-top: 2px solid #2dd4bf">
    <svg class="ring" viewBox="0 0 36 36">
      <path class="bg" d={PATH}/>
      <path class="fill" style="stroke:#2dd4bf" stroke-dasharray="{system.cpu.usage_percent},100" d={PATH}/>
      <text x="18" y="16.5" class="val">{system.cpu.usage_percent}%</text>
    </svg>
    <span class="label">CPU Usage</span>
  </div>

  <!-- RAM -->
  <div class="gauge-card" style="border-top: 2px solid #3fb950">
    <svg class="ring" viewBox="0 0 36 36">
      <path class="bg" d={PATH}/>
      <path class="fill" style="stroke:#3fb950" stroke-dasharray="{system.ram.usage_percent},100" d={PATH}/>
      <text x="18" y="16.5" class="val">{system.ram.usage_percent}%</text>
    </svg>
    <span class="label">RAM Usage</span>
    <span class="sub">{system.ram.used_mb} / {system.ram.total_mb} MB</span>
  </div>

  <!-- Disk -->
  <div class="gauge-card" style="border-top: 2px solid #d29922">
    <svg class="ring" viewBox="0 0 36 36">
      <path class="bg" d={PATH}/>
      <path class="fill" style="stroke:#d29922" stroke-dasharray="{system.disk.usage_percent},100" d={PATH}/>
      <text x="18" y="16.5" class="val">{system.disk.usage_percent}%</text>
    </svg>
    <span class="label">Disk Space</span>
    <span class="sub">{system.disk.used_gb} / {system.disk.total_gb} GB</span>
  </div>

  <!-- Temp -->
  {#if system.temp?.cpu_temp_c}
    <div class="gauge-card" style="border-top: 2px solid {tempColor(system.temp.cpu_temp_c)}">
      <svg class="ring" viewBox="0 0 36 36">
        <path class="bg" d={PATH}/>
        <path
          class="fill"
          style="stroke:{tempColor(system.temp.cpu_temp_c)}"
          stroke-dasharray="{system.temp.cpu_temp_c},100"
          d={PATH}
        />
        <text x="18" y="16.5" class="val">{system.temp.cpu_temp_c}°</text>
      </svg>
      <span class="label">Temperature</span>
      <span class="sub">{system.temp.cpu_temp_c}°C</span>
    </div>
  {/if}

  <!-- Network -->
  {#if system.net}
    <div class="gauge-card net-card" style="border-top: 2px solid #6366f1">
      <div class="net-rows">
        <div class="net-row">
          <svg viewBox="0 0 24 24" fill="none" stroke="#3fb950" stroke-width="2.5" width="16" height="16">
            <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
            <polyline points="17 6 23 6 23 12"/>
          </svg>
          <span class="net-val">{fmtBytes(system.net.tx_bytes_s)}</span>
        </div>
        <div class="net-row">
          <svg viewBox="0 0 24 24" fill="none" stroke="#2dd4bf" stroke-width="2.5" width="16" height="16">
            <polyline points="23 18 13.5 8.5 8.5 13.5 1 6"/>
            <polyline points="17 18 23 18 23 12"/>
          </svg>
          <span class="net-val">{fmtBytes(system.net.rx_bytes_s)}</span>
        </div>
      </div>
      <span class="label">Network I/O</span>
      <span class="sub">↑ upload &nbsp;·&nbsp; ↓ download</span>
    </div>
  {/if}
</div>

<style>
  .gauges {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(145px, 1fr));
    gap: 1rem;
  }

  .gauge-card {
    background: #161b22;
    border: 1px solid #21262d;
    border-radius: 14px;
    padding: 1.5rem 1rem 1.25rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.6rem;
    transition: border-color 0.2s;
  }

  .gauge-card:hover { border-color: #30363d; }

  .ring { width: 130px; height: 130px; }

  .bg {
    fill: none;
    stroke: #21262d;
    stroke-width: 2.8;
  }

  .fill {
    fill: none;
    stroke-width: 2.8;
    stroke-linecap: round;
    transition: stroke-dasharray 0.5s ease;
  }

  .val {
    fill: #f0f6fc;
    font-size: 7.5px;
    font-weight: 800;
    text-anchor: middle;
    dominant-baseline: middle;
  }

  .label {
    font-size: 0.82rem;
    color: #8b949e;
    font-weight: 500;
    text-align: center;
  }

  .sub {
    font-size: 0.7rem;
    color: #484f58;
    text-align: center;
  }

  .net-card { justify-content: center; gap: 1rem; }

  .net-rows {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    width: 100%;
    padding: 0 0.5rem;
  }

  .net-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .net-val {
    font-size: 1rem;
    font-weight: 700;
    color: #f0f6fc;
    letter-spacing: -0.01em;
    white-space: nowrap;
  }

  @media (max-width: 768px) {
    .gauges { grid-template-columns: repeat(2, 1fr); }
  }
</style>
