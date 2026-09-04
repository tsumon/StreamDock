// StreamDock API Client
const API = {
  authenticated: false,
  get username() { return sessionStorage.getItem('streamdock_user') || sessionStorage.getItem('meridian_user') || ''; },
  set username(v) {
    try { sessionStorage.removeItem('meridian_user'); } catch (e) {}
    v ? sessionStorage.setItem('streamdock_user', v) : sessionStorage.removeItem('streamdock_user');
  },

  async request(method, path, body) {
    const opts = {
      method,
      credentials: 'same-origin',
      headers: { 'Content-Type': 'application/json' },
    };
    if (body) opts.body = JSON.stringify(body);

    let res;
    try {
      res = await fetch(path, opts);
    } catch (err) {
      throw new Error('连不上面板。检查 StreamDock 是否在运行，然后重试。');
    }

    let data = {};
    try {
      data = await res.json();
    } catch (err) {
      if (res.status === 401) this.authenticated = false;
      if (!res.ok) throw new Error('请求失败（HTTP ' + res.status + '）');
      throw new Error('服务器返回了无法解析的响应。');
    }
    if (res.status === 401) this.authenticated = false;
    if (!res.ok) throw new Error(data.error || '请求失败（HTTP ' + res.status + '）');
    return data;
  },

  // Auth
  checkSetup() { return this.request('GET', '/api/auth/check'); },
  login(username, password) { return this.request('POST', '/api/auth/login', { username, password }); },
  setup(username, password, setupToken) {
    return this.request('POST', '/api/auth/setup', { username, password, setup_token: setupToken });
  },
  async logout() {
    try {
      await this.request('POST', '/api/auth/logout');
    } catch (e) {
      // Clear local session even if the server is unreachable.
    }
    this.authenticated = false;
    this.username = null;
  },

  // Dashboard
  dashboard() { return this.request('GET', '/api/dashboard'); },

  // Sites
  listSites() { return this.request('GET', '/api/sites'); },
  createSite(data) { return this.request('POST', '/api/sites', data); },
  updateSite(id, data) { return this.request('PUT', '/api/sites/' + id, data); },
  deleteSite(id) { return this.request('DELETE', '/api/sites/' + id); },
  toggleSite(id) { return this.request('POST', '/api/sites/' + id + '/toggle'); },
  diagSite(id) { return this.request('GET', '/api/sites/' + id + '/diag'); },

  // Config export/import
  exportSites() { return this.request('GET', '/api/sites/export'); },
  importSites(sites) { return this.request('POST', '/api/sites/import', { sites }); },

  // Traffic
  getTraffic(siteId, hours) { return this.request('GET', '/api/traffic/' + siteId + '?hours=' + (hours || 24)); },

  // UA Profiles
  getProfiles() { return this.request('GET', '/api/ua-profiles'); },
};

const uaClassMap = { infuse: 'pill-blue', web: 'pill-green', client: 'pill-orange' };
const uaNameMap = { infuse: 'Infuse', web: 'Web', client: '客户端' };

function esc(str) {
  const d = document.createElement('div');
  d.textContent = str == null ? '' : String(str);
  return d.innerHTML;
}

function formatBytes(bytes) {
  if (!bytes || bytes === 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return (bytes / Math.pow(1024, i)).toFixed(i > 1 ? 1 : 0) + '\u00a0' + units[i];
}

const UI = {
  empty(opts) {
    const actions = (opts.actions || []).map(a =>
      `<button type="button" class="${a.className || 'btn-secondary'}" data-empty-action="${esc(a.id)}">${a.label}</button>`
    ).join('');
    const icon = opts.icon || '<svg viewBox="0 0 24 24"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>';
    return `
      <div class="empty-state ${opts.inline ? 'inline' : ''}" ${opts.wide ? 'style="grid-column:1/-1"' : ''}>
        <div class="empty-icon" aria-hidden="true">${icon}</div>
        <div class="empty-title">${opts.title}</div>
        <p class="empty-body">${opts.body}</p>
        ${actions ? `<div class="empty-actions">${actions}</div>` : ''}
      </div>`;
  },

  error(opts) {
    const retry = opts.retry
      ? `<button type="button" class="btn-secondary" data-error-retry>${opts.retryLabel || '重试'}</button>`
      : '';
    return `
      <div class="error-banner ${opts.inline ? 'inline' : ''}" role="alert" ${opts.wide ? 'style="grid-column:1/-1"' : ''}>
        <div class="error-title">${opts.title || '加载失败'}</div>
        <p class="error-body">${esc(opts.body)}</p>
        ${retry ? `<div class="error-actions">${retry}</div>` : ''}
      </div>`;
  },

  skeletonCards(n) {
    return Array.from({ length: n }, () => `
      <div class="skeleton-card" aria-hidden="true">
        <div class="skeleton skeleton-line" style="width:40%;height:16px;margin-bottom:18px"></div>
        <div class="skeleton skeleton-line" style="width:88%;margin-bottom:10px"></div>
        <div class="skeleton skeleton-line" style="width:70%;margin-bottom:10px"></div>
        <div class="skeleton skeleton-line" style="width:54%"></div>
      </div>`).join('');
  },

  skeletonRows(n, cols) {
    const cells = Array.from({ length: cols || 6 }, () =>
      '<td><div class="skeleton skeleton-line" style="width:72%"></div></td>'
    ).join('');
    return Array.from({ length: n }, () => `<tr class="skeleton-table-row" aria-hidden="true">${cells}</tr>`).join('');
  },

  setFieldError(input, msg) {
    if (!input) return;
    const group = input.closest('.form-group') || input.parentElement;
    input.classList.toggle('invalid', !!msg);
    input.setAttribute('aria-invalid', msg ? 'true' : 'false');
    let err = group.querySelector('.field-error');
    if (msg) {
      if (!err) {
        err = document.createElement('div');
        err.className = 'field-error';
        err.id = (input.id || 'field') + '-error';
        group.appendChild(err);
      }
      input.setAttribute('aria-describedby', err.id);
      err.textContent = msg;
    } else if (err) {
      err.remove();
      input.removeAttribute('aria-describedby');
    }
  },

  clearFormErrors(root) {
    (root || document).querySelectorAll('.form-input.invalid, .form-select.invalid').forEach(el => {
      this.setFieldError(el, '');
    });
  }
};
window.UI = UI;

