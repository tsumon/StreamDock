const Router = {
  routes: {},
  current: null,
  initialized: false,

  register(path, handler) {
    this.routes[path] = handler;
  },

  navigate(path) {
    location.hash = path;
  },

  resolve() {
    let hash = location.hash.slice(1) || 'dashboard';
    if (!this.routes[hash] && hash !== 'dashboard') {
      hash = 'dashboard';
      if (location.hash.slice(1) !== hash) location.hash = hash;
    }
    const previous = this.current;

    if (previous === 'dashboard' && hash !== 'dashboard' && typeof stopDashSSE === 'function') {
      stopDashSSE();
    }

    this.current = hash;

    document.querySelectorAll('.topnav-link').forEach(link => {
      const on = link.dataset.page === hash;
      link.classList.toggle('active', on);
      if (on) link.setAttribute('aria-current', 'page');
      else link.removeAttribute('aria-current');
    });
    document.querySelectorAll('.mobile-tab').forEach(tab => {
      const on = tab.dataset.page === hash;
      tab.classList.toggle('active', on);
      if (on) tab.setAttribute('aria-current', 'page');
      else tab.removeAttribute('aria-current');
    });

    document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
    const target = document.getElementById('page-' + hash);
    if (target) target.classList.add('active');

    const handler = this.routes[hash];
    if (handler) handler();
  },

  init() {
    if (this.initialized) return;
    window.addEventListener('hashchange', () => this.resolve());
    this.initialized = true;
  }
};
