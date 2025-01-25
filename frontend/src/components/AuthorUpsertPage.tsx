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
import { authorApi, Author } from '../api/api';

function AuthorUpsertPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [author, setAuthor] = useState<Author>({
    name: '',
    bio: '',
    imageUrl: ''
  });
  const ref = React.useRef<MDXEditorMethods>(null)

  useEffect(() => {
    if (id) {
      fetchAuthor();
    }
  }, [id]);

  const fetchAuthor = async () => {
    try {
      if (id) {
        const data = await authorApi.getAuthor(id);
        setAuthor(data);
      }
    } catch (error) {
      console.error('Error fetching author:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (id) {
        await authorApi.updateAuthor(id, author);
      } else {
        await authorApi.createAuthor(author);
      }
      navigate('/list/authors');
    } catch (error) {
      console.error('Error saving:', error);
    }
  };

  const handleAuthorChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setAuthor(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleContentChange = (content: string) => {
    setAuthor(prev => ({
      ...prev,
      bio: content
    }));
  };

  useEffect(() => {
    if (ref.current) {
      ref.current.setMarkdown(author.bio || '');
    }
  }, [author.bio]);

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
      <h1>{id ? 'Edit Author' : 'Create New Author'}</h1>
      <form onSubmit={handleSubmit}>
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

        <div className="form-actions">
          <button type="button" className="cancel-button" onClick={() => navigate('/list/authors')}>
            Cancel
          </button>
          <button type="submit" className="save-button">
            {id ? 'Update' : 'Save'} Author
          </button>
        </div>
      </form>
    </div>
  );
}

export default AuthorUpsertPage; 