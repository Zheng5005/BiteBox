import React, { useEffect, useMemo, useState } from 'react';
import IFLButton from '../components/IFLButton';
import MealTypeFilter from '../components/MealTypeFilter';
import RecipeCard from '../components/RecipeCard';
import type { Recipe } from '../types';
import { useMealTypes } from '../hooks/useMealTypes';
import { getRecipes } from '../api/recipes';

const MainPage: React.FC = () => {
  const [search, setSearch] = useState<string>('');
  const [mealTypeSelected, setMealTypeSelected] = useState<string>('')
  const [recipesArray, setRecipesArray] = useState<Recipe[]>([])
  const mealTypes = useMealTypes()

  async function fetchRecipes(): Promise<void>{
    try {
      const res = await getRecipes();
      setRecipesArray(res.data);
    } catch (error) {
      console.error("Failed to fetch recipes:", error);
    }
  }

  useEffect(() => {
    fetchRecipes()
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
      <div className="flex flex-col md:flex-row justify-between items-center gap-4 mb-6">
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search recipes..."
          className="w-full md:w-1/2 p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
        />
        <MealTypeFilter mealTypes={mealTypes} selected={mealTypeSelected} onChange={setMealTypeSelected}/>
        <IFLButton recipes={recipesArray} />
      </div>

      <div className="grid gap-6">
        {filteredRecipes.map((recipe: Recipe) => (
          <RecipeCard key={recipe.id} recipe={recipe} />
        ))}
      </div>
    </div>
  );
};

export default MainPage;
