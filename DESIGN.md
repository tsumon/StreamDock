---
name: StreamDock
description: Fluorescent-paper topology worksheet for the Emby reverse-proxy admin panel.
colors:
  paper: "#f4f6f8"
  grid: "#d7e2ea"
  sheet: "#ffffff"
  ink: "#1c2b36"
  ink-secondary: "#334654"
  ink-muted: "#3e5360"
  hairline: "#c5d0d8"
  hairline-strong: "#9aafbc"
  primary: "#1c2b36"
  primary-hover: "#0f1a22"
  api: "#2f5f73"
  api-dim: "rgba(47, 95, 115, 0.12)"
  playback: "#c45c26"
  playback-dim: "rgba(196, 92, 38, 0.12)"
  green: "#2e7d4f"
  green-dim: "rgba(46, 125, 79, 0.12)"
  orange: "#a15c12"
  orange-dim: "rgba(161, 92, 18, 0.12)"
  red: "#b42318"
  red-dim: "rgba(180, 35, 24, 0.1)"
  surface-hover: "#eef2f5"
  surface-active: "#e4eaef"
typography:
  display:
    fontFamily: "Source Sans 3, Noto Sans SC, Source Han Sans SC, system-ui, sans-serif"
    fontSize: "1.7rem"
    fontWeight: 700
    lineHeight: 1.15
    letterSpacing: "-0.02em"
  headline:
    fontFamily: "Source Sans 3, Noto Sans SC, Source Han Sans SC, system-ui, sans-serif"
    fontSize: "1.45rem"
    fontWeight: 700
    lineHeight: 1.2
    letterSpacing: "-0.02em"
  title:
    fontFamily: "Source Sans 3, Noto Sans SC, Source Han Sans SC, system-ui, sans-serif"
    fontSize: "0.95rem"
    fontWeight: 700
    lineHeight: 1.3
    letterSpacing: "normal"
  figure:
    fontFamily: "Source Sans 3, Noto Sans SC, Source Han Sans SC, system-ui, sans-serif"
    fontSize: "1.15rem"
    fontWeight: 700
    lineHeight: 1.15
    letterSpacing: "-0.02em"
  body:
    fontFamily: "Source Sans 3, Noto Sans SC, Source Han Sans SC, system-ui, sans-serif"
    fontSize: "0.85rem"
    fontWeight: 400
    lineHeight: 1.5
    letterSpacing: "normal"
  label:
    fontFamily: "Source Sans 3, Noto Sans SC, Source Han Sans SC, system-ui, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 600
    lineHeight: 1.4
    letterSpacing: "normal"
  edge-key:
    fontFamily: "Source Sans 3, Noto Sans SC, Source Han Sans SC, system-ui, sans-serif"
    fontSize: "0.68rem"
    fontWeight: 700
    lineHeight: 1.2
    letterSpacing: "0.04em"
  mono:
    fontFamily: "ui-monospace, SF Mono, Cascadia Code, Noto Sans Mono, monospace"
    fontSize: "0.75rem"
    fontWeight: 400
    lineHeight: 1.4
    letterSpacing: "normal"
rounded:
  xs: "2px"
  sm: "3px"
  md: "4px"
  lg: "6px"
spacing:
  xs: "8px"
  sm: "12px"
  md: "16px"
  lg: "24px"
  xl: "32px"
  2xl: "40px"
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.sheet}"
    rounded: "{rounded.sm}"
    padding: "10px 16px"
    height: "40px"
  button-primary-hover:
    backgroundColor: "{colors.primary-hover}"
    textColor: "{colors.sheet}"
    rounded: "{rounded.sm}"
    padding: "10px 16px"
    height: "40px"
  button-ghost:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.ink-secondary}"
    rounded: "{rounded.xs}"
    padding: "8px 12px"
    height: "36px"
  button-ghost-hover:
    backgroundColor: "{colors.surface-hover}"
    textColor: "{colors.ink}"
    rounded: "{rounded.xs}"
    padding: "8px 12px"
    height: "36px"
  button-secondary:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.ink-secondary}"
    rounded: "{rounded.sm}"
    padding: "10px 16px"
    height: "40px"
  field:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "10px 12px"
    height: "44px"
  nav-title-block:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    height: "52px"
    padding: "0 28px"
  spec-strip:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.ink}"
    padding: "10px 16px 12px"
  topo-node:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "14px 16px 12px"
  pill:
    backgroundColor: "transparent"
    textColor: "{colors.api}"
    rounded: "{rounded.xs}"
    padding: "2px 8px"
  sheet:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "12px 16px"
  dialog:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.ink}"
    rounded: "{rounded.lg}"
    padding: "16px 20px"
    width: "480px"
