import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './SearchPage.css';


const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSearch = async (searchQuery: string) => {
    setIsLoading(true);
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/api/search?q=${encodeURIComponent(searchQuery)}`);
      const data = await response.json();
      if (data === null) {
        setResults([]);
      } else {
        setResults(data);
      }
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
        results.length > 0 ? (
        <div className="results-container">
          {results.map((result) => (
            <div key={result.id} className="result-item">
              <div className="result-header">
                <h3 dangerouslySetInnerHTML={{ __html: result.title }} />
                <button 
                  className="edit-button"
                  onClick={() => navigate(`/edit/${result.id}`)}
                >
                  Edit
                </button>
              </div>
              <div className="author">{result.author}</div>
              <p dangerouslySetInnerHTML={{ __html: result.content }} />
            </div>
          ))}
        </div>
      ) : (
        results === null ? (
          <div className="no-results">No results found</div>
        ) : (
          null
        )
      ))}
    </div>
  );
};

export default SearchPage; 