import { useState } from 'react';
import './SearchBar.css';

const SearchBar = ({ onSearch, onClear, currentSearchQuery = '' }) => {
  const [query, setQuery] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (query.trim()) {
      onSearch(query.trim());
    }
  };

  const handleClear = () => {
    setQuery('');
    onClear();
  };

  const isSearchActive = currentSearchQuery !== '';

  return (
    <div className="search-bar-container">
      <form onSubmit={handleSubmit} className="search-form">
        <input
          type="text"
          className="search-input"
          placeholder="Search posts..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button type="submit" className="search-btn" disabled={!query.trim()}>
          Search
        </button>
        {(query || isSearchActive) && (
          <button type="button" className="clear-btn" onClick={handleClear}>
            Clear
          </button>
        )}
      </form>
    </div>
  );
};

export default SearchBar;
