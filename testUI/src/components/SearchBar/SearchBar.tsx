import React from 'react';
import './SearchBar.css';

const SearchBar: React.FC = () => {
  return (
    <div className="search-container">
      <label htmlFor="search">Search</label>
      {/* Component 2: Search box */}
      <input 
        type="text" 
        id="search" 
        className="search-input" 
        placeholder="Find all of my repositories that use a Node version < 22"
      />
    </div>
  );
};

export default SearchBar;