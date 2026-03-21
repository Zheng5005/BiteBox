import React from 'react';
import RecipeCard from './RecipeCard';
import type { Recipe } from '../types';

interface RecipeListProps {
  recipes: Recipe[];
  title?: string;
  loading?: boolean;
}

const RecipeList: React.FC<RecipeListProps> = ({ recipes, title, loading }) => {
  if (loading) {
    return (
      <div className="py-8">
        {title && <h2 className="text-2xl font-bold mb-6">{title}</h2>}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="animate-pulse bg-white rounded-xl shadow-md overflow-hidden h-80">
              <div className="bg-gray-200 h-48 w-full" />
              <div className="p-4 space-y-3">
                <div className="h-6 bg-gray-200 rounded w-3/4" />
                <div className="h-4 bg-gray-200 rounded w-1/2" />
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (recipes.length === 0) return null;

  return (
    <div className="py-8">
      {title && <h2 className="text-2xl font-bold mb-6 text-gray-800 border-b-2 border-indigo-500 pb-2 inline-block">{title}</h2>}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {recipes.map((recipe) => (
          <RecipeCard key={recipe.id} recipe={recipe} />
        ))}
      </div>
    </div>
  );
};

export default RecipeList;
