# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

**[inferred from README, panel copy, and SKILL-PLAN-UI.md]** Homelab operators who already run Emby. They sit at a desk next to a NAS, a switch, and a monitor — often a spare room or closet-adjacent workspace — and need to keep reverse-proxy sites up, split playback away from API traffic, and restore a config from a JSON file without learning a new product.

They are not first-time SaaS buyers. They already know ports, upstreams, cookies, and systemd.

## Product Purpose

StreamDock is a single-binary Emby reverse-proxy admin panel. The operator adds sites in the browser; each site listens on its own port. API / WebSocket traffic goes to `target_url`; playback paths go to `playback_target_url` and `stream_hosts`. SQLite persists state. The four-page vanilla frontend is embedded in the same process.

Success is: sites stay up, playback actually splits, a JSON backup can be exported and re-imported, and first boot can create the only admin with `SETUP_TOKEN`.

## Positioning

**[inferred]** Neighboring products (upstream Meridian, Nginx Proxy Manager, a hand-written Caddyfile) can reverse-proxy Emby. StreamDock's claim that those cannot copy as a package: one Go binary with an embedded four-page panel, **plaintext** site JSON import/export (`GET /api/sites/export`, `POST /api/sites/import`, version `streamdock-v1`), an independent `#traffic` page, and a narrower playback split — without a JS framework or a build step.

Upstream Meridian already has `playback_target_url` / `playback_mode` / `stream_hosts`; those fields are preserved, not invented. The plaintext export and the standalone traffic page are the fork-unique capabilities.

## Operating Context

- Default panel bind: `127.0.0.1:9090`. Each site opens another port for Emby clients.
- Session: HttpOnly cookie `streamdock_session`. Scripts may still send `Authorization: Bearer`.
- First empty database: create the single admin with `SETUP_TOKEN` from the start log or `.env`.
- Backup ritual: Sites → 导出配置 downloads `streamdock_backup_YYYY-MM-DD.json`. Import only creates; port conflicts are skipped. Files do not contain users, traffic, or JWT.
- Hash routes: `#dashboard` `#sites` `#traffic` `#diagnostics`.
- Copy language of the panel is Chinese.
- GitHub: `tsumon/StreamDock`. Product name is StreamDock everywhere a human reads it.

## Capabilities and Constraints

Must keep:

- Site CRUD, start/stop, UA profiles, traffic quota.
- Playback fields: `playback_target_url`, `playback_mode`, `stream_hosts`.
- Plaintext JSON import/export.
- Independent traffic page (`#traffic`) with inbound/outbound series.
- Diagnostics against main upstream, playback upstream, certs, local listen.
- Cookie login and empty-db `SETUP_TOKEN` setup.
- Vanilla HTML/CSS/JS embedded by Go (`web/static/`, `web/embed.go`). No React, Vue, or build chain.

Must not:

- Change Go API behavior in this visual-world pass.
- Reintroduce the Meridian wordmark in the chrome, login, title, or DESIGN.md.
- Commit or push.

## Brand Commitments

- Name: **StreamDock**. GitHub `tsumon/StreamDock`.
- Joe (2026-09-04): the incumbent visual world **"The Dark Glass Console"** / Apple TV / cinema black / glass top bar / Inter / `#0A84FF` is an **anti-reference**. Replace it. Do not polish it. Never split the difference.
- Joe: management panel ≠ must be pure black. Do not land on near-black + one neon accent + glowing edges.
- Joe: expression must not hide the task, state, or familiar controls (Operate).

## Evidence on Hand

- Runnable panel: `web/static/index.html` + `css/style.css` + `js/{app,router,api,toast}.js` + `js/pages/{dashboard,sites,traffic,diag}.js`.
- README, `docs/*.png` (screenshots still show the old Meridian wordmark; they are evidence of the discarded world, not of the new one).
- No customer quotes, no usage metrics, no invented testimonials. Future work must not fabricate them.

## Product Principles

1. The operator came to move a site, not to watch a dashboard perform.
2. Playback split and plaintext backup are first-class, not hidden advanced fields.
3. Familiar controls beat invented widgets. Native `<select>`, native `<dialog>`, native file picker.
4. One binary, one panel, four pages. Density is a feature.
5. Empty, loading, and error are layouts, not afterthoughts.

## Accessibility & Inclusion

No product-specific standard was recorded beyond shipping a usable operator panel: skip link, `:focus-visible`, Chinese copy, 44px mobile tabs, `aria-live` toasts. **[inferred]** Keep those; do not regress contrast or keyboard access while replacing the visual world.
