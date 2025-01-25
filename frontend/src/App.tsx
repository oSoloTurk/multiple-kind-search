import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import SearchPage from './components/SearchPage';
import UpsertPage from './components/UpsertPage';
import ListPage from './components/ListPage';
import './App.css';

const NavBar: React.FC = () => {
  const location = useLocation();
  const isActive = (path: string) => location.pathname === path ? 'active' : '';

  return (
    <nav className="nav-bar">
      <ul>
        <li><Link to="/" className={isActive('/')}>Search</Link></li>
        <li><Link to="/list/news" className={isActive('/list/news')}>News</Link></li>
        <li><Link to="/list/authors" className={isActive('/list/authors')}>Authors</Link></li>
      </ul>
    </nav>
  );
};

function App() {
  return (
    <Router>
      <div className="App">
        <NavBar />
        <Routes>
          <Route path="/" element={<SearchPage />} />
          <Route path="/list/:type" element={<ListPage />} />
          <Route path="/edit/:type" element={<UpsertPage />} />
          <Route path="/edit/:type/:id" element={<UpsertPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App; 