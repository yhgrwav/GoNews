import React from 'react';

interface NewsItem {
    ID: number;
    Title: string;
    Content: string;
    PubTime: string;
    Link: string;
    Source: string;
}

interface NewsCardProps {
    news: NewsItem;
}

export const NewsCard: React.FC<NewsCardProps> = ({ news }) => {
    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow duration-300 border border-gray-200 dark:border-gray-700 flex flex-col h-full">
            <div className="flex justify-between items-start mb-2">
                <span className="text-xs font-semibold text-blue-600 dark:text-blue-400 uppercase tracking-wider">
                    {news.Source}
                </span>
                <span className="text-xs text-gray-500 dark:text-gray-400">
                    {new Date(news.PubTime).toLocaleString()}
                </span>
            </div>
            <h2 className="text-xl font-bold mb-3 text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
                <a href={news.Link} target="_blank" rel="noopener noreferrer" className="hover:underline">
                    {news.Title}
                </a>
            </h2>
            <p className="text-gray-700 dark:text-gray-300 mb-4 line-clamp-3 flex-grow">
                {news.Content}
            </p>
            <div className="mt-auto">
                <a
                    href={news.Link}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center text-blue-600 dark:text-blue-400 hover:underline font-medium"
                >
                    Read more
                    <svg className="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
                </a>
            </div>
        </div>
    );
};
