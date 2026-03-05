<script>
  let { pihole } = $props()

  function formatNumber(n) {
    return n.toLocaleString()
  }
</script>

<div class="stats-grid">
  <div class="stat-card">
    <div class="stat-header">
      <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
      </svg>
      <span class="stat-label">Total Queries</span>
    </div>
    <span class="stat-value">{formatNumber(pihole.total_queries)}</span>
  </div>

  <div class="stat-card">
    <div class="stat-header">
      <svg class="stat-icon blocked" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
      </svg>
      <span class="stat-label">Blocked</span>
    </div>
    <span class="stat-value blocked">{formatNumber(pihole.blocked_queries)}</span>
  </div>

  <div class="stat-card">
    <div class="stat-header">
      <svg class="stat-icon blocked" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
      </svg>
      <span class="stat-label">Blocked %</span>
    </div>
    <span class="stat-value blocked">{pihole.blocked_percent.toFixed(1)}%</span>
  </div>

  <div class="stat-card">
    <div class="stat-header">
      <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 2L2 7l10 5 10-5-10-5z"/>
        <path d="M2 17l10 5 10-5"/>
        <path d="M2 12l10 5 10-5"/>
      </svg>
      <span class="stat-label">Domains on Blocklist</span>
    </div>
    <span class="stat-value">{formatNumber(pihole.domains_blocked)}</span>
  </div>
</div>

<style>
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 1rem;
  }

  .stat-card {
    background: #111720;
    border: 1px solid #1c2333;
    border-radius: 12px;
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    transition: border-color 0.2s;
  }

  .stat-card:hover {
    border-color: #2dd4bf40;
  }

  .stat-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .stat-icon {
    width: 16px;
    height: 16px;
    color: #2dd4bf;
  }

  .stat-icon.blocked {
    color: #f85149;
  }

  .stat-label {
    font-size: 0.75rem;
    color: #484f58;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .stat-value {
    font-size: 1.75rem;
    font-weight: 700;
    color: #2dd4bf;
  }

  .stat-value.blocked {
    color: #f85149;
  }
</style>
