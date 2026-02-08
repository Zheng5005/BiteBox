import axiosInstance from './axiosInstance';
import type { Recipe } from '../types';

interface UserRecipeRaw {
  id: string;
  name_recipe: string;
  description: string;
  meal_type_id: string;
  img_url: string;
  rating: string;
}

export async function getUserRecipes(): Promise<Recipe[]> {
  const res = await axiosInstance.get<UserRecipeRaw[]>('/users');
  return res.data.map((r) => ({
    id: Number(r.id),
    name_recipe: r.name_recipe,
    description: r.description,
    meal_type_id: Number(r.meal_type_id),
    image: r.img_url,
    rating: Number(r.rating),
  }));
}
