let dashSSE = null;
let dashAbortController = null;
let dashRetryTimer = null;

function renderDashboard() {
  const page = document.getElementById('page-dashboard');
  page.innerHTML = `
    <header class="page-header fade-in">
      <h1 class="section-title">仪表盘</h1>
      <p class="section-sub">反代运行摘要 <span class="live-badge is-retry" id="sse-status">连接中</span></p>
    </header>
    <div class="stats-row" id="dash-stats">
      <div class="stat-card c-blue">
        <div class="stat-icon-wrap blue">
          <svg viewBox="0 0 24 24" aria-hidden="true"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
        </div>
        <div class="stat-number" id="s-total">—</div>
        <div class="stat-title">站点总数</div>
      </div>
      <div class="stat-card c-green">
        <div class="stat-icon-wrap green">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        </div>
        <div class="stat-number" id="s-running">—</div>
        <div class="stat-title">运行中</div>
      </div>
      <div class="stat-card c-teal">
        <div class="stat-icon-wrap teal">
          <svg viewBox="0 0 24 24" aria-hidden="true"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
        </div>
        <div class="stat-number" id="s-traffic">0 B</div>
        <div class="stat-title">总流量</div>
      </div>
      <div class="stat-card c-orange">
        <div class="stat-icon-wrap orange">
          <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
        </div>
        <div class="stat-number" id="s-uptime">—</div>
        <div class="stat-title">运行时长</div>
      </div>
    </div>
    <div class="glass-card fade-in">
      <div class="glass-card-header">
        <div class="glass-card-title"><span class="live-dot" id="table-live-dot"></span>站点实时状态</div>
        <div class="glass-card-title" style="font-size:.72rem;color:var(--ink-muted)" id="s-requests">0 请求</div>
      </div>
      <div class="table-wrap">
        <table>
          <thead><tr>
            <th>站点</th><th>状态</th><th>主回源</th><th>播放</th><th>UA</th><th>端口</th><th>已用流量</th>
          </tr></thead>
          <tbody id="dash-table">${UI.skeletonRows(3, 7)}</tbody>
        </table>
      </div>
    </div>
  `;

  startDashSSE();
  loadDashboardTable();
}

function startDashSSE() {
  stopDashSSE();
  startFetchSSE();
}

function setSseStatus(mode) {
  const statusEl = document.getElementById('sse-status');
  const dot = document.getElementById('table-live-dot');
  if (!statusEl) return;
  statusEl.classList.remove('is-live', 'is-retry', 'is-down');
  if (dot) dot.classList.toggle('is-live', mode === 'live');
  if (mode === 'live') {
    statusEl.classList.add('is-live');
    statusEl.textContent = '实时';
  } else if (mode === 'retry') {
    statusEl.classList.add('is-retry');
    statusEl.textContent = '重连中';
  } else {
    statusEl.classList.add('is-down');
    statusEl.textContent = '已断开';
  }
}

function queueDashSSERetry() {
  if (dashRetryTimer) clearTimeout(dashRetryTimer);
  setSseStatus('retry');
  dashRetryTimer = setTimeout(() => {
    if (Router.current === 'dashboard' && API.authenticated) startFetchSSE();
  }, 5000);
}

async function startFetchSSE() {
  const controller = new AbortController();
  dashAbortController = controller;
  setSseStatus('retry');

  try {
    const resp = await fetch('/api/events', {
      credentials: 'same-origin',
      signal: controller.signal,
    });

    if (!resp.ok) throw new Error('SSE failed');
    if (dashAbortController !== controller) return;

    setSseStatus('live');

    const reader = resp.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done || controller.signal.aborted) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop();

      for (const line of lines) {
        if (!line.startsWith('data: ')) continue;
        try {
          updateDashboardLive(JSON.parse(line.slice(6)));
        } catch (e) {
          // Skip malformed chunks and keep stream alive.
        }
      }
    }

    if (!controller.signal.aborted && dashAbortController === controller && Router.current === 'dashboard') {
      setSseStatus('down');
      queueDashSSERetry();
    }
  } catch (e) {
    if (controller.signal.aborted || dashAbortController !== controller) return;
    setSseStatus('down');
    queueDashSSERetry();
  }
}

