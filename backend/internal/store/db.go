package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

type SystemStat struct {
	Ts          int64   `json:"ts"`
	CPUPercent  float64 `json:"cpu_percent"`
	RAMPercent  float64 `json:"ram_percent"`
	RAMUsedMB   int64   `json:"ram_used_mb"`
	DiskPercent float64 `json:"disk_percent"`
	TempC       float64 `json:"temp_c"`
	NetRxBytesS float64 `json:"net_rx_bytes_s"`
	NetTxBytesS float64 `json:"net_tx_bytes_s"`
}

type ContainerStat struct {
	Ts          int64   `json:"ts"`
	Name        string  `json:"name,omitempty"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemPercent  float64 `json:"mem_percent"`
	MemUsedMB   float64 `json:"mem_used_mb"`
	NetRxBytesS float64 `json:"net_rx_bytes_s"`
	NetTxBytesS float64 `json:"net_tx_bytes_s"`
}

type ContainerStatInput struct {
	Name        string
	CPUPercent  float64
	MemPercent  float64
	MemUsedMB   float64
	NetRxBytesS float64
	NetTxBytesS float64
}

type HistoryResult struct {
	System     []SystemStat              `json:"system"`
	Containers map[string][]ContainerStat `json:"containers"`
}

func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("WAL mode: %w", err)
	}
	if _, err := db.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		return nil, fmt.Errorf("synchronous: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &DB{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS system_stats (
			ts             INTEGER PRIMARY KEY,
			cpu_percent    REAL,
			ram_percent    REAL,
			ram_used_mb    INTEGER,
			disk_percent   REAL,
			temp_c         REAL,
			net_rx_bytes_s REAL,
			net_tx_bytes_s REAL
		);
		CREATE TABLE IF NOT EXISTS container_stats (
			ts             INTEGER,
			name           TEXT,
			cpu_percent    REAL,
			mem_percent    REAL,
			mem_used_mb    REAL,
			net_rx_bytes_s REAL,
			net_tx_bytes_s REAL,
			PRIMARY KEY (ts, name)
		);
		CREATE INDEX IF NOT EXISTS idx_system_ts    ON system_stats(ts);
		CREATE INDEX IF NOT EXISTS idx_container_ts ON container_stats(ts);
	`)
	return err
}

func (d *DB) WriteSystemStat(ts int64, cpu, ram float64, ramMB int64, disk, temp, rx, tx float64) {
	_, err := d.db.Exec(
		`INSERT OR REPLACE INTO system_stats VALUES (?,?,?,?,?,?,?,?)`,
		ts, cpu, ram, ramMB, disk, temp, rx, tx,
	)
	if err != nil {
		log.Printf("store: write system stat: %v", err)
	}
}

func (d *DB) WriteContainerStats(ts int64, containers []ContainerStatInput) {
	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("store: begin tx: %v", err)
		return
	}
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO container_stats VALUES (?,?,?,?,?,?,?)`)
	if err != nil {
		tx.Rollback()
		log.Printf("store: prepare: %v", err)
		return
	}
	defer stmt.Close()
	for _, c := range containers {
		if _, err := stmt.Exec(ts, c.Name, c.CPUPercent, c.MemPercent, c.MemUsedMB, c.NetRxBytesS, c.NetTxBytesS); err != nil {
			log.Printf("store: write container %s: %v", c.Name, err)
		}
	}
	tx.Commit()
}

func (d *DB) Prune(retentionDays int) {
	cutoff := time.Now().Unix() - int64(retentionDays)*86400
	if _, err := d.db.Exec("DELETE FROM system_stats WHERE ts < ?", cutoff); err != nil {
		log.Printf("store: prune system_stats: %v", err)
	}
	if _, err := d.db.Exec("DELETE FROM container_stats WHERE ts < ?", cutoff); err != nil {
		log.Printf("store: prune container_stats: %v", err)
	}
}

// minResolution enforces a floor on resolution to avoid returning millions of rows.
func minResolution(windowSec int64) int64 {
	switch {
	case windowSec <= 6*3600:
		return 10
	case windowSec <= 24*3600:
		return 30
	case windowSec <= 7*86400:
		return 300
	default:
		return 900
	}
}

func (d *DB) QueryRange(from, to, resolution int64) (*HistoryResult, error) {
	window := to - from
	if floor := minResolution(window); resolution < floor {
		resolution = floor
	}

	result := &HistoryResult{
		System:     []SystemStat{},
		Containers: make(map[string][]ContainerStat),
	}

	rows, err := d.db.Query(`
		SELECT
			(ts / ?) * ? AS bucket,
			AVG(cpu_percent), AVG(ram_percent), CAST(AVG(ram_used_mb) AS INTEGER),
			AVG(disk_percent), AVG(temp_c), AVG(net_rx_bytes_s), AVG(net_tx_bytes_s)
		FROM system_stats
		WHERE ts >= ? AND ts <= ?
		GROUP BY bucket ORDER BY bucket`,
		resolution, resolution, from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("query system: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var s SystemStat
		if err := rows.Scan(&s.Ts, &s.CPUPercent, &s.RAMPercent, &s.RAMUsedMB,
			&s.DiskPercent, &s.TempC, &s.NetRxBytesS, &s.NetTxBytesS); err != nil {
			continue
		}
		result.System = append(result.System, s)
	}

	crows, err := d.db.Query(`
		SELECT
			name,
			(ts / ?) * ? AS bucket,
			AVG(cpu_percent), AVG(mem_percent), AVG(mem_used_mb),
			AVG(net_rx_bytes_s), AVG(net_tx_bytes_s)
		FROM container_stats
		WHERE ts >= ? AND ts <= ?
		GROUP BY name, bucket ORDER BY name, bucket`,
		resolution, resolution, from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("query containers: %w", err)
	}
	defer crows.Close()
	for crows.Next() {
		var c ContainerStat
		if err := crows.Scan(&c.Name, &c.Ts, &c.CPUPercent, &c.MemPercent, &c.MemUsedMB,
			&c.NetRxBytesS, &c.NetTxBytesS); err != nil {
			continue
		}
		result.Containers[c.Name] = append(result.Containers[c.Name], c)
	}

	return result, nil
}

func (d *DB) Close() {
	d.db.Close()
}
