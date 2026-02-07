import { Link } from 'react-router';
import { useAuth } from '../context/AuthContext';
import Avatar from './Avatar';

const Navbar: React.FC = () => {
  const { user, logout } = useAuth();
  return (
    <nav className="bg-white shadow-md w-full px-6 py-4">
      <div className="max-w-7xl mx-auto flex justify-between items-center flex-wrap">
        {/* Left links */}
        <div className="flex items-center gap-6">
          { /* any <a> tag should be replace by a Link component */}
          <Link to="/" className="text-lg font-medium hover:text-blue-600 transition">
            Recipes
          </Link>
          <Link to="/post" className="text-lg font-medium hover:text-blue-600 transition">
            Post
          </Link>
        </div>

        {/* Right controls */}
        <div className="flex items-center gap-4 mt-4 md:mt-0">
          {user ? (
            // If user is authenticated
            <>
              <Link to="/profile">
                <Avatar size={8} />
              </Link>
              <button onClick={logout} className="px-3 py-1 text-sm text-red-600 border border-red-600 rounded-md hover:bg-red-50 transition">
                  Logout
              </button>
            </>
          ) : (
            // If no user is logged in
            <>
              <Link to="/login" className="px-3 py-1 text-sm text-blue-600 border border-blue-600 rounded-md hover:bg-blue-50 transition">
                  Login
              </Link>
              <Link to="/signup" className="px-3 py-1 text-sm text-white bg-blue-600 rounded-md hover:bg-blue-700 transition">
                Register
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;

