import React from 'react';
import { useParams, Navigate } from 'react-router-dom';
import NewsUpsertPage from './NewsUpsertPage';
import AuthorUpsertPage from './AuthorUpsertPage';

type EntityType = 'news' | 'authors';

function UpsertPage() {
  const { type } = useParams<{ type: EntityType }>();

  if (type === 'news') {
    return <NewsUpsertPage />;
  } else if (type === 'authors') {
    return <AuthorUpsertPage />;
  }

  return <Navigate to="/" replace />;
}

export default UpsertPage;
