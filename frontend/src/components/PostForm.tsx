import React, { useState } from 'react';

interface Props {
  onSubmit: (title: string, content: string, tags: string[]) => void;
}

const PostForm: React.FC<Props> = ({ onSubmit }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [tags, setTags] = useState('');

  return (
    <form
      className="bg-white rounded shadow p-4 mb-6"
      onSubmit={e => {
        e.preventDefault();
        onSubmit(title, content, tags.split(',').map(t => t.trim()).filter(Boolean));
        setTitle('');
        setContent('');
        setTags('');
      }}
    >
      <input
        className="border p-2 w-full mb-2 rounded"
        placeholder="หัวข้อกระทู้"
        value={title}
        onChange={e => setTitle(e.target.value)}
        required
      />
      <textarea
        className="border p-2 w-full mb-2 rounded"
        placeholder="เนื้อหา"
        value={content}
        onChange={e => setContent(e.target.value)}
        required
      />
      <input
        className="border p-2 w-full mb-2 rounded"
        placeholder="แท็ก (คั่นด้วย , )"
        value={tags}
        onChange={e => setTags(e.target.value)}
      />
      <button className="bg-blue-600 text-white px-4 py-2 rounded" type="submit">
        โพสต์
      </button>
    </form>
  );
};

export default PostForm;