import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { newsApi, authorApi, News, Author } from '../api/api';
import { Button, CircularProgress } from '@mui/material';
import './ListPage.css';
import { ThemeProvider, createTheme } from "@mui/material/styles";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardActionArea from "@mui/material/CardActionArea";
import CardMedia from "@mui/material/CardMedia";
import Chip from "@mui/material/Chip";

type EntityType = 'news' | 'authors';

const darkTheme = createTheme({
  palette: {
    mode: "dark",
    background: {
      default: "#121212",
      paper: "#1e1e1e",
    },
    text: {
      primary: "#ffffff",
      secondary: "#b0b0b0",
    },
    primary: {
      main: "#bb86fc",
    },
    secondary: {
      main: "#03dac6",
    },
  },
});

const ListPage: React.FC = () => {
  const { type } = useParams<{ type: EntityType }>();
  const navigate = useNavigate();
  const [items, setItems] = useState<(News | Author)[]>([]);
  const [loading, setLoading] = useState(true);
  const searchTerm = new URLSearchParams(window.location.search).get('query') || '';

  useEffect(() => {
    loadItems();
  }, [type, searchTerm]);

  const loadItems = async () => {
    setItems([]);
    setLoading(true);
    try {
      let data: (News | Author)[] = [];
      if (type === 'news') {
        data = await newsApi.listNews();
      } else if (type === 'authors') {
        data = await authorApi.listAuthors();
      }

      if (searchTerm) {
        const lowercasedTerm = searchTerm.toLowerCase();
        data = data.filter(item => {
          if (type === 'news') {
            return (item as News).title.toLowerCase().includes(lowercasedTerm) ||
                   (item as News).content.toLowerCase().includes(lowercasedTerm);
          } else {
            return (item as Author).name.toLowerCase().includes(lowercasedTerm) ||
                   (item as Author).bio?.toLowerCase().includes(lowercasedTerm) || false;
          }
        });
      }

      setItems(data);
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
    <ThemeProvider theme={darkTheme}>
      <Box sx={{ p: 4, bgcolor: "background.default", color: "text.primary" }}>
        <div className="list-page">
          <Box className="list-header" display="flex" justifyContent="space-between" alignItems="center" mb={4}>
            <Typography variant="h4">
              {type.charAt(0).toUpperCase() + type.slice(1)} List
            </Typography>
            <Button variant="contained" color="primary" onClick={handleCreate}>
              Create New {type === 'news' ? 'Article' : 'Author'}
            </Button>
          </Box>

          {loading ? (
            <CircularProgress />
          ) : (
            <Grid container spacing={3}>
              {items.map((item) => (
                <Grid item xs={12} sm={6} md={4} key={item.id}>
                  <Card sx={{ bgcolor: "background.paper" }}>
                    <CardActionArea>
                      {type === 'news' ? (
                        <>
                          <CardMedia
                            component="img"
                            height="140"
                            image={(item as News).imageUrl || "https://via.placeholder.com/150"}
                            alt={(item as News).title}
                          />
                          <CardContent>
                            <Typography variant="h6" gutterBottom>
                              {(item as News).title}
                            </Typography>
                            <Typography variant="body2" color="text.secondary" gutterBottom>
                              By {/* Add author name if available */}
                            </Typography>
                            <Typography variant="body2" color="text.primary">
                              {(item as News).content?.substring(0, 150)}...
                            </Typography>
                            <Box mt={2} display="flex" flexWrap="wrap" gap={1}>
                              {/* Display tags if available */}
                            </Box>
                          </CardContent>
                        </>
                      ) : (
                        <>
                          <CardMedia
                            component="img"
                            height="140"
                            image={(item as Author).imageUrl || "https://via.placeholder.com/150"}
                            alt={(item as Author).name}
                          />
                          <CardContent>
                            <Typography variant="h6" gutterBottom>
                              {(item as Author).name}
                            </Typography>
                            <Typography variant="body2" color="text.primary">
                              {(item as Author).bio?.substring(0, 150)}...
                            </Typography>
                          </CardContent>
                        </>
                      )}
                    </CardActionArea>
                    <Box display="flex" justifyContent="space-between" p={2}>
                      <Button variant="contained" color="primary" onClick={() => handleEdit(item.id!)}>
                        Edit
                      </Button>
                      <Button variant="contained" color="secondary" onClick={() => handleDelete(item.id!)}>
                        Delete
                      </Button>
                    </Box>
                  </Card>
                </Grid>
              ))}
            </Grid>
          )}
        </div>
      </Box>
    </ThemeProvider>
  );
};

export default ListPage; 