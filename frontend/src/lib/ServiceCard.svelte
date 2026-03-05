<script>
  let { container } = $props()
</script>

<div class="card" class:running={container.running}>
  <div class="header">
    <span class="dot" class:up={container.running}></span>
    <span class="name">{container.name}</span>
  </div>
  <div class="meta">
    <span class="image">{container.image}</span>
    <span class="status">{container.status}</span>
  </div>
  {#if container.update_available}
    <div class="update">Update available</div>
  {/if}
</div>

<style>
  .card {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 1rem;
    transition: border-color 0.2s;
  }

  .card:hover {
    border-color: #58a6ff;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #f85149;
    flex-shrink: 0;
  }

  .dot.up {
    background: #3fb950;
  }

  .name {
    font-weight: 600;
    font-size: 0.95rem;
  }

  .meta {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .image {
    font-size: 0.75rem;
    color: #8b949e;
    word-break: break-all;
  }

  .status {
    font-size: 0.75rem;
    color: #8b949e;
  }

  .update {
    margin-top: 0.5rem;
    font-size: 0.75rem;
    color: #d29922;
    font-weight: 600;
  }
</style>
