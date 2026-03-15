<script>
  import { onMount, onDestroy } from 'svelte'
  import SystemStats from './lib/SystemStats.svelte'
  import ServiceCard from './lib/ServiceCard.svelte'
  import PiholeStats from './lib/PiholeStats.svelte'
  import ImmichStats from './lib/ImmichStats.svelte'
  import Bookmarks from './lib/Bookmarks.svelte'

  let system = $state(null)
  let history = $state([])
  let containers = $state([])
  let pihole = $state(null)
  let immich = $state(null)
  let connected = $state(false)
  let time = $state('')
  let date = $state('')
  let showBgPanel = $state(false)
  let bgInput = $state('')
  let bgBlur = $state(0)
  let bgUploading = $state(false)
  let darkMode = $state(false)

  // Config owned here — services + background persisted to backend
  let services = $state([])
  let ws

  // Whether a wallpaper image is active (for readability helpers)
  let hasWallpaper = $derived(
    bgInput && (bgInput.startsWith('http') || bgInput.startsWith('/') || bgInput.startsWith('data:'))
  )

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
      bgBlur = cfg.background_blur ?? 0
      if (bgInput) applyBg(bgInput)
      applyBlur(bgBlur)
    } catch (e) {
      console.error('Failed to load config:', e)
    }
  }

  async function saveConfig(newServices, newBg, newBlur) {
    try {
      await fetch('/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ services: newServices, background: newBg, background_blur: newBlur ?? bgBlur }),
      })
    } catch (e) {
      console.error('Failed to save config:', e)
    }
  }

  function applyDark(dark) {
    document.documentElement.classList.toggle('dark', dark)
  }

  function toggleDark() {
    darkMode = !darkMode
    localStorage.setItem('darkMode', String(darkMode))
    applyDark(darkMode)
  }

  function onServicesChange(updated) {
    services = updated
    saveConfig(services, bgInput, bgBlur)
  }

  function connect() {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    ws = new WebSocket(`${protocol}//${location.host}/ws`)
    ws.onopen = () => { connected = true }
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data)
      if (data.system) {
        system = data.system
        history = [...history.slice(-119), data.system]
      }
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
    const layer = document.getElementById('bg-layer')
    if (!layer) return
    if (!value) {
      layer.style.backgroundImage = ''
      layer.style.backgroundColor = ''
      return
    }
    if (value.startsWith('http') || value.startsWith('/') || value.startsWith('data:')) {
      layer.style.backgroundImage = `url(${value})`
      layer.style.backgroundColor = ''
    } else {
      layer.style.backgroundImage = ''
      layer.style.backgroundColor = value
    }
  }

  function applyBlur(value) {
    const layer = document.getElementById('bg-layer')
    if (!layer) return
    layer.style.filter = value > 0 ? `blur(${value}px)` : ''
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
      saveConfig(services, url, bgBlur)
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
    applyBlur(bgBlur)
    saveConfig(services, bgInput, bgBlur)
    showBgPanel = false
  }

  function clearBg() {
    bgInput = ''
    bgBlur = 0
    applyBg('')
    applyBlur(0)
    saveConfig(services, '', 0)
    showBgPanel = false
  }

  function handleBlurChange(e) {
    bgBlur = parseFloat(e.target.value)
    applyBlur(bgBlur)
  }

  function handleBlurCommit() {
    saveConfig(services, bgInput, bgBlur)
  }

  onMount(() => {
    darkMode = localStorage.getItem('darkMode') === 'true'
    applyDark(darkMode)
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

<!-- Fixed background layer — blur filter applied here so content stays sharp -->
<div id="bg-layer"></div>

<div class="app" class:has-wallpaper={hasWallpaper}>
  <header>
    <div class="header-left">
      <span class="brand">Pi Dashboard</span>
    </div>

    <div class="header-right">
      <div class="clock">
        <span class="clock-time">{time}</span>
        <span class="clock-date">{date}</span>
      </div>

      <!-- Dark mode toggle -->
      <button class="icon-btn" onclick={toggleDark} title={darkMode ? 'Switch to light mode' : 'Switch to dark mode'} aria-label="Toggle dark mode">
        {#if darkMode}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="17" height="17">
            <circle cx="12" cy="12" r="5"/>
            <line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
            <line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          </svg>
        {:else}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="17" height="17">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        {/if}
      </button>

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
          <div class="bg-panel card-panel">
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

            <!-- Blur slider -->
            <div class="blur-control">
              <label class="blur-label" for="blur-slider">
                Blur
                <span class="blur-value">{bgBlur.toFixed(0)}px</span>
              </label>
              <input
                id="blur-slider"
                type="range"
                min="0"
                max="30"
                step="1"
                value={bgBlur}
                oninput={handleBlurChange}
                onchange={handleBlurCommit}
                class="blur-slider"
              />
            </div>

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
        <SystemStats {system} {history} />
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
</div>

<style>
  /* ── Background layer ── */
  :global(#bg-layer) {
    position: fixed;
    inset: 0;
    z-index: -1;
    background-size: cover;
    background-position: center;
    background-attachment: fixed;
    background-color: var(--bg);
    transition: filter 0.3s ease;
    /* Slightly oversize to prevent blur edge artifacts */
    margin: -8px;
    padding: 8px;
  }

  .app { min-height: 100vh; display: flex; flex-direction: column; }

  /* ── Readability when wallpaper is active ── */
  .has-wallpaper .section-heading {
    text-shadow: 0 1px 8px rgba(0, 0, 0, 0.4), 0 0 2px rgba(0, 0, 0, 0.2);
  }

  .has-wallpaper .section-sub {
    text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
  }

  .has-wallpaper .uptime-badge {
    backdrop-filter: blur(12px) saturate(1.4);
    -webkit-backdrop-filter: blur(12px) saturate(1.4);
    background: rgba(0, 122, 255, 0.15);
  }

  .has-wallpaper .group-name,
  .has-wallpaper .group-count {
    text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
  }

  /* ── Header ── */
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0 2rem;
    height: 52px;
    background: var(--header-bg);
    backdrop-filter: saturate(180%) blur(20px);
    -webkit-backdrop-filter: saturate(180%) blur(20px);
    border-bottom: 0.5px solid var(--border);
    position: sticky;
    top: 0;
    z-index: 20;
    transition: background 0.2s ease;
  }

  .header-left { flex-shrink: 0; }

  .brand {
    font-size: 1rem;
    font-weight: 700;
    color: var(--text);
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
    color: var(--text);
    letter-spacing: -0.01em;
    line-height: 1.15;
  }

  .clock-date {
    font-size: 0.65rem;
    color: var(--text-2);
    line-height: 1.2;
  }

  /* Icon button */
  .icon-btn {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    border: none;
    background: var(--pill);
    color: var(--text-2);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.15s, color 0.15s;
  }

  .icon-btn:hover { background: var(--pill-hover); color: var(--text); }

  /* Background panel */
  .bg-wrap { position: relative; }

  .bg-panel {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    width: 280px;
    border-radius: 18px;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    z-index: 50;
  }

  .bg-label {
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--text-2);
    margin: 0;
  }

  .bg-upload-btn {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.45rem 0.85rem;
    border-radius: 8px;
    background: var(--pill);
    color: var(--text);
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
    width: fit-content;
  }

  .bg-upload-btn:hover { background: var(--pill-hover); }

  .bg-or {
    font-size: 0.72rem;
    color: var(--text-3);
  }

  .bg-input {
    background: var(--input-bg);
    border: 1px solid transparent;
    border-radius: 8px;
    padding: 0.5rem 0.75rem;
    font-size: 0.85rem;
    color: var(--text);
    outline: none;
    font-family: inherit;
    transition: border-color 0.15s;
    width: 100%;
  }

  .bg-input:focus { border-color: var(--accent); background: var(--card-bg); }
  .bg-input::placeholder { color: var(--text-3); }

  /* ── Blur slider ── */
  .blur-control {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .blur-label {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--text-2);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .blur-value {
    font-size: 0.7rem;
    font-weight: 500;
    color: var(--text-3);
    font-variant-numeric: tabular-nums;
  }

  .blur-slider {
    -webkit-appearance: none;
    appearance: none;
    width: 100%;
    height: 4px;
    border-radius: 2px;
    background: var(--ring-track);
    outline: none;
    cursor: pointer;
  }

  .blur-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--accent);
    border: 2px solid white;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
    cursor: pointer;
    transition: transform 0.15s ease;
  }

  .blur-slider::-webkit-slider-thumb:hover {
    transform: scale(1.15);
  }

  .blur-slider::-moz-range-thumb {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--accent);
    border: 2px solid white;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
    cursor: pointer;
  }

  .bg-actions {
    display: flex;
    gap: 0.5rem;
    justify-content: flex-end;
  }

  .bg-btn-clear {
    padding: 0.4rem 0.9rem;
    border-radius: 8px;
    border: none;
    background: var(--pill);
    color: var(--text);
    font-size: 0.82rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
  }

  .bg-btn-apply {
    padding: 0.4rem 0.9rem;
    border-radius: 8px;
    border: none;
    background: var(--accent);
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
    padding: 2rem 2rem 3rem;
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
    color: var(--text);
    margin: 0 0 0.2rem;
    letter-spacing: -0.02em;
  }

  .section-sub {
    margin: 0;
    font-size: 0.78rem;
    color: var(--text-2);
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

  .accent { color: var(--accent); font-weight: 500; }

  .uptime-badge {
    font-size: 0.68rem;
    font-weight: 600;
    padding: 0.25rem 0.75rem;
    border-radius: 999px;
    background: rgba(0,122,255,0.08);
    color: var(--accent);
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
    color: var(--text-3);
  }

  .group-count {
    font-size: 0.62rem;
    font-weight: 700;
    padding: 0.1rem 0.4rem;
    border-radius: 999px;
    background: var(--pill);
    color: var(--text-2);
  }

  .containers-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 0.75rem;
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
