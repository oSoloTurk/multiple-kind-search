import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './SearchPage.css';
import { searchApi, SearchResult } from '../api/api';

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSearch = async (searchQuery: string) => {
    setIsLoading(true);
    try {
      const data = await searchApi.search(searchQuery);
      setResults(data || []);
    } catch (error) {
      console.error('Error searching:', error);
      setResults([]);
    }
    setIsLoading(false);
  };

  const handleEdit = (result: SearchResult) => {
    const path = result.type === 'author' ? 'authors' : 'news';
    navigate(`/edit/${path}/${result.id}`);
  };

  return (
    <div className="search-container">
      <h1>Multiple Kind Search</h1>
      <div className="search-box">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for news or authors..."
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
                <h3>{result.title || result.author}</h3>
                <div className="result-type">{result.type}</div>
                <button 
                  className="edit-button"
                  onClick={() => handleEdit(result)}
                >
                  Edit
                </button>
              </div>
              {result.content && (
                <div className="content" dangerouslySetInnerHTML={{ 
                  __html: result.highlights?.content?.[0] || result.content 
                }} />
              )}
              {result.highlights && Object.entries(result.highlights).map(([field, highlights]) => (
                field !== 'content' && (
                  <div key={field} className="highlight">
                    <strong>{field}:</strong>
                    {highlights.map((highlight, i) => (
                      <div key={i} dangerouslySetInnerHTML={{ __html: highlight }} />
                    ))}
                  </div>
                )
              ))}
            </div>
          ))}
        </div>
      ) : (
        query && <div className="no-results">No results found</div>
      ))}
    </div>
  );
}

export default SearchPage; 