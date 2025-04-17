import React, { createContext, useState, useContext, ReactNode } from 'react';

export interface Thread {
  id: number;
  name: string;
  content: string;
  status: 'running' | 'failed' | 'needsFeedback' | 'completed';
}

interface AnimationContextType {
  animationStarted: boolean;
  setAnimationStarted: (started: boolean) => void;
  highlightNewPR: boolean;
  setHighlightNewPR: (highlight: boolean) => void;
  hideThread: boolean;
  setHideThread: (hide: boolean) => void;
  // Thread management for slider demo
  dynamicThreads: Thread[];
  setDynamicThreads: (threads: Thread[]) => void;
  addThreadFromPrompt: (prompt: string, iterationNumber: number) => void;
}

const AnimationContext = createContext<AnimationContextType>({
  animationStarted: false,
  setAnimationStarted: () => {},
  highlightNewPR: false,
  setHighlightNewPR: () => {},
  hideThread: false,
  setHideThread: () => {},
  dynamicThreads: [],
  setDynamicThreads: () => {},
  addThreadFromPrompt: () => {}
});

export const useAnimation = () => useContext(AnimationContext);

interface AnimationProviderProps {
  children: ReactNode;
}

export const AnimationProvider: React.FC<AnimationProviderProps> = ({ children }) => {
  const [animationStarted, setAnimationStarted] = useState(false);
  const [highlightNewPR, setHighlightNewPR] = useState(false);
  const [hideThread, setHideThread] = useState(false);
  const [dynamicThreads, setDynamicThreads] = useState<Thread[]>([]);
  
  // Function to create a new thread from a prompt
  const addThreadFromPrompt = (prompt: string, iterationNumber: number) => {
    const newThread: Thread = {
      id: 100 + iterationNumber, // Use high IDs to avoid conflicts with existing threads
      name: `Thread-${iterationNumber}`,
      content: `Starting to process Iteration ${iterationNumber}...\nAnalyzing prompt: ${prompt.substring(0, 50)}...\nGenerating plan...\nExecuting plan...\nCompiling...`,
      status: Math.random() > 0.7 ? 'completed' : (Math.random() > 0.5 ? 'needsFeedback' : 'running')
    };
    
    // Add the new thread to the existing threads
    setDynamicThreads(prev => {
      // Check if this thread already exists
      const exists = prev.some(t => t.id === newThread.id);
      if (exists) {
        // Replace the existing thread
        return prev.map(t => t.id === newThread.id ? newThread : t);
      } else {
        // Add new thread
        return [...prev, newThread];
      }
    });
  };
  
  return (
    <AnimationContext.Provider value={{
      animationStarted,
      setAnimationStarted,
      highlightNewPR,
      setHighlightNewPR,
      hideThread,
      setHideThread,
      dynamicThreads,
      setDynamicThreads,
      addThreadFromPrompt
    }}>
      {children}
    </AnimationContext.Provider>
  );
};