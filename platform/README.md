# Ethereum Validator API Platform

The frontend interface for the Ethereum Validator API, built with Next.js and React. This platform provides a modern, responsive interface for interacting with the validator API endpoints.

## Project Overview

This frontend application is designed to demonstrate:
- Clean architecture and component organization
- Modern React patterns and best practices
- Type-safe development with TypeScript
- Responsive design with TailwindCSS
- Production-ready deployment configuration

## Getting Started

### Prerequisites

- Node.js (v18 or higher)
- npm or yarn
- Docker (optional, for containerized deployment)

### Local Development

1. Clone the repository and install dependencies:
```bash
cd platform
npm install
```

2. Configure environment:
```bash
cp .env.example .env.local
# Update NEXT_PUBLIC_API_URL in .env.local
```

3. Start development server:
```bash
npm run dev
```

The application will be available at `http://localhost:3003`.

### Production Deployment

Build and start the production server:
```bash
npm run build
npm start
```

### Docker Deployment

```bash
docker build -t ethereum-validator-platform .
docker run -p 3003:3003 ethereum-validator-platform
```

## Architecture

The application follows a clean architecture pattern:

- `app/` - Next.js pages and components
- `components/` - Reusable UI components
- `lib/` - Utility functions and API clients
- `public/` - Static assets
- `styles/` - Global styles and Tailwind configuration

## API Integration

The platform integrates with two main API endpoints:

### Sync Committee Duties
```typescript
GET /api/syncduties/{slot}
```

Returns validator sync committee assignments for a given slot.

### Block Rewards
```typescript
GET /api/blockreward/{slot}
```

Returns block reward information including MEV status.

## Technology Stack

- **Framework**: Next.js 15.2.4
- **UI Library**: React 19
- **Styling**: TailwindCSS
- **Language**: TypeScript
- **Development Tools**: ESLint, PostCSS

## Testing

Run the test suite:
```bash
npm test
```