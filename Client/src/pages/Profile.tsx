import { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import Avatar from '../components/Avatar';
import RecipeCard, { type Recipe } from '../components/RecipeCard';
import axiosInstance from '../api/axiosInstance';

const Profile: React.FC = () => {
  const { user } = useAuth();
  const [recipes, setRecipes] = useState<Recipe[]>([]);

  useEffect(() => {
    async function fetchUserRecipes() {
      try {
        const res = await axiosInstance.get('/users');
        const mapped = res.data.map((r: { id: string; name_recipe: string; description: string; meal_type_id: string; img_url: string; rating: string }) => ({
          id: Number(r.id),
          name_recipe: r.name_recipe,
          description: r.description,
          meal_type_id: Number(r.meal_type_id),
          image: r.img_url,
          rating: Number(r.rating),
        }));
        setRecipes(mapped);
      } catch (error) {
        console.error('Failed to fetch user recipes:', error);
      }
    }

    if (user) {
      fetchUserRecipes();
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

      <h2 className="text-xl font-semibold mb-4">My Recipes</h2>
      {recipes.length === 0 ? (
        <p className="text-gray-500">You haven't posted any recipes yet.</p>
      ) : (
        <div className="grid gap-6">
          {recipes.map((recipe) => (
            <RecipeCard key={recipe.id} recipe={recipe} />
          ))}
        </div>
      )}
    </div>
  );
};

export default Profile;
