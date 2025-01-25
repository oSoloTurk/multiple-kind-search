import React from 'react';
import { useParams, Navigate } from 'react-router-dom';
import NewsUpsertPage from './NewsUpsertPage';
import AuthorUpsertPage from './AuthorUpsertPage';
import { Container } from '@mui/material';
import { ThemeProvider, createTheme } from "@mui/material/styles";
import Box from "@mui/material/Box";

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

function UpsertPage() {
  const { type } = useParams<{ type: EntityType }>();

  if (type === 'news') {
    return <NewsUpsertPage />;
  } else if (type === 'authors') {
    return <AuthorUpsertPage />;
  }

  return (
    <ThemeProvider theme={darkTheme}>
      <Box sx={{ p: 4, bgcolor: "background.default", color: "text.primary" }}>
        <div className="upsert-page">
          <Navigate to="/" replace />
        </div>
      </Box>
    </ThemeProvider>
  );
}

export default UpsertPage;