---

# Design System: StreamDock

## Overview

**Creative North Star: "The Topology Worksheet"**

StreamDock is a network worksheet, not a cinema console. Cool fluorescent paper carries a 20px drafting-blue grid; white sheets and navy-ink boxes sit on that field like annotations on a plot plan. Each site is a node with two outbound edges — drafting-blue for API, copper for playback — so the operator can see where traffic actually goes without decoding four hero tiles.

The panel is an Operate surface for a homelab desk: dense, bilingual (Source Sans 3 with Noto Sans SC), and built from native controls. Title-block chrome is a paper strip ruled with a 2px ink line. Depth is stroke, not glow. The discarded world — The Dark Glass Console, Inter, cinema black, system blue `#0A84FF`, glass bars — is an anti-reference, not a fallback.

**Key Characteristics:**
- Fluorescent paper canvas with a fixed 20px drafting grid
- Navy ink as type, chrome rule, and primary fill; copper reserved for the playback route and focus
- 2–6px drafting corners; square datasheet bands (spec strip, traffic totals)
- Two named route colors, never a third “accent”
- Source Sans 3 + Noto Sans SC at a 15px root; tabular/mono only for machine strings

## Colors

A cool paper field, navy ink, one drafting-blue route, one copper route, and muted traffic-light status.

### Primary
- **Navy Ink** (`{colors.ink}` / `{colors.primary}`): Body copy, title-block wordmark, 2px chrome rules, filled primary actions, login sheet stroke. The same ink is the type color and the button fill; it is not a neon accent.

### Secondary
- **Drafting API** (`{colors.api}`): The API / inbound edge, edge keys labeled `API`, inbound chart legend, in-page links, and the default toast marker. Dim wash (`{colors.api-dim}`) is selection and faint route fill, never a card background.

### Tertiary
- **Playback Copper** (`{colors.playback}`): The playback / outbound edge, active tab underline, mobile-tab active, input caret, and `:focus-visible` ring. It is a route and a focus signal, not a primary button.

### Neutral
- **Cool Paper** (`{colors.paper}`): Page and title-block canvas; table header wash.
- **Drafting Grid** (`{colors.grid}`): 1px lines on a 20px module; progress-track trough.
- **Sheet** (`{colors.sheet}`): Nodes, sheets, fields, dialogs, toasts — the paper you write on.
- **Secondary Ink** (`{colors.ink-secondary}`): Nav tabs at rest, ghost/secondary labels, mono URLs.
- **Muted Ink** (`{colors.ink-muted}`): Section subs, spec labels, table headers, empty-state body.
- **Hairline / Strong Hairline** (`{colors.hairline}`, `{colors.hairline-strong}`): Internal rules vs control strokes.
- **Hover / Press Wash** (`{colors.surface-hover}`, `{colors.surface-active}`): Row and ghost hover, not glass overlays.

### Status
- **Run Green** (`{colors.green}`): Live LED, running copy, good latency.
- **Quota Amber** (`{colors.orange}`): Warn fills, retry live-badge, mid latency.
- **Fault Red** (`{colors.red}`): Danger buttons, invalid fields, down LED, bad latency.

### Named Rules
**The Two Routes Rule.** API is drafting blue; playback is copper. Never invert them, never paint a third route, never use copper as a CTA fill.

**The Ink Primary Rule.** Primary actions fill with navy ink. Copper and API blue are routes and focus, not buttons.

**The Paper Field Rule.** The canvas is fluorescent paper plus the drafting grid. Sheets sit on the grid; they do not replace it with a solid dark field.

## Typography

**Display Font:** Source Sans 3 (with Noto Sans SC / Source Han Sans SC)
**Body Font:** Source Sans 3 (same pairing)
**Label/Mono Font:** ui-monospace / SF Mono / Cascadia Code / Noto Sans Mono

**Character:** An operate sans: tight but not compressed, bilingual, no display serif and no Inter. Chinese and Latin share one ramp. Machine strings (ports, URLs, bytes) drop to tabular mono.

Root size is 15px; body line-height is 1.5. Page titles track −0.02em.

