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

<div class="app">
  <header>
    <div class="header-left">
      <svg class="logo" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
        <line x1="8" y1="21" x2="16" y2="21"/>
        <line x1="12" y1="17" x2="12" y2="21"/>
      </svg>
      <h1>Pi Dashboard</h1>
    </div>
    <div class="header-right">
      <span class="connection-badge" class:online={connected}>
        <span class="conn-dot"></span>
        {connected ? 'Connected' : 'Reconnecting...'}
      </span>
    </div>
  </header>

  <main>
    {#if system}
      <section class="section">
        <h2 class="section-title">System Overview</h2>
        <SystemStats {system} />
      </section>
    {/if}

    {#if pihole}
      <section class="section">
        <h2 class="section-title">Pi-hole Statistics</h2>
        <PiholeStats {pihole} />
      </section>
    {/if}

    {#if containers.length > 0}
      <section class="section">
        <h2 class="section-title">Services</h2>
        <div class="services-grid">
          {#each containers as container (container.name)}
            <ServiceCard {container} />
          {/each}
        </div>
      </section>
    {/if}
  </main>

  <footer>
    <span class="footer-text">Pi Dashboard · {__APP_VERSION__}</span>
  </footer>
</div>

<style>
  :global(*) {
    box-sizing: border-box;
  }

  :global(body) {
    margin: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #0b0e14;
    color: #c9d1d9;
    min-height: 100vh;
  }

  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 2rem;
    background: #111720;
    border-bottom: 1px solid #1c2333;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .logo {
    width: 24px;
    height: 24px;
    color: #2dd4bf;
  }

  h1 {
    font-size: 1.25rem;
    font-weight: 700;
    margin: 0;
    color: #e6edf3;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .connection-badge {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.75rem;
    padding: 0.3rem 0.75rem;
    border-radius: 999px;
    background: rgba(248, 81, 73, 0.1);
    color: #f85149;
    border: 1px solid rgba(248, 81, 73, 0.2);
  }

  .connection-badge.online {
    background: rgba(45, 212, 191, 0.1);
    color: #2dd4bf;
    border-color: rgba(45, 212, 191, 0.2);
  }

  .conn-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: currentColor;
  }

  main {
    flex: 1;
    max-width: 1100px;
    width: 100%;
    margin: 0 auto;
    padding: 1.5rem 2rem;
  }

  .section {
    margin-bottom: 2rem;
  }

  .section-title {
    font-size: 1rem;
    font-weight: 600;
    color: #8b949e;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin: 0 0 1rem;
  }

  .services-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
  }

  footer {
    padding: 1rem 2rem;
    border-top: 1px solid #1c2333;
    text-align: center;
  }

  .footer-text {
    font-size: 0.75rem;
    color: #484f58;
  }

  @media (max-width: 640px) {
    header {
      padding: 1rem;
    }
    main {
      padding: 1rem;
    }
    .services-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
