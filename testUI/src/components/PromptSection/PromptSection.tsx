import React, { useState, useEffect } from 'react';
import './PromptSection.css';
import { useAnimation } from '../../AnimationContext';

const PromptSection: React.FC = () => {
  const { addThreadFromPrompt } = useAnimation();
  const [randomnessLevel, setRandomnessLevel] = useState<number>(0);
  
  // Save the previous level to detect changes
  const [prevLevel, setPrevLevel] = useState<number>(0);
  
  // Sample generated prompts based on randomness level
  const generatePrompts = (level: number, basePrompt: string) => {
    const variations = [
      'The Svelte5 docs are here: https://www. svelte.dev/docs. Read them, and upgrade this repo to use Svelte5. Do this in as few steps as possible, ensuring compilation, and implementing best practices according to Svelte 5 docs.',
      'Upgrade this repository to Svelte5 using minimal steps while ensuring proper compilation and following Svelte 5 documentation best practices.',
      'Your task: migrate this codebase to Svelte5, minimizing the conversion steps, maintaining build integrity, and adhering to recommended Svelte 5 patterns.',
      'Convert this project to use Svelte5 framework. Use the most efficient approach that ensures the code compiles and incorporates Svelte 5 best practices.',
      'Refactor this repository to implement Svelte5, taking the most direct path while ensuring successful compilation and following official Svelte 5 guidelines.'
    ];
    
    return variations.slice(0, level + 1);
  };

  const basePrompt = 'You are tasked with upgrading this repo to use Svelte5. Do this in as few steps as possible, ensuring compilation, and implementing best practices according to Svelte 5 docs.';
  const promptVariations = generatePrompts(randomnessLevel, basePrompt);
  
  // Create threads when the randomness level changes
  useEffect(() => {
    // Only create threads when level increases (not on initial load or reset)
    if (randomnessLevel > prevLevel) {
      // For each new variation, create a thread
      for (let i = prevLevel + 1; i <= randomnessLevel; i++) {
        // i will be 1-5, but we need to add 1 to get Iteration 2-6
        // (since Iteration 1 is the base prompt)
        addThreadFromPrompt(promptVariations[i - 1], i + 1);
      }
    }
    
    // Update the previous level
    setPrevLevel(randomnessLevel);
  }, [randomnessLevel, prevLevel, promptVariations, addThreadFromPrompt]);

  return (
    <div className="prompt-section">
      <label htmlFor="prompt">Iteration 1</label>
      {/* Component 3: Natural language prompt */}
      <textarea 
        id="prompt" 
        className="prompt-textarea" 
        placeholder="You are tasked with upgrading this repo to use Svelte5. Do this in as few steps as possible, ensuring compilation, and implementing best practices according to Svelte 5 docs."
        defaultValue={basePrompt}
      ></textarea>
      
      {/* Component 4: Randomness slider with marks */}
      <div className="slider-container">
        <div className="slider-label">
          <span>Iterations</span>
          <div className="tooltip">
            <span className="info-icon">ⓘ</span>
            <span className="tooltip-text">
              Slide to increase the variations of your prompt. Maximum 5 variations when slider is at 5.
            </span>
          </div>
        </div>
        <div className="slider-wrapper">
          <input 
            type="range" 
            min="0" 
            max="5" 
            value={randomnessLevel} 
            onChange={(e) => setRandomnessLevel(parseInt(e.target.value))}
            className="randomness-slider"
          />
          <div className="slider-marks">
            {[0, 1, 2, 3, 4, 5].map((mark) => (
              <div key={mark} className="slider-mark">
                <span>{mark}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
      
      {/* Display prompt variations based on slider position */}
      {randomnessLevel > 0 && (
        <div className="prompt-variations">
          {promptVariations.map((variation, index) => (
            <div key={index} className="prompt-variation">
              <span className="variation-label">Iteration {index + 2}</span>
              <p>{variation}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default PromptSection;