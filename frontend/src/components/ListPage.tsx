import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { newsApi, authorApi, News, Author } from '../api/api';
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
        <button className="create-button" onClick={handleCreate}>
          Create New {type === 'news' ? 'Article' : 'Author'}
        </button>
      </div>

      {loading ? (
        <div className="loading">Loading...</div>
      ) : (
        <div className="items-grid">
          {items.map((item) => (
            <div key={item.id} className="item-card">
              {type === 'news' ? (
                // News card
                <>
                  <h3>{(item as News).title}</h3>
                  <div className="item-preview">
                    {(item as News).content.substring(0, 150)}...
                  </div>
                </>
              ) : (
                // Author card
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
                <button 
                  className="edit-button"
                  onClick={() => handleEdit(item.id!)}
                >
                  Edit
                </button>
                <button 
                  className="delete-button"
                  onClick={() => handleDelete(item.id!)}
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ListPage; 