### Hierarchy
- **Display** (700, 1.7rem, 1.15 / −0.02em): Login wordmark only.
- **Headline** (700, 1.45rem → 1.25rem below 768px, −0.02em): Page titles (`仪表盘`, `站点管理`).
- **Title** (700, 0.95rem): Topo and site node names; sheet/chart heads at 0.9–1.05rem / 600–700.
- **Figure** (700, 1.15–1.25rem, tabular-nums, −0.02em): Spec-strip values and traffic totals — numbers as datasheet figures, not hero metrics.
- **Body** (400–500, 0.85–0.9rem, 1.5): UI copy, table cells, inputs. Empty-state body maxes around 42em.
- **Label** (600, 0.72–0.75rem): Field labels, spec `dt`, table headers, toolbar meta. Not uppercase except edge keys.
- **Edge key** (700, 0.68rem, 0.04em, uppercase): `API` / `播放` / `额外` on the node.
- **Mono** (400, 0.75rem, tabular-nums): Listen ports, upstream URLs, import previews.

### Named Rules
**The Operate Scale Rule.** Source Sans 3 + Noto Sans SC carry every human string. No Inter, no system display face, no all-caps chrome.

**The Tabular Port Rule.** Ports, bytes, uptime, and edge URLs are tabular/mono. Site names stay in the sans.

## Layout

The spatial model is a full-bleed worksheet, not a centered app column. `html` is 15px; the drafting module is 20px; the title block is 52px tall with 28px side padding and a 2px ink bottom rule. Main content starts `52px + 20px` from the top, 28px gutters (24px at 1024px, 16px at 768px), 48px bottom — 176px plus safe-area on small screens to clear the mobile bar.

Rhythm is 8 / 12 / 16 / 24 / 32 / 40. Page headers sit 16px above the first band. Toolbars and control rows use 12px gaps. Site cards auto-fit at `minmax(340px, 1fr)` with 12px gutters; the dashboard node list is a single column with 10px between nodes.

Datasheet bands (dashboard spec strip, traffic totals) are full-width ink-stroked rows, four or three equal columns, collapsing to a stack at 1024px. They belong to those pages. Other pages open with a title, a sub, and a toolbar or native `<select>` row — not a repeated spec strip.

At 768px the top tabs hide and a 44px paper tab bar with a 2px ink top rule takes the bottom edge. Active mobile tab is copper.

### Named Rules
**The Title Block Rule.** Chrome is a 52px paper strip ruled with 2px ink, not a floating glass bar.

**The Page Owns Its First Band Rule.** The four-cell spec strip is the dashboard’s datasheet header. Do not stamp it onto Sites, Traffic, or Diagnostics.

## Elevation & Depth

Rest-state surfaces are flat: white sheet, 1px hairline or ink stroke, no drop shadow. The drafting grid showing around the sheet is the depth cue. Shadows exist only for things that float off the worksheet — account menu, login card, toast, modal (`--shadow-menu` / `--shadow-modal`). Modal scrim is navy ink at 45% (`rgba(28, 43, 54, 0.45)`). Focus is a copper ring (`0 0 0 2px paper, 0 0 0 4px playback`), never a glow.

### Shadow Vocabulary
- **Menu / toast / login** (`0 1px 2px rgba(17,17,17,0.05), 0 4px 8px rgba(17,17,17,0.05), 0 8px 24px rgba(17,17,17,0.04)`): Overlay panels only.
- **Modal** (`0 2px 4px rgba(17,17,17,0.06), 0 8px 16px rgba(17,17,17,0.06), 0 16px 48px rgba(17,17,17,0.06)`): Native `<dialog>` sheet only.

### Named Rules
**The Stroke, Not Glow Rule.** Depth is a 1px ink or hairline stroke. No blur, no neon edge, no glass blur.

**The Overlay-Only Shadow Rule.** Rest-state nodes, sheets, and spec bands have no box-shadow. Shadow is a response to leaving the page plane.

## Shapes

Drafting geometry: tight corners, square datasheet bands, square LEDs-as-dots only for live state.

Buttons 3px, ghosts and pills 2px, nodes/sheets/fields 4px, dialog 6px. Nothing above 6px. Spec strip, traffic totals, and empty/error banners omit radius — they read as cut paper, not cards.

Borders are 1px. Title block and mobile bar use a 2px ink rule. Edge lines are 2px with a 6×6 square cap, not arrows. Progress is a 4px rectangular track. Status is an 8px circle LED; pills are outline rectangles, not stadium chips.

### Named Rules
**The Drafting Corner Rule.** Corners live between 2px and 6px. No 12px+ cards, no 100px pills.

