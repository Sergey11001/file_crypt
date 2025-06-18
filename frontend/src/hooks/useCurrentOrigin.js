import { useState, useEffect } from 'react';

export function useCurrentOrigin(): string {
    const [origin, setOrigin] = useState(() => {
        if (typeof window !== 'undefined') {
            return window.location.origin;
        }
        return '';
    });

    useEffect(() => {
        const handleLocationChange = () => {
            setOrigin(window.location.origin);
        };

        window.addEventListener('popstate', handleLocationChange);
        window.addEventListener('hashchange', handleLocationChange);

        return () => {
            window.removeEventListener('popstate', handleLocationChange);
            window.removeEventListener('hashchange', handleLocationChange);
        };
    }, []);

    return origin;
}