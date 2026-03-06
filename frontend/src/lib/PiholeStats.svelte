<script>
  let { pihole } = $props()
  const fmt = (n) => n.toLocaleString()
</script>

<div class="card">
  <div class="card-header">
    <div class="card-title">
      <svg viewBox="0 0 24 24" fill="none" stroke="#007AFF" stroke-width="2" width="16" height="16">
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
      </svg>
      Pi-hole
    </div>
    <span class="badge">LIVE</span>
  </div>
  <div class="stats-grid">
    <div class="stat">
      <span class="stat-label">QUERIES</span>
      <span class="stat-value">{fmt(pihole.total_queries)}</span>
    </div>
    <div class="stat">
      <span class="stat-label">BLOCKED</span>
      <span class="stat-value accent">{fmt(pihole.blocked_queries)}</span>
    </div>
    <div class="stat">
      <span class="stat-label">BLOCK RATE</span>
      <span class="stat-value accent">{pihole.blocked_percent.toFixed(1)}%</span>
    </div>
    <div class="stat">
      <span class="stat-label">BLOCKLIST</span>
      <span class="stat-value">{pihole.domains_blocked >= 1000 ? (pihole.domains_blocked / 1000).toFixed(0) + 'k' : fmt(pihole.domains_blocked)}</span>
    </div>
  </div>
</div>

<style>
  .card {
    background: #ffffff;
    border-radius: 16px;
    padding: 1.25rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.04);
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
    color: #1c1c1e;
  }

  .badge {
    font-size: 0.62rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    color: #aeaeb2;
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
    border-bottom: 1px solid rgba(60,60,67,0.1);
  }

  .stat:nth-last-child(-n+2) { border-bottom: none; padding-bottom: 0; }

  .stat-label {
    font-size: 0.62rem;
    font-weight: 700;
    letter-spacing: 0.07em;
    color: #aeaeb2;
    text-transform: uppercase;
  }

  .stat-value {
    font-size: 1.5rem;
    font-weight: 700;
    color: #1c1c1e;
    letter-spacing: -0.02em;
  }

  .stat-value.accent { color: #007AFF; }
</style>
