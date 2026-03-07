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
  <span class="name">{container.name}</span>
  <div class="tags">
    <span class="tag" class:tag-running={container.running} class:tag-stopped={!container.running}>
      {container.running ? 'Running' : 'Stopped'}
    </span>
    {#if container.running && container.uptime && container.uptime !== 'stopped'}
      <span class="tag tag-uptime">{container.uptime}</span>
    {/if}
    {#if container.update_available}
      <span class="tag tag-update">Update</span>
    {/if}
  </div>
</svelte:element>

<style>
  .card {
    background: var(--card-bg);
    border-radius: 14px;
    padding: 1rem 1.1rem;
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    box-shadow: 0 2px 10px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.04);
    text-decoration: none;
    color: inherit;
    transition: box-shadow 0.2s, transform 0.2s, background 0.2s;
  }

  .card.clickable { cursor: pointer; }
  .card.clickable:hover {
    box-shadow: 0 6px 20px rgba(0,0,0,0.1), 0 2px 6px rgba(0,0,0,0.06);
    transform: translateY(-1px);
  }

  .name {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--text);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.3rem;
    margin-top: 0.15rem;
  }

  .tag {
    font-size: 0.65rem;
    font-weight: 600;
    padding: 0.15rem 0.5rem;
    border-radius: 999px;
    letter-spacing: 0.02em;
  }

  .tag-running {
    background: rgba(52,199,89,0.12);
    color: #34c759;
  }

  .tag-stopped {
    background: rgba(255,59,48,0.1);
    color: #ff3b30;
  }

  .tag-uptime {
    background: var(--pill);
    color: var(--text-2);
  }

  .tag-update {
    background: rgba(255,149,0,0.12);
    color: #ff9500;
  }
</style>
