import { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import Avatar from '../components/Avatar';
import RecipeCard from '../components/RecipeCard';
import CreateCookbookModal from '../components/CreateCookbookModal';
import type { Recipe, Cookbook } from '../types';
import { getUserRecipes } from '../api/users';
import { getCookbooks } from '../api/cookbooks';

const Profile: React.FC = () => {
  const { user } = useAuth();
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [cookbooks, setCookbooks] = useState<Cookbook[]>([]);
  const [activeTab, setActiveTab] = useState<string>('recipes');
  const [showCreateModal, setShowCreateModal] = useState(false);

  const refreshCookbooks = async () => {
    try {
      setCookbooks(await getCookbooks());
    } catch {
      console.error('Failed to refresh cookbooks');
    }
  };

  useEffect(() => {
    async function fetchData() {
      const [recipesResult, cookbooksResult] = await Promise.allSettled([
        getUserRecipes(),
        getCookbooks(),
      ]);

      if (recipesResult.status === 'fulfilled') {
        setRecipes(recipesResult.value);
      } else {
      }

      if (cookbooksResult.status === 'fulfilled') {
        setCookbooks(cookbooksResult.value);
      } else {
      }
    }

    if (user) {
      fetchData();
    }
  }, [user]);

  if (!user) {
    return (
      <div className="max-w-4xl mx-auto p-6 text-center text-gray-500">
        Please log in to view your profile.
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto p-6 font-sans">
      <div className="bg-white shadow-md rounded-2xl p-6 flex items-center gap-6 mb-8">
        <Avatar size={20} />
        <div>
          <h1 className="text-2xl font-bold">{user.name}</h1>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex justify-center border-b mb-6">
        <button
          className={`px-4 py-2 font-semibold border-b-2 transition ${
            activeTab === 'recipes' ? 'text-green-400 border-green-600' : 'text-gray-600'
          }`}
          onClick={() => setActiveTab('recipes')}
        >
          My Recipes
        </button>
        <button
          className={`px-4 py-2 ml-4 font-semibold border-b-2 transition ${
            activeTab === 'cookbooks' ? 'text-green-400 border-green-600' : 'text-gray-600'
          }`}
          onClick={() => setActiveTab('cookbooks')}
        >
          My Cookbooks
        </button>
      </div>

      {activeTab === 'recipes' ? (
        recipes.length === 0 ? (
          <p className="text-gray-500">You haven't posted any recipes yet.</p>
        ) : (
          <div className="grid gap-6">
            {recipes.map((recipe) => (
              <RecipeCard key={recipe.id} recipe={recipe} />
            ))}
          </div>
        )
      ) : (
        <>
          <div className="flex justify-end mb-4">
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-4 py-2 text-sm bg-green-600 text-white rounded-md hover:bg-green-700 transition"
            >
              + Create Cookbook
            </button>
          </div>
          {cookbooks.length === 0 ? (
            <p className="text-gray-500">You haven't created any cookbooks yet.</p>
          ) : (
            <div className="grid gap-4">
              {cookbooks.map((cookbook) => (
                <div key={cookbook.id} className="bg-white shadow-md rounded-2xl p-4">
                  <div className="flex items-center justify-between mb-2">
                    <h3 className="text-lg font-bold">{cookbook.name}</h3>
                    <span className={`text-xs px-2 py-1 rounded-full ${
                      cookbook.is_public ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'
                    }`}>
                      {cookbook.is_public ? 'Public' : 'Private'}
                    </span>
                  </div>
                  {cookbook.description && (
                    <p className="text-gray-600 text-sm">{cookbook.description}</p>
                  )}
                </div>
              ))}
            </div>
          )}
        </>
      )}

      {showCreateModal && (
        <CreateCookbookModal
          onClose={() => setShowCreateModal(false)}
          onCreated={() => {
            setShowCreateModal(false);
            refreshCookbooks();
          }}
        />
      )}
    </div>
  );
};

export default Profile;