## Components

Controls feel like labeled instruments on a worksheet: ink fill for the one commit, hairline ghosts for everything else, native fields with a copper caret.

### Buttons
- **Shape:** Primary/secondary 3px (`{rounded.sm}`), 40px min height, 10px 16px, 0.85rem / 600. Ghost 2px, 36px, 8px 12px, 0.8rem / 500.
- **Primary:** Navy fill, sheet text. Hover `#0f1a22`. Used for login, add site, dialog confirm.
- **Ghost:** Sheet fill, hairline border, secondary ink. Hover wash + stronger hairline. Danger ghost tints red on hover.
- **Secondary:** Sheet fill, strong hairline, 40px — cancel and empty-state alternate.
- **Danger:** Fault red fill, sheet text.
- **Focus:** Copper ring. Disabled at 50% opacity.

### Chips
- **Pills:** Transparent, 2px corners, 1px currentColor, 2px 8px, 0.7rem / 600. API blue / green / orange / muted hairline. UA mode and playback-mode tags only.
- **Status LED:** 8px circle, green filled when running, hollow muted when stopped. Not a pill.

### Cards / Containers
- **Topo node / site card:** 4px, strong hairline (ink when traced), 14px 16px 12px, end-label rule under the name.
- **Sheet / chart wrap:** 4px, 1px ink, overflow clipped. Header 12px 16px then hairline. Tables: paper header cells, 0.72rem labels, 0.85rem body, hover wash.
- **Datasheet band:** Ink rectangle, no radius, internal hairline columns.
- **Empty / error:** Ink rectangle on paper, left-aligned, ~28px 20px, no illustration hero.

### Inputs / Fields
- **Style:** Sheet fill, 1px strong hairline, 3px corners, 44px min height, 10px 12px, 0.9rem.
- **Label:** 0.75rem / 600 secondary ink, 6px above the control.
- **Hover:** Border to secondary ink.
- **Focus:** Playback copper border + copper caret. No glow.
- **Invalid:** Red border, red-dim wash, 0.75rem red error line.
- **Select:** Same geometry; native control, ink chevron.

### Navigation
- **Title block:** 52px paper+grid, 2px ink bottom, StreamDock wordmark 1rem / 700, 22px mark.
- **Tabs:** 0.85rem / 500 secondary ink, 16px horizontal padding. Active: ink, 600, 2px copper underline that shares the chrome rule.
- **Account:** 32px 2px-radius ink square with sheet initial.
- **Mobile bar:** Same paper+grid, 2px ink top, 44px min tabs, 0.62rem labels, 22px stroke icons. Active is copper.

### Spec strip
Dashboard-only datasheet header: four equal cells (站点 / 运行 / 总流量 / 运行时长). Muted 0.72rem `dt`, 1.15rem tabular `dd`. Not four separate metric cards, not a page template.

### Topo node + two edges
Signature component. White node, name + `:{port}` + LED on an end-label. Two (or three) edges: uppercase key, 2px ruled line with a square cap, mono destination. API blue / playback copper. At rest the line is 22% opacity; traced or focused (and always on site cards) the line goes to 100% over 180ms. First dashboard node may ship traced; that is a page default, not a layout law.

### Dialog
Native `<dialog>`. 6px sheet, 1px ink, modal shadow, max 480px. Header 16px 20px 0 with 1.05rem title; body 16px 20px; footer right-aligned primary + secondary. Scrim is ink at 45%. Close control is a 32px paper square.

## Do's and Don'ts

### Do:
- **Do** put every screen on the fluorescent paper grid and name the product StreamDock.
- **Do** draw each site as a node with an API edge and a playback edge in the two named route colors.
- **Do** fill the one primary action with navy ink; keep add/export on Sites as primary + ghosts.
- **Do** use 2–6px corners, 1px strokes, and copper only for playback + focus + active tab.
- **Do** pair Source Sans 3 with Noto Sans SC at the 15px operate scale; mono for ports and URLs.

### Don't:
- **Don't** revive The Dark Glass Console: cinema black, Inter, `#0A84FF`, glass bars, glow, or blur.
- **Don't** print a Meridian wordmark anywhere in chrome, login, or examples.
- **Don't** replace the spec strip with four hero-metric cards, or clone the dashboard strip onto every page.
- **Don't** use copper or API blue as a primary button fill.
- **Don't** round past 6px, stadium-pill a chip, or drop a rest-state shadow on a node or sheet.
