import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import './App.css';

const API_BASE_URL = 'http://localhost:8080';

function App() {
  const [threads, setThreads] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const outputRefs = useRef({});

  useEffect(() => {
    const fetchThreads = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/threads`);
        if (response.data && response.data.threads) {
          // Compare thread IDs to see if anything has changed
          const newThreadIds = response.data.threads.map(t => t.thread_id).sort().join(',');
          const currentThreadIds = threads.map(t => t.thread_id).sort().join(',');
          
          // Only update if the list of threads changed
          if (newThreadIds !== currentThreadIds) {
            setThreads(response.data.threads);
          }
        }
        setLoading(false);
      } catch (err) {
        console.error('Error fetching threads:', err);
        setError('Failed to fetch threads. Please try again.');
        setLoading(false);
      }
    };

    // Fetch threads initially
    fetchThreads();

    // Set up polling every 1 second
    const intervalId = setInterval(fetchThreads, 1000);

    // Clean up interval on component unmount
    return () => clearInterval(intervalId);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Fetch thread outputs
  useEffect(() => {
    const fetchThreadOutputs = async () => {
      const updatedThreads = await Promise.all(
        threads.map(async (thread) => {
          try {
            const response = await axios.get(`${API_BASE_URL}/output?thread_id=${thread.thread_id}`);
            // Only update if there are changes to avoid unnecessary re-renders
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

      // Check if anything actually changed before updating state
      const hasChanges = updatedThreads.some((updated, index) => 
        updated.output !== threads[index].output || 
        updated.status !== threads[index].status ||
        updated.error !== threads[index].error
      );

      if (hasChanges) {
        setThreads(updatedThreads);
        
        // Scroll to bottom of output for processing threads
        setTimeout(() => {
          updatedThreads.forEach(thread => {
            if (thread.status === 'processing' && outputRefs.current[thread.thread_id]) {
              const outputElement = outputRefs.current[thread.thread_id];
              outputElement.scrollTop = outputElement.scrollHeight;
            }
          });
        }, 100);
      }
    };

    if (threads.length > 0) {
      fetchThreadOutputs();
    }
  }, [threads]);

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
                  <span className="loading-indicator">⏳ Updating...</span>
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