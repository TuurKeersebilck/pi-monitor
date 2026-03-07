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
  let time = $state('')
  let date = $state('')
  let showBgPanel = $state(false)
  let bgInput = $state('')
  let bgUploading = $state(false)

  // Config owned here — services + background persisted to backend
  let services = $state([])
  let ws

  // Group Docker containers
  let grouped = $derived.by(() => {
    const map = {}
    for (const c of containers) {
      const g = c.group || 'Containers'
      if (!map[g]) map[g] = []
      map[g].push(c)
    }
    return Object.entries(map).sort(([a], [b]) => a.localeCompare(b))
  })

  async function loadConfig() {
    try {
      const res = await fetch('/api/config')
      const cfg = await res.json()
      services = cfg.services ?? []
      bgInput = cfg.background ?? ''
      if (bgInput) applyBg(bgInput)
    } catch (e) {
      console.error('Failed to load config:', e)
    }
  }

  async function saveConfig(newServices, newBg) {
    try {
      await fetch('/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ services: newServices, background: newBg }),
      })
    } catch (e) {
      console.error('Failed to save config:', e)
    }
  }

  function onServicesChange(updated) {
    services = updated
    saveConfig(services, bgInput)
  }

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

  function applyBg(value) {
    const body = document.body
    if (!value) {
      body.style.backgroundImage = ''
      body.style.backgroundColor = '#f2f2f7'
      body.style.backgroundSize = ''
      body.style.backgroundPosition = ''
      body.style.backgroundAttachment = ''
      return
    }
    if (value.startsWith('http') || value.startsWith('/') || value.startsWith('data:')) {
      body.style.backgroundImage = `url(${value})`
      body.style.backgroundSize = 'cover'
      body.style.backgroundPosition = 'center'
      body.style.backgroundAttachment = 'fixed'
      body.style.backgroundColor = ''
    } else {
      body.style.backgroundImage = ''
      body.style.backgroundColor = value
    }
  }

  async function handleBgFile(e) {
    const file = e.target.files?.[0]
    if (!file) return
    bgUploading = true
    try {
      const fd = new FormData()
      fd.append('file', file)
      const res = await fetch('/api/upload', { method: 'POST', body: fd })
      const { url } = await res.json()
      bgInput = url
      applyBg(url)
      saveConfig(services, url)
      showBgPanel = false
    } catch {
      console.error('Background upload failed')
    } finally {
      bgUploading = false
      e.target.value = ''
    }
  }

  function saveBg() {
    applyBg(bgInput)
    saveConfig(services, bgInput)
    showBgPanel = false
  }

  function clearBg() {
    bgInput = ''
    applyBg('')
    saveConfig(services, '')
    showBgPanel = false
  }

  onMount(() => {
    connect()
    tickClock()
    clockId = setInterval(tickClock, 1000)
    loadConfig()
  })

  onDestroy(() => {
    if (ws) ws.close()
    clearInterval(clockId)
  })
</script>

