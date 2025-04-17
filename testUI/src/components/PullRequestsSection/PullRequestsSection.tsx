import React from 'react';
import './PullRequestsSection.css';

interface PullRequest {
  id: number;
  name: string;
  url: string;
}

const PullRequestsSection: React.FC = () => {
  // Mock pull request data
  const pullRequests: PullRequest[] = [
    { id: 1, name: 'PR-1', url: 'https://github.com/user/repo3/pull/123' },
    { id: 2, name: 'PR-2', url: 'https://github.com/user/repo4/pull/456' },
    { id: 3, name: 'PR-3', url: 'https://github.com/user/repo5/pull/789' },
    { id: 4, name: 'PR-4', url: 'https://github.com/user/repo6/pull/101' },
    { id: 5, name: 'PR-5', url: 'https://github.com/user/repo7/pull/112' },
  ];

  return (
    <div className="pull-requests-section">
      {/* Component 9: Status bar with thermometer */}
      <div className="pr-status-bar">
        <div className="pr-status-label">PRs issued</div>
        <div className="pr-thermometer">
          <div 
            className="pr-thermometer-fill" 
            style={{ width: `${(pullRequests.length / 10) * 100}%` }}
          ></div>
        </div>
        <div className="pr-count">{pullRequests.length}</div>
      </div>

      {/* Component 8: List of successful pull requests */}
      <div className="pr-list">
        {pullRequests.map(pr => (
          <a 
            key={pr.id} 
            href={pr.url} 
            target="_blank" 
            rel="noopener noreferrer"
            className="pr-item"
          >
            {pr.name}
          </a>
        ))}
      </div>
    </div>
  );
};

export default PullRequestsSection;