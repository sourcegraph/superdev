import React, { useState, useEffect } from 'react';
import './ThreadCard.css';
import { ThreadStatus } from './ThreadsSection';
import { useAnimation } from '../../AnimationContext';

interface ThreadCardProps {
  thread: {
    id: number;
    name: string;
    content: string;
    status: ThreadStatus;
  };
  feedback: string;
  shareWithOthers: boolean;
  onFeedbackChange: (value: string) => void;
  onShareChange: (checked: boolean) => void;
}

const ThreadCard: React.FC<ThreadCardProps> = ({ 
  thread, 
  feedback, 
  shareWithOthers, 
  onFeedbackChange, 
  onShareChange 
}) => {
  const [streamingComplete, setStreamingComplete] = useState(false);
  const { animationStarted } = useAnimation();

  const getStatusDecal = (status: ThreadStatus) => {
    switch (status) {
      case 'failed':
        return <div className="status-decal failed">✕</div>;
      case 'needsFeedback':
        return <div className="status-decal needs-feedback">!</div>;
      case 'completed':
        return <div className="status-decal completed">✓</div>;
      default:
        return null;
    }
  };

  // Show running status until streaming is complete
  const displayStatus = streamingComplete ? thread.status : 'running';
  
  // Add highlight class for completed Thread-104 when animation is active
  const isHighlighted = animationStarted && thread.id === 104 && streamingComplete;

  return (
    <div className={`thread-card ${displayStatus} ${isHighlighted ? 'highlight-success' : ''}`}>
      <div className="thread-header">
        <h3>{thread.name}</h3>
        {/* Show status decal only after streaming is complete */}
        {streamingComplete && getStatusDecal(thread.status)}
      </div>
      
      <div className="thread-content">
        {/* Display the streaming text with completion callback */}
        <StreamingContent 
          content={thread.content} 
          onComplete={() => setStreamingComplete(true)}
        />
      </div>
      
      {/* Component 7: Feedback section for failed or needs feedback threads - only show after streaming completes */}
      {streamingComplete && (thread.status === 'failed' || thread.status === 'needsFeedback') && (
        <div className="feedback-section">
          <textarea 
            className="feedback-input"
            placeholder="Add context or feedback..."
            value={feedback}
            onChange={(e) => onFeedbackChange(e.target.value)}
          ></textarea>
          
          <div className="feedback-actions">
            <label className="share-checkbox-label">
              <input 
                type="checkbox" 
                checked={shareWithOthers}
                onChange={(e) => onShareChange(e.target.checked)}
              />
              <span>Make available to all</span>
            </label>
            
            <button className="retry-button">Retry</button>
          </div>
        </div>
      )}
    </div>
  );
};

// StreamingContent component to simulate real-time computation
interface StreamingContentProps {
  content: string;
  onComplete?: () => void;
}

const StreamingContent: React.FC<StreamingContentProps> = ({ content, onComplete }) => {
  const [visibleLines, setVisibleLines] = useState<string[]>([]);
  const allLines = content.split('\n');
  
  useEffect(() => {
    // Reset visible lines when content changes
    setVisibleLines([]);
    
    // Add lines one by one with delays
    let currentLineIndex = 0;
    const timeouts: number[] = [];
    
    const addNextLine = () => {
      if (currentLineIndex < allLines.length) {
        setVisibleLines(prev => [...prev, allLines[currentLineIndex]]);
        currentLineIndex++;
        
        // Schedule next line with random delay for more natural feel
        const delay = Math.random() * 300 + 200; // 200-500ms
        const timeoutId = window.setTimeout(addNextLine, delay);
        timeouts.push(timeoutId);
      } else {
        // All lines have been displayed, call the onComplete callback
        if (onComplete) {
          onComplete();
        }
      }
    };
    
    // Start adding lines
    const initialTimeout = window.setTimeout(addNextLine, 300);
    timeouts.push(initialTimeout);
    
    // Cleanup function - cancel any pending timeouts
    return () => {
      timeouts.forEach(timeoutId => window.clearTimeout(timeoutId));
    };
  }, [content]);
  
  return (
    <>
      {visibleLines.map((line, index) => (
        <div key={index} className="content-line">{line}</div>
      ))}
      {visibleLines.length < allLines.length && (
        <div className="content-line typing-indicator">...</div>
      )}
    </>
  );
};

export default ThreadCard;