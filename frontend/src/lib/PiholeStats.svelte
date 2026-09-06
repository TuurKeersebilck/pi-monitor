<script>
  import StatCard from './StatCard.svelte'

  let { pihole } = $props()
  const fmt = (n) => n.toLocaleString()

  let stats = $derived([
    { label: 'QUERIES', value: fmt(pihole.total_queries) },
    { label: 'BLOCKED', value: fmt(pihole.blocked_queries), color: 'accent' },
    { label: 'BLOCK RATE', value: pihole.blocked_percent.toFixed(1) + '%', color: 'accent' },
    {
      label: 'BLOCKLIST',
      value: pihole.domains_blocked >= 1000
        ? (pihole.domains_blocked / 1000).toFixed(0) + 'k'
        : fmt(pihole.domains_blocked),
    },
  ])
</script>

<StatCard title="Pi-hole" badge="LIVE" {stats}>
  {#snippet icon()}
    <svg viewBox="0 0 24 24" fill="none" stroke="#007AFF" stroke-width="2" width="16" height="16">
      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
    </svg>
  {/snippet}
</StatCard>
