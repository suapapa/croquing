---
name: Croquing
description: Real-time croquis meetup UI — warm product surface, dark immersive drawing mode
colors:
  bg: "#fafafa"
  bg-accent: "#f4f4f5"
  surface: "#ffffff"
  text: "#1a1614"
  muted: "#5c534a"
  border: "#ddd4c8"
  accent: "#c2410c"
  accent-hover: "#9a3412"
  accent-soft: "#fff7ed"
  error: "#b91c1c"
  error-bg: "#fef2f2"
  success: "#15803d"
  drawing-bg: "#0c0a09"
  drawing-text: "#fafaf9"
  thumb-bg: "#e7e5e4"
  grid-bg: "#f5f5f4"
typography:
  body:
    fontFamily: "Manrope, Segoe UI, system-ui, sans-serif"
    fontSize: "0.9375rem"
    fontWeight: 400
    lineHeight: 1.5
  display:
    fontFamily: "Syne, Georgia, Times New Roman, serif"
    fontSize: "clamp(1.5rem, 5vw, 2.2rem)"
    fontWeight: 700
    lineHeight: 1.1
    letterSpacing: "-0.02em"
rounded:
  sm: "0.5rem"
  md: "0.75rem"
  lg: "1rem"
  pill: "999px"
spacing:
  touch-min: "2.75rem"
  content-max: "72rem"
components:
  button-primary:
    backgroundColor: "{colors.accent}"
    textColor: "#ffffff"
    rounded: "999px"
    height: "{spacing.touch-min}"
  button-secondary:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.text}"
    rounded: "999px"
    height: "{spacing.touch-min}"
  phase-panel:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.text}"
    rounded: "{rounded.lg}"
    padding: "1.5rem"
---

## Overview

Croquing uses a **product register** in the SPA and a **brand register** on the marketing landing (`docs/index*.html`). Both share terracotta accent and near-neutral backgrounds. The SPA uses `frontend/src/index.css`; the landing uses `docs/landing.css` with matching token names (`--color-bg`, `--color-accent`, etc.). Syne carries display headings; Manrope carries body text.

## Colors

| Role | Token | Usage |
|------|-------|--------|
| Page bg | `--color-bg` | Lobby, home shell, marketing landing |
| Surface | `--color-surface` | Panels, inputs |
| Ink | `--color-text` | Primary copy |
| Muted | `--color-muted` | Secondary copy (≥4.5:1 on surfaces) |
| Accent | `--color-accent` | Primary actions, admin badge, timer bar |
| Drawing | `--color-drawing-bg` / `--color-drawing-text` | Full-screen draw phase |

Accent is terracotta (`#c2410c`), used for actions and state—not decorative fills. Body backgrounds stay near-neutral (`#fafafa`), not saturated cream.

## Typography

- **Body:** Manrope 400–700, 15 px base, 1.5 line-height
- **Display:** Syne 600–700 for h1/h2 in lobby and home only—not buttons or data
- **Scale:** Fixed rem steps (product register); no fluid heading clamp in app chrome
- **Prose:** `text-wrap: pretty` on lead paragraphs; max ~30 rem on marketing lead

## Elevation

- **Panels:** 1 px border (`--color-border`) OR `--shadow-sm`, not both
- **Drawing timer:** Thin top bar with accent gradient; optional glow only in urgent state
- **Modals:** Solid surface + `--shadow-lg`; no bounce easing on open
- **Z-index:** drawing controls (10) → modal backdrop (native dialog) → selection dock (50)

## Components

- **Buttons:** Pill shape, `--touch-min` height, primary/secondary/icon-only variants
- **Phase panels:** `.phase-panel` — padded surface card for lobby phases
- **Pixabay grid:** Square aspect-ratio tiles, 2 px selection ring
- **Selection dock:** Fixed bottom bar when photos selected; solid surface, scrollable thumbs
- **Marketing landing:** `.btn-primary`, `.tab-nav` / `.tab-btn`, `.features-list` — border-first, no glass or metric hero bars
- **Drawing panel:** Fixed inset, black stage, `object-fit: contain` photo, footer attribution

## Do's and Don'ts

**Do**

- Use CSS variables for color; extend `:root` for new semantic roles
- Keep drawing mode minimal—photo, timer, attribution
- Localize all `aria-label`, `alt`, and visible strings via `t()`
- Honor `prefers-reduced-motion: reduce`

**Don't**

- Glassmorphism, cream paper backgrounds, or uppercase eyebrow kickers on every section
- Pair 1 px border with wide soft shadow on the same card
- Bounce/elastic easing on modals or decorative infinite animations in lobby chrome
- Custom scrollbars or touch targets below 44 px
