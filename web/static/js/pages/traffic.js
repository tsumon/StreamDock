// Traffic statistics page
function renderTraffic() {
  const page = document.getElementById('page-traffic');
  page.innerHTML = `
    <header class="page-header">
      <h1 class="section-title">流量统计</h1>
      <p class="section-sub">按站点查看入站 / 出站字节。计量是两端合计。</p>
    </header>
    <div class="controls-row">
      <select class="form-select" id="traffic-site-select" aria-label="选择站点">
        <option value="">加载中…</option>
      </select>
      <select class="form-select" id="traffic-hours-select" aria-label="时间范围">
        <option value="24">最近 24 小时</option>
        <option value="168">最近 7 天</option>
        <option value="720">最近 30 天</option>
      </select>
    </div>
    <div id="traffic-body"></div>
  `;

  loadTrafficSites();
  document.getElementById('traffic-site-select').onchange = loadTrafficChart;
  document.getElementById('traffic-hours-select').onchange = loadTrafficChart;
}

function trafficChartMarkup() {
  return `
    <div class="chart-wrap">
      <div class="chart-head">
        <h3>流量趋势</h3>
        <div class="chart-legend">
          <div class="legend-item"><div class="legend-dot in"></div>入站</div>
          <div class="legend-item"><div class="legend-dot out"></div>出站</div>
        </div>
      </div>
      <canvas id="trafficChart"></canvas>
    </div>
    <div class="traffic-totals" id="traffic-totals"></div>
  `;
}

async function loadTrafficSites() {
  const body = document.getElementById('traffic-body');
  const sel = document.getElementById('traffic-site-select');
  try {
    const sites = await API.listSites();
    if (!sites || sites.length === 0) {
      sel.innerHTML = '<option value="">暂无站点</option>';
      body.innerHTML = UI.empty({
        title: '还没有站点可统计',
        body: '先到站点管理添加一个反代，流量才会按站点记下来。',
        actions: [{ id: 'goto-sites', label: '前往站点管理', className: 'btn-primary' }],
      });
      body.querySelector('[data-empty-action="goto-sites"]')?.addEventListener('click', () => Router.navigate('sites'));
      return;
    }
    sel.innerHTML = sites.map(s => `<option value="${s.id}">${esc(s.name)}</option>`).join('');
    body.innerHTML = trafficChartMarkup();
    loadTrafficChart();
  } catch (e) {
    sel.innerHTML = '<option value="">加载失败</option>';
    body.innerHTML = UI.error({ body: e.message, retry: true });
    body.querySelector('[data-error-retry]')?.addEventListener('click', loadTrafficSites);
  }
}

async function loadTrafficChart() {
  const siteId = document.getElementById('traffic-site-select')?.value;
  const hours = parseInt(document.getElementById('traffic-hours-select')?.value, 10);
  const totals = document.getElementById('traffic-totals');
  if (!siteId) return;

  if (!document.getElementById('trafficChart')) {
    document.getElementById('traffic-body').innerHTML = trafficChartMarkup();
  }

  try {
    const [logs, sites] = await Promise.all([
      API.getTraffic(siteId, hours),
      API.listSites(),
    ]);
    const site = sites.find(s => s.id === parseInt(siteId, 10));

    const totalIn = logs.reduce((a, l) => a + (l.bytes_in || 0), 0);
    const totalOut = logs.reduce((a, l) => a + (l.bytes_out || 0), 0);

    if (totals) {
      totals.innerHTML = `
        <div class="total-card">
          <div class="total-label">入站流量</div>
          <div class="total-value">${formatBytes(totalIn)}</div>
        </div>
        <div class="total-card">
          <div class="total-label">出站流量</div>
          <div class="total-value">${formatBytes(totalOut)}</div>
        </div>
        <div class="total-card">
          <div class="total-label">累计使用</div>
          <div class="total-value">${formatBytes(site ? site.traffic_used : 0)}</div>
          ${site && site.traffic_quota > 0 ? `<div class="total-delta">额度 ${formatBytes(site.traffic_quota)}</div>` : ''}
        </div>
      `;
    }

    drawTrafficChart(logs, hours);
  } catch (e) {
    const body = document.getElementById('traffic-body');
    body.innerHTML = UI.error({ body: e.message, retry: true });
    body.querySelector('[data-error-retry]')?.addEventListener('click', loadTrafficChart);
  }
}

function trafficThemeColors() {
  const styles = getComputedStyle(document.documentElement);
  const read = function (name, fallback) {
    return styles.getPropertyValue(name).trim() || fallback;
  };
  return {
    grid: read('--grid', '#d7e2ea'),
    inkMuted: read('--ink-muted', '#3e5360'),
    api: read('--api', '#2f5f73'),
    playback: read('--playback', '#c45c26'),
    apiDim: read('--api-dim', 'rgba(47, 95, 115, 0.12)'),
    playbackDim: read('--playback-dim', 'rgba(196, 92, 38, 0.12)'),
  };
}

