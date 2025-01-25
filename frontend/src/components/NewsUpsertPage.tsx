import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './UpsertPage.css';
import '@mdxeditor/editor/style.css'
import {
  MDXEditor,
  headingsPlugin,
  listsPlugin,
  quotePlugin,
  thematicBreakPlugin,
  markdownShortcutPlugin,
  linkPlugin,
  linkDialogPlugin,
  imagePlugin,
  tablePlugin,
  codeBlockPlugin,
  frontmatterPlugin,
  AdmonitionDirectiveDescriptor,
  directivesPlugin,
  diffSourcePlugin,
  MDXEditorMethods,
  toolbarPlugin,
  BlockTypeSelect,
  BoldItalicUnderlineToggles,
  CodeToggle,
  CreateLink,
  InsertImage,
  InsertTable,
  InsertThematicBreak,
  ListsToggle,
  UndoRedo,
  InsertCodeBlock,
  InsertFrontmatter,
  sandpackPlugin,
  codeMirrorPlugin,
  InsertAdmonition,
} from '@mdxeditor/editor';
import { newsApi, authorApi, News, Author } from '../api/api';

function NewsUpsertPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [news, setNews] = useState<News>({
    title: '',
    content: '',
    authorID: '',
    tags: [],
    imageUrl: ''
  });
  const [availableAuthors, setAvailableAuthors] = useState<Author[]>([]);
  const [tagInput, setTagInput] = useState('');
  const ref = React.useRef<MDXEditorMethods>(null)

  useEffect(() => {
    loadAuthors();
    if (id) {
      fetchNews();
    }
  }, [id]);

  const loadAuthors = async () => {
    try {
      const authors = await authorApi.listAuthors();
      setAvailableAuthors(authors);
    } catch (error) {
      console.error('Error loading authors:', error);
    }
  };

  const fetchNews = async () => {
    try {
      if (id) {
        const data = await newsApi.getNews(id);
        setNews(data);
      }
    } catch (error) {
      console.error('Error fetching news:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (id) {
        await newsApi.updateNews(id, news);
      } else {
        await newsApi.createNews(news);
      }
      navigate('/list/news');
    } catch (error) {
      console.error('Error saving:', error);
    }
  };

  const handleNewsChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setNews(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleContentChange = (content: string) => {
    setNews(prev => ({
      ...prev,
      content
    }));
  };

  const handleAddTag = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && tagInput.trim()) {
      e.preventDefault();
      const newTag = tagInput.trim();
      const currentTags = news.tags || [];
      if (!currentTags.includes(newTag)) {
        setNews(prev => ({
          ...prev,
          tags: [...(prev.tags || []), newTag]
        }));
      }
      setTagInput('');
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setNews(prev => ({
      ...prev,
      tags: (prev.tags || []).filter(tag => tag !== tagToRemove)
    }));
  };

  useEffect(() => {
    if (ref.current) {
      ref.current.setMarkdown(news.content || '');
    }
  }, [news.content]);

  const plugins = [
    headingsPlugin(),
    listsPlugin(),
    quotePlugin(),
    thematicBreakPlugin(),
    markdownShortcutPlugin(),
    linkPlugin(),
    linkDialogPlugin(),
    imagePlugin(),
    tablePlugin(),
    codeBlockPlugin(),
    frontmatterPlugin(),
    directivesPlugin({
      directiveDescriptors: [AdmonitionDirectiveDescriptor],
    }),
    diffSourcePlugin(),
    toolbarPlugin({
      toolbarContents: () => (
        <>
          <UndoRedo />
          <BlockTypeSelect />
          <BoldItalicUnderlineToggles />
          <CodeToggle />
          <CreateLink />
          <InsertImage />
          <InsertTable />
          <InsertThematicBreak />
          <ListsToggle />
          <InsertCodeBlock />
          <InsertFrontmatter />
          <InsertAdmonition />
        </>
      )
    }),
    sandpackPlugin(),
    codeMirrorPlugin({
      codeBlockLanguages: {
        js: 'JavaScript',
        jsx: 'JavaScript React',
        ts: 'TypeScript',
        tsx: 'TypeScript React',
        python: 'Python',
        go: 'Go',
        rust: 'Rust',
        sql: 'SQL',
        json: 'JSON',
        html: 'HTML',
        css: 'CSS',
      }
    }),
  ];

  return (
    <div className="upsert-page">
      <h1>{id ? 'Edit News' : 'Create New News'}</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="title">Title</label>
          <input
            type="text"
            id="title"
            name="title"
            className="form-control"
            value={news.title}
            onChange={handleNewsChange}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="authorID">Author</label>
          <select
            id="authorID"
            name="authorID"
            className="form-control"
            value={news.authorID}
            onChange={handleNewsChange}
            required
          >
            <option value="">Select an author</option>
            {availableAuthors.map(author => (
              <option key={author.id} value={author.id}>
                {author.name}
              </option>
            ))}
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="imageUrl">Image URL</label>
          <input
            type="url"
            id="imageUrl"
            name="imageUrl"
            className="form-control"
            value={news.imageUrl}
            onChange={handleNewsChange}
            placeholder="https://example.com/image.jpg"
          />
          {news.imageUrl && (
            <img 
              src={news.imageUrl} 
              alt="Preview" 
              className="image-preview"
            />
          )}
        </div>

        <div className="form-group">
          <label htmlFor="tags">Tags</label>
          <div className="tags-input-container">
            <input
              type="text"
              id="tags"
              className="form-control"
              value={tagInput}
              onChange={(e) => setTagInput(e.target.value)}
              onKeyDown={handleAddTag}
              placeholder="Type a tag and press Enter"
            />
            <div className="tags-container">
              {(news.tags || []).map((tag) => (
                <span key={tag} className="tag">
                  {tag}
                  <button
                    type="button"
                    onClick={() => handleRemoveTag(tag)}
                    className="tag-remove"
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="content">Content</label>
          <MDXEditor
            markdown={news.content || ''}
            onChange={handleContentChange}
            plugins={plugins}
            contentEditableClassName="mdx-editor-content"
            ref={ref}
            className="mdxeditor"
          />
        </div>

        <div className="form-actions">
          <button type="button" className="cancel-button" onClick={() => navigate('/list/news')}>
            Cancel
          </button>
          <button type="submit" className="save-button">
            {id ? 'Update' : 'Save'} News
          </button>
        </div>
      </form>
    </div>
  );
}

export default NewsUpsertPage; 