import React, { useState, useEffect } from 'react';
import './PullRequestsSection.css';
import { useAnimation } from '../../AnimationContext';

interface PullRequest {
  id: number;
  title: string;
  repoName: string;
  baseBranch: string;
  prBranch: string;
  url: string;
}

const PullRequestsSection: React.FC = () => {
  const { highlightNewPR } = useAnimation();
  
  // Initial PR data that will be used for every page reload
  const initialPullRequests: (PullRequest & { isNew?: boolean })[] = [];
  
  // State for PR data including animation state
  const [pullRequests, setPullRequests] = useState<(PullRequest & { isNew?: boolean })[]>(initialPullRequests);
  
  // Reset to initial state on component mount
  useEffect(() => {
    setPullRequests(initialPullRequests);
  }, []);
  
  // For demo purposes, this is kept for reference - in reverse order (newest first)
  const _initialData = [
    { 
      id: 1, 
      title: '(feat/ui): Update repo to Svelte 5, attempt 1', 
      repoName: 'sourcegraph/superdev', 
      baseBranch: 'main', 
      prBranch: '(feat/ui): Update repo to Svelte 5, attempt 1', 
      url: 'https://github.com/sourcegraph/superdev/pull/123' 
    },
    {
      id: 2, 
      title: '(feat/ui): Update repo to Svelte 5, attempt 2', 
      repoName: 'sourcegraph/superdev', 
      baseBranch: 'main', 
      prBranch: '(feat/ui): Update repo to Svelte 5, attempt 2', 
      url: 'https://github.com/sourcegraph/superdev/pull/456' 
    },
    { 
      id: 3, 
      title: '(feat/ui): Update repo to Svelte 5, attempt 3', 
      repoName: 'sourcegraph/superdev', 
      baseBranch: 'main', 
      prBranch: '(feat/ui): Update repo to Svelte 5, attempt 3', 
      url: 'https://github.com/sourcegraph/superdev/pull/3' 
    },
    { 
      id: 4, 
      title: '(feat/ui): Update repo to Svelte 5, attempt 4', 
      repoName: 'sourcegraph/superdev', 
      baseBranch: 'main', 
      prBranch: '(feat/ui): Update repo to Svelte 5, attempt 4', 
      url: 'https://github.com/sourcegraph/superdev/pull/101' 
    },
    { 
      id: 5, 
      title: '(feat/ui): Update repo to Svelte 5, attempt 5', 
      repoName: 'sourcegraph/superdev', 
      baseBranch: 'main', 
      prBranch: '(feat/ui): Update repo to Svelte 5, attempt 5', 
      url: 'https://github.com/sourcegraph/superdev/pull/112' 
    },
  ];
  
  // Copy the initial data to our initialPullRequests in reverse order (newest first)
  // We're sorting by ID in descending order to show newest (highest ID) first
  initialPullRequests.push(..._initialData.sort((a, b) => b.id - a.id));
  
  // Track total PRs to display in status bar
  const [totalPRsIssued, setTotalPRsIssued] = useState<number>(5);
  
  // Add a new PR when highlightNewPR becomes true
  useEffect(() => {
    if (highlightNewPR) {
      // Increment total PRs when adding a new one
      setTotalPRsIssued(prev => prev + 1);
      const newPR = { 
        id: Math.floor(Math.random() * 1000) + 100, // Random ID to ensure uniqueness
        title: 'Implement Sveltekit Marketing Site, Attempt 6', 
        repoName: 'sourcegraph/superdev', 
        baseBranch: 'main', 
        prBranch: 'feat/sveltekit-marketing-site', 
        url: 'https://github.com/sourcegraph/superdev/pull/2',
        isNew: true
      };
      
      // Add the new PR to the beginning of the list (newest first)
      setPullRequests(prev => [newPR, ...prev]);
      
      // Store the new PR ID for cleanup
      const newPrId = newPR.id;
      
      // Stop blinking after 5 seconds
      const stopBlinkTimer = setTimeout(() => {
        setPullRequests(prev => 
          prev.map(pr => 
            pr.id === newPrId ? { ...pr, isNew: false } : pr
          )
        );
      }, 5000);
      
      return () => clearTimeout(stopBlinkTimer);
    }
  }, [highlightNewPR]);

  return (
    <div className="pull-requests-section">
      {/* Component 9: Status bar with thermometer */}
      <div className="pr-status-bar">
        <div className="pr-status-label">Pull requests ({totalPRsIssued})</div>
      </div>

      {/* Component 8: List of successful pull requests */}
      <div className="pr-list">
        {pullRequests.map(pr => (
          <a 
            key={pr.id} 
            href={pr.url} 
            target="_blank" 
            rel="noopener noreferrer"
            className={`pr-item ${pr.isNew ? 'new-pr' : ''}`}
          >
            <div className="pr-header">
              <div className="pr-title">{pr.title}</div>
              <div className="pr-outlink-icon">↗</div>
            </div>
            <div className="pr-meta">
              <div className="pr-repo">{pr.repoName}</div>
              <div className="pr-branches">
                <span className="pr-base-branch">{pr.baseBranch}</span>
                <span className="branch-separator">←</span>
                <span className="pr-head-branch">{pr.prBranch}</span>
              </div>
            </div>
          </a>
        ))}
      </div>
    </div>
  );
};

export default PullRequestsSection;