import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import { Header, ThreadList, ThreadDetail } from './components';
import { useThreads } from './hooks/useThreads';

function ThreadListPage() {
  const { threads, loading, error, outputRefs } = useThreads();

  if (loading) {
    return <div className="app">Loading threads...</div>;
  }

  if (error) {
    return <div className="app">Error: {error}</div>;
  }

  return (
    <>
      <Header />
      <ThreadList threads={threads} outputRefs={outputRefs} />
    </>
  );
}

function App() {
  return (
    <Router>
      <div className="app">
        <Routes>
          <Route path="/" element={<ThreadListPage />} />
          <Route path="/thread/:threadId" element={<ThreadDetail />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;