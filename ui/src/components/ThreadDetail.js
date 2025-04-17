import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useThreadDetail } from '../hooks/useThreadDetail';
import { processEscapeCodes } from '../utils/escapeCodeHandler';
import { copyToClipboard } from '../utils/clipboard';
import CodeWithLineNumbers from './CodeWithLineNumbers';

function ThreadDetail() {
  const { threadId } = useParams();
  const navigate = useNavigate();
  const { thread, loading, error } = useThreadDetail(threadId);
  const [copySuccess, setCopySuccess] = useState(false);

  const handleCopyId = async () => {
    if (!thread) return;
    
    const success = await copyToClipboard(thread.thread_id);
    if (success) {
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
    }
  };

  const handleBack = () => {
    navigate('/');
  };

  if (loading) {
    return (
      <div className="thread-detail-container">
        <div className="thread-detail-header skeleton">
          <div className="skeleton-title"></div>
          <div className="skeleton-meta"></div>
        </div>
        <div className="thread-detail-body skeleton">
          <div className="skeleton-line"></div>
          <div className="skeleton-line"></div>
          <div className="skeleton-line"></div>
          <div className="skeleton-line"></div>
          <div className="skeleton-line"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="thread-detail-container error">
        <div className="thread-detail-error">
          <h2>Error</h2>
          <p>{error}</p>
          <button onClick={handleBack} className="back-button">
            &larr; Back to Thread List
          </button>
        </div>
      </div>
    );
  }

  if (!thread) {
    return (
      <div className="thread-detail-container error">
        <div className="thread-detail-error">
          <h2>Thread Not Found</h2>
          <p>The requested thread could not be found.</p>
          <button onClick={handleBack} className="back-button">
            &larr; Back to Thread List
          </button>
        </div>
      </div>
    );
  }

  const processedOutput = thread.output ? processEscapeCodes(thread.output) : '';

  return (
    <div className="thread-detail-container">
      <button onClick={handleBack} className="back-button">
        &larr; Back to Thread List
      </button>
      
      <div className="thread-detail-header">
        <div className="thread-id-container">
          <h1>Thread Details</h1>
          <div className="thread-id">
            <span>Thread ID: {thread.thread_id}</span>
            <button 
              onClick={handleCopyId} 
              className="copy-button"
              title="Copy Thread ID"
            >
              {copySuccess ? '✓ Copied!' : 'Copy'}
            </button>
          </div>
        </div>
        
        <div className="thread-meta">
          <div className="meta-item">
            <strong>Created:</strong> {new Date(thread.created_at).toLocaleString()}
          </div>
          
          {thread.repository && (
            <div className="meta-item">
              <strong>Repository:</strong> 
              <a 
                href={`https://github.com/${thread.repository}`} 
                target="_blank" 
                rel="noopener noreferrer"
              >
                {thread.repository}
              </a>
            </div>
          )}
          
          <div className={`status status-${thread.status}`}>
            <strong>Status:</strong> {thread.status}
          </div>
        </div>
      </div>

      {thread.error && (
        <div className="thread-error-message">
          <h3>Error:</h3>
          <p>{thread.error}</p>
        </div>
      )}
      
      <div className="thread-detail-body">
        <h3>Output:</h3>
        <div className="thread-output-container">
          <CodeWithLineNumbers code={processedOutput} />
        </div>
      </div>
    </div>
  );
}

export default ThreadDetail;