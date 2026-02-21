# Frontend - Microservices E-commerce Demo

Next.js 15 frontend for the microservices e-commerce platform.

## Tech Stack

- **Next.js 15** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first CSS
- **shadcn/ui** - UI components built on Radix UI
- **TanStack Query** - Data fetching and caching
- **React Hook Form** - Form management
- **Zod** - Schema validation
- **Axios** - HTTP client

## Getting Started

### Install Dependencies

```bash
cd frontend
npm install
```

### Environment Variables

Create a `.env.local` file:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Development

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build

```bash
npm run build
npm start
```

### Docker

```bash
docker build -t frontend:latest .
docker run -p 3000:3000 frontend:latest
```

## Project Structure

```
frontend/
├── src/
│   ├── app/              # Next.js App Router pages
│   │   ├── (store)/      # Public store pages
│   │   ├── (dashboard)/  # User dashboard
│   │   ├── (admin)/      # Admin dashboard
│   │   ├── (auth)/       # Authentication pages
│   │   ├── layout.tsx    # Root layout
│   │   └── page.tsx      # Homepage
│   ├── components/
│   │   └── ui/           # shadcn/ui components
│   └── lib/
│       └── utils.ts      # Utility functions
├── public/               # Static files
└── package.json
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm start` - Start production server
- `npm run lint` - Run ESLint

## API Integration

The frontend communicates with the API Gateway at `http://localhost:8080`.

See `NEXT_STEPS.md` in the root directory for implementation details.
