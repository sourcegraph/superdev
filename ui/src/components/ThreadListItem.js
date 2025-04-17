import React from 'react';
import { formatDistance } from 'date-fns';

function ThreadListItem({ thread, onClick }) {
  const getStatusIndicator = () => {
    switch (thread.status) {
      case 'processing':
        return <span className="status-indicator processing">⏳</span>;
      case 'completed':
        return <span className="status-indicator completed">✅</span>;
      case 'error':
        return <span className="status-indicator error">❌</span>;
      default:
        return null;
    }
  };

  const formatTimeAgo = (date) => {
    try {
      return formatDistance(new Date(date), new Date(), { addSuffix: true });
    } catch (e) {
      return 'Invalid date';
    }
  };

  const getShortOutput = () => {
    if (!thread.output) {
      if (thread.error) return 'Error occurred';
      if (thread.status === 'processing') return 'Processing...';
      return 'No output';
    }
    
    // Take first line or first 60 chars
    const firstLine = thread.output.split('\n')[0];
    return firstLine.length > 60 
      ? firstLine.substring(0, 60) + '...' 
      : firstLine;
  };

  return (
    <div className="thread-list-item" onClick={onClick}>
      <div className="thread-item-left">
        {getStatusIndicator()}
        <div className="thread-item-info">
          <div className="thread-item-id">
            Thread: {thread.thread_id.substring(0, 8)}...
          </div>
          <div className="thread-item-summary">{getShortOutput()}</div>
        </div>
      </div>
      <div className="thread-item-time">
        {formatTimeAgo(thread.created_at)}
      </div>
    </div>
  );
}

export default ThreadListItem;