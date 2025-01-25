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
  ConditionalContents,
  InsertFrontmatter,
  sandpackPlugin,
  codeMirrorPlugin,
  KitchenSinkToolbar,
  InsertAdmonition,
} from '@mdxeditor/editor';
import { newsApi, authorApi, News, Author } from '../api/api';

type EntityType = 'news' | 'authors';

function UpsertPage() {
  const { type, id } = useParams<{ type: EntityType; id: string }>();
  const navigate = useNavigate();
  const [news, setNews] = useState<News>({
    title: '',
    content: '',
    authorID: '',
    tags: [],
    imageUrl: ''
  });
  const [author, setAuthor] = useState<Author>({
    name: '',
    bio: '',
    imageUrl: ''
  });
  const [availableAuthors, setAvailableAuthors] = useState<Author[]>([]);
  const [tagInput, setTagInput] = useState('');
  const ref = React.useRef<MDXEditorMethods>(null)

  useEffect(() => {
    if (type === 'news') {
      loadAuthors();
    }
    if (id) {
      fetchEntity();
    }
  }, [id, type]);

  const loadAuthors = async () => {
    try {
      const authors = await authorApi.listAuthors();
      setAvailableAuthors(authors);
    } catch (error) {
      console.error('Error loading authors:', error);
    }
  };

  const fetchEntity = async () => {
    try {
      if (type === 'news' && id) {
        const data = await newsApi.getNews(id);
        setNews(data);
      } else if (type === 'authors' && id) {
        const data = await authorApi.getAuthor(id);
        setAuthor(data);
      }
    } catch (error) {
      console.error('Error fetching entity:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      if (type === 'news') {
        if (id) {
          await newsApi.updateNews(id, news);
        } else {
          await newsApi.createNews(news);
        }
      } else if (type === 'authors') {
        if (id) {
          await authorApi.updateAuthor(id, author);
        } else {
          await authorApi.createAuthor(author);
        }
      }
      navigate('/list/' + type);
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

  const handleAuthorChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setAuthor(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleContentChange = (content: string) => {
    if (type === 'news') {
      setNews(prev => ({
        ...prev,
        content
      }));
    } else if (type === 'authors') {
      setAuthor(prev => ({
        ...prev,
        bio: content
      }));
    }
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

  // set the markdown value to the news.content
  useEffect(() => {
    if (ref.current) {
      ref.current.setMarkdown(news.content || '');
    }
  }, [news.content]);

  // set the markdown value to the author.bio
  useEffect(() => {
    if (ref.current) {
      ref.current.setMarkdown(author.bio || '');
    }
  }, [author.bio]);

  // Define markdown plugins
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

  if (!type) {
    return <div>Invalid entity type</div>;
  }

  return (
    <div className="upsert-page">
      <h1>{id ? `Edit ${type}` : `Create New ${type}`}</h1>
      <form onSubmit={handleSubmit}>
        {type === 'news' ? (
          <>
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
              <label htmlFor="authorId">Author</label>
              <select
                id="authorId"
                name="authorId"
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
          </>
        ) : (
          <>
            <div className="form-group">
              <label htmlFor="name">Name</label>
              <input
                type="text"
                id="name"
                name="name"
                className="form-control"
                value={author.name}
                onChange={handleAuthorChange}
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="imageUrl">Profile Image URL</label>
              <input
                type="url"
                id="imageUrl"
                name="imageUrl"
                className="form-control"
                value={author.imageUrl}
                onChange={handleAuthorChange}
                placeholder="https://example.com/profile.jpg"
              />
              {author.imageUrl && (
                <img 
                  src={author.imageUrl} 
                  alt="Profile Preview" 
                  className="image-preview profile"
                />
              )}
            </div>

            <div className="form-group">
              <label htmlFor="bio">Bio</label>
              <MDXEditor
                markdown={author.bio || ''}
                onChange={handleContentChange}
                plugins={plugins}
                contentEditableClassName="mdx-editor-content"
                className="mdxeditor"
                ref={ref}
              />
            </div>
          </>
        )}

        <div className="form-actions">
          <button type="button" className="cancel-button" onClick={() => navigate('/list/' + type)}>
            Cancel
          </button>
          <button type="submit" className="save-button">
            {id ? 'Update' : 'Save'} {type}
          </button>
        </div>
      </form>
    </div>
  );
}

export default UpsertPage;
