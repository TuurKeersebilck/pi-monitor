<script>
  import { onMount, onDestroy } from 'svelte'
  import SystemStats from './lib/SystemStats.svelte'
  import ServiceCard from './lib/ServiceCard.svelte'
  import PiholeStats from './lib/PiholeStats.svelte'

  let system = $state(null)
  let containers = $state([])
  let pihole = $state(null)
  let connected = $state(false)
  let ws

  function connect() {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    ws = new WebSocket(`${protocol}//${location.host}/ws`)

    ws.onopen = () => {
      connected = true
    }

    ws.onmessage = (e) => {
      const data = JSON.parse(e.data)
      if (data.system) system = data.system
      if (data.containers) containers = data.containers
      if (data.pihole) pihole = data.pihole
    }

    ws.onclose = () => {
      connected = false
      // Reconnect after 3 seconds
      setTimeout(connect, 3000)
    }
  }

  onMount(() => {
    connect()
  })

  onDestroy(() => {
    if (ws) ws.close()
  })
</script>

<main>
  <header>
    <h1>Pi Dashboard</h1>
    <span class="status" class:online={connected}>
      {connected ? 'Connected' : 'Reconnecting...'}
    </span>
  </header>

  {#if system}
    <SystemStats {system} />
  {/if}

  {#if pihole}
    <PiholeStats {pihole} />
  {/if}

  {#if containers.length > 0}
    <section class="services">
      <h2>Services</h2>
      <div class="grid">
        {#each containers as container (container.name)}
          <ServiceCard {container} />
        {/each}
      </div>
    </section>
  {/if}
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #0f1117;
    color: #e1e4e8;
  }

  main {
    max-width: 960px;
    margin: 0 auto;
    padding: 1.5rem;
  }

  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 2rem;
  }

  h1 {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
  }

  h2 {
    font-size: 1.1rem;
    font-weight: 600;
    margin: 0 0 1rem;
  }

  .status {
    font-size: 0.8rem;
    padding: 0.25rem 0.75rem;
    border-radius: 999px;
    background: #f8514933;
    color: #f85149;
  }

  .status.online {
    background: #23883233;
    color: #3fb950;
  }

  .services {
    margin-top: 1.5rem;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 0.75rem;
  }
</style>
