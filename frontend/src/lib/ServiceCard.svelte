<script>
  let { container } = $props()
</script>

<svelte:element
  this={container.url ? 'a' : 'div'}
  href={container.url || undefined}
  target={container.url ? '_blank' : undefined}
  rel={container.url ? 'noopener noreferrer' : undefined}
  class="card"
  class:clickable={!!container.url}
>
  {#if container.update_available}
    <div class="update-dot" title="Update available">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" width="10" height="10">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="8" x2="12" y2="12"/>
        <line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
    </div>
  {/if}

  <div class="card-top">
    <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75">
      <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
      <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
      <line x1="12" y1="22.08" x2="12" y2="12"/>
    </svg>
    <span class="health" class:healthy={container.running}>
      <span class="health-dot"></span>
      {container.running ? 'ONLINE' : 'STOPPED'}
    </span>
  </div>

  <div class="card-body">
    <span class="name">{container.name}</span>
    <span class="image">{container.image}</span>
  </div>

  <div class="card-footer">
    {#if container.running && container.uptime !== 'stopped'}
      <span class="uptime">UP {container.uptime.toUpperCase()}</span>
    {:else}
      <span class="uptime stopped">STOPPED</span>
    {/if}
    {#if container.url}
      <span class="open-hint">
        OPEN
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="11" height="11">
          <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
          <polyline points="15 3 21 3 21 9"/>
          <line x1="10" y1="14" x2="21" y2="3"/>
        </svg>
      </span>
    {/if}
  </div>
</svelte:element>

<style>
  .card {
    background: #161b22;
    border: 1px solid #21262d;
    border-radius: 14px;
    padding: 1.1rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    transition: border-color 0.2s, background 0.2s;
    position: relative;
    text-decoration: none;
    color: inherit;
  }

  .card.clickable { cursor: pointer; }
  .card.clickable:hover { border-color: #30363d; background: #1c2128; }

  .update-dot {
    position: absolute;
    top: 0.65rem;
    right: 0.65rem;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: rgba(210, 153, 34, 0.15);
    border: 1px solid rgba(210, 153, 34, 0.35);
    color: #d29922;
    display: flex;
    align-items: center;
    justify-content: center;
  }

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
    width: 5px;
    height: 5px;
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

  .uptime.stopped { color: #484f58; }

  .open-hint {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    color: #484f58;
    transition: color 0.15s;
  }

  .card.clickable:hover .open-hint { color: #2dd4bf; }
</style>
