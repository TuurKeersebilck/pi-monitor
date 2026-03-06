<script>
  import { onMount } from 'svelte'

  let services = $state([])
  let showDialog = $state(false)
  let form = $state({ name: '', url: '', icon: '' })
  let error = $state('')

  onMount(() => {
    try {
      const stored = localStorage.getItem('dashboard-services')
      if (stored) services = JSON.parse(stored)
    } catch {}
  })

  function persist() {
    localStorage.setItem('dashboard-services', JSON.stringify(services))
  }

  function addService() {
    const name = form.name.trim()
    const url = form.url.trim()
    if (!name || !url) { error = 'Name and URL are required.'; return }
    services = [...services, { name, url, icon: form.icon.trim() }]
    persist()
    form = { name: '', url: '', icon: '' }
    error = ''
    showDialog = false
  }

  function remove(i) {
    services = services.filter((_, idx) => idx !== i)
    persist()
  }

  function openDialog() {
    form = { name: '', url: '', icon: '' }
    error = ''
    showDialog = true
  }

  function iconSrc(s) {
    if (s.icon) return s.icon
    try { return new URL(s.url).origin + '/favicon.ico' } catch { return '' }
  }

  // Deterministic letter color from name
  const PALETTE = ['#007AFF','#34C759','#FF9500','#AF52DE','#FF2D55','#00C7BE','#FF6B35','#5856D6']
  function letterColor(name) {
    return PALETTE[(name.charCodeAt(0) + name.length) % PALETTE.length]
  }

  function handleKey(e) {
    if (e.key === 'Escape') showDialog = false
    if (e.key === 'Enter' && showDialog) addService()
  }
</script>

<svelte:window onkeydown={handleKey} />

