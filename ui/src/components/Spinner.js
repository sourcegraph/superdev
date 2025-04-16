import React from 'react';
import './Spinner.css';

const Spinner = ({ imageUrl, alt = 'Loading...' }) => {
  return (
    <div className="spinner-container">
      <img 
        src={imageUrl} 
        alt={alt} 
        className="spinner-image" 
      />
    </div>
  );
};

export default Spinner;