import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './UpsertPage.css';

interface Entry {
  id?: string;
  title: string;
  content: string;
  author: string;
}

function UpsertPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [entry, setEntry] = useState<Entry>({
    title: '',
    content: '',
    author: ''
  });

  useEffect(() => {
    if (id) {
      // Fetch entry if we're editing an existing one
      fetchEntry(id);
    }
  }, [id]);

  const fetchEntry = async (entryId: string) => {
    try {
      const response = await fetch(`/api/entries/${entryId}`);
      if (response.ok) {
        const data = await response.json();
        setEntry(data);
      }
    } catch (error) {
      console.error('Error fetching entry:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const url = id ? `/api/entries/${id}` : '/api/entries';
    const method = id ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(entry),
      });

      if (response.ok) {
        navigate('/');
      }
    } catch (error) {
      console.error('Error saving entry:', error);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setEntry(prev => ({
      ...prev,
      [name]: value
    }));
  };

  return (
    <div className="upsert-page">
      <h1>{id ? 'Edit Entry' : 'Create New Entry'}</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="title">Title</label>
          <input
            type="text"
            id="title"
            name="title"
            className="form-control"
            value={entry.title}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="author">Author</label>
          <input
            type="text"
            id="author"
            name="author"
            className="form-control"
            value={entry.author}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="content">Content</label>
          <textarea
            id="content"
            name="content"
            className="form-control editor-container"
            value={entry.content}
            onChange={handleChange}
            required
          />
        </div>

        <button type="submit" className="save-button">
          {id ? 'Update' : 'Save'} Entry
        </button>
      </form>
    </div>
  );
}

export default UpsertPage;
