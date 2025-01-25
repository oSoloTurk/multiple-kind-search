import axios from './config';

export interface Entry {
  id?: string;
  title: string;
  content: string;
  author: string;
}

export const entriesApi = {
  // Get a single entry by ID
  getEntry: async (id: string) => {
    const response = await axios.get<Entry>(`/api/entries/${id}`);
    return response.data;
  },

  // Create a new entry
  createEntry: async (entry: Entry) => {
    const response = await axios.post<Entry>('/api/entries', entry);
    return response.data;
  },

  // Update an existing entry
  updateEntry: async (id: string, entry: Entry) => {
    const response = await axios.put<Entry>(`/api/entries/${id}`, entry);
    return response.data;
  },

  // Search entries
  searchEntries: async (query: string) => {
    const response = await axios.get<Entry[]>(`/api/search`, {
      params: { q: query }
    });
    return response.data;
  }
};
