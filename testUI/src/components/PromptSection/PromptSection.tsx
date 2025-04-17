import React, { useState } from 'react';
import './PromptSection.css';
import { useAnimation, Thread } from '../../AnimationContext';

const PromptSection: React.FC = () => {
  const { setDynamicThreads, setAnimationStarted, setHighlightNewPR, setHideThread } = useAnimation();
  const [prompt, setPrompt] = useState('Act like a Senior Designer, your task is to generate a marketing site written in Sveltekit. This site will be the main vehicle through which you attract top talent and hype for your business, so it needs to be extremely beautiful, intuitive, and distinguished. Give me 4 versions I can evaluate.');
  const [isRunning, setIsRunning] = useState(false);

  const generateThreads = () => {
    // Set running state for visual feedback
    setIsRunning(true);
    // Clear any existing threads first
    setDynamicThreads([]);
    
    // Reset animation states
    setAnimationStarted(false);
    setHighlightNewPR(false);
    setHideThread(false);

    // Create 4 threads with different statuses
    const newThreads: Thread[] = [
      {
        id: 101,
        name: 'Thread-101',
        content: `Analyzing prompt: ${prompt.substring(0, 50)}...\nGenerating designs...\nInitializing Sveltekit project\nApplying styles...\nError: Could not resolve dependencies. Check package.json`,
        status: 'failed'
      },
      {
        id: 102,
        name: 'Thread-102',
        content: `Analyzing prompt: ${prompt.substring(0, 50)}...\nGenerating design concepts...\nNeed clarification: What color scheme would you prefer for this marketing site?`,
        status: 'needsFeedback'
      },
      {
        id: 103,
        name: 'Thread-103',
        content: `Analyzing prompt: ${prompt.substring(0, 50)}...\nGenerating alternative designs...\nIs there a specific industry focus for this marketing site?`,
        status: 'needsFeedback'
      },
      {
        id: 104,
        name: 'Thread-104',
        content: `Analyzing prompt: ${prompt.substring(0, 50)}...\nGenerating design concepts...\nCreating Sveltekit project structure\nImplementing responsive design\nAdding animations\nOptimizing for performance\nCompleted: Design implemented successfully!`,
        status: 'completed'
      }
    ];

    // Update the dynamic threads with a short delay
    setTimeout(() => {
      setDynamicThreads(newThreads);
      // Reset running state after threads are generated
      setIsRunning(false);
      
      // Start animation sequence after threads appear
      // First, highlight the successful thread after 5 seconds
      setTimeout(() => {
        setAnimationStarted(true);
        
        // Then, hide the thread and add a new PR after 2 seconds
        setTimeout(() => {
          setHideThread(true);
          
          // Add new PR immediately after hiding the thread
          setTimeout(() => {
            setHighlightNewPR(true);
          }, 200); // Small delay for visual effect
        }, 2000);
      }, 5000);
      
    }, 800); // Delay for visual effect
  };

  return (
    <div className="prompt-section">
      <label htmlFor="prompt">Prompt</label>
      <div className="prompt-container">
        {/* Component 3: Natural language prompt */}
        <textarea 
          id="prompt" 
          className="prompt-textarea" 
          placeholder="Enter your prompt here"
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
        ></textarea>
        <div className="button-container">
          <button 
            className={`amplify-button ${isRunning ? 'running' : ''}`} 
            onClick={generateThreads}
            disabled={isRunning}
          >
            {isRunning ? 'Amplifying...' : 'Amplify'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default PromptSection;