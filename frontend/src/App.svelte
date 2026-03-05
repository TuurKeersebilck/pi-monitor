<script>
  import { onMount, onDestroy } from 'svelte'
  import SystemStats from './lib/SystemStats.svelte'
  import ServiceCard from './lib/ServiceCard.svelte'
  import PiholeStats from './lib/PiholeStats.svelte'
  import ImmichStats from './lib/ImmichStats.svelte'

  let system = $state(null)
  let containers = $state([])
  let pihole = $state(null)
  let immich = $state(null)
  let connected = $state(false)
  let ws

  function connect() {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    ws = new WebSocket(`${protocol}//${location.host}/ws`)
    ws.onopen = () => { connected = true }
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data)
      if (data.system) system = data.system
      if (data.containers) containers = data.containers
      if (data.pihole) pihole = data.pihole
      if (data.immich) immich = data.immich
    }
    ws.onclose = () => { connected = false; setTimeout(connect, 3000) }
  }

  onMount(() => connect())
  onDestroy(() => { if (ws) ws.close() })
</script>

<div class="app">
  <header>
    <div class="header-left">
      <svg class="logo" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/>
        <rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/>
      </svg>
      <span class="brand">Pi Dashboard</span>
    </div>
    <div class="header-right">
      <span class="conn-badge" class:online={connected}>
        <span class="conn-dot"></span>
        {connected ? 'Live' : 'Reconnecting...'}
      </span>
    </div>
  </header>

  <main>
    {#if system}
      <section>
        <div class="section-header">
          <div>
            <h2 class="section-heading">System Overview</h2>
            <p class="system-sub">
              <span class="live-dot"></span>
              Raspberry Pi 5
              {#if system.info?.hostname}&nbsp;·&nbsp;Host: <span class="teal">{system.info.hostname}</span>{/if}
              {#if system.info?.ip}&nbsp;·&nbsp;IP: {system.info.ip}{/if}
            </p>
          </div>
          {#if system.info?.uptime}
            <span class="uptime-badge">UPTIME: {system.info.uptime}</span>
          {/if}
        </div>
        <SystemStats {system} />
      </section>
    {/if}

    {#if pihole || immich}
      <div class="middle-row">
        {#if pihole}<PiholeStats {pihole} />{/if}
        {#if immich}<ImmichStats {immich} />{/if}
      </div>
    {/if}

    {#if containers.length > 0}
      <section>
        <div class="section-header">
          <h2 class="section-heading" style="display:flex;align-items:center;gap:0.5rem">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20" style="color:#2dd4bf">
              <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/>
              <rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/>
            </svg>
            Services
          </h2>
        </div>
        <div class="services-grid">
          {#each containers as container (container.name)}
            <ServiceCard {container} />
          {/each}
        </div>
      </section>
    {/if}
  </main>

  <footer>
    <span class="footer-left">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="13" height="13">
        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
      </svg>
      Powered by Raspberry Pi 5 – Debian Bookworm
    </span>
    <span class="footer-right">Pi Dashboard v1.1.0</span>
  </footer>
</div>

<style>
  :global(*) { box-sizing: border-box; }

  :global(body) {
    margin: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    background: #0d1117;
    color: #c9d1d9;
    min-height: 100vh;
  }

  .app { min-height: 100vh; display: flex; flex-direction: column; }

  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.875rem 2rem;
    background: #161b22;
    border-bottom: 1px solid #21262d;
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .header-left { display: flex; align-items: center; gap: 0.6rem; }

  .logo { width: 20px; height: 20px; color: #2dd4bf; }

  .brand {
    font-size: 1.05rem;
    font-weight: 700;
    color: #f0f6fc;
    letter-spacing: -0.01em;
  }

  .conn-badge {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.72rem;
    font-weight: 500;
    padding: 0.28rem 0.7rem;
    border-radius: 999px;
    background: rgba(248, 81, 73, 0.12);
    color: #f85149;
    border: 1px solid rgba(248, 81, 73, 0.25);
  }

  .conn-badge.online {
    background: rgba(45, 212, 191, 0.1);
    color: #2dd4bf;
    border-color: rgba(45, 212, 191, 0.25);
  }

  .conn-dot {
    width: 6px; height: 6px;
    border-radius: 50%;
    background: currentColor;
  }

  main {
    flex: 1;
    max-width: 1200px;
    width: 100%;
    margin: 0 auto;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .section-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 1.25rem;
  }

  .section-heading {
    font-size: 1.5rem;
    font-weight: 800;
    color: #f0f6fc;
    margin: 0 0 0.3rem;
    letter-spacing: -0.02em;
  }

  .system-sub {
    margin: 0;
    font-size: 0.8rem;
    color: #8b949e;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
  }

  .live-dot {
    width: 7px; height: 7px;
    border-radius: 50%;
    background: #3fb950;
    display: inline-block;
    margin-right: 0.4rem;
    flex-shrink: 0;
  }

  .teal { color: #2dd4bf; font-weight: 500; }

  .uptime-badge {
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.05em;
    padding: 0.28rem 0.8rem;
    border-radius: 999px;
    background: rgba(45, 212, 191, 0.08);
    color: #2dd4bf;
    border: 1px solid rgba(45, 212, 191, 0.2);
    white-space: nowrap;
    align-self: flex-start;
    margin-top: 0.15rem;
  }

  .middle-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .services-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 1rem;
  }

  footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.875rem 2rem;
    border-top: 1px solid #21262d;
    background: #161b22;
  }

  .footer-left {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.73rem;
    color: #484f58;
  }

  .footer-right {
    font-size: 0.7rem;
    color: #484f58;
    font-weight: 600;
    letter-spacing: 0.04em;
    text-transform: uppercase;
  }

  @media (max-width: 768px) {
    main { padding: 1rem; }
    .middle-row { grid-template-columns: 1fr; }
    .services-grid { grid-template-columns: 1fr 1fr; }
  }

  @media (max-width: 480px) {
    header { padding: 0.75rem 1rem; }
    .services-grid { grid-template-columns: 1fr; }
  }
</style>
