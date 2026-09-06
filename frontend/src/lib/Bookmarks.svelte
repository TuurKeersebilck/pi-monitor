<script>
  /** @type {{ services: Array, onchange: (s: Array) => void }} */
  let { services, onchange } = $props()

  let showDialog = $state(false)
  let editIndex = $state(null)
  let form = $state({ name: '', url: '', icon: '' })
  let error = $state('')
  let uploading = $state(false)
  let openMenu = $state(null)

  function openDialog() {
    form = { name: '', url: '', icon: '' }
    error = ''
    editIndex = null
    showDialog = true
  }

  function openEditDialog(i) {
    const s = services[i]
    form = { name: s.name, url: s.url, icon: s.icon }
    error = ''
    editIndex = i
    showDialog = true
  }

  function saveService() {
    const name = form.name.trim()
    const url = form.url.trim()
    if (!name || !url) { error = 'Name and URL are required.'; return }
    if (editIndex !== null) {
      const updated = [...services]
      updated[editIndex] = { name, url, icon: form.icon.trim() }
      onchange(updated)
    } else {
      onchange([...services, { name, url, icon: form.icon.trim() }])
    }
    form = { name: '', url: '', icon: '' }
    error = ''
    showDialog = false
  }

  function remove(i) {
    onchange(services.filter((_, idx) => idx !== i))
  }

  async function handleIconFile(e) {
    const file = e.target.files?.[0]
    if (!file) return
    uploading = true
    try {
      const fd = new FormData()
      fd.append('file', file)
      const res = await fetch('/api/upload', { method: 'POST', body: fd })
      const { url } = await res.json()
      form.icon = url
    } catch {
      error = 'Upload failed.'
    } finally {
      uploading = false
      e.target.value = ''
    }
  }

  function iconSrc(s) {
    if (s.icon) return s.icon
    try { return new URL(s.url).origin + '/favicon.ico' } catch { return '' }
  }

  const PALETTE = ['#007AFF','#34C759','#FF9500','#AF52DE','#FF2D55','#00C7BE','#FF6B35','#5856D6']
  function letterColor(name) {
    if (!name) return PALETTE[0]
    return PALETTE[(name.charCodeAt(0) + name.length) % PALETTE.length]
  }

  function handleKey(e) {
    if (e.key === 'Escape') {
      if (openMenu !== null) { openMenu = null; return }
      showDialog = false
      return
    }
    if (e.key === 'Enter' && showDialog) saveService()
  }
</script>

<svelte:window onkeydown={handleKey} />

