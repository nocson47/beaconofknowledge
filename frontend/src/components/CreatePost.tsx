import React, { useState } from 'react';

interface Post {
    title: string;
    body: string;
    tags: string[];
}

interface Props {
    onSubmit: (post: Post) => void;
}

const CreatePost: React.FC<Props> = ({ onSubmit }) => {
    const [title, setTitle] = useState('');
    const [body, setBody] = useState('');
    const [tags, setTags] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const tagArray = tags.split(',').map(tag => tag.trim()).filter(tag => tag !== '');
        onSubmit({ title, body, tags: tagArray });
        setTitle('');
        setBody('');
        setTags('');
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full">
                <form 
                    className="bg-white rounded-lg shadow-xl p-8 space-y-6" 
                    onSubmit={handleSubmit}
                >
                    <h1 className="text-3xl font-bold text-center text-gray-900 mb-8">
                        Create New Thread
                    </h1>
                    <div className="space-y-4">
                        <input
                            className="appearance-none rounded-lg relative block w-full px-4 py-3 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-200"
                            placeholder="Title"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            required
                        />
                        <textarea
                            className="appearance-none rounded-lg relative block w-full px-4 py-3 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-200 min-h-[150px]"
                            placeholder="Body"
                            value={body}
                            onChange={(e) => setBody(e.target.value)}
                            required
                        />
                        <input
                            className="appearance-none rounded-lg relative block w-full px-4 py-3 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-200"
                            placeholder="Tags (separated by commas)"
                            value={tags}
                            onChange={(e) => setTags(e.target.value)}
                        />
                    </div>
                    <button 
                        className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition duration-200"
                        type="submit"
                    >
                        Create Thread
                    </button>
                </form>
            </div>
        </div>
    );
};

export default CreatePost;