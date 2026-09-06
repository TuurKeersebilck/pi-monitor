<script>
  import StatCard from './StatCard.svelte'

  let { immich } = $props()
  const fmt = (n) => n.toLocaleString()

  let stats = $derived([
    { label: 'PHOTOS', value: fmt(immich.photos) },
    { label: 'VIDEOS', value: fmt(immich.videos) },
    { label: 'STORAGE', value: `${immich.usage_gb} GB`, color: 'accent' },
    { label: 'HEALTH', value: immich.healthy ? 'Healthy' : 'Unhealthy', color: immich.healthy ? 'accent' : 'red' },
  ])
</script>

<StatCard title="Immich" badge="LIBRARY" {stats}>
  {#snippet icon()}
    <svg viewBox="0 0 24 24" fill="none" stroke="#007AFF" stroke-width="2" width="16" height="16">
      <rect x="3" y="3" width="18" height="18" rx="2"/>
      <circle cx="8.5" cy="8.5" r="1.5"/>
      <polyline points="21 15 16 10 5 21"/>
    </svg>
  {/snippet}
</StatCard>
