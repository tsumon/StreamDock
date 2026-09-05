// Sites management page
function renderSites() {
  const page = document.getElementById('page-sites');
  page.innerHTML = `
    <header class="page-header">
      <h1 class="section-title">站点管理</h1>
      <p class="section-sub">主回源、播放分流、监听端口与流量额度。配置可导出为明文 JSON。</p>
    </header>
    <div class="page-toolbar">
      <div class="toolbar-info" id="sites-count">加载中…</div>
      <div class="toolbar-actions">
        <button type="button" class="btn-ghost" id="btn-export-sites" title="导出所有站点配置为 JSON 备份文件">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          导出配置
        </button>
        <button type="button" class="btn-ghost" id="btn-import-sites" title="从 JSON 文件导入站点配置">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
          导入配置
        </button>
        <input type="file" id="import-file-input" accept=".json,application/json" hidden>
        <button type="button" class="btn-add" id="btn-add-site">
          <svg viewBox="0 0 24 24" aria-hidden="true"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          添加站点
        </button>
      </div>
    </div>
    <div class="sites-grid" id="sites-grid">${UI.skeletonCards(2)}</div>
  `;

  document.getElementById('btn-add-site').onclick = () => showSiteModal();
  document.getElementById('btn-export-sites').onclick = exportSitesConfig;
  document.getElementById('btn-import-sites').onclick = () => document.getElementById('import-file-input').click();
  document.getElementById('import-file-input').onchange = importSitesConfig;
  document.getElementById('sites-grid').addEventListener('click', onSitesGridClick);
  loadSites();
}

function onSitesGridClick(e) {
  const btn = e.target.closest('[data-site-action]');
  if (!btn) return;
  const id = parseInt(btn.dataset.id, 10);
  const action = btn.dataset.siteAction;
  if (action === 'toggle') toggleSiteAction(id);
  if (action === 'edit') editSiteAction(id);
  if (action === 'delete') deleteSiteAction(id, btn.dataset.name || '');
  if (action === 'add') showSiteModal();
  if (action === 'import') document.getElementById('import-file-input').click();
}

async function loadSites() {
  const grid = document.getElementById('sites-grid');
  const countEl = document.getElementById('sites-count');
  if (!grid) return;

  try {
    const sites = await API.listSites();
    if (countEl) countEl.innerHTML = `共 <strong>${sites.length}</strong> 个站点`;

    if (!sites || sites.length === 0) {
      grid.innerHTML = UI.empty({
        wide: true,
        title: '还没有站点',
        body: '添加反代，或导入 JSON。',
        actions: [
          { id: 'add', label: '添加站点', className: 'btn-primary' },
          { id: 'import', label: '导入配置', className: 'btn-secondary' },
        ],
      });
      grid.querySelector('[data-empty-action="add"]')?.addEventListener('click', () => showSiteModal());
      grid.querySelector('[data-empty-action="import"]')?.addEventListener('click', () => {
        document.getElementById('import-file-input').click();
      });
      return;
    }

    grid.innerHTML = sites.map((s) => {
      const pct = s.traffic_quota > 0 ? (s.traffic_used / s.traffic_quota * 100).toFixed(1) : 0;
      const pctClass = pct > 85 ? 'danger' : pct > 50 ? 'warn' : 'normal';
      const playbackRow = renderPlaybackRow(s);

      return `
      <article class="site-card">
        <div class="site-top">
          <div class="site-name">${esc(s.name)}</div>
          <div class="topo-endmeta">
            <span class="mono">:${s.listen_port}</span>
            <span class="status-badge">
              <span class="status-led ${s.running ? 'on' : 'off'}"></span>
              ${s.running ? '运行中' : '已停止'}
            </span>
          </div>
        </div>
        <div class="topo-edges" style="margin-bottom:12px">
          <div class="edge api">
            <span class="edge-key">API</span>
            <span class="edge-line" aria-hidden="true"></span>
            <span class="mono" title="${esc(s.target_url)}">${esc(s.target_url)}</span>
          </div>
          ${playbackRow}
        </div>
        <div class="site-rows">
          ${playbackModeRow(s)}
          <div class="site-row">
            <span class="site-row-label">UA 模式</span>
            <span class="pill ${uaClassMap[s.ua_mode] || 'pill-blue'}">${uaNameMap[s.ua_mode] || s.ua_mode}</span>
          </div>
          ${s.traffic_quota > 0 ? `
          <div class="progress-wrap">
            <div class="progress-labels">
              <span>已用 ${formatBytes(s.traffic_used)}</span>
              <span>${formatBytes(s.traffic_quota)}</span>
            </div>
            <div class="progress-track">
              <div class="progress-fill ${pctClass}" style="transform:scaleX(${Math.min(pct, 100) / 100})"></div>
            </div>
          </div>
          ` : `
          <div class="site-row">
            <span class="site-row-label">已用流量</span>
            <span>${formatBytes(s.traffic_used)}</span>
          </div>
          `}
        </div>
        <div class="site-actions">
          <button type="button" class="btn-ghost" data-site-action="toggle" data-id="${s.id}">${s.enabled ? '停用' : '启用'}</button>
          <button type="button" class="btn-ghost" data-site-action="edit" data-id="${s.id}">编辑</button>
          <button type="button" class="btn-ghost danger" data-site-action="delete" data-id="${s.id}" data-name="${esc(s.name)}">删除</button>
        </div>
      </article>`;
    }).join('');
  } catch (e) {
    if (countEl) countEl.textContent = '无法加载站点';
    grid.innerHTML = UI.error({
      wide: true,
      body: e.message,
      retry: true,
    });
    grid.querySelector('[data-error-retry]')?.addEventListener('click', loadSites);
  }
}

