let dashSSE = null;
let dashAbortController = null;
let dashRetryTimer = null;

function renderDashboard() {
  const page = document.getElementById('page-dashboard');
  page.innerHTML = `
    <header class="page-header">
      <h1 class="section-title">仪表盘</h1>
      <p class="section-sub">反代运行摘要 <span class="live-badge is-retry" id="sse-status">连接中</span></p>
    </header>
    <dl class="spec-strip" id="dash-stats">
      <div class="spec-cell"><dt>站点</dt><dd id="s-total">—</dd></div>
      <div class="spec-cell"><dt>运行</dt><dd id="s-running">—</dd></div>
      <div class="spec-cell"><dt>总流量</dt><dd id="s-traffic">0 B</dd></div>
      <div class="spec-cell"><dt>运行时长</dt><dd id="s-uptime">—</dd></div>
    </dl>
    <div class="topo-list" id="dash-topo">${UI.skeletonCards(2)}</div>
  `;

  document.getElementById('dash-topo').addEventListener('click', onTopoClick);
  document.getElementById('dash-topo').addEventListener('keydown', onTopoKey);
  startDashSSE();
  loadDashboardTable();
}

function onTopoClick(e) {
  const node = e.target.closest('.topo-node');
  if (!node) return;
  traceTopoNode(node);
}

function onTopoKey(e) {
  if (e.key !== 'Enter' && e.key !== ' ') return;
  const node = e.target.closest('.topo-node');
  if (!node) return;
  e.preventDefault();
  traceTopoNode(node);
}

function traceTopoNode(node) {
  const list = document.getElementById('dash-topo');
  if (!list) return;
  list.querySelectorAll('.topo-node.is-traced').forEach((el) => {
    if (el !== node) el.classList.remove('is-traced');
  });
  node.classList.toggle('is-traced');
}

function startDashSSE() {
  stopDashSSE();
  startFetchSSE();
}

function setSseStatus(mode) {
  const statusEl = document.getElementById('sse-status');
  if (!statusEl) return;
  statusEl.classList.remove('is-live', 'is-retry', 'is-down');
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
}

function formatUptime(seconds) {
  if (seconds < 60) return seconds + ' 秒';
  if (seconds < 3600) return Math.floor(seconds / 60) + ' 分';
  if (seconds < 86400) return Math.floor(seconds / 3600) + ' 时 ' + Math.floor((seconds % 3600) / 60) + ' 分';
  return Math.floor(seconds / 86400) + ' 天 ' + Math.floor((seconds % 86400) / 3600) + ' 时';
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

function sitePlaybackEdge(s) {
  const playback = (s.playback_target_url || '').trim();
  let extra = [];
  try { extra = JSON.parse(s.stream_hosts || '[]'); } catch (e) {}
  const total = (playback ? 1 : 0) + extra.length;
  if (total === 0) return '跟随主回源';
  if (total === 1 && playback === (s.target_url || '').trim()) return '与主回源相同';
  return playback || extra[0] || '已分流';
}

async function loadDashboardTable() {
  const topo = document.getElementById('dash-topo');
  if (!topo) return;

  try {
    const sites = await API.listSites();
    if (!sites || sites.length === 0) {
      topo.innerHTML = UI.empty({
        title: '还没有站点',
        body: '去站点管理添加第一个反代。',
        actions: [{ id: 'goto-sites', label: '前往站点管理', className: 'btn-primary' }],
      });
      topo.querySelector('[data-empty-action="goto-sites"]')?.addEventListener('click', () => Router.navigate('sites'));
      return;
    }

    topo.innerHTML = sites.map((s, i) => `
      <article class="topo-node${i === 0 ? ' is-traced' : ''}" tabindex="0" data-id="${s.id}">
        <div class="topo-endlabel">
          <span class="topo-name">${esc(s.name)}</span>
          <div class="topo-endmeta">
            <span class="mono">:${s.listen_port}</span>
            <span class="status-led ${s.running ? 'on' : 'off'}" title="${s.running ? '运行中' : '已停止'}"></span>
          </div>
        </div>
        <div class="topo-edges">
          <div class="edge api">
            <span class="edge-key">API</span>
            <span class="edge-line" aria-hidden="true"></span>
            <span class="mono" title="${esc(s.target_url)}">${esc(s.target_url)}</span>
          </div>
          <div class="edge playback">
            <span class="edge-key">播放</span>
            <span class="edge-line" aria-hidden="true"></span>
            <span class="mono">${esc(sitePlaybackEdge(s))}</span>
          </div>
        </div>
      </article>
    `).join('');
  } catch (e) {
    topo.innerHTML = UI.error({ body: e.message, retry: true });
    topo.querySelector('[data-error-retry]')?.addEventListener('click', loadDashboardTable);
  }
}

async function loadDashboardData() {
  loadDashboardTable();
}
