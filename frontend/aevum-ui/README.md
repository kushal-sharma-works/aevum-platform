# Aevum UI

Frontend application for deterministic event timelines, decision inspection, replay, diff, and audit workflows.

## Run

```bash
npm install
npm run dev
```

## Scripts

- `npm run dev` - start Vite dev server
- `npm run build` - type-check + production build
- `npm run test` - run unit tests
- `npm run test:coverage` - coverage report
- `npm run lint` - lint TS and Vue files
- `npm run type-check` - strict type validation

## Architecture

`router -> views -> components -> stores -> api -> backend services`

## Environment

- Vite dev server runs on `http://localhost:3000`
- Service proxy routes:
	- `/api/events` -> Event Timeline service (`:8080`)
	- `/api/decisions` -> Decision Engine (`:8081`)
	- `/api/query` -> Query & Audit (`:8082`)

## Component Catalog

- Timeline Viewer (virtualized list)
- Decision Inspector and Decision Trace
- Replay Console (SSE)
- Rule browser + condition builder
- Audit trail and diff views
- Dashboard overview with live metrics placeholders
