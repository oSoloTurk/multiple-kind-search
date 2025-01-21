import React, { useState, useEffect } from 'react';
import './SearchPage.css';

interface SearchSuggestion {
  text: string;
  type: 'title' | 'content' | 'author';
}

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
  const [results, setResults] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchSuggestions = async () => {
      if (query.length < 2) {
        setSuggestions([]);
        return;
      }

      try {
        const response = await fetch(`${process.env.REACT_APP_API_URL}/api/suggest?q=${encodeURIComponent(query)}`);
        const data = await response.json();
        setSuggestions(data);
      } catch (error) {
        console.error('Error fetching suggestions:', error);
      }
    };

    const debounceTimer = setTimeout(fetchSuggestions, 300);
    return () => clearTimeout(debounceTimer);
  }, [query]);

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

      {suggestions.length > 0 && (
        <div className="suggestions-container">
          {suggestions.map((suggestion, index) => (
            <div
              key={index}
              className="suggestion-item"
              onClick={() => {
                setQuery(suggestion.text);
                handleSearch(suggestion.text);
                setSuggestions([]);
              }}
            >
              <span>{suggestion.text}</span>
              <span className="suggestion-type">{suggestion.type}</span>
            </div>
          ))}
        </div>
      )}

      {isLoading ? (
        <div className="loading">Loading...</div>
      ) : (
        <div className="results-container">
          {results.map((result, index) => (
            <div key={index} className="result-item">
              <h3>{result.title}</h3>
              <p>{result.content}</p>
              <span className="author">By: {result.author}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default SearchPage; 