import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import SearchPage from './components/SearchPage';
import UpsertPage from './components/UpsertPage';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        <nav className="nav-bar">
          <ul>
            <li><Link to="/">Search</Link></li>
            <li><Link to="/edit">New Entry</Link></li>
          </ul>
        </nav>

        <Routes>
          <Route path="/" element={<SearchPage />} />
          <Route path="/edit" element={<UpsertPage />} />
          <Route path="/edit/:id" element={<UpsertPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App; 