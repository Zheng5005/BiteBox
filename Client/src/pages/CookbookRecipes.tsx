import { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router';
import RecipeCard from '../components/RecipeCard';
import StickyNote from '../components/StickyNote';
import EditCookbookModal from '../components/EditCookbookModal';
import DeleteCookbookModal from '../components/DeleteCookbookModal';
import type { CookbookRecipe, Cookbook } from '../types';
import { getCookbookRecipes, removeRecipeFromCookbook } from '../api/cookbooks';

const CookbookRecipes: React.FC = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [cookbook, setCookbook] = useState<Cookbook | null>(null);
  const [recipes, setRecipes] = useState<CookbookRecipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [recipeToRemove, setRecipeToRemove] = useState<CookbookRecipe | null>(null);
  const [removing, setRemoving] = useState(false);

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

  const handleRemoveRecipe = async () => {
    if (!recipeToRemove || !cookbook) return;
    setRemoving(true);
    try {
      await removeRecipeFromCookbook(cookbook.id, recipeToRemove.id);
      setRecipes((prev) => prev.filter((r) => r.id !== recipeToRemove.id));
      setRecipeToRemove(null);
    } catch {
      setError('Failed to remove recipe.');
    } finally {
      setRemoving(false);
    }
  };

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
            <div key={recipe.id}>
              <RecipeCard recipe={recipe} onRemove={() => setRecipeToRemove(recipe)} />
              <StickyNote note={recipe.notes} />
            </div>
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

      {recipeToRemove && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={() => setRecipeToRemove(null)}>
          <div
            className="bg-white rounded-2xl shadow-lg w-full max-w-md mx-4 p-6"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-bold">Remove Recipe</h2>
              <button onClick={() => setRecipeToRemove(null)} className="text-gray-400 hover:text-gray-600 text-2xl leading-none">
                &times;
              </button>
            </div>

            <p className="text-gray-700 mb-4">
              Are you sure you want to remove <span className="font-semibold">{recipeToRemove.name_recipe}</span> from this cookbook?
            </p>

            <div className="flex gap-3">
              <button
                onClick={() => setRecipeToRemove(null)}
                disabled={removing}
                className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                onClick={handleRemoveRecipe}
                disabled={removing}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition disabled:opacity-50"
              >
                {removing ? 'Removing...' : 'Remove'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CookbookRecipes;
