import { useState, useEffect } from 'react';
import { fetchThreadOutput } from '../services/api';

export const useThreadDetail = (threadId) => {
  const [thread, setThread] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    let isMounted = true;
    
    const fetchThread = async () => {
      if (!threadId) {
        setError('Thread ID is required');
        setLoading(false);
        return;
      }
      
      try {
        const data = await fetchThreadOutput(threadId);
        if (!isMounted) return;
        
        if (data) {
          setThread({
            thread_id: threadId,
            ...data
          });
        } else {
          setError('Thread not found');
        }
      } catch (err) {
        if (!isMounted) return;
        console.error('Error fetching thread details:', err);
        setError(err.response?.status === 404 
          ? 'Thread not found' 
          : 'Failed to fetch thread details. Please try again.');
      } finally {
        if (isMounted) setLoading(false);
      }
    };

    fetchThread();
    
    // Set up polling for thread updates if status is 'processing'
    let intervalId;
    if (thread && thread.status === 'processing') {
      intervalId = setInterval(fetchThread, 1000);
    }

    return () => {
      isMounted = false;
      if (intervalId) clearInterval(intervalId);
    };
  }, [threadId, thread]);  // Include thread as dependency

  return { thread, loading, error };
};