<div class="grid">
  {#each services as service, i (service.name)}
    <div class="wrap">
      <button
        class="menu-btn"
        onclick={(e) => { e.stopPropagation(); openMenu = openMenu === i ? null : i }}
        title="Options for {service.name}"
        aria-label="Options for {service.name}"
      >
        <svg viewBox="0 0 24 24" fill="currentColor" width="12" height="12">
          <circle cx="12" cy="5" r="1.8"/><circle cx="12" cy="12" r="1.8"/><circle cx="12" cy="19" r="1.8"/>
        </svg>
      </button>
      {#if openMenu === i}
        <div class="menu-backdrop" role="presentation" onclick={() => openMenu = null}></div>
        <div class="tile-menu">
          <button class="tile-menu-item" onclick={() => { openMenu = null; openEditDialog(i) }}>Edit</button>
          <button class="tile-menu-item danger" onclick={() => { openMenu = null; remove(i) }}>Delete</button>
        </div>
      {/if}
      <a class="card" href={service.url} target="_blank" rel="noopener noreferrer">
        <div class="icon-box" style={service.icon ? '' : `background:${letterColor(service.name)}`}>
          {#if !service.icon}
            <span class="letter">{(service.name?.[0] ?? '?').toUpperCase()}</span>
          {/if}
          <img
            src={iconSrc(service)}
            alt=""
            class="icon-img"
            onerror={(e) => e.currentTarget.style.display = 'none'}
          />
        </div>
        <span class="name">{service.name}</span>
      </a>
    </div>
  {/each}

  <button class="add-card" onclick={openDialog} title="Add service">
    <div class="add-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" width="22" height="22">
        <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
      </svg>
    </div>
    <span class="name muted">Add</span>
  </button>
</div>

{#if showDialog}
  <div class="overlay" role="presentation" onclick={() => showDialog = false} onkeydown={(e) => { if (e.key === 'Escape') showDialog = false }}>
    <div class="dialog" role="dialog" aria-modal="true" aria-label={editIndex !== null ? 'Edit service' : 'Add service'} tabindex="-1" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()}>
      <h3 class="dialog-title">{editIndex !== null ? 'Edit Service' : 'Add Service'}</h3>

      <div class="field">
        <label class="field-label" for="svc-name">Name</label>
        <input id="svc-name" class="field-input" type="text" bind:value={form.name} placeholder="Jellyfin" autocomplete="off" />
      </div>

      <div class="field">
        <label class="field-label" for="svc-url">URL</label>
        <input id="svc-url" class="field-input" type="url" bind:value={form.url} placeholder="http://pi.local:8096" autocomplete="off" />
      </div>

      <div class="field">
        <label class="field-label" for="svc-icon-url">Icon <span class="optional">optional — leave blank to auto-detect</span></label>
        <div class="icon-pick">
          <div class="icon-preview-box">
            {#if form.icon}
              <img src={form.icon} alt="" style="width:100%;height:100%;object-fit:contain;border-radius:10px" onerror={(e) => e.currentTarget.style.opacity='0.2'} />
            {:else}
              <svg viewBox="0 0 24 24" fill="none" stroke="#aeaeb2" stroke-width="1.5" width="22" height="22">
                <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>
              </svg>
            {/if}
          </div>
          <div class="icon-inputs">
            <label class="upload-btn" class:disabled={uploading}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="13" height="13">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
              {uploading ? 'Uploading…' : 'Upload PNG'}
              <input type="file" accept="image/*" onchange={handleIconFile} disabled={uploading} style="display:none" />
            </label>
            <span class="icon-or">or paste URL</span>
            <input id="svc-icon-url" class="field-input icon-url" type="text" bind:value={form.icon} placeholder="https://example.com/icon.png" autocomplete="off" />
            {#if form.icon}
              <button class="clear-icon" onclick={() => form.icon = ''}>Clear</button>
            {/if}
          </div>
        </div>
      </div>

      {#if error}<p class="error">{error}</p>{/if}

      <div class="actions">
        <button class="btn-cancel" onclick={() => showDialog = false}>Cancel</button>
        <button class="btn-add" onclick={saveService}>{editIndex !== null ? 'Save' : 'Add'}</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 12px;
  }

  .wrap { position: relative; }

  .menu-btn {
    position: absolute;
    top: 5px;
    right: 5px;
    z-index: 10;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: rgba(0,0,0,0.45);
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
    border: none;
    color: white;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    opacity: 0;
    transition: opacity 0.15s, background 0.15s;
    box-shadow: 0 1px 4px rgba(0,0,0,0.2);
  }

  .menu-btn:hover { background: rgba(0,0,0,0.65); }

  .wrap:hover .menu-btn, .menu-btn:focus-visible { opacity: 1; }

  .menu-backdrop {
    position: fixed;
    inset: 0;
    z-index: 40;
    background: transparent;
  }

  .tile-menu {
    position: absolute;
    top: 28px;
    right: 5px;
    z-index: 50;
    min-width: 110px;
    display: flex;
    flex-direction: column;
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: 12px;
    box-shadow: 0 8px 24px rgba(0,0,0,0.18);
    overflow: hidden;
    padding: 4px;
  }

  .tile-menu-item {
    background: none;
    border: none;
    text-align: left;
    padding: 0.45rem 0.6rem;
    font-size: 0.82rem;
    font-weight: 500;
    color: var(--text);
    cursor: pointer;
    border-radius: 8px;
    font-family: inherit;
    transition: background 0.15s;
  }

  .tile-menu-item:hover { background: var(--pill); }

  .tile-menu-item.danger { color: #ff3b30; }

  .card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 14px 8px 12px;
    background: var(--card-bg);
    backdrop-filter: blur(20px) saturate(1.6);
    -webkit-backdrop-filter: blur(20px) saturate(1.6);
    border: 1px solid var(--card-border);
    border-radius: 20px;
    box-shadow: var(--card-shadow);
    text-decoration: none;
    transition: transform 0.18s ease, box-shadow 0.18s ease, background 0.18s ease;
    width: 100%;
  }

  .card:hover {
    background: var(--card-bg-hover);
    transform: translateY(-3px);
    box-shadow: var(--card-shadow-hover);
  }

  .icon-box {
    width: 58px;
    height: 58px;
    border-radius: 15px;
    position: relative;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .letter {
    position: absolute;
    font-size: 24px;
    font-weight: 700;
    color: white;
    line-height: 1;
    user-select: none;
  }

  .icon-img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .name {
    font-size: 11px;
    font-weight: 500;
    color: var(--text);
    text-align: center;
    line-height: 1.3;
    max-width: 84px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .name.muted { color: var(--text-3); }

  .add-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 14px 8px 12px;
    background: transparent;
    border: 1.5px dashed var(--border);
    border-radius: 18px;
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
    width: 100%;
  }

  .add-card:hover { border-color: var(--accent); background: rgba(0,122,255,0.04); }

  .add-icon {
    width: 58px;
    height: 58px;
    border-radius: 15px;
    background: var(--pill);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-3);
    transition: background 0.15s, color 0.15s;
  }

  .add-card:hover .add-icon { background: rgba(0,122,255,0.1); color: var(--accent); }

  /* Dialog */
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.35);
    backdrop-filter: blur(6px);
    -webkit-backdrop-filter: blur(6px);
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }

  .dialog {
    background: var(--bg-elevated);
    border: 1px solid var(--card-border);
    border-radius: 24px;
    padding: 1.75rem;
    width: 100%;
    max-width: 380px;
    box-shadow: 0 24px 60px rgba(0, 0, 0, 0.25), 0 4px 16px rgba(0, 0, 0, 0.15);
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .dialog-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--text);
    margin: 0;
    letter-spacing: -0.02em;
  }

  .field { display: flex; flex-direction: column; gap: 0.35rem; }

  .field-label { font-size: 0.78rem; font-weight: 600; color: var(--text-2); }

  .optional { font-weight: 400; color: var(--text-3); margin-left: 0.3rem; }

  .field-input {
    background: var(--input-bg);
    border: 1px solid transparent;
    border-radius: 10px;
    padding: 0.6rem 0.875rem;
    font-size: 0.9rem;
    color: var(--text);
    outline: none;
    font-family: inherit;
    transition: border-color 0.15s;
  }

  .field-input:focus { border-color: var(--accent); background: var(--card-bg); }
  .field-input::placeholder { color: var(--text-3); }

  /* Icon picker */
  .icon-pick { display: flex; gap: 0.75rem; align-items: flex-start; }

  .icon-preview-box {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    background: var(--input-bg);
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }

  .icon-inputs { flex: 1; display: flex; flex-direction: column; gap: 0.4rem; }

  .upload-btn {
    display: inline-flex;
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
    width: fit-content;
    transition: background 0.15s;
  }

  .upload-btn:hover { background: var(--pill-hover); }
  .upload-btn.disabled { opacity: 0.5; pointer-events: none; }

  .icon-or { font-size: 0.72rem; color: var(--text-3); }

  .icon-url { font-size: 0.82rem; padding: 0.45rem 0.75rem; }

  .clear-icon {
    background: none;
    border: none;
    color: #ff3b30;
    font-size: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    padding: 0;
    font-family: inherit;
    width: fit-content;
  }

  .error { font-size: 0.8rem; color: #ff3b30; margin: 0; }

  .actions { display: flex; gap: 0.75rem; justify-content: flex-end; margin-top: 0.25rem; }

  .btn-cancel {
    padding: 0.55rem 1.1rem;
    border-radius: 10px;
    border: none;
    background: var(--pill);
    color: var(--text);
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
  }

  .btn-cancel:hover { background: var(--pill-hover); }

  .btn-add {
    padding: 0.55rem 1.4rem;
    border-radius: 10px;
    border: none;
    background: var(--accent);
    color: white;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
  }

  .btn-add:hover { background: var(--accent-hover); }

  @media (max-width: 600px) {
    .grid { grid-template-columns: repeat(auto-fill, minmax(88px, 1fr)); gap: 10px; }
  }
</style>
