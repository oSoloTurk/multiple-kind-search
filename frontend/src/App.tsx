import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation, Navigate } from 'react-router-dom';
import EditPage from './components/UpsertPage';
import './App.css';
import { Box, IconButton } from '@mui/material';
import { InputAdornment } from '@mui/material';
import { TextField } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import ListPage from './components/ListPage';

const NavBar: React.FC = () => {
  const location = useLocation();
  const isActive = (path: string) => location.pathname === path ? 'active' : '';

  const [searchTerm, setSearchTerm] = useState(new URLSearchParams(window.location.search).get('query') || '');

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const queryParams = new URLSearchParams(window.location.search);
    queryParams.set('query', searchTerm);
    window.history.pushState({}, '', `${location.pathname}?${queryParams}`);
  };

  useEffect(() => {
    const query = new URLSearchParams(window.location.search).get('query');
    if (query) {
      setSearchTerm(query);
    }
  }, [location.search]);


  return (
    <nav className="nav-bar">
      <ul>
        <li><Link to="/list/news" className={isActive('/list/news')}>News</Link></li>
        <li><Link to="/list/authors" className={isActive('/list/authors')}>Authors</Link></li>
      </ul>
      <ul className="search-box">
        <form onSubmit={handleSearchSubmit}>
          <TextField
            variant="outlined"
            size="small"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Search..."
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton type="submit" color="inherit">
                    <SearchIcon />
                  </IconButton>
                </InputAdornment>
              ),
            }}
            sx={{ backgroundColor: 'white', borderRadius: 1 }}
          />
        </form>
      </ul>
    </nav>
  );
};

const App: React.FC = () => {
  return (
    <Router>
      <div className="App">
        <NavBar />
        <Routes>
          <Route path="/" element={<Navigate to="/list/news" />} />
          <Route path="/list/:type" element={<ListPage />} />
          <Route path="/edit/:type/:id" element={<EditPage />} />
          <Route path="/edit/:type" element={<EditPage />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App; 