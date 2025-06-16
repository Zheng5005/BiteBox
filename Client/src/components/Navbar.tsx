import React, { useState } from 'react';

interface User {
  name: string;
  avatarUrl: string;
}

const Navbar: React.FC = () => {
  const [user, setUser] = useState<User | null>(null); // Simulate auth
  const [lang, setLang] = useState<'en' | 'es'>('en');

  const toggleLanguage = () => {
    setLang((prev) => (prev === 'en' ? 'es' : 'en'));
  };

  return (
    <nav className="bg-white shadow-md w-full px-6 py-4">
      <div className="max-w-7xl mx-auto flex justify-between items-center flex-wrap">
        {/* Left links */}
        <div className="flex items-center gap-6">
          { /* any <a> tag should be replace by a Link component */}
          <a href="/recipes" className="text-lg font-medium hover:text-blue-600 transition">
            Recipes
          </a>
          <a href="/post" className="text-lg font-medium hover:text-blue-600 transition">
            Post
          </a>
        </div>

        {/* Right controls */}
        <div className="flex items-center gap-4 mt-4 md:mt-0">
          {/* Language toggle */}
          <button
            onClick={toggleLanguage}
            className="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-100 transition"
          >
            {lang.toUpperCase()}
          </button>

          {user ? (
            // If user is authenticated
            <img
              src={user.avatarUrl}
              alt="User avatar"
              className="w-9 h-9 rounded-full object-cover border"
            />
          ) : (
            // If no user is logged in
            <>
              <button className="px-3 py-1 text-sm text-blue-600 border border-blue-600 rounded-md hover:bg-blue-50 transition">
                Login
              </button>
              <button className="px-3 py-1 text-sm text-white bg-blue-600 rounded-md hover:bg-blue-700 transition">
                Register
              </button>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;

