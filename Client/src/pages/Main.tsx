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
      const matchCategory = mealTypeSelected === '' || String(recipe.meal_type_id) === mealTypeSelected;
      return matchSearch && matchCategory;
    });
  }, [search, mealTypeSelected, recipesArray]);

  const isSearching = search !== '' || mealTypeSelected !== '';

  const selectedMealTypeName = mealTypes.find((mt) => String(mt.id) === mealTypeSelected)?.name;

  const clearFilters = () => {
    setSearch('');
    setMealTypeSelected('');
  };

  return (
    <div className="max-w-7xl mx-auto p-6 font-sans">
      <div className="bg-gradient-to-br from-indigo-600 via-blue-500 to-purple-500 rounded-2xl p-8 mb-8 text-white shadow-lg relative overflow-hidden">
        <div className="absolute -top-6 -right-6 text-[120px] opacity-10 select-none pointer-events-none rotate-12">
          🍽️
        </div>

        <div className="relative z-10">
          <h1 className="text-4xl font-extrabold mb-2 text-center md:text-left">
            Discover Your Next Favorite Meal
          </h1>
          <p className="text-indigo-100 text-lg mb-6 text-center md:text-left">
            Explore thousands of recipes shared by the BiteBox community.
          </p>

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
                className="w-full pl-10 pr-4 py-3 bg-white/95 backdrop-blur border-0 rounded-xl focus:ring-4 focus:ring-blue-300 text-gray-900 transition-all shadow-inner"
              />
            </div>
            <div className="flex w-full md:w-auto gap-3 items-center">
              <MealTypeFilter mealTypes={mealTypes} selected={mealTypeSelected} onChange={setMealTypeSelected}/>
              <IFLButton recipes={recipesArray} />
            </div>
          </div>
        </div>
      </div>

      {isSearching && (
        <div className="flex flex-wrap items-center gap-2 mb-6">
          <span className="text-sm text-gray-500">Active filters:</span>
          {search && (
            <span className="inline-flex items-center gap-1 px-3 py-1 bg-indigo-100 text-indigo-700 rounded-full text-sm font-medium">
              "{search}"
              <button onClick={() => setSearch('')} className="ml-1 hover:text-indigo-900 transition" aria-label="Clear search">
                ✕
              </button>
            </span>
          )}
          {selectedMealTypeName && (
            <span className="inline-flex items-center gap-1 px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm font-medium">
              {selectedMealTypeName}
              <button onClick={() => setMealTypeSelected('')} className="ml-1 hover:text-purple-900 transition" aria-label="Clear meal type">
                ✕
              </button>
            </span>
          )}
          <button
            onClick={clearFilters}
            className="text-sm text-gray-400 hover:text-gray-600 underline underline-offset-2 transition ml-2"
          >
            Clear all
          </button>
        </div>
      )}

      {isSearching ? (
        filteredRecipes.length > 0 ? (
          <RecipeList title={`Search Results (${filteredRecipes.length})`} recipes={filteredRecipes} loading={loading} />
        ) : (
          !loading && (
            <div className="text-center py-16">
              <p className="text-5xl mb-4">🔍</p>
              <h3 className="text-xl font-semibold text-gray-700 mb-2">No recipes found</h3>
              <p className="text-gray-500 mb-4">
                Try adjusting your search or filters to find what you're looking for.
              </p>
              <button
                onClick={clearFilters}
                className="px-5 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition font-medium"
              >
                Clear filters
              </button>
            </div>
          )
        )
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
