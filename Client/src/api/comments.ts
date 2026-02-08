import axiosInstance from './axiosInstance';
import type { Comment } from '../types';

export function getComments(recipeId: string) {
  return axiosInstance.get<Comment[]>(`/comments/${recipeId}`);
}

export function postComment(recipeId: string, comment: string, rating: number) {
  return axiosInstance.post(`/comments/post/${recipeId}`, { comment, rating });
}
