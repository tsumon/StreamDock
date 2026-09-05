package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type UAProfile struct {
	Name      string `json:"name"`
	UserAgent string `json:"user_agent"`
	Client    string `json:"client"`
	Version   string `json:"version"`
}

var uaProfiles = map[string]UAProfile{
	"infuse": {Name: "Infuse", UserAgent: "Infuse/7.8.1", Client: "Infuse", Version: "7.8.1"},
	"web":    {Name: "Web", UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Emby Theater", Client: "Emby Web", Version: "4.9.0.42"},
	"client": {Name: "Client", UserAgent: "Emby-Theater/4.7.0", Client: "Emby Theater", Version: "4.7.0"},
}

func getUAProfile(mode string) UAProfile {
	if p, ok := uaProfiles[strings.ToLower(mode)]; ok {
		return p
	}
	return uaProfiles["infuse"]
}

type redirectFollowTransport struct {
	base          http.RoundTripper
	playbackHosts map[string]bool
	profile       UAProfile
}

func (t *redirectFollowTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 3; i++ {
		if resp.StatusCode != 301 && resp.StatusCode != 302 && resp.StatusCode != 307 && resp.StatusCode != 308 {
			break
		}
		loc := resp.Header.Get("Location")
		if loc == "" {
			break
		}
		locURL, err := url.Parse(loc)
		if err != nil {
			break
		}
		if locURL.Host == "" {
			locURL = req.URL.ResolveReference(locURL)
		}
		if !t.playbackHosts[strings.ToLower(locURL.Host)] {
			break
		}
		resp.Body.Close()
		newReq, err := http.NewRequestWithContext(req.Context(), req.Method, locURL.String(), nil)
		if err != nil {
			break
		}
		for k, v := range req.Header {
			newReq.Header[k] = v
		}
		newReq.Host = locURL.Host
		applyUAProfileHeaders(newReq.Header, t.profile)
		resp, err = t.base.RoundTrip(newReq)
		if err != nil {
			return nil, err
		}
		req = newReq
	}
	return resp, nil
}

var embyAuthClientRe = regexp.MustCompile(`(?i)(Client=")[^"]*"`)
var embyAuthVersionRe = regexp.MustCompile(`(?i)(Version=")[^"]*"`)

type ProxyInstance struct {
	Site             Site
	server           *http.Server
	listener         net.Listener
	bytesIn          atomic.Int64
	bytesOut         atomic.Int64
	reqCount         atomic.Int64
	persistedTraffic atomic.Int64
}

type ProxyManager struct {
	mu       sync.RWMutex
	proxies  map[int64]*ProxyInstance
	database *DB
}

func NewProxyManager(db *DB) *ProxyManager {
	return &ProxyManager{
		proxies:  make(map[int64]*ProxyInstance),
		database: db,
	}
}

// metered response writer
type meteredWriter struct {
	http.ResponseWriter
	written *atomic.Int64
}

func (m *meteredWriter) Write(b []byte) (int, error) {
	n, err := m.ResponseWriter.Write(b)
	m.written.Add(int64(n))
	return n, err
}

// Flush support for streaming
func (m *meteredWriter) Flush() {
	if f, ok := m.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack support for WebSocket upgrade
func (m *meteredWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := m.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijack not supported")
}

// metered request body reader
type meteredReader struct {
	io.ReadCloser
	read *atomic.Int64
}

func (m *meteredReader) Read(p []byte) (int, error) {
	n, err := m.ReadCloser.Read(p)
	m.read.Add(int64(n))
	return n, err
}

// 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?Rate-limited writer (token bucket) 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?
type rateLimitedWriter struct {
	http.ResponseWriter
	bytesPerSec int64
	written     *atomic.Int64
	start       time.Time
}

func (w *rateLimitedWriter) Write(b []byte) (int, error) {
	if w.bytesPerSec <= 0 {
		n, err := w.ResponseWriter.Write(b)
		w.written.Add(int64(n))
		return n, err
	}
	totalWritten := 0
	for len(b) > 0 {
		elapsed := time.Since(w.start).Seconds()
		if elapsed < 0.001 {
			elapsed = 0.001
		}
		allowed := int64(elapsed*float64(w.bytesPerSec)) - w.written.Load()
		if allowed <= 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		chunk := b
		if int64(len(chunk)) > allowed {
			chunk = b[:allowed]
		}
		n, err := w.ResponseWriter.Write(chunk)
		w.written.Add(int64(n))
		totalWritten += n
		b = b[n:]
		if err != nil {
			return totalWritten, err
		}
	}
	return totalWritten, nil
}

func (w *rateLimitedWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *rateLimitedWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijack not supported")
}

// 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?WebSocket tunnel 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?
func isWebSocketUpgrade(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}

func normalizeTargetURL(addr string) (*url.URL, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, fmt.Errorf("target URL is required")
	}
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}
	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid target URL")
	}
	return parsed, nil
}

