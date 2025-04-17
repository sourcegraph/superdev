import React, { useState } from 'react';
import './ThreadsSection.css';
import ThreadCard from './ThreadCard';

export type ThreadStatus = 'running' | 'failed' | 'needsFeedback' | 'completed';

interface Thread {
  id: number;
  name: string;
  content: string;
  status: ThreadStatus;
}

const ThreadsSection: React.FC = () => {
  // Mock thread data
  const [threads, setThreads] = useState<Thread[]>([
    { 
      id: 1, 
      name: 'Thread-1', 
      content: 'Scanning repository: user/repo1\nChecking package.json\nFound Node v16\nUpdating to Node v24\nUpdating dependencies...\nFailed: Incompatible dependency found', 
      status: 'failed' 
    },
    { 
      id: 2, 
      name: 'Thread-2', 
      content: 'Scanning repository: user/repo2\nChecking package.json\nFound Node v20\nNeeds clarification on handling custom build scripts\nWaiting for feedback...', 
      status: 'needsFeedback' 
    },
    { 
      id: 3, 
      name: 'Thread-3', 
      content: 'Scanning repository: user/repo3\nChecking package.json\nFound Node v18\nUpdating to Node v24\nUpdating dependencies...\nRunning tests...\nAll tests passed\nCreating pull request...', 
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

  return (
    <div className="threads-section">
      <h2>Threads</h2>
      {/* Component 5: Threads container */}
      <div className="threads-container">
        {threads.map(thread => (
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
    </div>
  );
};

export default ThreadsSection;