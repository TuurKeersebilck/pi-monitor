<script>
  let { immich } = $props()
  const fmt = (n) => n.toLocaleString()
</script>

<div class="card">
  <div class="card-header">
    <div class="card-title">
      <svg viewBox="0 0 24 24" fill="none" stroke="#007AFF" stroke-width="2" width="16" height="16">
        <rect x="3" y="3" width="18" height="18" rx="2"/>
        <circle cx="8.5" cy="8.5" r="1.5"/>
        <polyline points="21 15 16 10 5 21"/>
      </svg>
      Immich
    </div>
    <span class="badge">LIBRARY</span>
  </div>
  <div class="stats-grid">
    <div class="stat">
      <span class="stat-label">PHOTOS</span>
      <span class="stat-value">{fmt(immich.photos)}</span>
    </div>
    <div class="stat">
      <span class="stat-label">VIDEOS</span>
      <span class="stat-value">{fmt(immich.videos)}</span>
    </div>
    <div class="stat">
      <span class="stat-label">STORAGE</span>
      <span class="stat-value accent">{immich.usage_gb} GB</span>
    </div>
    <div class="stat">
      <span class="stat-label">HEALTH</span>
      <span class="stat-value" class:accent={immich.healthy} class:red={!immich.healthy}>
        {immich.healthy ? 'Healthy' : 'Unhealthy'}
      </span>
    </div>
  </div>
</div>

<style>
  .card {
    background: var(--card-bg);
    border-radius: 16px;
    padding: 1.25rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.04);
    transition: background 0.2s ease;
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
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--text);
  }

  .badge {
    font-size: 0.62rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    color: var(--text-3);
  }

  .stats-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.1rem 1rem;
  }

  .stat {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--border-subtle);
  }

  .stat:nth-last-child(-n+2) { border-bottom: none; padding-bottom: 0; }

  .stat-label {
    font-size: 0.62rem;
    font-weight: 700;
    letter-spacing: 0.07em;
    color: var(--text-3);
    text-transform: uppercase;
  }

  .stat-value {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--text);
    letter-spacing: -0.02em;
  }

  .stat-value.accent { color: #007AFF; }
  .stat-value.red { color: #ff3b30; }
</style>
