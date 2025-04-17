import React, { useState } from 'react';
import './PromptSection.css';

const PromptSection: React.FC = () => {
  const [randomnessLevel, setRandomnessLevel] = useState<number>(0);
  
  // Sample generated prompts based on randomness level
  const generatePrompts = (level: number, basePrompt: string) => {
    const variations = [
      'You are tasked with upgrading all repositories to Node 24. Do this in as few steps as possible, ensuring compilation.',
      'Upgrade all repos to Node 24 with minimal steps while maintaining compilation integrity.',
      'Your task: update all repositories to Node version 24, using the fewest possible steps and ensuring successful compilation.',
      'Convert all repositories to Node 24. Minimize steps and ensure code compiles correctly.',
      'Migrate all repositories to Node 24 with minimal effort while maintaining build integrity.'
    ];
    
    return variations.slice(0, level + 1);
  };

  const basePrompt = 'You are tasked with upgrading all repositories to Node 24. Do this in as few steps as possible, ensuring compilation.';
  const promptVariations = generatePrompts(randomnessLevel, basePrompt);

  return (
    <div className="prompt-section">
      <label htmlFor="prompt">Prompt</label>
      {/* Component 3: Natural language prompt */}
      <textarea 
        id="prompt" 
        className="prompt-textarea" 
        placeholder="You are tasked with upgrading all repositories to Node 24. Do this in as few steps as possible, ensuring compilation."
        defaultValue={basePrompt}
      ></textarea>
      
      {/* Component 4: Randomness slider */}
      <div className="slider-container">
        <div className="slider-label">
          <span>Randomness</span>
          <div className="tooltip">
            <span className="info-icon">ⓘ</span>
            <span className="tooltip-text">
              Slide to increase the variations of your prompt. Maximum 5 variations when slider is at 5.
            </span>
          </div>
        </div>
        <input 
          type="range" 
          min="0" 
          max="5" 
          value={randomnessLevel} 
          onChange={(e) => setRandomnessLevel(parseInt(e.target.value))}
          className="randomness-slider"
        />
      </div>
      
      {/* Display prompt variations based on slider position */}
      {randomnessLevel > 0 && (
        <div className="prompt-variations">
          {promptVariations.map((variation, index) => (
            <div key={index} className="prompt-variation">
              <span className="variation-label">Prompt-{index + 1}</span>
              <p>{variation}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default PromptSection;