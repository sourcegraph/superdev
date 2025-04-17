import React, { useState, useEffect } from 'react';
import './ThreadsSection.css';
import ThreadCard from './ThreadCard';
import ampLogo from '../../assets/amp-logo.png';
import { useAnimation } from '../../AnimationContext';

export type ThreadStatus = 'running' | 'failed' | 'needsFeedback' | 'completed';

interface Thread {
  id: number;
  name: string;
  content: string;
  status: ThreadStatus;
}

const ThreadsSection: React.FC = () => {
  const { animationStarted, setAnimationStarted, hideThread, setHideThread, setHighlightNewPR, dynamicThreads } = useAnimation();
  
  // Start animation sequence after 5 seconds
  useEffect(() => {
    const timer = setTimeout(() => {
      // Step 1: Highlight Thread 3
      setAnimationStarted(true);
      
      // Step 2: Hide Thread 3 after 2 seconds
      setTimeout(() => {
        setHideThread(true);
        
        // Step 3: Add new PR immediately after
        setTimeout(() => {
          setHighlightNewPR(true);
        }, 200); // Just a small delay for visual effect
      }, 3000);
    }, 5000);
    
    return () => clearTimeout(timer);
  }, [setAnimationStarted, setHighlightNewPR, setHideThread]);
  // Mock thread data
  const [threads, setThreads] = useState<Thread[]>([
    { 
      id: 1, 
      name: 'Thread-1', 
      content: 'Scanning repository: someFile.tsx\nReading file\nFound Svelte 3 usage\nUpdating to Svelte 5\nChecking compilation...\nFailed: Unable to resolve compilation errors.', 
      status: 'failed' 
    },
    { 
      id: 2, 
      name: 'Thread-2', 
      content: 'Scanning repository: anotherFile.tsx\nReading file\nFound Svelte 4 usage\nNeeds clarification on handling custom function\nWaiting for feedback...', 
      status: 'needsFeedback' 
    },
    { 
      id: 3, 
      name: 'Thread-3', 
      content: 'Scanning repository: andYetAnotherFile.tsx\nReading file\nFound Svelte 4 usage\nUpdating to Svelte 5\nCompiling...\nRunning CI...\nCI was successful!\nCreating pull request...', 
      status: 'completed' 
    },
  ]);

  const [feedback, setFeedback] = useState<{[key: number]: string}>({
    1: '',
    2: ''
  });

  const [shareWithOthers, setShareWithOthers] = useState<{[key: number]: boolean}>({
    1: false,
    2: false
  });

  const handleFeedbackChange = (threadId: number, value: string) => {
    setFeedback(prev => ({
      ...prev,
      [threadId]: value
    }));
  };

  const handleShareChange = (threadId: number, checked: boolean) => {
    setShareWithOthers(prev => ({
      ...prev,
      [threadId]: checked
    }));
  };

  // Filter out Thread-3 when hideThread is true
  const visibleThreads = threads.filter(thread => {
    if (hideThread && thread.id === 3) return false;
    return true;  
  });
  
  // Combine the original threads with dynamic threads from the slider
  const allThreads = [...visibleThreads, ...dynamicThreads];
  
  // Organize threads into rows of 3
  const organizeThreadsIntoRows = (threads: Thread[]) => {
    const rows: Thread[][] = [];
    for (let i = 0; i < threads.length; i += 3) {
      rows.push(threads.slice(i, i + 3));
    }
    return rows;
  };
  
  const threadRows = organizeThreadsIntoRows(allThreads);

  return (
    <div className="threads-section">
      <div className="threads-header">
        <img src={ampLogo} alt="Amp Logo" className="thread-icon" />
        <h2>Threads</h2>
      </div>
      {/* Component 5: Threads container */}
      <div className="threads-container">
        {threadRows.map((row, rowIndex) => (
          <div key={`row-${rowIndex}`} className="thread-row">
            {row.map(thread => (
              <ThreadCard 
                key={thread.id} 
                thread={thread} 
                feedback={feedback[thread.id] || ''}
                shareWithOthers={shareWithOthers[thread.id] || false}
                onFeedbackChange={(value) => handleFeedbackChange(thread.id, value)}
                onShareChange={(checked) => handleShareChange(thread.id, checked)}
              />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default ThreadsSection;