<script>
  import { services } from './services.js'

  let { query = '' } = $props()

  // SVG inner-content for each named icon
  const ICONS = {
    film:     `<rect x="2" y="2" width="20" height="20" rx="2.18"/><line x1="7" y1="2" x2="7" y2="22"/><line x1="17" y1="2" x2="17" y2="22"/><line x1="2" y1="12" x2="22" y2="12"/><line x1="2" y1="7" x2="7" y2="7"/><line x1="2" y1="17" x2="7" y2="17"/><line x1="17" y1="17" x2="22" y2="17"/><line x1="17" y1="7" x2="22" y2="7"/>`,
    image:    `<rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/>`,
    music:    `<path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/>`,
    shield:   `<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>`,
    cloud:    `<path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"/>`,
    monitor:  `<rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/>`,
    box:      `<path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/>`,
    globe:    `<circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>`,
    home:     `<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>`,
    database: `<ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>`,
    download: `<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>`,
    code:     `<polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/>`,
    lock:     `<rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>`,
    rss:      `<path d="M4 11a9 9 0 0 1 9 9"/><path d="M4 4a16 16 0 0 1 16 16"/><circle cx="5" cy="19" r="1"/>`,
    tv:       `<rect x="2" y="7" width="20" height="15" rx="2"/><polyline points="17 2 12 7 7 2"/>`,
    activity: `<polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>`,
    server:   `<rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/>`,
    wifi:     `<path d="M5 12.55a11 11 0 0 1 14.08 0"/><path d="M1.42 9a16 16 0 0 1 21.16 0"/><path d="M8.53 16.11a6 6 0 0 1 6.95 0"/><line x1="12" y1="20" x2="12.01" y2="20"/>`,
    router:   `<rect x="2" y="9" width="20" height="8" rx="2"/><line x1="8" y1="9" x2="8" y2="17"/><line x1="12" y1="9" x2="12" y2="17"/><line x1="16" y1="9" x2="16" y2="17"/><line x1="6" y1="3" x2="6" y2="9"/><line x1="12" y1="3" x2="12" y2="9"/><line x1="18" y1="3" x2="18" y2="9"/>`,
    camera:   `<path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/>`,
    book:     `<path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>`,
    settings: `<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>`,
    zap:      `<polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>`,
    layers:   `<polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/>`,
  }

  const FALLBACK = `<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>`

  function icon(name) {
    return ICONS[name] ?? FALLBACK
  }

  function hexToRgba(hex, alpha) {
    const r = parseInt(hex.slice(1, 3), 16)
    const g = parseInt(hex.slice(3, 5), 16)
    const b = parseInt(hex.slice(5, 7), 16)
    return `rgba(${r}, ${g}, ${b}, ${alpha})`
  }

  // Group and filter
  let grouped = $derived.by(() => {
    const q = query.trim().toLowerCase()
    const filtered = q
      ? services.filter(s =>
          s.name.toLowerCase().includes(q) ||
          (s.description ?? '').toLowerCase().includes(q) ||
          (s.group ?? '').toLowerCase().includes(q)
        )
      : services

    const map = {}
    for (const s of filtered) {
      const g = s.group || 'Other'
      if (!map[g]) map[g] = []
      map[g].push(s)
    }
    return Object.entries(map).sort(([a], [b]) => a.localeCompare(b))
  })
</script>

{#if grouped.length === 0}
  <p class="no-results">No services match "{query}"</p>
{:else}
  {#each grouped as [group, items]}
    <div class="service-group">
      <div class="group-label">
        <span class="group-dot" style="background:{items[0].color ?? '#8b949e'}"></span>
        <span class="group-name">{group}</span>
        <span class="group-count">{items.length}</span>
      </div>
      <div class="grid">
        {#each items as service (service.name)}
          {@const color = service.color ?? '#8b949e'}
          <a
            class="card"
            href={service.url}
            target="_blank"
            rel="noopener noreferrer"
            style="--color:{color}; --color-bg:{hexToRgba(color, 0.1)}; --color-border:{hexToRgba(color, 0.3)}"
          >
            <div class="icon-wrap" style="background:{hexToRgba(color, 0.12)}">
              <svg viewBox="0 0 24 24" fill="none" stroke={color} stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" width="20" height="20">
                {@html icon(service.icon)}
              </svg>
            </div>
            <div class="info">
              <span class="name">{service.name}</span>
              {#if service.description}
                <span class="desc">{service.description}</span>
              {/if}
            </div>
            <svg class="arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
              <polyline points="15 3 21 3 21 9"/>
              <line x1="10" y1="14" x2="21" y2="3"/>
            </svg>
          </a>
        {/each}
      </div>
    </div>
  {/each}
{/if}

<style>
  .no-results {
    font-size: 0.875rem;
    color: #484f58;
    margin: 0;
    padding: 2rem 0;
  }

  .service-group {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1.75rem;
  }
  .service-group:last-child { margin-bottom: 0; }

  .group-label {
    display: flex;
    align-items: center;
    gap: 0.55rem;
  }

  .group-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .group-name {
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.07em;
    text-transform: uppercase;
    color: #8b949e;
  }

  .group-count {
    font-size: 0.65rem;
    font-weight: 700;
    padding: 0.1rem 0.45rem;
    border-radius: 999px;
    background: #21262d;
    color: #484f58;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(230px, 1fr));
    gap: 0.75rem;
  }

  /* Card */
  .card {
    display: flex;
    align-items: center;
    gap: 0.875rem;
    padding: 0.875rem 1rem;
    background: #161b22;
    border: 1px solid #21262d;
    border-radius: 12px;
    text-decoration: none;
    color: inherit;
    transition: border-color 0.2s, background 0.2s;
    cursor: pointer;
  }

  .card:hover {
    border-color: var(--color-border);
    background: #1c2128;
  }

  .icon-wrap {
    width: 42px;
    height: 42px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    transition: background 0.2s;
  }

  .card:hover .icon-wrap {
    background: var(--color-bg) !important;
  }

  .info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
    min-width: 0;
  }

  .name {
    font-size: 0.9rem;
    font-weight: 700;
    color: #f0f6fc;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .desc {
    font-size: 0.72rem;
    color: #484f58;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .arrow {
    color: #30363d;
    flex-shrink: 0;
    transition: color 0.2s;
  }

  .card:hover .arrow { color: var(--color); }

  @media (max-width: 600px) {
    .grid { grid-template-columns: 1fr 1fr; }
  }

  @media (max-width: 400px) {
    .grid { grid-template-columns: 1fr; }
  }
</style>
