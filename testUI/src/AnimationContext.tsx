import React, { createContext, useState, useContext, ReactNode } from 'react';

export interface Thread {
  id: number;
  name: string;
  content: string;
  status: 'running' | 'failed' | 'needsFeedback' | 'completed';
  onClick?: () => void; // Add optional click handler
}

interface AnimationContextType {
  animationStarted: boolean;
  setAnimationStarted: (started: boolean) => void;
  highlightNewPR: boolean;
  setHighlightNewPR: (highlight: boolean) => void;
  hideThread: boolean;
  setHideThread: (hide: boolean) => void;
  // Thread management
  dynamicThreads: Thread[];
  setDynamicThreads: (threads: Thread[]) => void;
}

const AnimationContext = createContext<AnimationContextType>({
  animationStarted: false,
  setAnimationStarted: () => {},
  highlightNewPR: false,
  setHighlightNewPR: () => {},
  hideThread: false,
  setHideThread: () => {},
  dynamicThreads: [],
  setDynamicThreads: () => {}
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
  
  return (
    <AnimationContext.Provider value={{
      animationStarted,
      setAnimationStarted,
      highlightNewPR,
      setHighlightNewPR,
      hideThread,
      setHideThread,
      dynamicThreads,
      setDynamicThreads
    }}>
      {children}
    </AnimationContext.Provider>
  );
};