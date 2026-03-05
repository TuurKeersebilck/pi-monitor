<script>
  let { pihole } = $props()

  function formatNumber(n) {
    return n.toLocaleString()
  }
</script>

<section class="pihole">
  <h2>Pi-hole</h2>
  <div class="stats">
    <div class="stat">
      <span class="value">{formatNumber(pihole.total_queries)}</span>
      <span class="label">Total Queries</span>
    </div>
    <div class="stat">
      <span class="value">{formatNumber(pihole.blocked_queries)}</span>
      <span class="label">Blocked</span>
    </div>
    <div class="stat">
      <span class="value blocked">{pihole.blocked_percent.toFixed(1)}%</span>
      <span class="label">Blocked %</span>
    </div>
    <div class="stat">
      <span class="value">{formatNumber(pihole.domains_blocked)}</span>
      <span class="label">Domains on Blocklist</span>
    </div>
  </div>
</section>

<style>
  .pihole {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1.25rem;
    margin-top: 1rem;
  }

  h2 {
    font-size: 1.1rem;
    margin: 0 0 1rem;
  }

  .stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 1rem;
  }

  .stat {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
  }

  .value {
    font-size: 1.5rem;
    font-weight: 700;
    color: #58a6ff;
  }

  .value.blocked {
    color: #f85149;
  }

  .label {
    font-size: 0.75rem;
    color: #8b949e;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
</style>
