<script>
  let { system } = $props()
</script>

<section class="system-stats">
  <h2>System</h2>
  <div class="gauges">
    <div class="gauge">
      <svg viewBox="0 0 36 36">
        <path
          class="bg"
          d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <path
          class="fill"
          stroke-dasharray="{system.cpu.usage_percent}, 100"
          d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <text x="18" y="20.5" class="label">{system.cpu.usage_percent}%</text>
      </svg>
      <span class="name">CPU</span>
    </div>

    <div class="gauge">
      <svg viewBox="0 0 36 36">
        <path
          class="bg"
          d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <path
          class="fill"
          stroke-dasharray="{system.ram.usage_percent}, 100"
          d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <text x="18" y="20.5" class="label">{system.ram.usage_percent}%</text>
      </svg>
      <span class="name">RAM</span>
      <span class="detail">{system.ram.used_mb} / {system.ram.total_mb} MB</span>
    </div>

    <div class="gauge">
      <svg viewBox="0 0 36 36">
        <path
          class="bg"
          d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <path
          class="fill"
          stroke-dasharray="{system.disk.usage_percent}, 100"
          d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
        />
        <text x="18" y="20.5" class="label">{system.disk.usage_percent}%</text>
      </svg>
      <span class="name">Disk</span>
      <span class="detail">{system.disk.used_gb} / {system.disk.total_gb} GB</span>
    </div>

    {#if system.temp?.cpu_temp_c}
      <div class="gauge">
        <svg viewBox="0 0 36 36">
          <path
            class="bg"
            d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
          />
          <path
            class="fill temp"
            class:warm={system.temp.cpu_temp_c >= 60}
            class:hot={system.temp.cpu_temp_c >= 75}
            stroke-dasharray="{system.temp.cpu_temp_c}, 100"
            d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
          />
          <text x="18" y="20.5" class="label">{system.temp.cpu_temp_c}°</text>
        </svg>
        <span class="name">Temp</span>
        <span class="detail">{system.temp.cpu_temp_c}°C</span>
      </div>
    {/if}
  </div>
</section>

<style>
  .system-stats {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.25rem;
  }

  h2 {
    font-size: 1.1rem;
    margin: 0 0 1rem;
  }

  .gauges {
    display: flex;
    gap: 2rem;
    justify-content: center;
  }

  .gauge {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
  }

  svg {
    width: 100px;
    height: 100px;
  }

  .bg {
    fill: none;
    stroke: #30363d;
    stroke-width: 3;
  }

  .fill {
    fill: none;
    stroke: #58a6ff;
    stroke-width: 3;
    stroke-linecap: round;
    transition: stroke-dasharray 0.6s ease;
  }

  .fill.temp {
    stroke: #3fb950;
  }

  .fill.warm {
    stroke: #d29922;
  }

  .fill.hot {
    stroke: #f85149;
  }

  .label {
    fill: #e1e4e8;
    font-size: 7px;
    text-anchor: middle;
    font-weight: 600;
  }

  .name {
    font-size: 0.85rem;
    font-weight: 600;
  }

  .detail {
    font-size: 0.7rem;
    color: #8b949e;
  }
</style>
