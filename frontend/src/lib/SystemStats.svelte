<script>
  let { system } = $props()

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
</script>

<div class="gauges">
  <!-- CPU -->
  <div class="gauge-card">
    <svg class="ring" viewBox="0 0 36 36">
      <path class="bg" d={PATH}/>
      <path class="fill" style="stroke:#007AFF" stroke-dasharray="{system.cpu.usage_percent},100" d={PATH}/>
      <text x="18" y="16.5" class="val">{system.cpu.usage_percent}%</text>
    </svg>
    <span class="label">CPU</span>
  </div>

  <!-- RAM -->
  <div class="gauge-card">
    <svg class="ring" viewBox="0 0 36 36">
      <path class="bg" d={PATH}/>
      <path class="fill" style="stroke:#34c759" stroke-dasharray="{system.ram.usage_percent},100" d={PATH}/>
      <text x="18" y="16.5" class="val">{system.ram.usage_percent}%</text>
    </svg>
    <span class="label">RAM</span>
    <span class="sub">{system.ram.used_mb} / {system.ram.total_mb} MB</span>
  </div>

  <!-- Disk -->
  <div class="gauge-card">
    <svg class="ring" viewBox="0 0 36 36">
      <path class="bg" d={PATH}/>
      <path class="fill" style="stroke:#ff9500" stroke-dasharray="{system.disk.usage_percent},100" d={PATH}/>
      <text x="18" y="16.5" class="val">{system.disk.usage_percent}%</text>
    </svg>
    <span class="label">Disk</span>
    <span class="sub">{system.disk.used_gb} / {system.disk.total_gb} GB</span>
  </div>

  <!-- Temp -->
  {#if system.temp?.cpu_temp_c}
    <div class="gauge-card">
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

<style>
  .gauges {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(145px, 1fr));
    gap: 1rem;
  }

  .gauge-card {
    background: #ffffff;
    border-radius: 16px;
    padding: 1.5rem 1rem 1.25rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.04);
  }

  .ring { width: 120px; height: 120px; }

  .bg {
    fill: none;
    stroke: #f2f2f7;
    stroke-width: 2.8;
  }

  .fill {
    fill: none;
    stroke-width: 2.8;
    stroke-linecap: round;
    transition: stroke-dasharray 0.5s ease;
  }

  .val {
    fill: #1c1c1e;
    font-size: 7.5px;
    font-weight: 700;
    text-anchor: middle;
    dominant-baseline: middle;
  }

  .label {
    font-size: 0.8rem;
    color: #6c6c70;
    font-weight: 500;
    text-align: center;
  }

  .sub {
    font-size: 0.68rem;
    color: #aeaeb2;
    text-align: center;
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
    color: #1c1c1e;
    letter-spacing: -0.01em;
    white-space: nowrap;
  }

  @media (max-width: 768px) {
    .gauges { grid-template-columns: repeat(2, 1fr); }
  }
</style>
