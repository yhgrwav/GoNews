import { NewsCard } from './components/NewsCard';
import { RefreshButton } from './components/RefreshButton';

interface NewsItem {
  ID: number;
  Title: string;
  Content: string;
  PubTime: string;
  Link: string;
  Source: string;
}

async function getNews(): Promise<NewsItem[]> {
  try {
    // Fetch from backend directly since this is server-side
    const apiUrl = process.env.API_URL || 'http://localhost:8080';
    const res = await fetch(`${apiUrl}/news?limit=50`, {
      cache: 'no-store',
      next: { revalidate: 0 }
    });

    if (!res.ok) {
      console.error('Failed to fetch news, status:', res.status);
      return [];
    }
    return res.json();
  } catch (error) {
    console.error('Error fetching news:', error);
    return [];
  }
}

export default async function Home() {
  const news = await getNews();

  return (
    <main className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-extrabold text-gray-900 dark:text-white">
            GoNews Feed
          </h1>
          <RefreshButton />
        </div>

        {news.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500 dark:text-gray-400 text-lg">
              No news available. Make sure the backend is running and try refreshing the feed.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {news.map((item) => (
              <NewsCard key={item.ID} news={item} />
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
