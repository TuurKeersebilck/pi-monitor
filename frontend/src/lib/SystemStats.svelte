<script>
  let { system } = $props()

  const circlePath = "M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"

  function tempColor(temp) {
    if (temp >= 75) return '#f85149'
    if (temp >= 60) return '#d29922'
    return '#2dd4bf'
  }
</script>

<div class="gauges">
  <div class="gauge-card">
    <div class="gauge-header">
      <svg class="gauge-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="4" y="4" width="16" height="16" rx="2"/>
        <rect x="9" y="9" width="6" height="6"/>
        <line x1="9" y1="1" x2="9" y2="4"/><line x1="15" y1="1" x2="15" y2="4"/>
        <line x1="9" y1="20" x2="9" y2="23"/><line x1="15" y1="20" x2="15" y2="23"/>
        <line x1="20" y1="9" x2="23" y2="9"/><line x1="20" y1="14" x2="23" y2="14"/>
        <line x1="1" y1="9" x2="4" y2="9"/><line x1="1" y1="14" x2="4" y2="14"/>
      </svg>
      <span class="gauge-label">CPU</span>
    </div>
    <div class="gauge-ring">
      <svg viewBox="0 0 36 36">
        <path class="ring-bg" d={circlePath} />
        <path class="ring-fill" stroke-dasharray="{system.cpu.usage_percent}, 100" d={circlePath} />
        <text x="18" y="19" class="ring-value">{system.cpu.usage_percent}</text>
        <text x="18" y="23.5" class="ring-unit">%</text>
      </svg>
    </div>
  </div>

  <div class="gauge-card">
    <div class="gauge-header">
      <svg class="gauge-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="3" y="3" width="7" height="18" rx="1"/>
        <rect x="14" y="3" width="7" height="18" rx="1"/>
      </svg>
      <span class="gauge-label">RAM</span>
    </div>
    <div class="gauge-ring">
      <svg viewBox="0 0 36 36">
        <path class="ring-bg" d={circlePath} />
        <path class="ring-fill" stroke-dasharray="{system.ram.usage_percent}, 100" d={circlePath} />
        <text x="18" y="19" class="ring-value">{system.ram.usage_percent}</text>
        <text x="18" y="23.5" class="ring-unit">%</text>
      </svg>
    </div>
    <span class="gauge-detail">{system.ram.used_mb} / {system.ram.total_mb} MB</span>
  </div>

  <div class="gauge-card">
    <div class="gauge-header">
      <svg class="gauge-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <ellipse cx="12" cy="5" rx="9" ry="3"/>
        <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
        <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
      </svg>
      <span class="gauge-label">Disk</span>
    </div>
    <div class="gauge-ring">
      <svg viewBox="0 0 36 36">
        <path class="ring-bg" d={circlePath} />
        <path class="ring-fill" stroke-dasharray="{system.disk.usage_percent}, 100" d={circlePath} />
        <text x="18" y="19" class="ring-value">{system.disk.usage_percent}</text>
        <text x="18" y="23.5" class="ring-unit">%</text>
      </svg>
    </div>
    <span class="gauge-detail">{system.disk.used_gb} / {system.disk.total_gb} GB</span>
  </div>

  {#if system.temp?.cpu_temp_c}
    <div class="gauge-card">
      <div class="gauge-header">
        <svg class="gauge-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 14.76V3.5a2.5 2.5 0 0 0-5 0v11.26a4.5 4.5 0 1 0 5 0z"/>
        </svg>
        <span class="gauge-label">Temp</span>
      </div>
      <div class="gauge-ring">
        <svg viewBox="0 0 36 36">
          <path class="ring-bg" d={circlePath} />
          <path
            class="ring-fill"
            style="stroke: {tempColor(system.temp.cpu_temp_c)}"
            stroke-dasharray="{system.temp.cpu_temp_c}, 100"
            d={circlePath}
          />
          <text x="18" y="19" class="ring-value">{system.temp.cpu_temp_c}</text>
          <text x="18" y="23.5" class="ring-unit">°C</text>
        </svg>
      </div>
    </div>
  {/if}
</div>

<style>
  .gauges {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
    gap: 1rem;
  }

  .gauge-card {
    background: #111720;
    border: 1px solid #1c2333;
    border-radius: 12px;
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    transition: border-color 0.2s;
  }

  .gauge-card:hover {
    border-color: #2dd4bf40;
  }

  .gauge-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    align-self: flex-start;
  }

  .gauge-icon {
    width: 16px;
    height: 16px;
    color: #2dd4bf;
  }

  .gauge-label {
    font-size: 0.8rem;
    font-weight: 600;
    color: #8b949e;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .gauge-ring {
    width: 110px;
    height: 110px;
  }

  .gauge-ring svg {
    width: 100%;
    height: 100%;
  }

  .ring-bg {
    fill: none;
    stroke: #1c2333;
    stroke-width: 2.5;
  }

  .ring-fill {
    fill: none;
    stroke: #2dd4bf;
    stroke-width: 2.5;
    stroke-linecap: round;
    transition: stroke-dasharray 0.6s ease;
  }

  .ring-value {
    fill: #e6edf3;
    font-size: 8px;
    text-anchor: middle;
    font-weight: 700;
  }

  .ring-unit {
    fill: #484f58;
    font-size: 4.5px;
    text-anchor: middle;
  }

  .gauge-detail {
    font-size: 0.7rem;
    color: #484f58;
  }
</style>
