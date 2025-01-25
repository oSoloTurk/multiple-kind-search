import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { newsApi, authorApi, News, Author } from '../api/api';
import { Button, CircularProgress } from '@mui/material';
import './ListPage.css';

type EntityType = 'news' | 'authors';

const ListPage: React.FC = () => {
  const { type } = useParams<{ type: EntityType }>();
  const navigate = useNavigate();
  const [items, setItems] = useState<(News | Author)[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadItems();
  }, [type]);

  const loadItems = async () => {
    setItems([]);
    setLoading(true);
    try {
      if (type === 'news') {
        const data = await newsApi.listNews();
        setItems(data);
      } else if (type === 'authors') {
        const data = await authorApi.listAuthors();
        setItems(data);
      }
    } catch (error) {
      console.error('Error loading items:', error);
    }
    setLoading(false);
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this item?')) {
      return;
    }

    try {
      if (type === 'news') {
        await newsApi.deleteNews(id);
      } else if (type === 'authors') {
        await authorApi.deleteAuthor(id);
      }
      await loadItems();
    } catch (error) {
      console.error('Error deleting item:', error);
    }
  };

  const handleEdit = (id: string) => {
    navigate(`/edit/${type}/${id}`);
  };

  const handleCreate = () => {
    navigate(`/edit/${type}`);
  };

  if (!type) {
    return <div>Invalid entity type</div>;
  }

  return (
    <div className="list-page">
      <div className="list-header">
        <h1>{type.charAt(0).toUpperCase() + type.slice(1)} List</h1>
        <Button variant="contained" color="primary" onClick={handleCreate}>
          Create New {type === 'news' ? 'Article' : 'Author'}
        </Button>
      </div>

      {loading ? (
        <CircularProgress />
      ) : (
        <div className="items-grid">
          {items.map((item) => (
            <div key={item.id} className="item-card">
              {type === 'news' ? (
                <>
                  <h3>{(item as News).title}</h3>
                  <div className="item-preview">
                    {(item as News).content?.substring(0, 150)}...
                  </div>
                </>
              ) : (
                <>
                  <h3>{(item as Author).name}</h3>
                  {(item as Author).imageUrl && (
                    <img 
                      src={(item as Author).imageUrl} 
                      alt={(item as Author).name}
                      className="author-image"
                    />
                  )}
                  <div className="item-preview">
                    {(item as Author).bio?.substring(0, 150)}...
                  </div>
                </>
              )}
              <div className="item-actions">
                <Button 
                  variant="contained" 
                  color="primary" 
                  onClick={() => handleEdit(item.id!)}
                >
                  Edit
                </Button>
                <Button 
                  variant="contained" 
                  color="secondary" 
                  onClick={() => handleDelete(item.id!)}
                >
                  Delete
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ListPage; 