---
name: StreamDock
description: Dark glass operator console for the Emby reverse-proxy admin panel.
colors:
  primary: "#0A84FF"
  primary-glow: "rgba(10,132,255,0.35)"
  cinema-black: "#000000"
  glass: "rgba(30,30,30,0.72)"
  surface: "rgba(255,255,255,0.06)"
  surface-hover: "rgba(255,255,255,0.10)"
  surface-active: "rgba(255,255,255,0.14)"
  hairline: "rgba(255,255,255,0.08)"
  ink: "rgba(255,255,255,0.87)"
  ink-secondary: "rgba(255,255,255,0.60)"
  ink-muted: "rgba(255,255,255,0.38)"
  status-green: "#30D158"
  status-teal: "#64D2FF"
  status-orange: "#FF9F0A"
  status-red: "#FF453A"
  status-purple: "#BF5AF2"
  white: "#FFFFFF"
typography:
  display:
    fontFamily: "Inter, -apple-system, system-ui, sans-serif"
    fontSize: "2.4rem"
    fontWeight: 700
    lineHeight: 1.1
    letterSpacing: "-0.04em"
  headline:
    fontFamily: "Inter, -apple-system, system-ui, sans-serif"
    fontSize: "1.5rem"
    fontWeight: 700
    lineHeight: 1.2
    letterSpacing: "-0.03em"
  title:
    fontFamily: "Inter, -apple-system, system-ui, sans-serif"
    fontSize: "1.05rem"
    fontWeight: 600
    lineHeight: 1.3
    letterSpacing: "-0.01em"
  body:
    fontFamily: "Inter, -apple-system, system-ui, sans-serif"
    fontSize: "0.9rem"
    fontWeight: 400
    lineHeight: 1.5
    letterSpacing: "0em"
  label:
    fontFamily: "Inter, -apple-system, system-ui, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 500
    lineHeight: 1.4
    letterSpacing: "0.06em"
  mono:
    fontFamily: "ui-monospace, SF Mono, Fira Code, Cascadia Code, monospace"
    fontSize: "0.75rem"
    fontWeight: 400
    lineHeight: 1.4
    letterSpacing: "0em"
rounded:
  xs: "8px"
  sm: "12px"
  md: "16px"
  lg: "22px"
  pill: "100px"
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
    textColor: "{colors.white}"
    rounded: "{rounded.sm}"
    padding: "10px 20px"
  button-primary-hover:
    backgroundColor: "#0077ED"
    textColor: "{colors.white}"
    rounded: "{rounded.sm}"
    padding: "10px 20px"
  button-secondary:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "10px 20px"
  button-ghost:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink-secondary}"
    rounded: "{rounded.xs}"
    padding: "9px 12px"
  button-danger:
    backgroundColor: "{colors.status-red}"
    textColor: "{colors.white}"
    rounded: "{rounded.sm}"
    padding: "10px 20px"
  input:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "12px 16px"
    height: "44px"
  card:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "20px 24px"
  chip:
    backgroundColor: "rgba(10,132,255,0.15)"
    textColor: "{colors.primary}"
    rounded: "{rounded.pill}"
    padding: "4px 12px"
  nav-link:
    backgroundColor: "transparent"
    textColor: "{colors.ink-secondary}"
    rounded: "{rounded.xs}"
    padding: "8px 16px"
  nav-link-active:
    backgroundColor: "{colors.surface-hover}"
    textColor: "{colors.white}"
    rounded: "{rounded.xs}"
    padding: "8px 16px"
---

# Design System: StreamDock

## Overview

**Creative North Star: "The Dark Glass Console"**

StreamDock is an Operate surface: a homelab operator sits in a dark room and drives reverse-proxy sites. The incumbent stylesheet already named this world ("Apple TV Design System") — cinema black, hairline glass, Inter, one system-blue accent. This file carbonizes that world. It does not replace it.

The console should disappear into the task. Familiar controls, dense data, and status color do the work. Brand lives in the wordmark, the favicon stroke, and the restraint of a single accent — not in invented widgets, gradient type, or page-load choreography. Product name is **StreamDock**. GitHub repository remains `tsumon/Meridian-merged`.

