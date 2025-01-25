import axios from './config';

export interface Author {
  id?: string;
  name: string;
  bio?: string;
  imageUrl?: string;
}

export interface News {
  id?: string;
  title: string;
  content: string;
  authorId: string;
  tags?: string[];
  imageUrl?: string;
}

export interface SearchResult {
  id: string;
  title?: string;
  content?: string;
  author?: string;
  type: string;
  highlights?: Record<string, string[]>;
}

export const newsApi = {
  getNews: async (id: string) => {
    const response = await axios.get<News>(`/api/news/${id}`);
    return response.data;
  },

  createNews: async (news: News) => {
    const response = await axios.post<News>('/api/news', news);
    return response.data;
  },

  updateNews: async (id: string, news: News) => {
    const response = await axios.put<News>(`/api/news/${id}`, news);
    return response.data;
  },

  listNews: async () => {
    const response = await axios.get<News[]>('/api/news');
    return response.data;
  },

  deleteNews: async (id: string) => {
    await axios.delete(`/api/news/${id}`);
  }
};

export const authorApi = {
  getAuthor: async (id: string) => {
    const response = await axios.get<Author>(`/api/authors/${id}`);
    return response.data;
  },

  createAuthor: async (author: Author) => {
    const response = await axios.post<Author>('/api/authors', author);
    return response.data;
  },

  updateAuthor: async (id: string, author: Author) => {
    const response = await axios.put<Author>(`/api/authors/${id}`, author);
    return response.data;
  },

  listAuthors: async () => {
    const response = await axios.get<Author[]>('/api/authors');
    return response.data;
  },

  deleteAuthor: async (id: string) => {
    await axios.delete(`/api/authors/${id}`);
  }
};

export const searchApi = {
  search: async ({ q, username }: { q: string; username: string }) => {
    const response = await axios.get<SearchResult[]>('/api/search', {
      params: { q, username }
    });
    return response.data;
  }
};
