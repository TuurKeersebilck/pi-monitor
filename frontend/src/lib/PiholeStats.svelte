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
    background: var(--glass-bg);
    backdrop-filter: blur(20px) saturate(1.8);
    -webkit-backdrop-filter: blur(20px) saturate(1.8);
    border: 1px solid var(--glass-border);
    border-radius: 20px;
    padding: 1.25rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    box-shadow:
      0 4px 24px var(--glass-shadow),
      0 1px 3px var(--glass-shadow),
      inset 0 1px 0 var(--glass-shine);
    transform: translateZ(0);
    will-change: transform;
    isolation: isolate;
    transition: background 0.2s ease, box-shadow 0.2s ease;
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
</style>
