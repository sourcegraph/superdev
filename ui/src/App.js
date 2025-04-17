import React, { useState, useEffect } from 'react';
import './App.css';
import { Header, ThreadList, ThreadDetail } from './components';
import { useThreads } from './hooks/useThreads';

function App() {
  const { threads, loading, error, getUpdatedThread } = useThreads();
  const [selectedThread, setSelectedThread] = useState(null);

  // Keep the selected thread updated
  useEffect(() => {
    if (!selectedThread) return;
    
    let isMounted = true;
    const intervalId = setInterval(async () => {
      const updatedThreadData = await getUpdatedThread(selectedThread.thread_id);
      if (!isMounted || !updatedThreadData) return;
      
      // Only update if content has changed to avoid unnecessary re-renders
      if (updatedThreadData.status !== selectedThread.status || 
          updatedThreadData.output !== selectedThread.output || 
          updatedThreadData.error !== selectedThread.error) {
        setSelectedThread({
          ...selectedThread,
          output: updatedThreadData.output || '',
          error: updatedThreadData.error || '',
          status: updatedThreadData.status,
        });
      }
    }, 1000);
    
    return () => {
      isMounted = false;
      clearInterval(intervalId);
    };
  }, [selectedThread, getUpdatedThread]);

  // Update selected thread if it exists in updated threads list
  useEffect(() => {
    if (!selectedThread) return;
    
    const updatedThread = threads.find(t => t.thread_id === selectedThread.thread_id);
    if (updatedThread) {
      // Check if there are changes before updating to avoid unnecessary re-renders
      if (updatedThread.status !== selectedThread.status || 
          updatedThread.output !== selectedThread.output || 
          updatedThread.error !== selectedThread.error) {
        setSelectedThread(updatedThread);
      }
    }
  }, [threads, selectedThread]);

  if (loading) {
    return <div className="app">Loading threads...</div>;
  }

  if (error) {
    return <div className="app">Error: {error}</div>;
  }

  const handleThreadSelect = (thread) => {
    setSelectedThread(thread);
  };

  const handleBackToList = () => {
    setSelectedThread(null);
  };

  return (
    <div className="app">
      <Header />
      {selectedThread ? (
        <ThreadDetail thread={selectedThread} onBack={handleBackToList} />
      ) : (
        <ThreadList threads={threads} onThreadSelect={handleThreadSelect} />
      )}
    </div>
  );
}

export default App;