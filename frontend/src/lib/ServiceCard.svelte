<script>
  let { container } = $props()
</script>

<div class="card">
  <div class="card-top">
    <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
      <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
      <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
      <line x1="12" y1="22.08" x2="12" y2="12"/>
    </svg>
    <span class="health" class:healthy={container.running}>
      <span class="health-dot"></span>
      {container.running ? 'HEALTHY' : 'STOPPED'}
    </span>
  </div>

  <div class="card-body">
    <span class="name">{container.name}</span>
    <span class="image">{container.image}</span>
  </div>

  <div class="card-footer">
    <span class="uptime">UPTIME: {container.uptime.toUpperCase()}</span>
    <svg class="ext-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
      <polyline points="15 3 21 3 21 9"/>
      <line x1="10" y1="14" x2="21" y2="3"/>
    </svg>
  </div>

  {#if container.update_available}
    <div class="update-banner">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="13" height="13">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      Update available
    </div>
  {/if}
</div>

<style>
  .card {
    background: #161b22;
    border: 1px solid #21262d;
    border-radius: 14px;
    padding: 1.1rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    transition: border-color 0.2s;
  }

  .card:hover { border-color: #30363d; }

  .card-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .icon {
    width: 28px;
    height: 28px;
    color: #2dd4bf;
  }

  .health {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    padding: 0.2rem 0.55rem;
    border-radius: 999px;
    background: rgba(248, 81, 73, 0.1);
    color: #f85149;
    border: 1px solid rgba(248, 81, 73, 0.2);
  }

  .health.healthy {
    background: rgba(63, 185, 80, 0.1);
    color: #3fb950;
    border-color: rgba(63, 185, 80, 0.25);
  }

  .health-dot {
    width: 5px; height: 5px;
    border-radius: 50%;
    background: currentColor;
  }

  .card-body {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }

  .name {
    font-size: 0.95rem;
    font-weight: 700;
    color: #f0f6fc;
  }

  .image {
    font-size: 0.72rem;
    color: #484f58;
    word-break: break-all;
  }

  .card-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-top: 0.5rem;
    border-top: 1px solid #21262d;
  }

  .uptime {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    color: #8b949e;
  }

  .ext-icon {
    width: 14px;
    height: 14px;
    color: #484f58;
    cursor: pointer;
    transition: color 0.15s;
  }

  .ext-icon:hover { color: #8b949e; }

  .update-banner {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.4rem 0.6rem;
    border-radius: 8px;
    background: rgba(210, 153, 34, 0.1);
    border: 1px solid rgba(210, 153, 34, 0.2);
    font-size: 0.7rem;
    font-weight: 600;
    color: #d29922;
  }
</style>