function renderPlaybackRow(site) {
  const playback = (site.playback_target_url || '').trim();
  let extraHosts = [];
  try { extraHosts = JSON.parse(site.stream_hosts || '[]'); } catch (e) {}
  const totalHosts = (playback ? 1 : 0) + extraHosts.length;

  if (totalHosts === 0) {
    return `
      <div class="edge playback">
        <span class="edge-key">播放</span>
        <span class="edge-line" aria-hidden="true"></span>
        <span class="mono mono-subtle">跟随主回源</span>
      </div>
    `;
  }

  if (totalHosts === 1 && playback === (site.target_url || '').trim()) {
    return `
      <div class="edge playback">
        <span class="edge-key">播放</span>
        <span class="edge-line" aria-hidden="true"></span>
        <span class="mono mono-subtle">与主回源相同</span>
      </div>
    `;
  }

  let rows = '';
  if (playback) {
    rows += `
    <div class="edge playback">
      <span class="edge-key">播放</span>
      <span class="edge-line" aria-hidden="true"></span>
      <span class="mono" title="${esc(playback)}">${esc(playback)}</span>
    </div>`;
  }
  for (const h of extraHosts) {
    rows += `
    <div class="edge playback">
      <span class="edge-key">额外</span>
      <span class="edge-line" aria-hidden="true"></span>
      <span class="mono" title="${esc(h)}">${esc(h)}</span>
    </div>`;
  }
  return rows;
}

function playbackModeRow(site) {
  const playback = (site.playback_target_url || '').trim();
  let extraHosts = [];
  try { extraHosts = JSON.parse(site.stream_hosts || '[]'); } catch (e) {}
  const totalHosts = (playback ? 1 : 0) + extraHosts.length;
  if (totalHosts === 0) return '';
  if (totalHosts === 1 && playback === (site.target_url || '').trim()) return '';
  const modeLabel = site.playback_mode === 'redirect' ? '重定向跟随' : '直连分流';
  return `
    <div class="site-row">
      <span class="site-row-label">播放模式</span>
      <span class="pill pill-orange">${modeLabel}</span>
    </div>`;
}

