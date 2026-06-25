# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## ⚠️ This is Next.js 16 — NOT the Next.js you know

Next.js 16 (`next@16.2.9` per `package.json`) has breaking changes from the version in most training data. **Before writing any Next.js code, read the relevant guide in `node_modules/next/dist/docs/`.** APIs, conventions, file structure, and metadata may all differ. Heed deprecation notices — do not rely on memorized patterns.

## What this is

The frontend for the shop. Currently a stock `create-next-app` scaffold (App Router, TypeScript, Tailwind v4) — no product UI yet. The backend lives in `../server/` and exposes `GET /products` returning the catalogue from PostgreSQL.

## Build / run / lint

All commands run from this `client/` directory.

```sh
# Dev server (http://localhost:3000)
npm run dev

# Production build
npm run build

# Run the production build
npm start

# Lint (ESLint via eslint-config-next — core-web-vitals + TypeScript)
npm run lint
```

There are no tests yet.

## Layout

```
client/
├── next.config.ts          # empty NextConfig; no custom config options
├── tsconfig.json           # strict TS, path alias `@/*` → ./src/*
├── eslint.config.mjs       # uses next/core-web-vitals + next/typescript
├── postcss.config.mjs      # @tailwindcss/postcss plugin only
├── .prettierrc             # 80 cols, single quotes, semi, trailingComma es5
├── src/app/                # App Router — layout.tsx, page.tsx
└── public/                 # static assets (file.svg, globe.svg, next.svg, vercel.svg, window.svg)
```

## Architectural notes

- **Next.js 16 + React 19.** Run `node_modules/next/dist/docs/` checks before any feature work — the App Router entry points, route handlers, and metadata APIs are not what older docs describe.
- **App Router only.** All routes go under `src/app/`. No `pages/` directory. Server Components by default; mark Client Components with `"use client"` only when needed.
- **Path alias:** `@/*` resolves to `./src/*` (see `tsconfig.json`).
- **Styling is Tailwind v4** via the `@tailwindcss/postcss` plugin. There is no `tailwind.config.*` file — Tailwind v4 uses CSS-first config. No `globals.css` exists in `src/app/` yet beyond what `create-next-app` set up; check before adding theme tokens.
- **Fonts** are loaded via `next/font/google` (Geist + Geist Mono) in `src/app/layout.tsx` and exposed as `--font-geist-sans` / `--font-geist-mono` CSS variables.
- **Linting:** `eslint-config-next` (core-web-vitals + typescript). `eslint-plugin-prettier` is wired in via dev deps but not yet added to `eslint.config.mjs` — add it there if you want `npm run lint` to surface Prettier issues.
- **Formatting:** Prettier config in `.prettierrc` (single quotes, semi, 80 cols, LF, `arrowParens: "avoid"`).
- **Env files are gitignored** (`.env*` in `.gitignore`). Do not commit secrets. When wiring the backend, decide on a variable like `NEXT_PUBLIC_API_URL` and document it here.

## Connecting to the backend

The Go server (`../server/`) is the source of truth for products. CORS is **not yet set up** server-side — before the client can call `/products` from the browser, middleware for CORS needs to be added in `server/internal/routes/routes.go` (or the client can proxy through a Next.js Route Handler / rewrite to avoid CORS). Pick one approach and update this file with the decision.