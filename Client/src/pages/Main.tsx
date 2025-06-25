import React, { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router';

interface Recipe {
  id: number;
  name_recipe: string;
  description: string;
  meal_type_id: number,
  image: string;
  rating: number
}

interface Meal_type {
  id: number,
  name: string
}

const MainPage: React.FC = () => {
  const [search, setSearch] = useState<string>('');
  const [mealTypeSelected, setMealTypeSelected] = useState<string>('')
  const [recipesArray, setRecipesArray] = useState<Recipe[]>([])
  const [mealTypes, setMealTypes] = useState<Meal_type[]>([])

  async function fetchRecipes(): Promise<void>{
    const res = await fetch('http://localhost:8080/api/recipes')
    const data = await res.json()
    setRecipesArray(data)
  }

  async function fetchMealTypes(): Promise<void>{
    const res = await fetch('http://localhost:8080/api/mealtypes')
    const data = await res.json()
    setMealTypes(data)
  }  

  useEffect(() => {
    fetchRecipes()
    fetchMealTypes()
  }, [])

  const filteredRecipes = useMemo(() => {
    return recipesArray.filter((recipe) => {
      const matchSearch = recipe.name_recipe.toLowerCase().includes(search.toLowerCase());
      const matchCategory = mealTypeSelected === '' || recipe.meal_type_id === Number(mealTypeSelected);
      return matchSearch && matchCategory;

    });
  }, [search, mealTypeSelected, recipesArray]);

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
          <select
            value={mealTypeSelected}
            onChange={(e) => setMealTypeSelected(e.target.value)}
            className="w-full md:w-1/4 px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring focus:border-red-300"
          >
            <option value="">All meals</option>
            {mealTypes.map((mt) => (
              <option key={mt.id} value={mt.id}>
                {mt.name}
              </option>
            ))}
          </select>
      </div>

      {/* Recipes */}
      <div className="grid gap-6">
        {filteredRecipes.map((recipe: Recipe) => (
          //this should be a Link
          <Link
            to={`/details/${recipe.id}`}
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
                <span>{recipe.rating > 0.0 ? recipe.rating : "BE THE FIRST ONE TO RATE IT!"}</span>
              </div>
              <p className="text-gray-600">{recipe.description}</p>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default MainPage;
