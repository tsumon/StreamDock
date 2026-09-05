(function() {
  'use strict';

  var STORAGE_KEY = 'streamdock-theme';
  var THEME_COLOR = { light: '#f4f6f8', dark: '#151c21' };

  function storedTheme() {
    try {
      var value = localStorage.getItem(STORAGE_KEY);
      if (value === 'light' || value === 'dark') return value;
    } catch (e) {}
    return null;
  }

  function systemTheme() {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function resolveTheme() {
    return storedTheme() || systemTheme();
  }

  function applyTheme(theme, persist) {
    var next = theme === 'dark' ? 'dark' : 'light';
    var root = document.documentElement;
    root.setAttribute('data-theme', next);
    root.style.colorScheme = next;
    var meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.setAttribute('content', THEME_COLOR[next]);
    if (persist) {
      try { localStorage.setItem(STORAGE_KEY, next); } catch (e) {}
    }
    syncToggles(next);
    document.dispatchEvent(new CustomEvent('streamdock:theme', { detail: { theme: next } }));
  }

  function syncToggles(theme) {
    document.querySelectorAll('[data-theme-choice]').forEach(function(btn) {
      btn.setAttribute('aria-pressed', btn.getAttribute('data-theme-choice') === theme ? 'true' : 'false');
    });
  }

  document.addEventListener('click', function(e) {
    var btn = e.target.closest('[data-theme-choice]');
    if (!btn) return;
    e.preventDefault();
    applyTheme(btn.getAttribute('data-theme-choice'), true);
  });

  var media = window.matchMedia('(prefers-color-scheme: dark)');
  function onSystemChange(e) {
    if (storedTheme()) return;
    applyTheme(e.matches ? 'dark' : 'light', false);
  }
  if (typeof media.addEventListener === 'function') {
    media.addEventListener('change', onSystemChange);
  } else if (typeof media.addListener === 'function') {
    media.addListener(onSystemChange);
  }

  applyTheme(resolveTheme(), false);

  window.Theme = {
    resolve: resolveTheme,
    apply: applyTheme,
  };
})();
