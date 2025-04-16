import React from 'react';
import './App.css';
import { Header, ThreadList } from './components';
import { useThreads } from './hooks/useThreads';

function App() {
  const { threads, loading, error, outputRefs } = useThreads();

  if (loading) {
    return <div className="app">Loading threads...</div>;
  }

  if (error) {
    return <div className="app">Error: {error}</div>;
  }

  return (
    <div className="app">
      <Header />
      <ThreadList threads={threads} outputRefs={outputRefs} />
    </div>
  );
}

export default App;