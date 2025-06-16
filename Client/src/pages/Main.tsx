
import React, { useEffect, useMemo, useState } from 'react';

interface Recipe {
  id: number;
  name_recipe: string;
  description: string;
  meal_type_id: number,
  image: string;
  rating: number
}

const MainPage: React.FC = () => {
  const [search, setSearch] = useState<string>('');
  const [recipesArray, setRecipesArray] = useState<Recipe[]>([])

  async function fetchRecipes(): Promise<void>{
    const res = await fetch('http://localhost:8080/api/recipes')
    const data = await res.json()
    setRecipesArray(data)
  }

  useEffect(() => {
    fetchRecipes()
  }, [])

  const filteredRecipes = useMemo(() => {
    return recipesArray.filter((recipe) =>
      recipe.name_recipe.toLowerCase().includes(search.toLowerCase())
    );
  }, [search, recipesArray]);


  return (
    <div className="max-w-6xl mx-auto p-6 font-sans">
      {/* Search and Filter */}
      <div className="flex flex-col md:flex-row justify-between items-center gap-4 mb-6">
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search recipes..."
          className="w-full md:w-1/2 p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
        />
        <button className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition">
          Filter
        </button>
      </div>

      {/* Recipes */}
      <div className="grid gap-6">
        {filteredRecipes.map((recipe: Recipe) => (
          <div
            key={recipe.id}
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
                <span>{recipe.rating}</span>
              </div>
              <p className="text-gray-600">{recipe.description}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default MainPage;