function showSiteModal(site) {
  const isEdit = !!site;
  const title = isEdit ? '编辑站点' : '添加站点';

  document.getElementById('modal-title').textContent = title;
  document.getElementById('modal-body').innerHTML = `
    <div class="form-group">
      <label for="m-name">站点名称</label>
      <input type="text" class="form-input" id="m-name" value="${isEdit ? esc(site.name) : ''}" placeholder="如：Emby-US-01" required>
    </div>
    <div class="form-group">
      <label for="m-target">主回源地址</label>
      <input type="text" class="form-input" id="m-target" value="${isEdit ? esc(site.target_url) : ''}" placeholder="如：192.168.1.10:8096 或 https://emby.example.com" required>
      <div class="form-help">网页、API 和默认回源都走这里。</div>
    </div>
    <div class="form-group">
      <label>播放回源列表（可选）</label>
      <div id="m-playback-list"></div>
      <button type="button" class="btn-ghost" id="m-add-playback" style="margin-top:6px">+ 添加播放回源</button>
      <div class="form-help">留空则跟随主回源。可添加多个，用于多推流 / 播放节点。</div>
    </div>
    <div class="form-group" id="playback-mode-group" hidden>
      <label for="m-playback-mode">播放模式</label>
      <select class="form-select modal-select" id="m-playback-mode">
        <option value="direct" ${(!isEdit || site.playback_mode !== 'redirect') ? 'selected' : ''}>直连分流</option>
        <option value="redirect" ${isEdit && site.playback_mode === 'redirect' ? 'selected' : ''}>重定向跟随</option>
      </select>
      <div class="form-help">直连分流：播放请求发到首个播放回源。重定向跟随：请求经主回源，只允许跳到播放回源集合里的 host。</div>
    </div>
    <div class="form-group">
      <label for="m-port">监听端口</label>
      <input type="number" class="form-input" id="m-port" value="${isEdit ? site.listen_port : ''}" placeholder="如：8001" min="1" max="65535" required>
    </div>
    <div class="form-group">
      <label for="m-ua">UA 模式</label>
      <select class="form-select modal-select" id="m-ua">
        <option value="infuse" ${(!isEdit || site.ua_mode === 'infuse') ? 'selected' : ''}>Infuse</option>
        <option value="web" ${isEdit && site.ua_mode === 'web' ? 'selected' : ''}>Web</option>
        <option value="client" ${isEdit && site.ua_mode === 'client' ? 'selected' : ''}>客户端</option>
      </select>
    </div>
    <div class="form-group">
      <label for="m-quota">流量额度（GB，0 = 不限）</label>
      <input type="number" class="form-input" id="m-quota" value="${isEdit ? Math.round((site.traffic_quota || 0) / 1073741824) : 0}" min="0" placeholder="0">
    </div>
  `;

  document.getElementById('modal-footer').innerHTML = `
    <button type="button" class="btn-modal secondary" id="m-cancel">取消</button>
    <button type="button" class="btn-modal primary" id="m-submit">${isEdit ? '保存' : '创建'}</button>
  `;
  document.getElementById('m-cancel').onclick = closeModal;

  const listContainer = document.getElementById('m-playback-list');
  const modeGroup = document.getElementById('playback-mode-group');
  let existingHosts = [];
  if (isEdit) {
    if ((site.playback_target_url || '').trim()) existingHosts.push(site.playback_target_url.trim());
    try {
      const extra = JSON.parse(site.stream_hosts || '[]');
      for (const h of extra) if (h && h.trim()) existingHosts.push(h.trim());
    } catch (e) {}
  }
  if (existingHosts.length === 0) existingHosts = [''];

  function renderPlaybackInputs() {
    listContainer.innerHTML = existingHosts.map((val, idx) => `
      <div class="playback-row">
        <input type="text" class="form-input m-pb-input" value="${esc(val)}" placeholder="${idx === 0 ? '主播放回源地址' : '额外播放回源地址'}">
        ${existingHosts.length > 1 ? `<button type="button" class="btn-ghost danger m-pb-remove" data-idx="${idx}">删除</button>` : ''}
      </div>
    `).join('');
    listContainer.querySelectorAll('.m-pb-remove').forEach(btn => {
      btn.onclick = () => {
        existingHosts.splice(parseInt(btn.dataset.idx, 10), 1);
        renderPlaybackInputs();
        toggleModeGroup();
      };
    });
    listContainer.querySelectorAll('.m-pb-input').forEach((inp, idx) => {
      inp.oninput = () => { existingHosts[idx] = inp.value; toggleModeGroup(); };
    });
  }
  renderPlaybackInputs();

  document.getElementById('m-add-playback').onclick = () => {
    existingHosts.push('');
    renderPlaybackInputs();
    const inputs = listContainer.querySelectorAll('.m-pb-input');
    if (inputs.length) inputs[inputs.length - 1].focus();
  };

  function toggleModeGroup() {
    const hasAny = existingHosts.some(h => h.trim());
    modeGroup.hidden = !hasAny;
  }
  toggleModeGroup();

  document.getElementById('m-submit').onclick = async () => {
    const nameEl = document.getElementById('m-name');
    const targetEl = document.getElementById('m-target');
    const portEl = document.getElementById('m-port');
    const quotaEl = document.getElementById('m-quota');
    UI.clearFormErrors(document.getElementById('modal-body'));

    const name = nameEl.value.trim();
    const target = targetEl.value.trim();
    const port = parseInt(portEl.value, 10);
    const quotaRaw = quotaEl.value === '' ? 0 : parseInt(quotaEl.value, 10);

    let firstInvalid = null;
    if (!name) {
      UI.setFieldError(nameEl, '填写站点名称。');
      firstInvalid = firstInvalid || nameEl;
    }
    if (!target) {
      UI.setFieldError(targetEl, '填写主回源地址，例如 192.168.1.10:8096。');
      firstInvalid = firstInvalid || targetEl;
    }
    if (!port || port < 1 || port > 65535) {
      UI.setFieldError(portEl, '监听端口必须是 1–65535。');
      firstInvalid = firstInvalid || portEl;
    }
    if (Number.isNaN(quotaRaw) || quotaRaw < 0) {
      UI.setFieldError(quotaEl, '额度必须是 0 或正整数 GB。');
      firstInvalid = firstInvalid || quotaEl;
    }
    if (firstInvalid) {
      firstInvalid.focus();
      return;
    }

    const allHosts = existingHosts.map(h => h.trim()).filter(Boolean);
    const data = {
      name,
      target_url: target,
      playback_target_url: allHosts.length > 0 ? allHosts[0] : '',
      playback_mode: document.getElementById('m-playback-mode').value,
      stream_hosts: allHosts.length > 1 ? allHosts.slice(1) : [],
      listen_port: port,
      ua_mode: document.getElementById('m-ua').value,
      traffic_quota: quotaRaw * 1073741824,
    };

    const submit = document.getElementById('m-submit');
    submit.disabled = true;
    submit.setAttribute('aria-busy', 'true');
    submit.textContent = isEdit ? '正在保存…' : '正在创建…';

    try {
      if (isEdit) {
        await API.updateSite(site.id, data);
        Toast.success('站点已更新');
      } else {
        await API.createSite(data);
        Toast.success('站点已创建');
      }
      closeModal();
      loadSites();
    } catch (e) {
      Toast.error(e.message);
      submit.disabled = false;
      submit.removeAttribute('aria-busy');
      submit.textContent = isEdit ? '保存' : '创建';
    }
  };

  openModal();
}

