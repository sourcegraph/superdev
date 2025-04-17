import React, { useEffect, useRef } from 'react';
import { processEscapeCodes } from '../utils/escapeCodeHandler';

function ThreadDetail({ thread, onBack }) {
  const outputRef = useRef(null);

  useEffect(() => {
    // Scroll to bottom of output when it changes
    if (outputRef.current && thread.status === 'processing') {
      outputRef.current.scrollTop = outputRef.current.scrollHeight;
    }
  }, [thread.output, thread.status]);

  const getStatusClass = () => {
    switch (thread.status) {
      case 'processing': return 'status-processing';
      case 'completed': return 'status-completed';
      case 'error': return 'status-error';
      default: return '';
    }
  };

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
    <div className="thread-detail">
      <div className="thread-detail-header">
        <button className="back-button" onClick={onBack}>
          &larr; Back to list
        </button>
        <div className="thread-detail-title">
          Thread: {thread.thread_id}
        </div>
        <div className={`thread-detail-status ${getStatusClass()}`}>
          Status: {thread.status}
          {thread.status === 'processing' && <span className="loading-indicator"> ⏳</span>}
        </div>
      </div>
      <div className="thread-detail-meta">
        <div title={new Date(thread.created_at).toLocaleString()}>
          Created: {new Date(thread.created_at).toLocaleString()}
        </div>
      </div>
      <div 
        className="thread-detail-content" 
        ref={outputRef}
        dangerouslySetInnerHTML={{ __html: getProcessedOutput() }}
      />
    </div>
  );
}

export default ThreadDetail;