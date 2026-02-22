import { Link } from 'react-router';
import type { Recipe } from '../types';

interface RecipeCardProps {
  recipe: Recipe;
  onRemove?: () => void;
}

const RecipeCard: React.FC<RecipeCardProps> = ({ recipe, onRemove }) => {
  return (
    <div className="relative">
      {onRemove && (
        <button
          onClick={onRemove}
          className="absolute top-2 right-2 z-10 text-gray-400 hover:text-red-500 text-xl leading-none transition bg-white rounded-full w-7 h-7 flex items-center justify-center shadow"
          title="Remove from cookbook"
        >
          &times;
        </button>
      )}
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
          <div className="flex items-center text-red-500 mb-2">
            <span className="mr-1">❤</span>
            <span>{recipe.likes ? recipe.likes : 0}</span>
          </div>
          <p className="text-gray-600">{recipe.description}</p>
        </div>
      </Link>
    </div>
  );
};

export default RecipeCard;
