import axiosInstance from './axiosInstance';
import type { Cookbook, Recipe } from '../types';

interface CookbookRecipeRaw {
  id: string;
  name_recipe: string;
  description: string;
  meal_type_id: string;
  img_url: string;
  rating: string;
  likes: number;
}

export async function getCookbooks(): Promise<Cookbook[]> {
  const res = await axiosInstance.get<Cookbook[]>('/cookbooks');
  return res.data ?? [];
}

export async function createCookbook(name: string, description: string, is_public: boolean): Promise<void> {
  await axiosInstance.post('/cookbooks', { name, description, is_public });
}

export async function addRecipeToCookbook(cookbookId: number, recipeId: number, notes = ''): Promise<void> {
  await axiosInstance.post(`/cookbooks/recipes/add/${cookbookId}`, { recipe_id: Number(recipeId), notes });
}

export async function updateCookbook(
  cookbookId: number,
  fields: { name?: string; description?: string; is_public?: boolean },
): Promise<void> {
  await axiosInstance.patch(`/cookbooks/edit/${cookbookId}`, fields);
}

export async function deleteCookbook(cookbookId: number): Promise<void> {
  await axiosInstance.delete(`/cookbooks/delete/${cookbookId}`);
}

export async function removeRecipeFromCookbook(cookbookId: number, recipeId: number): Promise<void> {
  await axiosInstance.delete(`/cookbooks/recipes/remove/${cookbookId}`, {
    data: { recipe_id: recipeId },
  });
}

export async function getCookbookRecipes(cookbookId: string): Promise<{ cookbook: Cookbook; recipes: Recipe[] }> {
  const [cookbooksRes, recipesRes] = await Promise.all([
    axiosInstance.get<Cookbook[]>('/cookbooks'),
    axiosInstance.get<CookbookRecipeRaw[]>(`/cookbooks/recipes/${cookbookId}`),
  ]);

  const cookbook = (cookbooksRes.data ?? []).find((c) => c.id === Number(cookbookId));
  if (!cookbook) throw new Error('Cookbook not found');

  const recipes = (recipesRes.data ?? []).map((r) => ({
    id: Number(r.id),
    name_recipe: r.name_recipe,
    description: r.description,
    meal_type_id: Number(r.meal_type_id),
    image: r.img_url,
    rating: Number(r.rating),
    likes: r.likes,
  }));

  return { cookbook, recipes };
}
