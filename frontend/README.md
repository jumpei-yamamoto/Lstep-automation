# Frontend - Lstep Automation

This is the frontend application built with Next.js, TypeScript, and Tailwind CSS according to the specifications in CLAUDE.md.

## Getting Started

1. Install dependencies:
```bash
npm install
```

2. Copy environment variables:
```bash
cp .env.local.example .env.local
```

3. Run the development server:
```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Project Structure

- `app/` - Next.js App Router pages and layouts
- `components/` - Reusable React components
- `lib/` - Utility functions and API client
- `services/` - API service layer (UI components should use these, not direct fetch)
- `styles/` - Global CSS styles

## Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run type-check` - Run TypeScript type checking

## Architecture Guidelines

This project follows the guidelines specified in CLAUDE.md:

- **UI Components**: No direct API calls, use services layer
- **State Management**: Recoil for client state, TanStack Query for server state
- **Styling**: Tailwind CSS utilities
- **API Communication**: HTTP-only cookies for authentication

## Technologies

- **Next.js 15** with App Router
- **TypeScript** for type safety
- **Tailwind CSS** for styling
- **Recoil** for state management
- **TanStack Query** for server state management
- **ESLint** for code linting