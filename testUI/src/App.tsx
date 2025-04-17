import React, { useEffect, useState } from 'react';
import './App.css';
import Header from './components/Header/Header';
import SearchBar from './components/SearchBar/SearchBar';
import PromptSection from './components/PromptSection/PromptSection';
import ThreadsSection from './components/ThreadsSection/ThreadsSection';
import PullRequestsSection from './components/PullRequestsSection/PullRequestsSection';

import { AnimationProvider } from './AnimationContext';

function App() {
  return (
    <AnimationProvider>
      <div className="App">
        <Header />
        <div className="main-content">
          <div className="left-panel">
            <SearchBar />
            <PromptSection />
            <ThreadsSection />
          </div>
          <div className="right-panel">
            <PullRequestsSection />
          </div>
        </div>
      </div>
    </AnimationProvider>
  );
}

export default App;