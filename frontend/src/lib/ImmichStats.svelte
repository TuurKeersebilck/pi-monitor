<script>
  let { immich } = $props()
  const fmt = (n) => n.toLocaleString()
</script>

<div class="card">
  <div class="card-header">
    <div class="card-title">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="18" height="18" style="color:#2dd4bf">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
        <circle cx="8.5" cy="8.5" r="1.5"/>
        <polyline points="21 15 16 10 5 21"/>
      </svg>
      Immich Stats
    </div>
    <span class="badge">LIBRARY STATS</span>
  </div>
  <div class="stats-grid">
    <div class="stat">
      <span class="stat-label">TOTAL PHOTOS</span>
      <span class="stat-value">{fmt(immich.photos)}</span>
    </div>
    <div class="stat">
      <span class="stat-label">TOTAL VIDEOS</span>
      <span class="stat-value">{fmt(immich.videos)}</span>
    </div>
    <div class="stat">
      <span class="stat-label">STORAGE USED</span>
      <span class="stat-value teal">{immich.usage_gb} GB</span>
    </div>
    <div class="stat">
      <span class="stat-label">LIBRARY HEALTH</span>
      <span class="stat-value" class:teal={immich.healthy} class:red={!immich.healthy}>
        {immich.healthy ? 'Healthy' : 'Unhealthy'}
        {#if immich.healthy}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" width="16" height="16">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
        {:else}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" width="16" height="16">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
        {/if}
      </span>
    </div>
  </div>
</div>

<style>
  .card {
    background: #161b22;
    border: 1px solid #21262d;
    border-radius: 14px;
    padding: 1.25rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .card-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1rem;
    font-weight: 700;
    color: #f0f6fc;
  }

  .badge {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    color: #8b949e;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem 1rem;
  }

  .stat {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid #21262d;
  }

  .stat:nth-last-child(-n+2) { border-bottom: none; padding-bottom: 0; }

  .stat-label {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.07em;
    color: #484f58;
    text-transform: uppercase;
  }

  .stat-value {
    font-size: 1.6rem;
    font-weight: 800;
    color: #f0f6fc;
    letter-spacing: -0.02em;
    display: flex;
    align-items: center;
    gap: 0.3rem;
  }

  .stat-value.teal { color: #2dd4bf; }
  .stat-value.red { color: #f85149; }
</style>
