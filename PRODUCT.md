# Product

## Register

product

The **SPA** (`frontend/`) is the primary product surface. **`docs/index*.html`** is the marketing landing (brand register): longer-form copy and screenshots, but it shares the same palette and restraint as the app.

## Users

Weekly croquis (figure-drawing) meetup participants: an admin who runs the session and joiners who follow along from their own devices. They are in a relaxed group setting—studio, café, or home call—and need synchronized reference photos and timers without screen sharing or video overhead.

## Product Purpose

Croquing lets a drawing group share one lobby URL, pick reference images together (via Pixabay), and draw in timed rounds with a server-authoritative countdown. Success means everyone sees the same photo at the same time, trusts the timer, and stays in flow from lobby setup through the last round.

## Brand Personality

Warm, focused, communal. The tool should feel like a calm studio assistant—not a flashy SaaS landing page. Confidence through reliability; delight in the drawing moment, not in decorative chrome. The marketing landing (`docs/`) follows the same restraint: no hero-metric bars, glass chrome, or cream-paper aesthetics.

## Anti-references

- Generic AI warm-cream marketing pages with glass cards and eyebrow kickers
- Hero-metric dashboards and identical icon-card grids
- Over-animated onboarding that delays getting into a session
- Modal-first flows for actions that belong inline
- Screen-share or video-first drawing tools (out of scope)

## Design Principles

1. **Task first** — Every screen serves the current phase (lobby, select, draw, break). Decoration never competes with the reference photo or timer.
2. **Server truth** — Timers and phase state come from the server; the UI reflects snapshots, it does not invent timing.
3. **Earned familiarity** — Standard controls, predictable layout, consistent button and form vocabulary across phases.
4. **Inclusive by default** — Touch-friendly targets, i18n for all user-facing and assistive strings, reduced-motion support.
5. **Quiet between rounds, bold during draw** — Lobby chrome stays restrained; the drawing surface gets maximum space and contrast.

## Accessibility & Inclusion

- Target WCAG 2.1 AA for lobby and drawing flows
- Minimum 44×44 px touch targets on interactive controls
- `prefers-reduced-motion` respected globally
- Five locales: en, ko, ja, pl, zh — all visible and aria copy localized
- High contrast on the drawing surface (dark stage, light timer chrome)
