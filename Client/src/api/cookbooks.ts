import axiosInstance from './axiosInstance';
import type { Cookbook } from '../types';

export async function getCookbooks(): Promise<Cookbook[]> {
  const res = await axiosInstance.get<Cookbook[]>('/cookbooks');
  return res.data ?? [];
}

export async function addRecipeToCookbook(cookbookId: number, recipeId: number, notes = ''): Promise<void> {
  await axiosInstance.post(`/cookbooks/recipes/add/${cookbookId}`, { recipe_id: Number(recipeId), notes });
}
