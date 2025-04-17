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
        placeholder="Get me all the files"
      />
    </div>
  );
};

export default SearchBar;