import React, { useState } from 'react';
import './ThreadList.css';

const ThreadList = ({ threads, outputRefs }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const threadsPerPage = 10;

  // Calculate pagination
  const indexOfLastThread = currentPage * threadsPerPage;
  const indexOfFirstThread = indexOfLastThread - threadsPerPage;
  const currentThreads = threads.slice(indexOfFirstThread, indexOfLastThread);
  const totalPages = Math.ceil(threads.length / threadsPerPage);

  // Format relative time (e.g., "2 minutes ago")
  const formatRelativeTime = (dateString) => {
    const now = new Date();
    const date = new Date(dateString);
    const diffInSeconds = Math.floor((now - date) / 1000);

    if (diffInSeconds < 60) return `${diffInSeconds} seconds ago`;
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} minutes ago`;
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    return `${Math.floor(diffInSeconds / 86400)} days ago`;
  };

  // Get status color class
  const getStatusClass = (status) => {
    switch (status.toLowerCase()) {
      case 'processing': return 'status-processing';
      case 'completed': return 'status-completed';
      case 'error': return 'status-error';
      default: return '';
    }
  };

  // Truncate thread ID and add copy functionality
  const TruncatedId = ({ id }) => {
    const truncated = id.length > 8 ? `${id.substring(0, 8)}...` : id;
    
    const copyToClipboard = (e) => {
      e.stopPropagation(); // Prevent row click
      navigator.clipboard.writeText(id)
        .then(() => {
          // Could add a toast notification here
          console.log('Thread ID copied to clipboard');
        })
        .catch(err => {
          console.error('Failed to copy: ', err);
        });
    };

    return (
      <div className="thread-id">
        <span>{truncated}</span>
        <button className="copy-button" onClick={copyToClipboard} aria-label="Copy thread ID">
          <span className="icon">📋</span>
        </button>
      </div>
    );
  };

  // Navigate to thread detail page
  const navigateToThread = (threadId) => {
    window.location.href = `/thread/${threadId}`;
  };

  // Render empty state
  if (threads.length === 0) {
    return <div className="empty-state">No threads found.</div>;
  }

  return (
    <div className="thread-list-container">
      <table className="thread-list">
        <thead>
          <tr>
            <th>Thread ID</th>
            <th>Status</th>
            <th>Creation Time</th>
          </tr>
        </thead>
        <tbody>
          {currentThreads.map((thread) => (
            <tr 
              key={thread.thread_id} 
              onClick={() => navigateToThread(thread.thread_id)}
              className="thread-row"
            >
              <td><TruncatedId id={thread.thread_id} /></td>
              <td>
                <span className={`status-badge ${getStatusClass(thread.status)}`}>
                  {thread.status}
                </span>
              </td>
              <td>{formatRelativeTime(thread.created_at)}</td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="pagination">
          <button 
            onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
            disabled={currentPage === 1}
          >
            Previous
          </button>
          <span className="page-info">{currentPage} of {totalPages}</span>
          <button 
            onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
            disabled={currentPage === totalPages}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
};

export default ThreadList;