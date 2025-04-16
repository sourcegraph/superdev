import { useState, useEffect, useRef } from 'react';
import { fetchThreads, fetchThreadOutput } from '../services/api';

export const useThreads = () => {
  const [threads, setThreads] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const outputRefs = useRef({});

  // Sort threads by creation date (newest first)
  const sortThreads = (threads) => {
    return [...threads].sort((a, b) => 
      new Date(b.created_at) - new Date(a.created_at)
    );
  };

  // Initial fetch of threads
  useEffect(() => {
    let isMounted = true;
    
    const getThreads = async () => {
      try {
        const data = await fetchThreads();
        if (!isMounted) return;
        
        if (data && data.threads) {
          setThreads(sortThreads(data.threads));
        }
        setLoading(false);
      } catch (err) {
        if (!isMounted) return;
        console.error('Error fetching threads:', err);
        setError('Failed to fetch threads. Please try again.');
        setLoading(false);
      }
    };

    // Fetch threads only once on component mount
    getThreads();

    return () => {
      isMounted = false;
    };
  }, []);

  // Fetch thread outputs and poll for updates
  useEffect(() => {
    let isMounted = true;
    
    const fetchThreadOutputs = async () => {
      // First check if there are any new threads
      try {
        const data = await fetchThreads();
        if (!isMounted) return;
        
        if (data && data.threads) {
          const serverThreads = data.threads;
          
          // Compare thread lists
          const serverIds = serverThreads.map(t => t.thread_id).sort();
          const clientIds = threads.map(t => t.thread_id).sort();
          
          // Check if the lists are different
          if (JSON.stringify(serverIds) !== JSON.stringify(clientIds)) {
            // We have new threads or threads were removed - update the entire list
            setThreads(sortThreads(serverThreads));
            return; // Exit this polling cycle, next one will fetch outputs
          }
        }
      } catch (err) {
        console.error('Error checking for new threads:', err);
      }
      
      // Then fetch outputs for existing threads
      const updatedThreads = await Promise.all(
        threads.map(async (thread) => {
          try {
            const threadData = await fetchThreadOutput(thread.thread_id);
            // Only update if there are content changes to avoid unnecessary re-renders
            if (threadData.status !== thread.status || 
                threadData.output !== thread.output || 
                threadData.error !== thread.error) {
              return {
                ...thread,
                output: threadData.output || '',
                error: threadData.error || '',
                status: threadData.status,
              };
            }
            return thread; // No changes, return original
          } catch (err) {
            // If there's an error, just return the thread as is
            return thread;
          }
        })
      );

      if (!isMounted) return;

      // Find threads with content changes
      let contentChanged = false;
      const updatedWithChanges = updatedThreads.map((updated, index) => {
        const original = threads[index];
        const hasChanged = 
          updated.output !== original.output || 
          updated.status !== original.status || 
          updated.error !== original.error;
        
        if (hasChanged) contentChanged = true;
        return updated;
      });

      if (contentChanged) {
        const sortedThreads = sortThreads(updatedWithChanges);
        
        setThreads(sortedThreads);
        
        // Scroll to bottom of output for processing threads
        setTimeout(() => {
          if (!isMounted) return;
          
          sortedThreads.forEach(thread => {
            if (thread.status === 'processing' && outputRefs.current[thread.thread_id]) {
              const outputElement = outputRefs.current[thread.thread_id];
              outputElement.scrollTop = outputElement.scrollHeight;
            }
          });
        }, 100);
      }
    };

    // Initial fetch
    if (threads.length > 0) {
      fetchThreadOutputs();
    }
    
    // Set up polling
    const intervalId = setInterval(fetchThreadOutputs, 1000);
    
    // Cleanup function
    return () => {
      isMounted = false;
      clearInterval(intervalId);
    };
  }, [threads]);

  return { threads, loading, error, outputRefs };
};