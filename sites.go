package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Site struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	ListenPort        int    `json:"listen_port"`
	TargetURL         string `json:"target_url"`
	PlaybackTargetURL string `json:"playback_target_url"`
	PlaybackMode      string `json:"playback_mode"`
	StreamHosts       string `json:"stream_hosts"`
	UAMode            string `json:"ua_mode"`
	Enabled           bool   `json:"enabled"`
	TrafficQuota      int64  `json:"traffic_quota"`
	TrafficUsed       int64  `json:"traffic_used"`
	SpeedLimit        int    `json:"speed_limit"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func (d *DB) ListSites() ([]Site, error) {
	rows, err := d.db.Query("SELECT id, name, listen_port, target_url, playback_target_url, playback_mode, stream_hosts, ua_mode, enabled, traffic_quota, traffic_used, speed_limit, created_at, updated_at FROM sites ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sites []Site
	for rows.Next() {
		var s Site
		var enabled int
		rows.Scan(&s.ID, &s.Name, &s.ListenPort, &s.TargetURL, &s.PlaybackTargetURL, &s.PlaybackMode, &s.StreamHosts, &s.UAMode, &enabled, &s.TrafficQuota, &s.TrafficUsed, &s.SpeedLimit, &s.CreatedAt, &s.UpdatedAt)
		s.Enabled = enabled == 1
		sites = append(sites, s)
	}
	if sites == nil {
		sites = []Site{}
	}
	return sites, nil
}

func (d *DB) GetSite(id int64) (*Site, error) {
	var s Site
	var enabled int
	err := d.db.QueryRow("SELECT id, name, listen_port, target_url, playback_target_url, playback_mode, stream_hosts, ua_mode, enabled, traffic_quota, traffic_used, speed_limit, created_at, updated_at FROM sites WHERE id=?", id).
		Scan(&s.ID, &s.Name, &s.ListenPort, &s.TargetURL, &s.PlaybackTargetURL, &s.PlaybackMode, &s.StreamHosts, &s.UAMode, &enabled, &s.TrafficQuota, &s.TrafficUsed, &s.SpeedLimit, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	s.Enabled = enabled == 1
	return &s, nil
}

func (d *DB) CreateSite(name string, port int, targetURL, playbackTargetURL, playbackMode, streamHosts, uaMode string, quota int64, speedLimit int) (*Site, error) {
	if streamHosts == "" {
		streamHosts = "[]"
	}
	res, err := d.db.Exec(
		"INSERT INTO sites (name, listen_port, target_url, playback_target_url, playback_mode, stream_hosts, ua_mode, traffic_quota, speed_limit) VALUES (?,?,?,?,?,?,?,?,?)",
		name, port, targetURL, playbackTargetURL, playbackMode, streamHosts, uaMode, quota, speedLimit,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return d.GetSite(id)
}

func (d *DB) UpdateSite(id int64, name string, port int, targetURL, playbackTargetURL, playbackMode, streamHosts, uaMode string, quota int64, speedLimit int) error {
	if streamHosts == "" {
		streamHosts = "[]"
	}
	_, err := d.db.Exec(
		"UPDATE sites SET name=?, listen_port=?, target_url=?, playback_target_url=?, playback_mode=?, stream_hosts=?, ua_mode=?, traffic_quota=?, speed_limit=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		name, port, targetURL, playbackTargetURL, playbackMode, streamHosts, uaMode, quota, speedLimit, id,
	)
	return err
}

func (d *DB) DeleteSite(id int64) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec("DELETE FROM traffic_logs WHERE site_id=?", id); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM sites WHERE id=?", id); err != nil {
		return err
	}
	return tx.Commit()
}

func (d *DB) ToggleSite(id int64) (bool, error) {
	var enabled int
	d.db.QueryRow("SELECT enabled FROM sites WHERE id=?", id).Scan(&enabled)
	newVal := 1 - enabled
	_, err := d.db.Exec("UPDATE sites SET enabled=?, updated_at=CURRENT_TIMESTAMP WHERE id=?", newVal, id)
	return newVal == 1, err
}

type ExportSiteRecord struct {
	Name              string   `json:"name"`
	ListenPort        int      `json:"listen_port"`
	TargetURL         string   `json:"target_url"`
	PlaybackTargetURL string   `json:"playback_target_url"`
	PlaybackMode      string   `json:"playback_mode"`
	StreamHosts       []string `json:"stream_hosts"`
	UAMode            string   `json:"ua_mode"`
	TrafficQuota      int64    `json:"traffic_quota"`
	SpeedLimit        int      `json:"speed_limit"`
}

// GET /api/sites/export — download all sites as JSON
func (a *App) handleSitesExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.jsonErr(w, 405, "method not allowed")
		return
	}
	sites, err := a.db.ListSites()
	if err != nil {
		a.jsonErr(w, 500, err.Error())
		return
	}
	records := make([]ExportSiteRecord, 0, len(sites))
	for _, s := range sites {
		var streamHosts []string
		if s.StreamHosts != "" && s.StreamHosts != "[]" {
			_ = json.Unmarshal([]byte(s.StreamHosts), &streamHosts)
		}
		if streamHosts == nil {
			streamHosts = []string{}
		}
		records = append(records, ExportSiteRecord{
			Name:              s.Name,
			ListenPort:        s.ListenPort,
			TargetURL:         s.TargetURL,
			PlaybackTargetURL: s.PlaybackTargetURL,
			PlaybackMode:      s.PlaybackMode,
			StreamHosts:       streamHosts,
			UAMode:            s.UAMode,
			TrafficQuota:      s.TrafficQuota,
			SpeedLimit:        s.SpeedLimit,
		})
	}
	out := map[string]interface{}{
		"version": backupFormatVersion,
		"sites":   records,
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+backupDownloadName+"\"")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(out)
}

// POST /api/sites/import — restore sites from exported JSON
func (a *App) handleSitesImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		a.jsonErr(w, 405, "method not allowed")
		return
	}
	var payload struct {
		Version   string             `json:"version"`
		Overwrite bool               `json:"overwrite"`
		Sites     []ExportSiteRecord `json:"sites"`
	}
	if err := decodeJSONBody(w, r, &payload); err != nil {
		a.jsonErr(w, 400, "invalid JSON: "+err.Error())
		return
	}
	if !acceptedBackupVersion(payload.Version) {
		a.jsonErr(w, 400, "unsupported backup version")
		return
	}
	if len(payload.Sites) == 0 {
		a.jsonErr(w, 400, "no sites in import payload")
		return
	}

	created := 0
	skipped := 0
	skippedItems := make([]ImportSkip, 0)
	for _, rec := range payload.Sites {
		if rec.Name == "" || rec.TargetURL == "" || rec.ListenPort == 0 {
			skipped++
			skippedItems = append(skippedItems, ImportSkip{
				Name:       rec.Name,
				ListenPort: rec.ListenPort,
				Reason:     importSkipReason(rec, nil),
			})
			continue
		}
		if rec.UAMode == "" {
			rec.UAMode = "infuse"
		}
		if rec.PlaybackMode == "" {
			rec.PlaybackMode = "direct"
		}
		streamHostsJSON, _ := json.Marshal(rec.StreamHosts)
		site, err := a.db.CreateSite(
			rec.Name, rec.ListenPort, rec.TargetURL, rec.PlaybackTargetURL,
			rec.PlaybackMode, string(streamHostsJSON), rec.UAMode,
			rec.TrafficQuota, rec.SpeedLimit,
		)
		if err != nil {
			skipped++
			skippedItems = append(skippedItems, ImportSkip{
				Name:       rec.Name,
				ListenPort: rec.ListenPort,
				Reason:     importSkipReason(rec, err),
			})
			continue
		}
		if site.Enabled {
			_ = a.pm.StartSite(*site)
		}
		created++
	}
	a.jsonOK(w, map[string]interface{}{
		"created":       created,
		"skipped":       skipped,
		"skipped_items": skippedItems,
	})
}

type ImportSkip struct {
	Name       string `json:"name"`
	ListenPort int    `json:"listen_port"`
	Reason     string `json:"reason"`
}

func importSkipReason(rec ExportSiteRecord, err error) string {
	var missing []string
	if rec.Name == "" {
		missing = append(missing, "名称")
	}
	if rec.ListenPort == 0 {
		missing = append(missing, "端口")
	}
	if rec.TargetURL == "" {
		missing = append(missing, "回源 URL")
	}
	if len(missing) > 0 {
		return "缺少" + strings.Join(missing, "、")
	}
	if err == nil {
		return "已跳过"
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "unique") && strings.Contains(msg, "listen_port") {
		return "监听端口已被占用"
	}
	if strings.Contains(lower, "unique") {
		return "与现有站点冲突"
	}
	return msg
}

// GET/POST /api/sites
func (a *App) handleSites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sites, err := a.db.ListSites()
		if err != nil {
			a.jsonErr(w, 500, err.Error())
			return
		}
		// Add running status
		type SiteWithStatus struct {
			Site
			Running bool `json:"running"`
		}
		result := make([]SiteWithStatus, len(sites))
		for i, s := range sites {
			result[i] = SiteWithStatus{Site: s, Running: a.pm.IsRunning(s.ID)}
		}
		a.jsonOK(w, result)

	case "POST":
		var req struct {
			Name              string   `json:"name"`
			ListenPort        int      `json:"listen_port"`
			TargetURL         string   `json:"target_url"`
			PlaybackTargetURL string   `json:"playback_target_url"`
			PlaybackMode      string   `json:"playback_mode"`
			StreamHosts       []string `json:"stream_hosts"`
			UAMode            string   `json:"ua_mode"`
			Quota             int64    `json:"traffic_quota"`
			SpeedLimit        int      `json:"speed_limit"`
		}
		if err := decodeJSONBody(w, r, &req); err != nil {
			a.jsonErr(w, 400, "invalid request")
			return
		}
		if req.Name == "" || req.ListenPort == 0 || req.TargetURL == "" {
			a.jsonErr(w, 400, "name, listen_port, and target_url are required")
			return
		}
		if req.UAMode == "" {
			req.UAMode = "infuse"
		}
		if req.PlaybackMode == "" {
			req.PlaybackMode = "direct"
		}
		streamHostsJSON, _ := json.Marshal(req.StreamHosts)
		if req.StreamHosts == nil {
			streamHostsJSON = []byte("[]")
		}
		site, err := a.db.CreateSite(req.Name, req.ListenPort, req.TargetURL, req.PlaybackTargetURL, req.PlaybackMode, string(streamHostsJSON), req.UAMode, req.Quota, req.SpeedLimit)
		if err != nil {
			a.jsonErr(w, 500, err.Error())
			return
		}
		// Auto start
		if site.Enabled {
			if err := a.pm.StartSite(*site); err != nil {
				if deleteErr := a.db.DeleteSite(site.ID); deleteErr != nil {
					a.jsonErr(w, 500, fmt.Sprintf("start site: %v; rollback create: %v", err, deleteErr))
					return
				}
				a.jsonErr(w, 500, err.Error())
				return
			}
		}
		w.WriteHeader(201)
		a.jsonOK(w, site)

	default:
		a.jsonErr(w, 405, "method not allowed")
	}
}

// PUT/DELETE /api/sites/{id}, POST /api/sites/{id}/toggle, GET /api/sites/{id}/diag
func (a *App) handleSiteByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sites/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		a.jsonErr(w, 400, "invalid site id")
		return
	}

	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch {
	case action == "toggle" && r.Method == "POST":
		newState, err := a.db.ToggleSite(id)
		if err != nil {
			a.jsonErr(w, 500, err.Error())
			return
		}
		if newState {
			site, err := a.db.GetSite(id)
			if err != nil {
				if _, revertErr := a.db.ToggleSite(id); revertErr != nil {
					a.jsonErr(w, 500, fmt.Sprintf("load site: %v; rollback toggle: %v", err, revertErr))
					return
				}
				a.jsonErr(w, 500, err.Error())
				return
			}
			if err := a.pm.StartSite(*site); err != nil {
				if _, revertErr := a.db.ToggleSite(id); revertErr != nil {
					a.jsonErr(w, 500, fmt.Sprintf("start site: %v; rollback toggle: %v", err, revertErr))
					return
				}
				a.jsonErr(w, 500, err.Error())
				return
			}
		} else {
			a.pm.StopSite(id)
		}
		a.jsonOK(w, map[string]interface{}{"enabled": newState})

	case action == "diag" && r.Method == "GET":
		site, err := a.db.GetSite(id)
		if err != nil {
			a.jsonErr(w, 404, "site not found")
			return
		}
		result := diagnoseSite(site, a.pm)
		a.jsonOK(w, result)

	case action == "" && r.Method == "PUT":
		oldSite, err := a.db.GetSite(id)
		if err != nil {
			a.jsonErr(w, 404, "site not found")
			return
		}
		var req struct {
			Name              string    `json:"name"`
			ListenPort        int       `json:"listen_port"`
			TargetURL         string    `json:"target_url"`
			PlaybackTargetURL *string   `json:"playback_target_url"`
			PlaybackMode      *string   `json:"playback_mode"`
			StreamHosts       *[]string `json:"stream_hosts"`
			UAMode            string    `json:"ua_mode"`
			Quota             int64     `json:"traffic_quota"`
			SpeedLimit        int       `json:"speed_limit"`
		}
		if err := decodeJSONBody(w, r, &req); err != nil {
			a.jsonErr(w, 400, "invalid request")
			return
		}
		name := oldSite.Name
		if req.Name != "" {
			name = req.Name
		}
		listenPort := oldSite.ListenPort
		if req.ListenPort != 0 {
			listenPort = req.ListenPort
		}
		targetURL := oldSite.TargetURL
		if req.TargetURL != "" {
			targetURL = req.TargetURL
		}
		if name == "" || listenPort == 0 || targetURL == "" {
			a.jsonErr(w, 400, "name, listen_port, and target_url are required")
			return
		}
		playbackTargetURL := oldSite.PlaybackTargetURL
		if req.PlaybackTargetURL != nil {
			playbackTargetURL = *req.PlaybackTargetURL
		}
		playbackMode := oldSite.PlaybackMode
		if req.PlaybackMode != nil {
			playbackMode = *req.PlaybackMode
		}
		streamHosts := oldSite.StreamHosts
		if req.StreamHosts != nil {
			sh, _ := json.Marshal(*req.StreamHosts)
			streamHosts = string(sh)
		}
		if req.UAMode == "" {
			req.UAMode = oldSite.UAMode
		}
		if err := a.db.UpdateSite(id, name, listenPort, targetURL, playbackTargetURL, playbackMode, streamHosts, req.UAMode, req.Quota, req.SpeedLimit); err != nil {
			a.jsonErr(w, 500, err.Error())
			return
		}
		site, err := a.db.GetSite(id)
		if err != nil {
			a.jsonErr(w, 500, err.Error())
			return
		}
		if site.Enabled {
			needsPreStop := oldSite.Enabled && oldSite.ListenPort == site.ListenPort
			if needsPreStop {
				a.pm.StopSite(id)
			}
			if err := a.pm.StartSite(*site); err != nil {
				if rollbackErr := a.db.UpdateSite(oldSite.ID, oldSite.Name, oldSite.ListenPort, oldSite.TargetURL, oldSite.PlaybackTargetURL, oldSite.PlaybackMode, oldSite.StreamHosts, oldSite.UAMode, oldSite.TrafficQuota, oldSite.SpeedLimit); rollbackErr != nil {
					a.jsonErr(w, 500, fmt.Sprintf("start updated site: %v; rollback update: %v", err, rollbackErr))
					return
				}
				restoredSite, getErr := a.db.GetSite(id)
				if getErr != nil {
					a.jsonErr(w, 500, fmt.Sprintf("start updated site: %v; reload rollback site: %v", err, getErr))
					return
				}
				if oldSite.Enabled && !a.pm.IsRunning(id) {
					if restartErr := a.pm.StartSite(*restoredSite); restartErr != nil {
						a.jsonErr(w, 500, fmt.Sprintf("start updated site: %v; restore previous site: %v", err, restartErr))
						return
					}
				}
				a.jsonErr(w, 500, err.Error())
				return
			}
		}
		a.jsonOK(w, site)

	case action == "" && r.Method == "DELETE":
		a.pm.StopSite(id)
		if err := a.db.DeleteSite(id); err != nil {
			a.jsonErr(w, 500, err.Error())
			return
		}
		a.jsonOK(w, map[string]string{"status": "deleted"})

	default:
		a.jsonErr(w, 405, "method not allowed")
	}
}

func acceptedBackupVersion(v string) bool {
	switch strings.TrimSpace(v) {
	case "", backupFormatVersion, legacyBackupVersion:
		return true
	default:
		return false
	}
}
