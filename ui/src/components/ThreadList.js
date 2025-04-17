import React from 'react';
import ThreadListItem from './ThreadListItem';

function ThreadList({ threads, onThreadSelect }) {
  return (
    <div className="threads-list">
      {threads.length > 0 ? threads.map((thread) => (
        <ThreadListItem 
          key={thread.thread_id} 
          thread={thread} 
          onClick={() => onThreadSelect(thread)} 
        />
      )) : (
        <div className="no-threads-message">
          No threads found. Start a new thread to see output.
        </div>
      )}
    </div>
  );
}

export default ThreadList;