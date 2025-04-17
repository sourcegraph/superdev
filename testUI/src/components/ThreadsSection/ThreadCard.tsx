import React from 'react';
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

  // Get animation context to control highlighting
  const { animationStarted } = useAnimation();
  
  return (
    <div className={`thread-card ${thread.status} ${(thread.id === 3 && animationStarted) ? 'highlight-success' : ''}`}>
      <div className="thread-header">
        <h3>{thread.name}</h3>
        {/* Component 6: Status decal */}
        {getStatusDecal(thread.status)}
      </div>
      
      <div className="thread-content">
        {/* Display the streaming text */}
        {thread.content.split('\n').map((line, index) => (
          <div key={index} className="content-line">{line}</div>
        ))}
      </div>
      
      {/* Component 7: Feedback section for failed or needs feedback threads */}
      {(thread.status === 'failed' || thread.status === 'needsFeedback') && (
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

export default ThreadCard;