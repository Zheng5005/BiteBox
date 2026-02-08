import axiosInstance from './axiosInstance';
import type { MealType } from '../types';

export function getMealTypes() {
  return axiosInstance.get<MealType[]>('/mealtypes');
}
