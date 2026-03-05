<script>
  let { container } = $props()
</script>

<div class="card">
  <div class="card-header">
    <div class="name-row">
      <svg class="container-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
        <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
        <line x1="12" y1="22.08" x2="12" y2="12"/>
      </svg>
      <span class="name">{container.name}</span>
    </div>
    <span class="health-badge" class:healthy={container.running}>
      {container.running ? 'HEALTHY' : 'STOPPED'}
    </span>
  </div>

  <div class="card-body">
    <div class="info-row">
      <span class="info-label">Image</span>
      <span class="info-value">{container.image}</span>
    </div>
    <div class="info-row">
      <span class="info-label">Status</span>
      <span class="info-value">{container.status}</span>
    </div>
  </div>

  {#if container.update_available}
    <div class="update-banner">
      <svg class="update-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
    background: #111720;
    border: 1px solid #1c2333;
    border-radius: 12px;
    padding: 1.25rem;
    transition: border-color 0.2s;
  }

  .card:hover {
    border-color: #2dd4bf40;
  }

  .card-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .name-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 0;
  }

  .container-icon {
    width: 18px;
    height: 18px;
    color: #2dd4bf;
    flex-shrink: 0;
  }

  .name {
    font-weight: 600;
    font-size: 0.95rem;
    color: #e6edf3;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .health-badge {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    padding: 0.2rem 0.6rem;
    border-radius: 999px;
    white-space: nowrap;
    flex-shrink: 0;
    background: rgba(248, 81, 73, 0.1);
    color: #f85149;
    border: 1px solid rgba(248, 81, 73, 0.2);
  }

  .health-badge.healthy {
    background: rgba(45, 212, 191, 0.1);
    color: #2dd4bf;
    border-color: rgba(45, 212, 191, 0.2);
  }

  .card-body {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 1rem;
  }

  .info-label {
    font-size: 0.75rem;
    color: #484f58;
    flex-shrink: 0;
  }

  .info-value {
    font-size: 0.75rem;
    color: #8b949e;
    text-align: right;
    word-break: break-all;
  }

  .update-banner {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-top: 0.75rem;
    padding: 0.5rem 0.75rem;
    border-radius: 8px;
    background: rgba(210, 153, 34, 0.1);
    border: 1px solid rgba(210, 153, 34, 0.2);
    font-size: 0.75rem;
    font-weight: 600;
    color: #d29922;
  }

  .update-icon {
    width: 14px;
    height: 14px;
    flex-shrink: 0;
  }
</style>
