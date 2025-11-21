'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';

export const RefreshButton: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const handleRefresh = async () => {
        setLoading(true);
        try {
            const res = await fetch('/api/news/refresh');
            if (res.ok) {
                router.refresh();
            } else {
                console.error('Failed to refresh news');
            }
        } catch (error) {
            console.error('Failed to refresh news:', error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <button
            onClick={handleRefresh}
            disabled={loading}
            className={`px-4 py-2 rounded-md text-white font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 ${loading
                    ? 'bg-blue-400 cursor-not-allowed'
                    : 'bg-blue-600 hover:bg-blue-700'
                }`}
        >
            {loading ? (
                <span className="flex items-center">
                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Refreshing...
                </span>
            ) : (
                'Refresh News'
            )}
        </button>
    );
};