<div class="grid">
  {#each services as service, i (service.name + i)}
    <div class="wrap">
      <button class="remove-btn" onclick={() => remove(i)} title="Remove {service.name}" aria-label="Remove {service.name}">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" width="10" height="10">
          <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
      <a class="card" href={service.url} target="_blank" rel="noopener noreferrer">
        <div class="icon-box" style="background:{letterColor(service.name)}">
          <span class="letter">{service.name[0].toUpperCase()}</span>
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

  <!-- Add button — same grid slot as cards -->
  <button class="add-card" onclick={openDialog} title="Add service">
    <div class="add-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" width="22" height="22">
        <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
      </svg>
    </div>
    <span class="name muted">Add</span>
  </button>
</div>

<!-- Add Service Dialog -->
{#if showDialog}
  <div class="overlay" role="presentation" onclick={() => showDialog = false}>
    <div class="dialog" role="dialog" aria-modal="true" aria-label="Add service" onclick={(e) => e.stopPropagation()}>
      <h3 class="dialog-title">Add Service</h3>

      <div class="field">
        <label class="field-label" for="svc-name">Name</label>
        <input id="svc-name" class="field-input" type="text" bind:value={form.name} placeholder="Jellyfin" autocomplete="off" />
      </div>

      <div class="field">
        <label class="field-label" for="svc-url">URL</label>
        <input id="svc-url" class="field-input" type="url" bind:value={form.url} placeholder="http://pi.local:8096" autocomplete="off" />
      </div>

      <div class="field">
        <label class="field-label" for="svc-icon">
          Icon URL
          <span class="optional">optional — leave blank to auto-detect</span>
        </label>
        <input id="svc-icon" class="field-input" type="url" bind:value={form.icon} placeholder="http://pi.local:8096/favicon.ico" autocomplete="off" />
      </div>

      {#if error}
        <p class="error">{error}</p>
      {/if}

      <div class="actions">
        <button class="btn-cancel" onclick={() => showDialog = false}>Cancel</button>
        <button class="btn-add" onclick={addService}>Add</button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* ── Grid ── */
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 12px;
  }

  /* ── Card wrap (for remove button positioning) ── */
  .wrap { position: relative; }

  .remove-btn {
    position: absolute;
    top: 5px;
    left: 5px;
    z-index: 10;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #ff3b30;
    border: 2px solid white;
    color: white;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    opacity: 0;
    transition: opacity 0.15s;
    box-shadow: 0 1px 4px rgba(0,0,0,0.2);
  }

  .wrap:hover .remove-btn { opacity: 1; }

  /* ── Service card ── */
  .card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 14px 8px 12px;
    background: white;
    border-radius: 18px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.04);
    text-decoration: none;
    transition: transform 0.15s ease, box-shadow 0.15s ease;
    width: 100%;
  }

  .card:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 24px rgba(0,0,0,0.1), 0 2px 6px rgba(0,0,0,0.06);
  }

  /* ── Icon ── */
  .icon-box {
    width: 58px;
    height: 58px;
    border-radius: 15px; /* iOS icon corner ratio */
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

  /* ── Name ── */
  .name {
    font-size: 11px;
    font-weight: 500;
    color: #1c1c1e;
    text-align: center;
    line-height: 1.3;
    max-width: 84px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .name.muted { color: #aeaeb2; }

  /* ── Add button (same shape as card) ── */
  .add-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 14px 8px 12px;
    background: transparent;
    border: 1.5px dashed rgba(60,60,67,0.2);
    border-radius: 18px;
    cursor: pointer;
    transition: border-color 0.15s, background 0.15s;
    width: 100%;
  }

  .add-card:hover {
    border-color: #007AFF;
    background: rgba(0,122,255,0.04);
  }

  .add-icon {
    width: 58px;
    height: 58px;
    border-radius: 15px;
    background: rgba(60,60,67,0.07);
    display: flex;
    align-items: center;
    justify-content: center;
    color: #aeaeb2;
    transition: background 0.15s, color 0.15s;
  }

  .add-card:hover .add-icon {
    background: rgba(0,122,255,0.1);
    color: #007AFF;
  }

  /* ── Dialog overlay ── */
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
    background: white;
    border-radius: 20px;
    padding: 1.75rem;
    width: 100%;
    max-width: 380px;
    box-shadow: 0 24px 60px rgba(0,0,0,0.18), 0 4px 16px rgba(0,0,0,0.1);
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .dialog-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: #1c1c1e;
    margin: 0;
    letter-spacing: -0.02em;
  }

  /* ── Form fields ── */
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .field-label {
    font-size: 0.78rem;
    font-weight: 600;
    color: #6c6c70;
  }

  .optional {
    font-weight: 400;
    color: #aeaeb2;
    margin-left: 0.3rem;
  }

  .field-input {
    background: #f2f2f7;
    border: 1px solid transparent;
    border-radius: 10px;
    padding: 0.6rem 0.875rem;
    font-size: 0.9rem;
    color: #1c1c1e;
    outline: none;
    font-family: inherit;
    transition: border-color 0.15s;
  }

  .field-input:focus { border-color: #007AFF; background: white; }
  .field-input::placeholder { color: #aeaeb2; }

  .error {
    font-size: 0.8rem;
    color: #ff3b30;
    margin: 0;
  }

  /* ── Dialog actions ── */
  .actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    margin-top: 0.25rem;
  }

  .btn-cancel {
    padding: 0.55rem 1.1rem;
    border-radius: 10px;
    border: none;
    background: rgba(60,60,67,0.1);
    color: #1c1c1e;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
  }

  .btn-cancel:hover { background: rgba(60,60,67,0.16); }

  .btn-add {
    padding: 0.55rem 1.4rem;
    border-radius: 10px;
    border: none;
    background: #007AFF;
    color: white;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s;
  }

  .btn-add:hover { background: #0071e3; }

  /* ── Responsive ── */
  @media (max-width: 600px) {
    .grid { grid-template-columns: repeat(auto-fill, minmax(88px, 1fr)); gap: 10px; }
  }
</style>
