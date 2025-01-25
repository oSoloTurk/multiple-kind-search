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
import { entriesApi, Entry } from '../api/api';

function UpsertPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [entry, setEntry] = useState<Entry>({
    title: '',
    content: '',
    author: ''
  });
  const ref = React.useRef<MDXEditorMethods>(null)

  useEffect(() => {
    if (id) {
      // Fetch entry if we're editing an existing one
      fetchEntry(id);
    }
  }, [id]);

  const fetchEntry = async (entryId: string) => {
    try {
      const data = await entriesApi.getEntry(entryId);
      setEntry(data);
    } catch (error) {
      console.error('Error fetching entry:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      if (id) {
        await entriesApi.updateEntry(id, entry);
      } else {
        await entriesApi.createEntry(entry);
      }
      navigate('/');
    } catch (error) {
      console.error('Error saving entry:', error);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setEntry(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleEditorChange = (content: string) => {
    setEntry(prev => ({
      ...prev,
      content
    }));
  };

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

  return (
    <div className="upsert-page">
      <h1>{id ? 'Edit Entry' : 'Create New Entry'}</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="title">Title</label>
          <input
            type="text"
            id="title"
            name="title"
            className="form-control"
            value={entry.title}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="author">Author</label>
          <input
            type="text"
            id="author"
            name="author"
            className="form-control"
            value={entry.author}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
            <label htmlFor="content">Content</label>
            <MDXEditor
                markdown={entry.content}
                onChange={handleEditorChange}
                plugins={plugins}
                contentEditableClassName="mdx-editor-content"
                ref={ref}
                className="mdxeditor"
            />
        </div>

        <button type="submit" className="save-button">
          {id ? 'Update' : 'Save'} Entry
        </button>
      </form>
    </div>
  );
}

export default UpsertPage;
