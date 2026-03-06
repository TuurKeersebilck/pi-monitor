<script>
  import { onMount, onDestroy } from 'svelte'
  import SystemStats from './lib/SystemStats.svelte'
  import ServiceCard from './lib/ServiceCard.svelte'
  import PiholeStats from './lib/PiholeStats.svelte'
  import ImmichStats from './lib/ImmichStats.svelte'
  import Bookmarks from './lib/Bookmarks.svelte'

  let system = $state(null)
  let containers = $state([])
  let pihole = $state(null)
  let immich = $state(null)
  let connected = $state(false)
  let query = $state('')
  let time = $state('')
  let date = $state('')
  let ws

  // Group + filter containers
  let grouped = $derived.by(() => {
    const q = query.toLowerCase()
    const filtered = q
      ? containers.filter(c => c.name.toLowerCase().includes(q))
      : containers
    const map = {}
    for (const c of filtered) {
      const g = c.group || 'Other'
      if (!map[g]) map[g] = []
      map[g].push(c)
    }
    return Object.entries(map).sort(([a], [b]) => a.localeCompare(b))
  })

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

  let clockId
  function tickClock() {
    const now = new Date()
    time = now.toLocaleTimeString('en', { hour: '2-digit', minute: '2-digit' })
    date = now.toLocaleDateString('en', { weekday: 'short', month: 'short', day: 'numeric' })
  }

  onMount(() => {
    connect()
    tickClock()
    clockId = setInterval(tickClock, 1000)
  })
  onDestroy(() => {
    if (ws) ws.close()
    clearInterval(clockId)
  })
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

    <div class="search-wrap">
      <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="15" height="15">
        <circle cx="11" cy="11" r="8"/>
        <line x1="21" y1="21" x2="16.65" y2="16.65"/>
      </svg>
      <input
        class="search"
        type="search"
        placeholder="Search services…"
        bind:value={query}
      />
    </div>

    <div class="header-right">
      <div class="clock">
        <span class="clock-time">{time}</span>
        <span class="clock-date">{date}</span>
      </div>
      <span class="conn-badge" class:online={connected}>
        <span class="conn-dot"></span>
        {connected ? 'Live' : 'Reconnecting…'}
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

    <section>
      <div class="section-header">
        <div>
          <h2 class="section-heading">Services</h2>
          <p class="section-sub">Quick access to your homelab</p>
        </div>
        {#if query}
          <span class="result-count">filtering: "{query}"</span>
        {/if}
      </div>
      <Bookmarks {query} />
    </section>

    {#if pihole || immich}
      <div class="middle-row">
        {#if pihole}<PiholeStats {pihole} />{/if}
        {#if immich}<ImmichStats {immich} />{/if}
      </div>
    {/if}

    {#if containers.length > 0}
      <section>
        <div class="section-header">
          <div>
            <h2 class="section-heading" style="display:flex;align-items:center;gap:0.5rem">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20" style="color:#6366f1">
                <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/>
                <rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/>
              </svg>
              Docker Containers
            </h2>
            <p class="section-sub">Live container status</p>
          </div>
        </div>

        {#each grouped as [group, items]}
          <div class="service-group">
            <div class="group-label">
              <span class="group-name">{group}</span>
              <span class="group-count">{items.length}</span>
            </div>
            <div class="services-grid">
              {#each items as container (container.name)}
                <ServiceCard {container} />
              {/each}
            </div>
          </div>
        {/each}
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

  /* ── Header ── */
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0 2rem;
    height: 56px;
    background: #161b22;
    border-bottom: 1px solid #21262d;
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-shrink: 0;
  }

  .logo { width: 20px; height: 20px; color: #2dd4bf; }

  .brand {
    font-size: 1.05rem;
    font-weight: 700;
    color: #f0f6fc;
    letter-spacing: -0.01em;
    white-space: nowrap;
  }

  /* Search */
  .search-wrap {
    position: relative;
    flex: 1;
    max-width: 340px;
  }

  .search-icon {
    position: absolute;
    left: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    color: #484f58;
    pointer-events: none;
  }

  .search {
    width: 100%;
    background: #0d1117;
    border: 1px solid #21262d;
    border-radius: 8px;
    padding: 0.4rem 0.75rem 0.4rem 2.25rem;
    font-size: 0.85rem;
    color: #c9d1d9;
    outline: none;
    transition: border-color 0.2s;
    font-family: inherit;
  }

  .search::placeholder { color: #484f58; }
  .search:focus { border-color: #30363d; }
  .search::-webkit-search-cancel-button { display: none; }

  /* Right side */
  .header-right {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-shrink: 0;
  }

  .clock {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 1px;
  }

  .clock-time {
    font-size: 0.95rem;
    font-weight: 700;
    color: #f0f6fc;
    letter-spacing: -0.01em;
    line-height: 1;
  }

  .clock-date {
    font-size: 0.68rem;
    color: #8b949e;
    line-height: 1;
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
    white-space: nowrap;
  }

  .conn-badge.online {
    background: rgba(45, 212, 191, 0.1);
    color: #2dd4bf;
    border-color: rgba(45, 212, 191, 0.25);
  }

  .conn-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: currentColor;
  }

  /* ── Main ── */
  main {
    flex: 1;
    max-width: 1400px;
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

  .section-sub {
    margin: 0;
    font-size: 0.8rem;
    color: #484f58;
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
    width: 7px;
    height: 7px;
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

  .result-count {
    font-size: 0.75rem;
    color: #8b949e;
    align-self: center;
  }

  /* ── Widget row ── */
  .middle-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  /* ── Service groups ── */
  .service-group { display: flex; flex-direction: column; gap: 0.75rem; margin-bottom: 1.75rem; }
  .service-group:last-child { margin-bottom: 0; }

  .group-label {
    display: flex;
    align-items: center;
    gap: 0.6rem;
  }

  .group-name {
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.07em;
    text-transform: uppercase;
    color: #484f58;
  }

  .group-count {
    font-size: 0.65rem;
    font-weight: 700;
    padding: 0.1rem 0.45rem;
    border-radius: 999px;
    background: #21262d;
    color: #8b949e;
  }

  .services-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    gap: 1rem;
  }

  /* ── Footer ── */
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

  /* ── Responsive ── */
  @media (max-width: 900px) {
    .search-wrap { max-width: 220px; }
    .clock { display: none; }
  }

  @media (max-width: 768px) {
    main { padding: 1rem; }
    .middle-row { grid-template-columns: 1fr; }
    .services-grid { grid-template-columns: 1fr 1fr; }
  }

  @media (max-width: 600px) {
    header { padding: 0 1rem; }
    .search-wrap { display: none; }
  }

  @media (max-width: 480px) {
    .services-grid { grid-template-columns: 1fr; }
  }
</style>