func isPlaybackRequest(path string) bool {
	path = strings.ToLower(path)
	switch {
	case strings.HasPrefix(path, "/videos/"),
		strings.HasPrefix(path, "/emby/videos/"),
		strings.HasPrefix(path, "/audio/"),
		strings.HasPrefix(path, "/emby/audio/"),
		strings.HasPrefix(path, "/livetv/"),
		strings.HasPrefix(path, "/emby/livetv/"):
		return true
	case strings.HasPrefix(path, "/items/"),
		strings.HasPrefix(path, "/emby/items/"):
		return strings.Contains(path, "/download") || strings.Contains(path, "/file")
	default:
		return false
	}
}

func upstreamTargetForRequest(r *http.Request, apiTarget, playbackTarget *url.URL) *url.URL {
	if playbackTarget != nil && isPlaybackRequest(r.URL.Path) {
		return playbackTarget
	}
	return apiTarget
}

func applyUAProfileHeaders(header http.Header, profile UAProfile) {
	header.Set("User-Agent", profile.UserAgent)
	if auth := header.Get("X-Emby-Authorization"); auth != "" {
		if embyAuthClientRe.MatchString(auth) {
			auth = embyAuthClientRe.ReplaceAllString(auth, `${1}`+profile.Client+`"`)
		}
		if embyAuthVersionRe.MatchString(auth) {
			auth = embyAuthVersionRe.ReplaceAllString(auth, `${1}`+profile.Version+`"`)
		}
		header.Set("X-Emby-Authorization", auth)
	}
	if auth := header.Get("Authorization"); auth != "" {
		if embyAuthClientRe.MatchString(auth) {
			auth = embyAuthClientRe.ReplaceAllString(auth, `${1}`+profile.Client+`"`)
		}
		if embyAuthVersionRe.MatchString(auth) {
			auth = embyAuthVersionRe.ReplaceAllString(auth, `${1}`+profile.Version+`"`)
		}
		header.Set("Authorization", auth)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, target *url.URL, profile UAProfile, inst *ProxyInstance) {
	// Build target WS URL
	scheme := "ws"
	if target.Scheme == "https" {
		scheme = "wss"
	}
	targetURL := scheme + "://" + target.Host + r.URL.RequestURI()

	// Hijack client connection
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "WebSocket not supported", 500)
		return
	}
	clientConn, clientBuf, err := hj.Hijack()
	if err != nil {
		log.Printf("[WS] hijack error: %v", err)
		return
	}
	defer clientConn.Close()

	// Connect to upstream
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	var upstreamConn net.Conn
	host := target.Host
	if !strings.Contains(host, ":") {
		if scheme == "wss" {
			host += ":443"
		} else {
			host += ":80"
		}
	}
	if scheme == "wss" {
		serverName := host
		if h, _, splitErr := net.SplitHostPort(host); splitErr == nil {
			serverName = h
		}
		upstreamConn, err = tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: serverName,
		})
	} else {
		upstreamConn, err = dialer.Dial("tcp", host)
	}
	if err != nil {
		log.Printf("[WS] upstream dial error: %v", err)
		clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer upstreamConn.Close()

	// Send upgrade request to upstream
	reqLine := fmt.Sprintf("%s %s HTTP/1.1\r\n", r.Method, r.URL.RequestURI())
	upstreamConn.Write([]byte(reqLine))
	r.Header.Set("Host", target.Host)
	applyUAProfileHeaders(r.Header, profile)
	r.Header.Write(upstreamConn)
	upstreamConn.Write([]byte("\r\n"))

	_ = targetURL
	log.Printf("[WS] tunnel established: client <-> %s", target.Host)

	// Bidirectional copy
	done := make(chan struct{}, 2)
	go func() {
		n, _ := io.Copy(upstreamConn, clientBuf)
		inst.bytesIn.Add(n)
		done <- struct{}{}
	}()
	go func() {
		n, _ := io.Copy(clientConn, upstreamConn)
		inst.bytesOut.Add(n)
		done <- struct{}{}
	}()
	<-done
}