<div class="app">
  <header>
    <div class="header-left">
      <span class="brand">Pi Dashboard</span>
    </div>

    <div class="header-right">
      <div class="clock">
        <span class="clock-time">{time}</span>
        <span class="clock-date">{date}</span>
      </div>

      <!-- Background setting -->
      <div class="bg-wrap">
        <button class="icon-btn" onclick={() => showBgPanel = !showBgPanel} title="Set background" aria-label="Set background">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="17" height="17">
            <rect x="3" y="3" width="18" height="18" rx="2"/>
            <circle cx="8.5" cy="8.5" r="1.5"/>
            <polyline points="21 15 16 10 5 21"/>
          </svg>
        </button>
        {#if showBgPanel}
          <div class="bg-panel">
            <p class="bg-label">Background</p>
            <label class="bg-upload-btn">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
              Upload image
              <input type="file" accept="image/*" onchange={handleBgFile} style="display:none" />
            </label>
            <span class="bg-or">or paste URL / color</span>
            <input
              class="bg-input"
              type="text"
              bind:value={bgInput}
              placeholder="https://… or #1a1a2e"
              onkeydown={(e) => { if (e.key === 'Enter') saveBg(); if (e.key === 'Escape') showBgPanel = false }}
            />
            <div class="bg-actions">
              <button class="bg-btn-clear" onclick={clearBg}>Clear</button>
              <button class="bg-btn-apply" onclick={saveBg}>Apply</button>
            </div>
          </div>
        {/if}
      </div>

      <span class="conn-badge" class:online={connected}>
        <span class="conn-dot"></span>
        {connected ? 'Live' : 'Reconnecting…'}
      </span>
    </div>
  </header>

  <main>
    <!-- System Overview -->
    {#if system}
      <section>
        <div class="section-header">
          <div>
            <h2 class="section-heading">System Overview</h2>
            <p class="section-sub">
              <span class="live-dot"></span>
              Raspberry Pi 5
              {#if system.info?.hostname}&nbsp;·&nbsp;<span class="accent">{system.info.hostname}</span>{/if}
              {#if system.info?.ip}&nbsp;·&nbsp;{system.info.ip}{/if}
            </p>
          </div>
          {#if system.info?.uptime}
            <span class="uptime-badge">↑ {system.info.uptime}</span>
          {/if}
        </div>
        <SystemStats {system} />
      </section>
    {/if}

    <!-- Services bookmarks -->
    <section>
      <div class="section-header">
        <div>
          <h2 class="section-heading">Services</h2>
          <p class="section-sub">Hover a card to remove it</p>
        </div>
      </div>
      <Bookmarks services={services} onchange={onServicesChange} />
    </section>

    <!-- Widgets row -->
    {#if pihole || immich}
      <div class="widget-row">
        {#if pihole}<PiholeStats {pihole} />{/if}
        {#if immich}<ImmichStats {immich} />{/if}
      </div>
    {/if}

    <!-- Docker containers -->
    {#if containers.length > 0}
      <section>
        <div class="section-header">
          <div>
            <h2 class="section-heading">Docker</h2>
            <p class="section-sub">Live container status</p>
          </div>
        </div>
        {#each grouped as [group, items]}
          {#if grouped.length > 1}
            <div class="group-label">
              <span class="group-name">{group}</span>
              <span class="group-count">{items.length}</span>
            </div>
          {/if}
          <div class="containers-grid">
            {#each items as container (container.name)}
              <ServiceCard {container} />
            {/each}
          </div>
        {/each}
      </section>
    {/if}
  </main>

  <footer>
    <span>Raspberry Pi 5 · Debian Bookworm</span>
    <span>Pi Dashboard</span>
  </footer>
</div>

<style>
  :global(*) { box-sizing: border-box; }

  :global(body) {
    margin: 0;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Helvetica Neue', sans-serif;
    background-color: #f2f2f7;
    color: #1c1c1e;
    min-height: 100vh;
    -webkit-font-smoothing: antialiased;
  }

  .app { min-height: 100vh; display: flex; flex-direction: column; }

  /* ── Header ── */
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0 2rem;
    height: 52px;
    background: rgba(255,255,255,0.72);
    backdrop-filter: saturate(180%) blur(20px);
    -webkit-backdrop-filter: saturate(180%) blur(20px);
    border-bottom: 0.5px solid rgba(60,60,67,0.18);
    position: sticky;
    top: 0;
    z-index: 20;
  }

  .header-left { flex-shrink: 0; }

  .brand {
    font-size: 1rem;
    font-weight: 700;
    color: #1c1c1e;
    letter-spacing: -0.02em;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-shrink: 0;
  }

  /* Clock */
  .clock {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 0;
  }

  .clock-time {
    font-size: 0.92rem;
    font-weight: 600;
    color: #1c1c1e;
    letter-spacing: -0.01em;
    line-height: 1.15;
  }

  .clock-date {
    font-size: 0.65rem;
    color: #6c6c70;
    line-height: 1.2;
  }

  /* Icon button */
  .icon-btn {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    border: none;
    background: rgba(60,60,67,0.08);
    color: #6c6c70;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.15s, color 0.15s;
  }

  .icon-btn:hover { background: rgba(60,60,67,0.14); color: #1c1c1e; }

  /* Background panel */
  .bg-wrap { position: relative; }

  .bg-panel {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    width: 280px;
    background: white;
    border-radius: 14px;
    padding: 1rem;
    box-shadow: 0 8px 30px rgba(0,0,0,0.14), 0 2px 8px rgba(0,0,0,0.08);
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    z-index: 50;
  }

  .bg-label {
    font-size: 0.78rem;
    font-weight: 600;
    color: #6c6c70;
    margin: 0;
  }

  .bg-upload-btn {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.45rem 0.85rem;
    border-radius: 8px;
    background: rgba(60,60,67,0.08);
    color: #1c1c1e;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
    width: fit-content;
  }

  .bg-upload-btn:hover { background: rgba(60,60,67,0.14); }

  .bg-or {
    font-size: 0.72rem;
    color: #aeaeb2;
  }

  .bg-input {
    background: #f2f2f7;
    border: 1px solid transparent;
    border-radius: 8px;
    padding: 0.5rem 0.75rem;
    font-size: 0.85rem;
    color: #1c1c1e;
    outline: none;
    font-family: inherit;
    transition: border-color 0.15s;
    width: 100%;
  }

  .bg-input:focus { border-color: #007AFF; background: white; }
  .bg-input::placeholder { color: #aeaeb2; }

  .bg-actions {
    display: flex;
    gap: 0.5rem;
    justify-content: flex-end;
  }

  .bg-btn-clear {
    padding: 0.4rem 0.9rem;
    border-radius: 8px;
    border: none;
    background: rgba(60,60,67,0.1);
    color: #1c1c1e;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
  }

  .bg-btn-apply {
    padding: 0.4rem 0.9rem;
    border-radius: 8px;
    border: none;
    background: #007AFF;
    color: white;
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
  }

  /* Connection badge */
  .conn-badge {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.7rem;
    font-weight: 600;
    padding: 0.25rem 0.65rem;
    border-radius: 999px;
    background: rgba(255,59,48,0.1);
    color: #ff3b30;
    white-space: nowrap;
  }

  .conn-badge.online {
    background: rgba(52,199,89,0.12);
    color: #34c759;
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
    max-width: 1300px;
    width: 100%;
    margin: 0 auto;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
  }

  section { display: flex; flex-direction: column; }

  .section-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 1rem;
  }

  .section-heading {
    font-size: 1.3rem;
    font-weight: 700;
    color: #1c1c1e;
    margin: 0 0 0.2rem;
    letter-spacing: -0.02em;
  }

  .section-sub {
    margin: 0;
    font-size: 0.78rem;
    color: #aeaeb2;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.2rem;
  }

  .live-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #34c759;
    display: inline-block;
    flex-shrink: 0;
  }

  .accent { color: #007AFF; font-weight: 500; }

  .uptime-badge {
    font-size: 0.68rem;
    font-weight: 600;
    padding: 0.25rem 0.75rem;
    border-radius: 999px;
    background: rgba(0,122,255,0.08);
    color: #007AFF;
    white-space: nowrap;
    align-self: flex-start;
    margin-top: 0.1rem;
  }

  /* Widget row */
  .widget-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  /* Docker groups */
  .group-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.6rem;
    margin-top: 1rem;
  }

  .group-label:first-child { margin-top: 0; }

  .group-name {
    font-size: 0.68rem;
    font-weight: 700;
    letter-spacing: 0.07em;
    text-transform: uppercase;
    color: #aeaeb2;
  }

  .group-count {
    font-size: 0.62rem;
    font-weight: 700;
    padding: 0.1rem 0.4rem;
    border-radius: 999px;
    background: rgba(60,60,67,0.08);
    color: #6c6c70;
  }

  .containers-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 0.75rem;
  }

  /* ── Footer ── */
  footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 2rem;
    font-size: 0.7rem;
    color: #aeaeb2;
    border-top: 0.5px solid rgba(60,60,67,0.12);
  }

  /* ── Responsive ── */
  @media (max-width: 900px) {
    .clock { display: none; }
  }

  @media (max-width: 768px) {
    main { padding: 1rem; gap: 2rem; }
    .widget-row { grid-template-columns: 1fr; }
    .containers-grid { grid-template-columns: 1fr 1fr; }
  }

  @media (max-width: 600px) {
    header { padding: 0 1rem; }
  }

  @media (max-width: 480px) {
    .containers-grid { grid-template-columns: 1fr; }
  }
</style>
