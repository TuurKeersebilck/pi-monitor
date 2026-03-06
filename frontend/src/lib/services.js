// ┌─────────────────────────────────────────────────────────────────────────┐
// │  DASHBOARD SERVICES CONFIG                                              │
// │  Add, remove, or reorder entries to manage what appears on the         │
// │  dashboard. Changes here are reflected immediately after a page reload. │
// └─────────────────────────────────────────────────────────────────────────┘
//
// Fields:
//   name        — display name shown on the card (required)
//   url         — where clicking the card opens (required)
//   group       — category label; cards with the same group appear together
//   description — short subtitle shown below the name
//   icon        — icon name (see list below) or omit for default
//   color       — accent color as a hex string, e.g. '#8b5cf6'
//
// Available icons:
//   film · image · music · shield · cloud · monitor · box · globe · home
//   database · download · code · lock · rss · tv · activity · server · wifi
//   router · camera · book · settings · zap · layers
// ─────────────────────────────────────────────────────────────────────────────

export const services = [

  // ── Media ──────────────────────────────────────────────────────────────────
  {
    name: 'Jellyfin',
    url: 'http://pi.local:8096',
    group: 'Media',
    description: 'Media server',
    icon: 'film',
    color: '#8b5cf6',
  },
  {
    name: 'Immich',
    url: 'http://pi.local:2283',
    group: 'Media',
    description: 'Photo & video library',
    icon: 'image',
    color: '#ec4899',
  },
  {
    name: 'Navidrome',
    url: 'http://pi.local:4533',
    group: 'Media',
    description: 'Music streaming',
    icon: 'music',
    color: '#f97316',
  },

  // ── Networking ─────────────────────────────────────────────────────────────
  {
    name: 'Pi-hole',
    url: 'http://pi.local/admin',
    group: 'Networking',
    description: 'DNS ad blocker',
    icon: 'shield',
    color: '#ef4444',
  },
  {
    name: 'Nginx Proxy Manager',
    url: 'http://pi.local:81',
    group: 'Networking',
    description: 'Reverse proxy',
    icon: 'globe',
    color: '#3b82f6',
  },

  // ── Storage ────────────────────────────────────────────────────────────────
  {
    name: 'Nextcloud',
    url: 'http://pi.local:8080',
    group: 'Storage',
    description: 'File sync & share',
    icon: 'cloud',
    color: '#0ea5e9',
  },

  // ── Monitoring ─────────────────────────────────────────────────────────────
  {
    name: 'Grafana',
    url: 'http://pi.local:3000',
    group: 'Monitoring',
    description: 'Metrics & dashboards',
    icon: 'activity',
    color: '#f97316',
  },
  {
    name: 'Portainer',
    url: 'http://pi.local:9000',
    group: 'Monitoring',
    description: 'Docker management',
    icon: 'box',
    color: '#2dd4bf',
  },
  {
    name: 'Uptime Kuma',
    url: 'http://pi.local:3001',
    group: 'Monitoring',
    description: 'Service uptime monitor',
    icon: 'monitor',
    color: '#3fb950',
  },

  // ── Automation ─────────────────────────────────────────────────────────────
  {
    name: 'Home Assistant',
    url: 'http://homeassistant.local:8123',
    group: 'Automation',
    description: 'Home automation',
    icon: 'home',
    color: '#d29922',
  },

  // ── Development ────────────────────────────────────────────────────────────
  {
    name: 'Gitea',
    url: 'http://pi.local:3030',
    group: 'Development',
    description: 'Git repository',
    icon: 'code',
    color: '#6366f1',
  },

]
