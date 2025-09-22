import React, { useState } from "react";
import { Search } from "lucide-react"; // or any icon library you like

interface SearchBarProps {
  onSearch: (query: string) => void;
}

const SearchBar: React.FC<SearchBarProps> = ({ onSearch }) => {
  const [query, setQuery] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(query.trim());
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="mx-auto flex w-full max-w-lg items-center rounded-md border border-gray-600 bg-gray-800 px-3 py-2 shadow-sm focus-within:ring-2 focus-within:ring-indigo-500"
    >
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search articles or topics..."
        className="flex-1 bg-transparent text-gray-100 placeholder-gray-400 focus:outline-none"
      />
      <button
        type="submit"
        aria-label="Search"
        className="ml-2 text-gray-300 hover:text-white"
      >
        <Search size={20} />
      </button>
    </form>
  );
};

export default SearchBar;
