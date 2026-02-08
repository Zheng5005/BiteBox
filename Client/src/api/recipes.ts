import axiosInstance from './axiosInstance';
import type { Recipe, RecipeDetail } from '../types';

export function getRecipes() {
  return axiosInstance.get<Recipe[]>('/recipes');
}

export function getRecipeById(id: string) {
  return axiosInstance.get<RecipeDetail>(`/recipes/${id}`);
}

export function postRecipe(formData: FormData) {
  return axiosInstance.postForm('/recipes/post', formData);
}
