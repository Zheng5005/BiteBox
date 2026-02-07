import { Link } from 'react-router';
import type { Recipe } from '../types';

interface RecipeCardProps {
  recipe: Recipe;
}

const RecipeCard: React.FC<RecipeCardProps> = ({ recipe }) => {
  return (
    <Link
      to={`/details/${recipe.id}`}
      className="bg-white shadow-md rounded-2xl overflow-hidden grid grid-cols-1 md:grid-cols-3"
    >
      <img
        src={recipe.image}
        alt={recipe.name_recipe}
        className="object-cover w-full h-full md:col-span-1"
      />
      <div className="col-span-2 p-4">
        <h2 className="text-xl font-bold mb-1">{recipe.name_recipe}</h2>
        <div className="flex items-center text-yellow-500 mb-2">
          <span className="mr-1">⭐</span>
          <span>{recipe.rating > 0.0 ? recipe.rating : "BE THE FIRST ONE TO RATE IT!"}</span>
        </div>
        <p className="text-gray-600">{recipe.description}</p>
      </div>
    </Link>
  );
};

export default RecipeCard;
