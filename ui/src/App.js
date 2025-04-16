import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import './App.css';
import Spinner from './components/Spinner';

const API_BASE_URL = 'http://localhost:8080';

function App() {
  const [threads, setThreads] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const outputRefs = useRef({});

  // Initial fetch of threads
  useEffect(() => {
    let isMounted = true;
    
    const fetchThreads = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/threads`);
        if (!isMounted) return;
        
        if (response.data && response.data.threads) {
          // Get the new threads and sort them consistently
          const newThreads = [...response.data.threads].sort((a, b) => 
            new Date(b.created_at) - new Date(a.created_at)
          );
          
          setThreads(newThreads);
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
    fetchThreads();

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
        const threadsResponse = await axios.get(`${API_BASE_URL}/threads`);
        if (!isMounted) return;
        
        if (threadsResponse.data && threadsResponse.data.threads) {
          const serverThreads = threadsResponse.data.threads;
          
          // Compare thread lists
          const serverIds = serverThreads.map(t => t.thread_id).sort();
          const clientIds = threads.map(t => t.thread_id).sort();
          
          // Check if the lists are different
          if (JSON.stringify(serverIds) !== JSON.stringify(clientIds)) {
            // We have new threads or threads were removed - update the entire list
            const sortedThreads = [...serverThreads].sort((a, b) => 
              new Date(b.created_at) - new Date(a.created_at)
            );
            setThreads(sortedThreads);
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
            const response = await axios.get(`${API_BASE_URL}/output?thread_id=${thread.thread_id}`);
            // Only update if there are content changes to avoid unnecessary re-renders
            if (response.data.status !== thread.status || 
                response.data.output !== thread.output || 
                response.data.error !== thread.error) {
              return {
                ...thread,
                output: response.data.output || '',
                error: response.data.error || '',
                status: response.data.status,
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
        // Keep the same order by using thread ID for stable sort
        const sortedThreads = [...updatedWithChanges].sort((a, b) => 
          new Date(b.created_at) - new Date(a.created_at)
        );
        
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
  }, [threads]);  // Re-create this effect when threads list changes

  if (loading) {
    return <div className="app">Loading threads...</div>;
  }

  if (error) {
    return <div className="app">Error: {error}</div>;
  }

  return (
    <div className="app">
      <header className="header">
        <div className="logo">
          <span>SuperDev</span>
        </div>
        <div className="nav-item active">Threads</div>
      </header>

      <div className="threads-grid">
        {threads.map((thread) => (
          <div key={thread.thread_id} className="thread-card">
            <div className="thread-header">
              <h3 className="thread-title">Thread ID: {thread.thread_id.substring(0, 8)}...</h3>
              <div className="thread-meta">
                <div title={new Date(thread.created_at).toLocaleString()}>
                  Created: {new Date(thread.created_at).toLocaleTimeString()}
                </div>
              </div>
            </div>
            <div 
              className="thread-body" 
              ref={el => outputRefs.current[thread.thread_id] = el}
            >
              {thread.output ? thread.output : 
               thread.error ? `Error: ${thread.error}` : 
               thread.status === 'processing' ? 'Processing...' : 'No output yet'}
            </div>
            <div className="thread-footer">
              <div className={`status status-${thread.status}`}>
                Status: {thread.status}
              </div>
              <div>
                {thread.status === 'processing' && (
                  <div className="loading-indicator">
                    <Spinner imageUrl="/ampster.png" alt="Loading..." />
                    <span style={{ marginLeft: '8px' }}>Updating...</span>
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}

        {threads.length === 0 && (
          <div>No threads found. Start a new thread to see output.</div>
        )}
      </div>
    </div>
  );
}

export default App;