**Key Characteristics:**
- Cinema-black canvas with translucent glass chrome
- One accent (system blue) on primary action, current nav, and focus
- Semantic status colors only — green / teal / orange / red / purple never decorate
- Inter for UI; tabular/mono only for addresses, ports, and bytes
- Surfaces stay flat at rest; motion is 150–250ms and state-only
- Empty, loading, and error are first-class layouts

## Colors

A restrained dark palette: black sea, translucent ice, one system-blue voice, and a small status vocabulary.

### Primary
- **System Blue**: used for the single primary action on a screen, the current navigation mark, focus rings, inbound-traffic series, and "on" affordances. Rarity is the point.

### Neutral
- **Cinema Black**: page canvas. This product is a night console; the incumbent chose true black and it stays.
- **Glass**: frosted top and bottom chrome over scrolling content (`backdrop-filter` belongs here, not on every card).
- **Surface / Surface Hover / Surface Active**: translucent white lifts for cards, rows, and rest/hover/press.
- **Hairline**: 1px separators and card strokes.
- **Ink / Ink Secondary / Ink Muted**: primary text, secondary labels, placeholders and meta. Body text is ink, never full white.

### Named Rules
**The One Accent Rule.** System blue is for primary action, current selection, and focus. If a screen has two solid blue buttons, one of them is wrong.

**The Status Vocabulary Rule.** Green is running/healthy/success. Teal is outbound/secondary series. Orange is warning or quota pressure. Red is error or destructive. Purple is diagnostic commentary only. A color that means nothing gets deleted.

## Typography

**Display Font:** Inter (with -apple-system, system-ui)
**Body Font:** Inter (same family — product UI does not pair a display face)
**Label/Mono Font:** SF Mono, Fira Code, Cascadia Code for addresses, ports, probes, and JSON

**Character:** One well-tuned grotesque. Hierarchy comes from size and weight, not from a second family. Operate copy is Chinese; do not Title-Case or uppercase body sentences.

### Hierarchy
- **Display** (700, 2.4rem, -0.04em): login wordmark only.
- **Headline** (700, 1.5rem, -0.03em): page titles (`仪表盘`, `站点管理`).
- **Title** (600, 1.05rem): card and modal titles.
- **Body** (400, 0.9rem, 1.5): forms, table cells, help text. UI copy can run denser than 65ch.
- **Label** (500, 0.75rem, 0.06em, uppercase English keys only): form labels and table headers. Chinese labels stay sentence case; do not force `letter-spacing` on CJK.

### Named Rules
**The One Family Rule.** Inter carries headings, buttons, labels, and body. Mono is for machine strings, never as a "technical" costume on human copy.

## Layout

Fixed 15px root. Content column max-width 1400px, centered. Desktop chrome is a 56px glass top bar; pages pad 24px 32px under it. Mobile keeps the top bar for brand + account, hides the text links, and uses a 4-tab bottom bar with 44px hit targets. Breakpoints: 1024px (stats 2-up, diag 1-up), 768px (bottom tabs, sites 1-up), 480px (stats 1-up).

Density is structural: tighter card padding (20px 24px), compact tables, no decorative vertical gaps. Responsive behavior collapses columns and chrome — it does not fluid-scale type.

URL hash is the page state (`#dashboard` `#sites` `#traffic` `#diagnostics`).

## Elevation & Depth

Dark mode elevates by tint and hairline, not by drop shadow. Cards sit on translucent surface; they do not lift on hover. Shadows are reserved for layers that leave the page plane: modal, toast, account menu.

Glass blur is a chrome material (nav sitting over scrolling content), not a card effect.

### Shadow Vocabulary
- **Menu** (`0 8px 32px rgba(0,0,0,0.45), 0 0 0 1px rgba(255,255,255,0.08)`): account panel.
- **Toast** (same hairline + darker wash): transient status.
- **Modal panel** (`0 16px 48px rgba(0,0,0,0.55), 0 0 0 1px rgba(255,255,255,0.08)`): dialog surface over the dimmed canvas.
- **Focus ring** (`0 0 0 3px rgba(10,132,255,0.35)`): keyboard only (`:focus-visible`).

### Named Rules
**The Flat Rest Rule.** Surfaces do not translate or grow a shadow on hover. Hover is a background-color or border-color shift, gated by `@media (hover: hover)`.

## Shapes

Soft Apple radii, used on purpose: 8px for compact controls and nav pills, 12px for buttons and inputs, 16px for cards, 22px for the dialog panel, full pill only for status chips. Nested radius steps inward.

