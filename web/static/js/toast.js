// Toast notification system
const Toast = {
  show(msg, type, duration) {
    const container = document.getElementById('toast-container');
    if (!container) return;
    const el = document.createElement('div');
    el.className = 'toast ' + (type || 'info');
    el.textContent = msg;
    container.appendChild(el);
    setTimeout(() => {
      el.classList.add('out');
      setTimeout(() => el.remove(), 180);
    }, duration || 3000);
  },
  success(msg, duration) { this.show(msg, 'success', duration); },
  error(msg, duration) { this.show(msg, 'error', duration || 5000); },
  info(msg, duration) { this.show(msg, 'info', duration); },
};