window.toggleSiteAction = async function(id) {
  try {
    const res = await API.toggleSite(id);
    Toast.success(res.enabled ? '站点已启用' : '站点已停用');
    loadSites();
  } catch (e) {
    Toast.error(e.message);
  }
};

window.editSiteAction = async function(id) {
  try {
    const sites = await API.listSites();
    const site = sites.find(s => s.id === id);
    if (site) showSiteModal(site);
  } catch (e) {
    Toast.error(e.message);
  }
};

window.deleteSiteAction = function(id, name) {
  document.getElementById('modal-title').textContent = '删除站点';
  document.getElementById('modal-body').innerHTML = `<p style="color:var(--ink-secondary)">删除 <strong>${esc(name)}</strong> 后无法撤销。该监听端口会停止。</p>`;
  document.getElementById('modal-footer').innerHTML = `
    <button type="button" class="btn-modal secondary" onclick="closeModal()">取消</button>
    <button type="button" class="btn-modal danger" id="confirm-delete">删除</button>
  `;
  document.getElementById('confirm-delete').onclick = () => confirmDelete(id);
  openModal({ closeOnBackdrop: true });
};

window.confirmDelete = async function(id) {
  try {
    await API.deleteSite(id);
    Toast.success('站点已删除');
    closeModal();
    loadSites();
  } catch (e) {
    Toast.error(e.message);
  }
};

