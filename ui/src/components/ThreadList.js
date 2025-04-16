import React from 'react';
import ThreadCard from './ThreadCard';

function ThreadList({ threads, outputRefs }) {
  return (
    <div className="threads-grid">
      {threads.map((thread) => (
        <ThreadCard key={thread.thread_id} thread={thread} outputRefs={outputRefs} />
      ))}

      {threads.length === 0 && (
        <div>No threads found. Start a new thread to see output.</div>
      )}
    </div>
  );
}

export default ThreadList;