Hairline 1px strokes, never a thick accent bar on the side of a card. Status LEDs are 7px circles. Avatars are circles.

## Components

Familiar admin controls. Native `<select>`, native `<dialog>`, native file picker. Do not invent a scrollbar, a custom dropdown, or a non-standard modal.

### Buttons
- **Shape:** 12px radius, shared across primary / secondary / danger. Ghost actions on cards use 8px.
- **Primary:** system blue fill, white label, 10px 20px, font-weight 600. One per screen region.
- **Secondary / Ghost:** surface fill, ink-secondary label. Danger ghost turns red on hover.
- **Hover / Focus:** background darkens; `:active` scales to 0.98. `:focus-visible` uses the focus ring. Disabled is 0.5 opacity and `not-allowed`.
- **Loading:** keep the button width, replace label with `…` / `正在保存`, `aria-busy="true"`. No spinner under 800ms.

### Chips
- **Style:** tinted pill, 4px 12px, 0.7rem 600. Blue = Infuse, green = Web, orange = 客户端. Playback mode uses the same chip, not a second shape.
- **State:** chips are readouts, not filters.

### Cards / Containers
- **Corner Style:** 16px
- **Background:** surface
- **Shadow Strategy:** none at rest (Flat Rest Rule)
- **Border:** hairline
- **Internal Padding:** 20px 24px
- Stat cards, site cards, diag cards, and chart wells share this shell. Nested cards are wrong.

### Inputs / Fields
- **Style:** surface fill, hairline stroke, 12px radius, 44px min height.
- **Focus:** border shifts toward primary at 50% alpha + focus ring. No scale.
- **Error:** 1px status-red border, status-red message under the field (13px). Validate on blur and submit, not on each keystroke.
- **Disabled:** darker wash, ink-muted text.
- Labels sit above fields. Placeholder is not a label. Help text is ink-muted, 0.75rem.

### Navigation
- Desktop: glass top bar, wordmark + four text links + account. Active link is surface-hover + white.
- Mobile: same top bar without the links; four icon+label tabs on a glass bottom bar. Active tab is system blue.
- Account is a circle avatar that opens a small menu (username + 退出登录). Avatar is not itself a logout button.

### Tables
- Uppercase 0.7rem headers in ink-muted. Rows 16px 24px, hairline dividers, hover wash. Tabular nums on bytes and latency. Empty, loading (skeleton rows), and error (retry) are real table states.

### Dialog
- Native `<dialog>`. 22px panel, 480px max, dimmed backdrop. Esc closes. Destructive confirms use this dialog, never `window.confirm`.

### Toast
- Top-right under the nav. Glass panel, 12px radius, 3px semantic left edge. Copy has no emoji. `aria-live="polite"`.

### Empty / Loading / Error
- Empty names the gap and the next action (添加站点, 导入配置).
- Loading uses layout-matching skeletons after 800ms; button labels change immediately.
- Error names the failure and offers 重试. Empty is not error.

### Signature: site card
The site card is the product: name, running LED, main upstream, playback upstreams + mode, listen port, UA chip, quota bar, then 停用 / 编辑 / 删除. Playback fields and 导出配置 / 导入配置 stay visible — they are fork-unique, not optional chrome.

## Do's and Don'ts

### Do:
- **Do** keep the cinema-black canvas, Inter, system blue, and glass chrome.
- **Do** use one primary button per region and put 导出 / 导入 next to 添加站点.
- **Do** write Chinese empty/error copy that names the next step (`还没有站点。添加一个反代，或导入 JSON 备份。`).
- **Do** keep hash routes, site CRUD, `playback_target_url` / `playback_mode` / `stream_hosts`, and plaintext JSON import/export.
- **Do** theme selection, caret, focus rings, and scrollbars from this palette.
- **Do** honor `prefers-reduced-motion`: drop translate/scale, keep short opacity.

### Don't:
- **Don't** introduce a second visual world, a new brand name, or a marketing landing layout.
- **Don't** use gradient text, emoji as icons, or Lucide-as-a-package.
- **Don't** lift cards, scale from 0, or run page-load stagger.
- **Don't** hide scrollbars or restyle native select into a custom combobox.
- **Don't** add React, Vue, or a build step.
- **Don't** treat "暂无数据" as both empty and error.
- **Don't** apply Duolingo / Taste / Frontend Design defaults. This file wins.
