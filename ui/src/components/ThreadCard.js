import React, { useRef, useEffect } from 'react';
import { processEscapeCodes } from '../utils/escapeCodeHandler';

function ThreadCard({ thread, outputRefs }) {
  const getProcessedOutput = () => {
    if (!thread.output) {
      if (thread.error) return `Error: ${thread.error}`;
      if (thread.status === 'processing') return 'Processing...';
      return 'No output yet';
    }
    
    // Process any escape codes in the output
    return processEscapeCodes(thread.output);
  };
  
  return (
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
        dangerouslySetInnerHTML={{ __html: getProcessedOutput() }}
      />
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
  );
}

export default ThreadCard;