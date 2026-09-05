(function() {
  'use strict';

  var STORAGE_KEY = 'streamdock-theme';
  var THEME_COLOR = { light: '#f4f6f8', dark: '#151c21' };

  function storedPreference() {
    try {
      var value = localStorage.getItem(STORAGE_KEY);
      if (value === 'light' || value === 'dark' || value === 'system') return value;
    } catch (e) {}
    return 'system';
  }

  function systemTheme() {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function resolveTheme(preference) {
    var pref = preference || storedPreference();
    if (pref === 'dark' || pref === 'light') return pref;
    return systemTheme();
  }

  function applyResolved(resolved, preference) {
    var next = resolved === 'dark' ? 'dark' : 'light';
    var root = document.documentElement;
    root.setAttribute('data-theme', next);
    root.style.colorScheme = next;
    var meta = document.querySelector('meta[name="theme-color"]');
    if (meta) meta.setAttribute('content', THEME_COLOR[next]);
    syncToggles(preference || storedPreference());
    document.dispatchEvent(new CustomEvent('streamdock:theme', { detail: { theme: next, preference: preference || storedPreference() } }));
  }

  function persistPreference(preference) {
    try { localStorage.setItem(STORAGE_KEY, preference); } catch (e) {}
  }

  function syncToggles(preference) {
    document.querySelectorAll('[data-theme-choice]').forEach(function(btn) {
      btn.setAttribute('aria-pressed', btn.getAttribute('data-theme-choice') === preference ? 'true' : 'false');
    });
  }

  document.addEventListener('click', function(e) {
    var btn = e.target.closest('[data-theme-choice]');
    if (!btn) return;
    e.preventDefault();
    var preference = btn.getAttribute('data-theme-choice');
    if (preference !== 'light' && preference !== 'dark' && preference !== 'system') return;
    persistPreference(preference);
    applyResolved(resolveTheme(preference), preference);
  });

  var media = window.matchMedia('(prefers-color-scheme: dark)');
  function onSystemChange() {
    if (storedPreference() !== 'system') return;
    applyResolved(systemTheme(), 'system');
  }
  if (typeof media.addEventListener === 'function') {
    media.addEventListener('change', onSystemChange);
  } else if (typeof media.addListener === 'function') {
    media.addListener(onSystemChange);
  }

  applyResolved(resolveTheme(), storedPreference());

  window.Theme = {
    preference: storedPreference,
    resolve: resolveTheme,
    apply: function(preference) {
      persistPreference(preference);
      applyResolved(resolveTheme(preference), preference);
    },
  };
})();