async function exportSitesConfig() {
  try {
    const res = await fetch('/api/sites/export', {
      credentials: 'same-origin',
    });
    if (!res.ok) throw new Error('导出失败：HTTP ' + res.status);
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'streamdock_backup_' + new Date().toISOString().split('T')[0] + '.json';
    a.click();
    URL.revokeObjectURL(url);
    Toast.success('已下载 ' + a.download);
  } catch (e) {
    Toast.error('导出失败：' + e.message);
  }
}

function renderImportSkipList(items) {
  const existing = document.getElementById('import-skip-list');
  if (!items || !items.length) {
    if (existing) existing.remove();
    return;
  }
  const rows = items.map((item) => {
    const port = item.listen_port ? `:${esc(String(item.listen_port))}` : '';
    return `<li><span class="mono">${port}</span> ${esc(item.name || '(未命名)')} — ${esc(item.reason || '已跳过')}</li>`;
  }).join('');
  const html = `<div class="import-skip-list" id="import-skip-list"><p>跳过 ${items.length} 个</p><ul>${rows}</ul></div>`;
  if (existing) {
    existing.outerHTML = html;
    return;
  }
  document.querySelector('#page-sites .page-toolbar')?.insertAdjacentHTML('afterend', html);
}

async function importSitesConfig(e) {
  const file = e.target.files[0];
  if (!file) return;
  e.target.value = '';

  let parsed;
  try {
    const text = await file.text();
    parsed = JSON.parse(text);
  } catch (err) {
    Toast.error('文件不是有效 JSON。');
    return;
  }

  let sites = [];
  if (Array.isArray(parsed)) {
    sites = parsed;
  } else if (parsed.sites && Array.isArray(parsed.sites)) {
    sites = parsed.sites;
  } else {
    Toast.error('配置文件需要是站点数组，或包含 sites 数组的对象。');
    return;
  }

  if (sites.length === 0) {
    Toast.error('文件里没有可导入的站点。');
    return;
  }

  const names = sites.map(s => `• ${s.name || '(未命名)'} → :${s.listen_port || '?'} → ${s.target_url || '?'}`).join('\n');
  document.getElementById('modal-title').textContent = '导入配置';
  document.getElementById('modal-body').innerHTML = `
    <p style="color:var(--ink-secondary);margin-bottom:12px">将导入 <strong>${sites.length}</strong> 个站点。只新建，不覆盖已有配置。</p>
    <pre class="import-preview">${esc(names)}</pre>
    <p class="form-help" style="margin-top:10px">端口冲突或字段不全会跳过，并列出原因。</p>
  `;
  document.getElementById('modal-footer').innerHTML = `
    <button type="button" class="btn-modal secondary" onclick="closeModal()">取消</button>
    <button type="button" class="btn-modal primary" id="m-confirm-import">确认导入</button>
  `;
  document.getElementById('m-confirm-import').onclick = async () => {
    const btn = document.getElementById('m-confirm-import');
    btn.disabled = true;
    btn.textContent = '正在导入…';
    try {
      const res = await fetch('/api/sites/import', {
        method: 'POST',
        credentials: 'same-origin',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ sites, version: parsed && parsed.version })
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || '导入失败');
      closeModal();
      const skippedItems = Array.isArray(data.skipped_items) ? data.skipped_items : [];
      Toast.success('导入 ' + data.created + ' 个站点' + (data.skipped > 0 ? '，跳过 ' + data.skipped + ' 个' : ''));
      if (skippedItems.length) {
        const lines = skippedItems.map((item) => {
          const name = item.name || '(未命名)';
          const port = item.listen_port ? ':' + item.listen_port : '';
          return name + port + ' — ' + (item.reason || '已跳过');
        });
        Toast.info('跳过：\n' + lines.join('\n'), 8000);
      }
      loadSites();
      renderImportSkipList(skippedItems);
    } catch (err) {
      Toast.error('导入失败：' + err.message);
      btn.disabled = false;
      btn.textContent = '确认导入';
    }
  };
  openModal({ closeOnBackdrop: true });
}
