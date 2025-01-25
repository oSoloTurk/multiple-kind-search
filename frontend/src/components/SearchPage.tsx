import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { TextField, Button, CircularProgress } from '@mui/material';
import './SearchPage.css';
import { searchApi, SearchResult } from '../api/api';

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [username, setUsername] = useState('');
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSearch = async (searchQuery: string) => {
    if (!username) {
      alert('Username is required');
      return;
    }
    
    setIsLoading(true);
    try {
      const data = await searchApi.search({ q: searchQuery, username });
      // Filter only news results
      setResults((data || []))
    } catch (error) {
      console.error('Error searching:', error);
      setResults([]);
    }
    setIsLoading(false);
  };

  const handleEdit = (result: SearchResult) => {
    navigate(`/edit/news/${result.id}`);
  };

  return (
    <div className="search-container">
      <h1>News Search</h1>
      <div className="search-box">
        <TextField
          variant="outlined"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Enter username..."
          className="search-input"
        />
        <TextField
          variant="outlined"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for news..."
          className="search-input"
        />
        <Button 
          variant="contained" 
          color="primary" 
          onClick={() => handleSearch(query)}
        >
          Search
        </Button>
      </div>

      {isLoading ? (
        <CircularProgress />
      ) : (
        results.length > 0 ? (
        <div className="results-container">
          {results.map((result: SearchResult) => (
            <div key={result.id} className="result-card">
              {result.type === 'author' ? (
                <div className="author-card">
                  <h2 dangerouslySetInnerHTML={{ __html: result.title }} />
                  <p  dangerouslySetInnerHTML={{ __html: result.content }} />
                  <Button onClick={() => handleEdit(result)}>Edit Author</Button>
                </div>
              ) : (
                <div className="news-card">
                  <h2 dangerouslySetInnerHTML={{ __html: result.title }} />
                  <p  dangerouslySetInnerHTML={{ __html: result.content }} />
                  <Button onClick={() => handleEdit(result)}>Edit News</Button>
                </div>
              )}
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