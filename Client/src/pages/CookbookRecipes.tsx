import { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router';
import RecipeCard from '../components/RecipeCard';
import EditCookbookModal from '../components/EditCookbookModal';
import DeleteCookbookModal from '../components/DeleteCookbookModal';
import type { Recipe, Cookbook } from '../types';
import { getCookbookRecipes } from '../api/cookbooks';

const CookbookRecipes: React.FC = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [cookbook, setCookbook] = useState<Cookbook | null>(null);
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  useEffect(() => {
    async function fetchCookbookRecipes() {
      try {
        const data = await getCookbookRecipes(id!);
        setCookbook(data.cookbook);
        setRecipes(data.recipes);
      } catch {
        setError('Failed to load cookbook recipes.');
      } finally {
        setLoading(false);
      }
    }

    fetchCookbookRecipes();
  }, [id]);

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto p-6 text-center text-gray-500">
        Loading...
      </div>
    );
  }

  if (error || !cookbook) {
    return (
      <div className="max-w-4xl mx-auto p-6 text-center text-red-500">
        {error ?? 'Cookbook not found.'}
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto p-6 font-sans">
      <Link to="/profile" className="text-green-600 hover:underline text-sm mb-4 inline-block">
        ← Back to Profile
      </Link>

      <div className="bg-white shadow-md rounded-2xl p-6 mb-8">
        <div className="flex items-center justify-between mb-2">
          <h1 className="text-2xl font-bold">{cookbook.name}</h1>
          <div className="flex items-center gap-2">
            <span className={`text-xs px-2 py-1 rounded-full ${
              cookbook.is_public ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'
            }`}>
              {cookbook.is_public ? 'Public' : 'Private'}
            </span>
            <button
              onClick={() => setShowEditModal(true)}
              className="px-3 py-1 text-sm bg-green-600 text-white rounded-md hover:bg-green-700 transition"
            >
              Edit
            </button>
            <button
              onClick={() => setShowDeleteModal(true)}
              className="px-3 py-1 text-sm bg-red-600 text-white rounded-md hover:bg-red-700 transition"
            >
              Delete
            </button>
          </div>
        </div>
        {cookbook.description && (
          <p className="text-gray-600">{cookbook.description}</p>
        )}
      </div>

      <h2 className="text-xl font-semibold mb-4">Recipes</h2>

      {recipes.length === 0 ? (
        <p className="text-gray-500">No recipes in this cookbook yet.</p>
      ) : (
        <div className="grid gap-6">
          {recipes.map((recipe) => (
            <RecipeCard key={recipe.id} recipe={recipe} />
          ))}
        </div>
      )}

      {showEditModal && (
        <EditCookbookModal
          cookbook={cookbook}
          onClose={() => setShowEditModal(false)}
          onUpdated={(updated) => {
            setCookbook(updated);
            setShowEditModal(false);
          }}
        />
      )}

      {showDeleteModal && (
        <DeleteCookbookModal
          cookbook={cookbook}
          onClose={() => setShowDeleteModal(false)}
          onDeleted={() => navigate('/profile')}
        />
      )}
    </div>
  );
};

export default CookbookRecipes;
