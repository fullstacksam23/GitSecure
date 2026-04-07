# GitSecure Frontend

Production-oriented React frontend for the GitSecure backend. It provides:

- Dashboard overview for vulnerability totals and recent scans
- Scan creation flow for GitHub repositories
- Vulnerability explorer with filtering and pagination
- Detailed vulnerability record view

## Stack

- React + Vite
- JavaScript
- Tailwind CSS
- shadcn-style UI primitives with Radix Dialog

## Setup

1. Install dependencies:

```bash
npm install
```

2. Create an environment file:

```bash
cp .env.example .env
```

3. Start the app:

```bash
npm run dev
```

The frontend expects the Go API at `http://localhost:3000` by default.

## Backend Routes Used

Implemented today:

- `POST /scan`
- `GET /health`

Required for the full dashboard experience:

- `GET /dashboard/summary`
- `GET /scans`
- `GET /scans/:jobId`
- `GET /vulnerabilities`
- `GET /vulnerabilities/:id`

The exact request and response shapes expected by the frontend are documented in [docs/frontend-api-contract.md](./docs/frontend-api-contract.md).
