import axiosInstance from './axiosInstance';
import type { Cookbook } from '../types';

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