function drawTrafficChart(logs, hours) {
  const canvas = document.getElementById('trafficChart');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  const dpr = window.devicePixelRatio || 1;
  const w = canvas.parentElement.clientWidth;
  const h = 260;
  canvas.width = w * dpr;
  canvas.height = h * dpr;
  canvas.style.width = w + 'px';
  canvas.style.height = h + 'px';
  ctx.setTransform(1, 0, 0, 1, 0, 0);
  ctx.scale(dpr, dpr);

  const colors = trafficThemeColors();
  const pad = { top: 24, right: 24, bottom: 36, left: 54 };
  const cw = w - pad.left - pad.right;
  const ch = h - pad.top - pad.bottom;

  const numPoints = Math.min(hours, 24);
  const inbound = new Array(numPoints).fill(0);
  const outbound = new Array(numPoints).fill(0);

  if (logs.length > 0) {
    const now = Date.now();
    logs.forEach(l => {
      const t = new Date(l.recorded_at).getTime();
      const hoursAgo = (now - t) / 3600000;
      const idx = numPoints - 1 - Math.floor(hoursAgo * numPoints / hours);
      if (idx >= 0 && idx < numPoints) {
        inbound[idx] += l.bytes_in / (1024 * 1024);
        outbound[idx] += l.bytes_out / (1024 * 1024);
      }
    });
  }

  const rawMax = Math.max(0, ...inbound, ...outbound);
  const maxV = rawMax > 0 ? rawMax * 1.2 : 1;
  const x = i => pad.left + (i / (numPoints - 1 || 1)) * cw;
  const y = v => pad.top + (1 - v / maxV) * ch;

  ctx.clearRect(0, 0, w, h);

  function axisLabel(v) {
    if (v === 0) return '0';
    if (v >= 10) return v.toFixed(0);
    if (v >= 1) return v.toFixed(1);
    return v.toFixed(2);
  }

  ctx.strokeStyle = colors.grid;
  ctx.lineWidth = 1;
  for (let i = 0; i <= 4; i++) {
    const yy = pad.top + (i / 4) * ch;
    ctx.beginPath(); ctx.moveTo(pad.left, yy); ctx.lineTo(w - pad.right, yy); ctx.stroke();
    ctx.fillStyle = colors.inkMuted;
    ctx.font = '11px "Source Sans 3", "Noto Sans SC", system-ui';
    ctx.textAlign = 'right';
    const value = (4 - i) / 4 * (logs.length === 0 ? 0 : maxV);
    ctx.fillText(axisLabel(value) + ' MB', pad.left - 12, yy + 4);
  }

  if (logs.length === 0) {
    ctx.fillStyle = colors.inkMuted;
    ctx.font = '14px "Source Sans 3", "Noto Sans SC", system-ui';
    ctx.textAlign = 'center';
    ctx.fillText('这段时间没有流量记录', w / 2, h / 2);
    return;
  }

  function smoothLine(data, color, fillFrom) {
    ctx.save();
    ctx.beginPath();
    ctx.moveTo(x(0), y(data[0]));
    for (let i = 1; i < data.length; i++) {
      const xc = (x(i - 1) + x(i)) / 2;
      const yc = (y(data[i - 1]) + y(data[i])) / 2;
      ctx.quadraticCurveTo(x(i - 1), y(data[i - 1]), xc, yc);
    }
    ctx.lineTo(x(data.length - 1), y(data[data.length - 1]));
    ctx.strokeStyle = color;
    ctx.lineWidth = 2;
    ctx.stroke();
    ctx.lineTo(x(data.length - 1), pad.top + ch);
    ctx.lineTo(x(0), pad.top + ch);
    ctx.closePath();
    const grad = ctx.createLinearGradient(0, pad.top, 0, pad.top + ch);
    grad.addColorStop(0, fillFrom);
    grad.addColorStop(1, 'transparent');
    ctx.fillStyle = grad;
    ctx.fill();
    ctx.restore();
  }

  smoothLine(outbound, colors.playback, colors.playbackDim);
  smoothLine(inbound, colors.api, colors.apiDim);
}

function redrawTrafficIfVisible() {
  if (typeof Router === 'undefined' || Router.current !== 'traffic') return;
  const canvas = document.getElementById('trafficChart');
  if (canvas) loadTrafficChart();
}

window.addEventListener('resize', redrawTrafficIfVisible);
document.addEventListener('streamdock:theme', redrawTrafficIfVisible);
new MutationObserver(redrawTrafficIfVisible).observe(document.documentElement, {
  attributes: true,
  attributeFilter: ['data-theme'],
});