func (pm *ProxyManager) StartSite(site Site) error {
	target, err := normalizeTargetURL(site.TargetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL: %w", err)
	}
	var playbackTarget *url.URL
	if strings.TrimSpace(site.PlaybackTargetURL) != "" {
		playbackTarget, err = normalizeTargetURL(site.PlaybackTargetURL)
		if err != nil {
			return fmt.Errorf("invalid playback target URL: %w", err)
		}
	}

	// Build playback hosts set from playback_target_url + stream_hosts
	playbackHostsSet := make(map[string]bool)
	if playbackTarget != nil {
		playbackHostsSet[strings.ToLower(playbackTarget.Host)] = true
	}
	var extraHosts []string
	if site.StreamHosts != "" && site.StreamHosts != "[]" {
		json.Unmarshal([]byte(site.StreamHosts), &extraHosts)
	}
	for _, raw := range extraHosts {
		if parsed, e := normalizeTargetURL(raw); e == nil {
			playbackHostsSet[strings.ToLower(parsed.Host)] = true
			if playbackTarget == nil {
				playbackTarget = parsed
			}
		}
	}

	profile := getUAProfile(site.UAMode)
	inst := &ProxyInstance{Site: site}
	inst.persistedTraffic.Store(site.TrafficUsed)

	isRedirectMode := playbackTarget != nil && site.PlaybackMode == "redirect"

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			var upstream *url.URL
			if isRedirectMode {
				upstream = target
			} else {
				upstream = upstreamTargetForRequest(req, target, playbackTarget)
			}
			req.URL.Scheme = upstream.Scheme
			req.URL.Host = upstream.Host
			req.Host = upstream.Host
			applyUAProfileHeaders(req.Header, profile)
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[%s] proxy error: %v", site.Name, err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error":"upstream unavailable"}`))
		},
	}

	if isRedirectMode {
		proxy.Transport = &redirectFollowTransport{
			base:          http.DefaultTransport,
			playbackHosts: playbackHostsSet,
			profile:       profile,
		}
	}

	// Speed limit in bytes/sec (field is in Mbps, 0 = unlimited)
	speedLimitBytes := int64(site.SpeedLimit) * 125000 // Mbps -> bytes/sec

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inst.reqCount.Add(1)

		// 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?Traffic quota check 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?
		if site.TrafficQuota > 0 {
			currentUsed := inst.persistedTraffic.Load() + inst.bytesIn.Load() + inst.bytesOut.Load()
			if currentUsed >= site.TrafficQuota {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"traffic quota exceeded"}`))
				return
			}
		}

		// 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?WebSocket upgrade 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?
		if isWebSocketUpgrade(r) {
			wsTarget := upstreamTargetForRequest(r, target, playbackTarget)
			if isRedirectMode {
				wsTarget = target
			}
			handleWebSocket(w, r, wsTarget, profile, inst)
			return
		}

		// 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?Normal proxy with metering 闂傚倸鍊搁崐椋庣矆娓氣偓瀹曨垶宕稿Δ鈧崒銊︾節婵犲倻澧曠痪鎯ь煼閺岀喖宕滆鐢盯鏌ｉ幘鍐叉殻闁哄本绋栫粻娑㈠箼閸愨敩锔界箾?
		if r.Body != nil {
			r.Body = &meteredReader{ReadCloser: r.Body, read: &inst.bytesIn}
		}

		var rw http.ResponseWriter
		if speedLimitBytes > 0 {
			rw = &rateLimitedWriter{
				ResponseWriter: w,
				bytesPerSec:    speedLimitBytes,
				written:        &inst.bytesOut,
				start:          time.Now(),
			}
		} else {
			rw = &meteredWriter{ResponseWriter: w, written: &inst.bytesOut}
		}
		proxy.ServeHTTP(rw, r)
	})

	listenAddr := fmt.Sprintf(":%d", site.ListenPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", listenAddr, err)
	}

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	inst.server = server
	inst.listener = listener

	pm.mu.Lock()
	if existing, ok := pm.proxies[site.ID]; ok {
		if existing.server != nil {
			existing.server.Close()
		}
		delete(pm.proxies, site.ID)
	}
	pm.proxies[site.ID] = inst
	pm.mu.Unlock()

	go func() {
		if len(playbackHostsSet) > 0 {
			hosts := make([]string, 0, len(playbackHostsSet))
			for h := range playbackHostsSet {
				hosts = append(hosts, h)
			}
			log.Printf("[%s] proxy :%d -> %s (playback hosts: %s, mode: %s, UA: %s)", site.Name, site.ListenPort, site.TargetURL, strings.Join(hosts, ", "), site.PlaybackMode, site.UAMode)
		} else {
			log.Printf("[%s] proxy :%d -> %s (UA: %s)", site.Name, site.ListenPort, site.TargetURL, site.UAMode)
		}
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("[%s] server error: %v", site.Name, err)
		}
	}()

	return nil
}

