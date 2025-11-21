# GoNews Frontend

This is the frontend application for the GoNews project, built with [Next.js](https://nextjs.org) and [Tailwind CSS](https://tailwindcss.com).

## Features

- **News Feed**: Displays a grid of news items fetched from the Go backend.
- **News Detail**: Click on any news item to read the full content on a dedicated page.
- **Refresh**: Manually trigger an RSS feed update from the UI.
- **Responsive Design**: Optimized for both desktop and mobile devices.
- **Dark Mode**: Supports system preference for dark/light mode.

## Prerequisites

- Node.js 18+ installed.
- GoNews Backend running on `http://localhost:8080`.

## Getting Started

1.  **Install dependencies:**

    ```bash
    npm install
    ```

2.  **Run the development server:**

    ```bash
    npm run dev
    ```

3.  **Open the application:**

    Visit [http://localhost:3000](http://localhost:3000) in your browser.

## Configuration

The application is configured to proxy API requests to the backend.
See `next.config.ts` for the proxy configuration:

```typescript
async rewrites() {
  return [
    {
      source: '/api/:path*',
      destination: 'http://localhost:8080/:path*',
    },
  ];
}
```

## Project Structure

- `src/app/page.tsx`: Main news feed page.
- `src/app/news/[id]/page.tsx`: News detail page.
- `src/app/components/`: Reusable UI components (`NewsCard`, `RefreshButton`).
