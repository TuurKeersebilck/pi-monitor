<script>
  /**
   * @typedef {{ label: string, value: string, color?: 'accent' | 'red' }} Stat
   * @type {{ title: string, badge: string, stats: Stat[], icon: import('svelte').Snippet }}
   */
  let { title, badge, stats, icon } = $props()
</script>

<div class="card card-surface">
  <div class="card-header">
    <div class="card-title">
      {@render icon()}
      {title}
    </div>
    <span class="badge">{badge}</span>
  </div>
  <div class="stats-grid">
    {#each stats as s (s.label)}
      <div class="stat">
        <span class="stat-label">{s.label}</span>
        <span class="stat-value" class:accent={s.color === 'accent'} class:red={s.color === 'red'}>{s.value}</span>
      </div>
    {/each}
  </div>
</div>

<style>
  .card {
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

  .stat-value.accent { color: var(--accent); }
  .stat-value.red { color: #ff3b30; }
</style>