func (pm *ProxyManager) StopSite(id int64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if inst, ok := pm.proxies[id]; ok {
		pm.flushProxyTraffic(inst)
		if inst.server != nil {
			inst.server.Close()
		}
		delete(pm.proxies, id)
	}
}

func (pm *ProxyManager) IsRunning(id int64) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	_, ok := pm.proxies[id]
	return ok
}

func (pm *ProxyManager) StartAllEnabled() {
	sites, _ := pm.database.ListSites()
	for _, s := range sites {
		if s.Enabled {
			if err := pm.StartSite(s); err != nil {
				log.Printf("[%s] failed to start: %v", s.Name, err)
			}
		}
	}
}

// Flush traffic counters to DB periodically
func (pm *ProxyManager) FlushTraffic() {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	for _, inst := range pm.proxies {
		pm.flushProxyTraffic(inst)
	}
}

func (pm *ProxyManager) flushProxyTraffic(inst *ProxyInstance) {
	in := inst.bytesIn.Swap(0)
	out := inst.bytesOut.Swap(0)
	if in == 0 && out == 0 {
		return
	}
	if err := pm.database.addTraffic(inst.Site.ID, in, out); err != nil {
		inst.bytesIn.Add(in)
		inst.bytesOut.Add(out)
		log.Printf("[%s] failed to flush traffic: %v", inst.Site.Name, err)
		return
	}
	delta := in + out
	inst.persistedTraffic.Add(delta)
	inst.Site.TrafficUsed += delta
}

func (pm *ProxyManager) GetRunningCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.proxies)
}

// GracefulShutdown stops all proxies gracefully
func (pm *ProxyManager) GracefulShutdown(ctx context.Context) {
	pm.FlushTraffic()
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for id, inst := range pm.proxies {
		log.Printf("[%s] shutting down...", inst.Site.Name)
		inst.server.Shutdown(ctx)
		delete(pm.proxies, id)
	}
}

// GetTotalRequests returns total request count across all proxies
func (pm *ProxyManager) GetTotalRequests() int64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	var total int64
	for _, inst := range pm.proxies {
		total += inst.reqCount.Load()
	}
	return total
}
