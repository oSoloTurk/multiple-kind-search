import React, { useState, useEffect } from 'react';
import './SearchPage.css';


const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const handleSearch = async (searchQuery: string) => {
    setIsLoading(true);
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/api/search?q=${encodeURIComponent(searchQuery)}`);
      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error('Error searching:', error);
    }
    setIsLoading(false);
  };

  return (
    <div className="search-container">
      <h1>Multiple Kind Search</h1>
      <div className="search-box">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for news..."
          className="search-input"
        />
        <button 
          onClick={() => handleSearch(query)}
          className="search-button"
        >
          Search
        </button>
      </div>

      {isLoading ? (
        <div className="loading">Loading...</div>
      ) : (
        <div className="results-container">
          {results.map((result, index) => (
            <div key={index} className="result-item">
              {result.title && (
                <h3 dangerouslySetInnerHTML={{ __html: result.title }} />
              )}
              {result.content && (
                <p dangerouslySetInnerHTML={{ __html: result.content }} />
              )}
              {result.author && (
                <span className="author">By: {result.author}</span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default SearchPage; 