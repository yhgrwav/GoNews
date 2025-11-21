import React from 'react';

interface NewsItem {
    ID: number;
    Title: string;
    Content: string;
    PubTime: string;
    Link: string;
    Source: string;
}

async function getNewsItem(id: string): Promise<NewsItem | null> {
    try {
        const res = await fetch(`http://localhost:8080/news/${id}`, {
            cache: 'no-store',
            next: { revalidate: 0 }
        });
        if (!res.ok) return null;
        return res.json();
    } catch (error) {
        console.error('Error fetching news item:', error);
        return null;
    }
}

export default async function NewsPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const news = await getNewsItem(id);

    if (!news) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
                <div className="text-center">
                    <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">News not found</h1>
                    <a href="/" className="text-blue-600 hover:underline">Back to feed</a>
                </div>
            </div>
        );
    }

    return (
        <main className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
            <article className="max-w-3xl mx-auto bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
                <div className="p-8">
                    <div className="flex items-center justify-between mb-6">
                        <span className="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-sm font-semibold rounded-full">
                            {news.Source}
                        </span>
                        <time className="text-gray-500 dark:text-gray-400 text-sm">
                            {new Date(news.PubTime).toLocaleString()}
                        </time>
                    </div>

                    <h1 className="text-3xl sm:text-4xl font-extrabold text-gray-900 dark:text-white mb-6 leading-tight">
                        {news.Title}
                    </h1>

                    <div className="prose dark:prose-invert max-w-none mb-8 text-gray-700 dark:text-gray-300 leading-relaxed whitespace-pre-wrap">
                        {news.Content}
                    </div>

                    <div className="border-t border-gray-200 dark:border-gray-700 pt-6 flex justify-between items-center">
                        <a
                            href="/"
                            className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white font-medium transition-colors flex items-center"
                        >
                            <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
                            Back to Feed
                        </a>

                        <a
                            href={news.Link}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
                        >
                            Read Original
                            <svg className="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path></svg>
                        </a>
                    </div>
                </div>
            </article>
        </main>
    );
}
