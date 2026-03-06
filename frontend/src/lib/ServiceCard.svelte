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
  <div class="top">
    <span class="dot" class:running={container.running}></span>
    <span class="name">{container.name}</span>
    {#if container.update_available}
      <span class="update-dot" title="Update available"></span>
    {/if}
  </div>
  <span class="uptime" class:stopped={!container.running}>
    {container.running && container.uptime !== 'stopped' ? container.uptime : 'Stopped'}
  </span>
</svelte:element>

<style>
  .card {
    background: #ffffff;
    border-radius: 14px;
    padding: 1rem 1.1rem;
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.04);
    text-decoration: none;
    color: inherit;
    transition: box-shadow 0.2s, transform 0.2s;
  }

  .card.clickable { cursor: pointer; }
  .card.clickable:hover {
    box-shadow: 0 6px 20px rgba(0,0,0,0.1), 0 2px 6px rgba(0,0,0,0.06);
    transform: translateY(-1px);
  }

  .top {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: #ff3b30;
    flex-shrink: 0;
  }

  .dot.running { background: #34c759; }

  .update-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #ff9500;
    margin-left: auto;
    flex-shrink: 0;
  }

  .name {
    font-size: 0.9rem;
    font-weight: 600;
    color: #1c1c1e;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .uptime {
    font-size: 0.75rem;
    color: #6c6c70;
    padding-left: calc(7px + 0.5rem); /* align with name */
  }

  .uptime.stopped { color: #aeaeb2; }
</style>
