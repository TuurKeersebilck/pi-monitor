<script>
  import { onMount, onDestroy } from 'svelte'
  import SystemStats from './lib/SystemStats.svelte'
  import DockerList from './lib/DockerList.svelte'
  import PiholeStats from './lib/PiholeStats.svelte'
  import ImmichStats from './lib/ImmichStats.svelte'
  import Bookmarks from './lib/Bookmarks.svelte'

  let history = $state([])
  // system is always exactly history's last element -- never set directly.
  let system = $derived(history.at(-1) ?? null)
  let containers = $state([])
  let containerHistory = $state({}) // name -> [{ts, cpu_percent, mem_percent, mem_used_mb}]
  let pihole = $state(null)
  let immich = $state(null)
  let connected = $state(false)
  let time = $state('')
  let date = $state('')
  let showSettings = $state(false)
  let bgInput = $state('')
  let bgBlur = $state(0)
  let bgUploading = $state(false)
  let darkMode = $state(false)

  let services = $state([])
  let ws

  let hasWallpaper = $derived(
    bgInput && (bgInput.startsWith('http') || bgInput.startsWith('/') || bgInput.startsWith('data:'))
  )

  // ── Config ────────────────────────────────────────────────────────────────

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
        body: JSON.stringify({
          services: newServices,
          background: newBg,
          background_blur: newBlur ?? bgBlur,
        }),
      })
    } catch (e) {
      console.error('Failed to save config:', e)
    }
  }

  // ── WebSocket ─────────────────────────────────────────────────────────────

  function connect() {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
    ws = new WebSocket(`${protocol}//${location.host}/ws`)
    ws.onopen = () => { connected = true }
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data)
      if (data.system) {
        history = [...history.slice(-9999), data.system]
      }
      if (data.containers) {
        const now = Math.floor(Date.now() / 1000)
        const next = { ...containerHistory }
        const activeNames = new Set(data.containers.map(c => c.name))

        // Prune removed containers
        for (const name of Object.keys(next)) {
          if (!activeNames.has(name)) delete next[name]
        }
        // Append live stats
        for (const c of data.containers) {
          const arr = next[c.name] ?? []
          next[c.name] = [...arr.slice(-119), {
            ts: now,
            cpu_percent: c.cpu_percent,
            mem_percent: c.mem_percent,
            mem_used_mb: c.mem_used_mb,
          }]
        }
        containerHistory = next
        containers = data.containers
      }
      if (data.pihole) pihole = data.pihole
      if (data.immich) immich = data.immich
    }
    ws.onclose = () => { connected = false; setTimeout(connect, 3000) }
  }

  // ── Dark mode ─────────────────────────────────────────────────────────────

  function applyDark(dark) {
    document.documentElement.classList.toggle('dark', dark)
  }

  function toggleDark() {
    darkMode = !darkMode
    localStorage.setItem('darkMode', String(darkMode))
    applyDark(darkMode)
  }

  function handleDarkToggle() {
    toggleDark()
    saveConfig(services, bgInput, bgBlur)
  }

  // ── Services ──────────────────────────────────────────────────────────────

  function onServicesChange(updated) {
    services = updated
    saveConfig(services, bgInput, bgBlur)
  }

  // ── Background ────────────────────────────────────────────────────────────

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
      showSettings = false
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
    showSettings = false
  }

  function clearBg() {
    bgInput = ''
    bgBlur = 0
    applyBg('')
    applyBlur(0)
    saveConfig(services, '', 0)
    showSettings = false
  }

  function handleBlurChange(e) {
    bgBlur = parseFloat(e.target.value)
    applyBlur(bgBlur)
  }

  function handleBlurCommit() {
    saveConfig(services, bgInput, bgBlur)
  }

  // ── Clock ─────────────────────────────────────────────────────────────────

  let clockId
  function tickClock() {
    const now = new Date()
    time = now.toLocaleTimeString('en', { hour: '2-digit', minute: '2-digit' })
    date = now.toLocaleDateString('en', { weekday: 'short', month: 'short', day: 'numeric' })
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

      <div class="settings-wrap">
        <button class="icon-btn" onclick={() => showSettings = !showSettings} title="Settings" aria-label="Settings">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="17" height="17">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </button>

        {#if showSettings}
          <div class="settings-panel card-panel">

            <!-- Appearance -->
            <p class="settings-section-label">Appearance</p>

            <div class="settings-row">
              <span class="settings-row-label">Dark mode</span>
              <button
                class="toggle"
                class:on={darkMode}
                onclick={handleDarkToggle}
                role="switch"
                aria-checked={darkMode}
                aria-label="Dark mode"
              >
                <span class="toggle-thumb"></span>
              </button>
            </div>

            <div class="settings-row col">
              <span class="settings-row-label">Background</span>
              <label class="bg-upload-btn">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="13" height="13">
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
                onkeydown={(e) => { if (e.key === 'Enter') saveBg(); if (e.key === 'Escape') showSettings = false }}
              />
              <div class="bg-actions">
                <button class="bg-btn-clear" onclick={clearBg}>Clear</button>
                <button class="bg-btn-apply" onclick={saveBg}>Apply</button>
              </div>
            </div>

            <div class="slider-row">
              <label class="slider-label" for="blur-slider">
                Blur <span class="slider-val">{bgBlur.toFixed(0)}px</span>
              </label>
              <input
                id="blur-slider"
                type="range" min="0" max="30" step="1"
                value={bgBlur}
                oninput={handleBlurChange}
                onchange={handleBlurCommit}
                class="dash-slider"
              />
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
          <div class="header-right-controls">
            {#if system.info?.uptime}
              <span class="uptime-badge">↑ {system.info.uptime}</span>
            {/if}
          </div>
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
            <p class="section-sub">Click a container to see its CPU history</p>
          </div>
        </div>
        <DockerList {containers} {containerHistory} />
      </section>
    {/if}
  </main>
</div>

<style>
  :global(#bg-layer) {
    position: fixed;
    inset: 0;
    z-index: -1;
    background-size: cover;
    background-position: center;
    background-attachment: fixed;
    background-color: var(--bg);
    transition: filter 0.3s ease;
    margin: -8px;
    padding: 8px;
  }

  .app { min-height: 100vh; display: flex; flex-direction: column; }

  .has-wallpaper .brand,
  .has-wallpaper .clock-time,
  .has-wallpaper .section-heading {
    text-shadow: 0 1px 8px rgba(0,0,0,0.5), 0 0 2px rgba(0,0,0,0.3);
  }

  .has-wallpaper .clock-date,
  .has-wallpaper .section-sub,
  .has-wallpaper :global(.group-name),
  .has-wallpaper :global(.group-count) {
    text-shadow: 0 1px 4px rgba(0,0,0,0.4);
  }

  .has-wallpaper .uptime-badge {
    backdrop-filter: blur(12px) saturate(1.4);
    -webkit-backdrop-filter: blur(12px) saturate(1.4);
    background: rgba(0,122,255,0.15);
    text-shadow: 0 1px 3px rgba(0,0,0,0.3);
  }

  .has-wallpaper .conn-badge { text-shadow: 0 1px 3px rgba(0,0,0,0.3); }

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

  /* ── Settings panel ── */
  .settings-wrap { position: relative; }

  .settings-panel {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    width: 300px;
    border-radius: 18px;
    padding: 1.1rem 1.1rem 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
    z-index: 50;
    /* Override card-panel's near-transparent background so text is always readable */
    background: var(--bg-elevated);
    border-color: var(--border);
  }

  .settings-section-label {
    font-size: 0.68rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-3);
    margin: 0;
  }

  /* Row with label + control side by side */
  .settings-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
  }

  /* Row that stacks vertically (background sub-section) */
  .settings-row.col {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.55rem;
  }

  .settings-row-label {
    font-size: 0.82rem;
    font-weight: 500;
    color: var(--text);
  }

  /* iOS-style toggle */
  .toggle {
    position: relative;
    width: 40px;
    height: 24px;
    border-radius: 999px;
    border: none;
    background: var(--ring-track);
    cursor: pointer;
    transition: background 0.2s ease;
    flex-shrink: 0;
    padding: 0;
  }

  .toggle.on { background: var(--accent); }

  .toggle-thumb {
    position: absolute;
    top: 3px;
    left: 3px;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: white;
    box-shadow: 0 1px 4px rgba(0,0,0,0.25);
    transition: transform 0.2s cubic-bezier(0.34,1.2,0.64,1);
  }

  .toggle.on .toggle-thumb { transform: translateX(16px); }

  /* Background fields */
  .bg-upload-btn {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.4rem 0.8rem;
    border-radius: 8px;
    background: var(--pill);
    color: var(--text);
    font-size: 0.8rem;
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
    margin-top: -0.2rem;
  }

  .bg-input {
    background: var(--input-bg);
    border: 1px solid transparent;
    border-radius: 8px;
    padding: 0.45rem 0.7rem;
    font-size: 0.83rem;
    color: var(--text);
    outline: none;
    font-family: inherit;
    transition: border-color 0.15s;
    width: 100%;
  }

  .bg-input:focus { border-color: var(--accent); background: var(--card-bg); }
  .bg-input::placeholder { color: var(--text-3); }

  .bg-actions {
    display: flex;
    gap: 0.5rem;
    align-self: flex-end;
  }

  .bg-btn-clear {
    padding: 0.35rem 0.85rem;
    border-radius: 8px;
    border: none;
    background: var(--pill);
    color: var(--text);
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
  }

  .bg-btn-apply {
    padding: 0.35rem 0.85rem;
    border-radius: 8px;
    border: none;
    background: var(--accent);
    color: white;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
  }

  /* Sliders */
  .slider-row {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
  }

  .slider-label {
    font-size: 0.78rem;
    font-weight: 500;
    color: var(--text-2);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .slider-val {
    font-size: 0.72rem;
    font-weight: 500;
    color: var(--text-3);
    font-variant-numeric: tabular-nums;
  }

  .dash-slider {
    -webkit-appearance: none;
    appearance: none;
    width: 100%;
    height: 4px;
    border-radius: 2px;
    background: var(--ring-track);
    outline: none;
    cursor: pointer;
  }

  .dash-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--accent);
    border: 2px solid white;
    box-shadow: 0 1px 4px rgba(0,0,0,0.2);
    cursor: pointer;
    transition: transform 0.15s ease;
  }

  .dash-slider::-webkit-slider-thumb:hover { transform: scale(1.15); }

  .dash-slider::-moz-range-thumb {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--accent);
    border: 2px solid white;
    box-shadow: 0 1px 4px rgba(0,0,0,0.2);
    cursor: pointer;
  }

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
    gap: 1rem;
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

  .header-right-controls {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-shrink: 0;
  }

  .uptime-badge {
    font-size: 0.68rem;
    font-weight: 600;
    padding: 0.25rem 0.75rem;
    border-radius: 999px;
    background: rgba(0,122,255,0.08);
    color: var(--accent);
    white-space: nowrap;
  }

  .widget-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  @media (max-width: 900px) { .clock { display: none; } }

  @media (max-width: 768px) {
    main { padding: 1rem; gap: 2rem; }
    .widget-row { grid-template-columns: 1fr; }
  }

  @media (max-width: 600px) {
    header { padding: 0 1rem; }
    .header-right-controls { gap: 0.4rem; }
  }
</style>