function updateDashboardLive(stats) {
  animateValue('s-total', stats.total_sites || 0);
  animateValue('s-running', stats.running_sites || 0);

  const trafficEl = document.getElementById('s-traffic');
  if (trafficEl) trafficEl.textContent = formatBytes(stats.total_traffic || 0);

  const uptimeEl = document.getElementById('s-uptime');
  if (uptimeEl) uptimeEl.textContent = formatUptime(stats.uptime_seconds || 0);

  const requestsEl = document.getElementById('s-requests');
  if (requestsEl) requestsEl.textContent = formatNumber(stats.total_requests || 0) + ' 请求';
}

function formatUptime(seconds) {
  if (seconds < 60) return seconds + ' 秒';
  if (seconds < 3600) return Math.floor(seconds / 60) + ' 分';
  if (seconds < 86400) return Math.floor(seconds / 3600) + ' 时 ' + Math.floor((seconds % 3600) / 60) + ' 分';
  return Math.floor(seconds / 86400) + ' 天 ' + Math.floor((seconds % 86400) / 3600) + ' 时';
}

function formatNumber(n) {
  return n.toLocaleString();
}

function animateValue(id, newVal) {
  const el = document.getElementById(id);
  if (!el) return;
  const current = parseInt(el.textContent, 10) || 0;
  if (current === newVal) {
    el.textContent = newVal;
    return;
  }
  el.textContent = newVal;
  el.style.transition = 'transform 150ms var(--ease)';
  el.style.transform = 'scale(1.04)';
  setTimeout(() => { el.style.transform = ''; }, 150);
}

function stopDashSSE() {
  if (dashRetryTimer) {
    clearTimeout(dashRetryTimer);
    dashRetryTimer = null;
  }
  if (dashAbortController) {
    dashAbortController.abort();
    dashAbortController = null;
  }
  if (dashSSE) {
    dashSSE.close();
    dashSSE = null;
  }
}

function sitePlaybackLabel(s) {
  const playback = (s.playback_target_url || '').trim();
  let extra = [];
  try { extra = JSON.parse(s.stream_hosts || '[]'); } catch (e) {}
  const total = (playback ? 1 : 0) + extra.length;
  if (total === 0) return '<span class="pill pill-muted">跟随主回源</span>';
  if (total === 1 && playback === (s.target_url || '').trim()) {
    return '<span class="pill pill-muted">与主回源相同</span>';
  }
  const mode = s.playback_mode === 'redirect' ? '重定向跟随' : '直连分流';
  return `<span class="pill pill-blue">${mode}</span>`;
}

async function loadDashboardTable() {
  const tbody = document.getElementById('dash-table');
  if (!tbody) return;

  try {
    const sites = await API.listSites();
    if (!sites || sites.length === 0) {
      tbody.innerHTML = `<tr><td colspan="7">${UI.empty({
        inline: true,
        title: '还没有站点',
        body: '添加一个反代后，这里会显示运行状态、回源和流量。',
        actions: [{ id: 'goto-sites', label: '前往站点管理', className: 'btn-primary' }],
      })}</td></tr>`;
      tbody.querySelector('[data-empty-action="goto-sites"]')?.addEventListener('click', () => Router.navigate('sites'));
      return;
    }

    tbody.innerHTML = sites.map(s => `
      <tr>
        <td style="font-weight:600">${esc(s.name)}</td>
        <td><span class="status-badge"><span class="status-led ${s.running ? 'on' : 'off'}"></span>${s.running ? '运行中' : '已停止'}</span></td>
        <td class="mono">${esc(s.target_url)}</td>
        <td>${sitePlaybackLabel(s)}</td>
        <td><span class="pill ${uaClassMap[s.ua_mode] || 'pill-blue'}">${uaNameMap[s.ua_mode] || s.ua_mode}</span></td>
        <td class="mono">:${s.listen_port}</td>
        <td>${formatBytes(s.traffic_used)}</td>
      </tr>
    `).join('');
  } catch (e) {
    tbody.innerHTML = `<tr><td colspan="7">${UI.error({
      inline: true,
      body: e.message,
      retry: true,
    })}</td></tr>`;
    tbody.querySelector('[data-error-retry]')?.addEventListener('click', loadDashboardTable);
  }
}

async function loadDashboardData() {
  loadDashboardTable();
}
