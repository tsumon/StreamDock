(function() {
  'use strict';

  if (typeof window.UI === 'undefined') {
    window.UI = {
      clearFormErrors: function() {},
      setFieldError: function(_el, msg) { Toast.error(msg); },
    };
  }

  const loginEl = document.getElementById('page-login');
  const shellEl = document.getElementById('app-shell');
  const loginFooterEl = document.getElementById('login-footer');
  const loginButtonEl = document.getElementById('btn-login');
  const modalEl = document.getElementById('modal');
  const accountPanel = document.getElementById('account-panel');
  const avatarBtn = document.getElementById('avatar-btn');
  let dashboardRefreshTimer = null;
  let appBootstrapped = false;
  let modalBackdropClosable = false;
  let lastModalOpener = null;
  const setupTokenGroupEl = document.getElementById('setup-token-group');
  const setupTokenInputEl = document.getElementById('inp-setup-token');
  let authStatus = {
    needs_setup: false,
    mode: 'single_admin',
    jwt_secret_ephemeral: false,
    setup_token_required: false,
    authenticated: false,
  };

  window.openModal = function(options) {
    modalBackdropClosable = !!(options && options.closeOnBackdrop);
    lastModalOpener = document.activeElement;
    if (typeof modalEl.showModal === 'function') {
      if (!modalEl.open) modalEl.showModal();
    } else {
      modalEl.setAttribute('open', '');
    }
  };

  window.closeModal = function() {
    modalBackdropClosable = false;
    if (typeof modalEl.close === 'function') {
      if (modalEl.open) modalEl.close();
    } else {
      modalEl.removeAttribute('open');
    }
    if (lastModalOpener && typeof lastModalOpener.focus === 'function') {
      lastModalOpener.focus();
    }
    lastModalOpener = null;
  };

  modalEl.addEventListener('click', function(e) {
    if (e.target === modalEl && modalBackdropClosable) closeModal();
  });
  modalEl.addEventListener('cancel', function(e) {
    e.preventDefault();
    closeModal();
  });
  document.getElementById('modal-close').addEventListener('click', closeModal);

  function closeAccountMenu() {
    accountPanel.hidden = true;
    avatarBtn.setAttribute('aria-expanded', 'false');
  }

  function toggleAccountMenu() {
    const open = accountPanel.hidden;
    accountPanel.hidden = !open;
    avatarBtn.setAttribute('aria-expanded', open ? 'true' : 'false');
  }

  avatarBtn.addEventListener('click', function(e) {
    e.stopPropagation();
    toggleAccountMenu();
  });
  document.addEventListener('click', function(e) {
    if (!accountPanel.hidden && !e.target.closest('.account-menu')) closeAccountMenu();
  });
  document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape' && !accountPanel.hidden) closeAccountMenu();
  });

  async function checkAuth() {
    try { localStorage.removeItem('meridian_token'); } catch (e) {}
    try { localStorage.removeItem('streamdock_token'); } catch (e) {}
    try { sessionStorage.removeItem('meridian_user'); } catch (e) {}

    try {
      const res = await API.checkSetup();
      authStatus = Object.assign({}, authStatus, res || {});
      if (res.authenticated) {
        API.authenticated = true;
        if (res.username) API.username = res.username;
        enterApp();
        return;
      }
      if (res.needs_setup) {
        showSetupMode();
        return;
      }
    } catch (e) {
      loginFooterEl.textContent = '连不上面板。确认服务已启动后再刷新。';
    }

    showLoginMode();
  }

  function renderLoginFooter(isSetup) {
    const lines = [];
    if (authStatus.mode === 'single_admin') {
      lines.push(isSetup
        ? '当前为单管理员模式，创建唯一的管理员账号。'
        : '当前为单管理员模式。还没有账号？<a href="#" id="link-register">创建管理员</a>');
    } else {
      lines.push(isSetup
        ? '首次使用，创建管理员账号。'
        : '还没有账号？<a href="#" id="link-register">创建管理员</a>');
    }

    if (authStatus.jwt_secret_ephemeral) {
      lines.push('<span class="login-note warn">未固定 JWT_SECRET，服务重启后需要重新登录。</span>');
    }

    return lines.join('');
  }

  function showSetupMode() {
    loginButtonEl.textContent = '创建管理员';
    loginButtonEl.disabled = false;
    loginFooterEl.innerHTML = renderLoginFooter(true);
    loginEl._isSetup = true;
    document.getElementById('inp-password').autocomplete = 'new-password';
    if (setupTokenGroupEl) setupTokenGroupEl.hidden = false;
    if (setupTokenInputEl) setupTokenInputEl.required = true;
  }

  function showLoginMode() {
    loginButtonEl.textContent = '登录';
    loginButtonEl.disabled = false;
    loginFooterEl.innerHTML = renderLoginFooter(false);
    loginEl._isSetup = false;
    document.getElementById('inp-password').autocomplete = 'current-password';
    if (setupTokenGroupEl) setupTokenGroupEl.hidden = true;
    if (setupTokenInputEl) {
      setupTokenInputEl.required = false;
      setupTokenInputEl.value = '';
    }
  }

  function startDashboardRefresh() {
    if (dashboardRefreshTimer) clearInterval(dashboardRefreshTimer);
    dashboardRefreshTimer = setInterval(() => {
      if (Router.current === 'dashboard') loadDashboardData();
    }, 15000);
  }

  function stopDashboardRefresh() {
    if (!dashboardRefreshTimer) return;
    clearInterval(dashboardRefreshTimer);
    dashboardRefreshTimer = null;
  }

  function teardownAppRuntime() {
    stopDashboardRefresh();
    if (typeof stopDashSSE === 'function') stopDashSSE();
  }

  document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const usernameEl = document.getElementById('inp-username');
    const passwordEl = document.getElementById('inp-password');
    const username = usernameEl.value.trim();
    const password = passwordEl.value;
    UI.clearFormErrors(this);

    let firstInvalid = null;
    if (!username) {
      UI.setFieldError(usernameEl, '填写用户名。');
      firstInvalid = firstInvalid || usernameEl;
    }
    if (!password) {
      UI.setFieldError(passwordEl, '填写密码。');
      firstInvalid = firstInvalid || passwordEl;
    } else if (password.length < 12) {
      UI.setFieldError(passwordEl, '密码至少 12 位。');
      firstInvalid = firstInvalid || passwordEl;
    }
    if (loginEl._isSetup && setupTokenInputEl && !setupTokenInputEl.value.trim()) {
      UI.setFieldError(setupTokenInputEl, '填写初始化令牌 SETUP_TOKEN。');
      firstInvalid = firstInvalid || setupTokenInputEl;
    }
    if (firstInvalid) {
      firstInvalid.focus();
      return;
    }

    loginButtonEl.disabled = true;
    loginButtonEl.setAttribute('aria-busy', 'true');
    const idleLabel = loginEl._isSetup ? '创建管理员' : '登录';
    loginButtonEl.textContent = loginEl._isSetup ? '正在创建…' : '正在登录…';

    try {
      let res;
      if (loginEl._isSetup) {
        const setupToken = setupTokenInputEl ? setupTokenInputEl.value.trim() : '';
        res = await API.setup(username, password, setupToken);
        Toast.success('管理员已创建');
        if (setupTokenInputEl) setupTokenInputEl.value = '';
      } else {
        res = await API.login(username, password);
        Toast.success('已登录');
      }
      API.authenticated = true;
      API.username = res.username;
      enterApp();
    } catch (err) {
      Toast.error(err.message);
      loginButtonEl.disabled = false;
      loginButtonEl.removeAttribute('aria-busy');
      loginButtonEl.textContent = idleLabel;
    }
  });

  loginFooterEl.addEventListener('click', function(e) {
    const registerLink = e.target.closest('#link-register');
    if (!registerLink) return;
    e.preventDefault();
    showSetupMode();
  });

  function enterApp() {
    loginEl.classList.add('hidden');
    shellEl.classList.add('active');
    closeAccountMenu();

    const initial = (API.username || 'A')[0].toUpperCase();
    avatarBtn.textContent = initial;
    document.getElementById('account-name').textContent = API.username || '管理员';

    if (!appBootstrapped) {
      Router.register('dashboard', renderDashboard);
      Router.register('sites', renderSites);
      Router.register('traffic', renderTraffic);
      Router.register('diagnostics', renderDiag);
      Router.init();
      appBootstrapped = true;
    }

    Router.resolve();
    startDashboardRefresh();
  }

  function confirmLogout() {
    closeAccountMenu();
    document.getElementById('modal-title').textContent = '退出登录';
    document.getElementById('modal-body').innerHTML = '<p style="color:var(--ink-secondary)">当前会话会结束，需要重新登录才能管理站点。</p>';
    document.getElementById('modal-footer').innerHTML = `
      <button type="button" class="btn-modal secondary" id="logout-cancel">取消</button>
      <button type="button" class="btn-modal primary" id="logout-confirm">退出</button>
    `;
    document.getElementById('logout-cancel').onclick = closeModal;
    document.getElementById('logout-confirm').onclick = async function() {
      closeModal();
      teardownAppRuntime();
      await API.logout();
      loginEl.classList.remove('hidden');
      shellEl.classList.remove('active');
      showLoginMode();
      document.getElementById('inp-password').value = '';
      document.getElementById('inp-setup-token').value = '';
      Toast.info('已退出登录');
    };
    openModal({ closeOnBackdrop: true });
  }

  document.getElementById('btn-logout').addEventListener('click', confirmLogout);

  checkAuth();
})();
