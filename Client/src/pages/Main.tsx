import React, { useEffect, useMemo, useState } from 'react';
import IFLButton from '../components/IFLButton';
import MealTypeFilter from '../components/MealTypeFilter';
import RecipeList from '../components/RecipeList';
import type { Recipe } from '../types';
import { useMealTypes } from '../hooks/useMealTypes';
import { getRecipes, getPopularRecipes, getRecommendedRecipes } from '../api/recipes';

const MainPage: React.FC = () => {
  const [search, setSearch] = useState<string>('');
  const [mealTypeSelected, setMealTypeSelected] = useState<string>('')
  const [recipesArray, setRecipesArray] = useState<Recipe[]>([])
  const [popularRecipes, setPopularRecipes] = useState<Recipe[]>([])
  const [recommendedRecipes, setRecommendedRecipes] = useState<Recipe[]>([])
  const [loading, setLoading] = useState<boolean>(true)
  const mealTypes = useMealTypes()

  async function fetchData(): Promise<void>{
    setLoading(true);
    try {
      const [recipesRes, popularRes, recommendedRes] = await Promise.all([
        getRecipes(),
        getPopularRecipes(),
        getRecommendedRecipes()
      ]);
      setRecipesArray(recipesRes.data);
      setPopularRecipes(popularRes.data);
      setRecommendedRecipes(recommendedRes.data);
    } catch (error) {
      console.error("Failed to fetch data:", error);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const filteredRecipes = useMemo(() => {
    return recipesArray.filter((recipe) => {
      const matchSearch = recipe.name_recipe.toLowerCase().includes(search.toLowerCase());
      const matchCategory = mealTypeSelected === '' || recipe.meal_type_id === Number(mealTypeSelected);
      return matchSearch && matchCategory;
    });
  }, [search, mealTypeSelected, recipesArray]);

  const isSearching = search !== '' || mealTypeSelected !== '';

  return (
    <div className="max-w-7xl mx-auto p-6 font-sans">
      <div className="bg-gradient-to-r from-indigo-600 to-blue-500 rounded-2xl p-8 mb-8 text-white shadow-lg">
        <h1 className="text-4xl font-extrabold mb-2 text-center md:text-left">Discover Your Next Favorite Meal</h1>
        <p className="text-indigo-100 text-lg mb-6 text-center md:text-left">Explore thousands of recipes shared by the BiteBox community.</p>
        
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <div className="relative w-full md:w-1/2">
            <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">
              <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </span>
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search recipes by name..."
              className="w-full pl-10 pr-4 py-3 bg-white border-0 rounded-xl focus:ring-4 focus:ring-blue-300 text-gray-900 transition-all shadow-inner"
            />
          </div>
          <div className="flex w-full md:w-auto gap-4">
            <MealTypeFilter mealTypes={mealTypes} selected={mealTypeSelected} onChange={setMealTypeSelected}/>
            <IFLButton recipes={recipesArray} />
          </div>
        </div>
      </div>

      {isSearching ? (
        <RecipeList title={`Search Results (${filteredRecipes.length})`} recipes={filteredRecipes} loading={loading} />
      ) : (
        <>
          <RecipeList title="Recommended for You" recipes={recommendedRecipes} loading={loading} />
          <RecipeList title="Most Popular" recipes={popularRecipes} loading={loading} />
          <RecipeList title="All Recipes" recipes={recipesArray} loading={loading} />
        </>
      )}
    </div>
  );
};

export default MainPage;
