import React from 'react';

function CodeWithLineNumbers({ code }) {
  if (!code) return null;
  
  const lines = code.split('\n');
  
  return (
    <div className="code-container">
      <div className="line-numbers">
        {lines.map((_, i) => (
          <div key={`line-${i}`} className="line-number">{i + 1}</div>
        ))}
      </div>
      <pre className="code-content">
        <code dangerouslySetInnerHTML={{ __html: code }} />
      </pre>
    </div>
  );
}

export default CodeWithLineNumbers;