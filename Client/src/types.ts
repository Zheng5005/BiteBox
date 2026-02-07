export interface User {
  id: number;
  name: string;
  url_photo: string;
}

export interface Recipe {
  id: number;
  name_recipe: string;
  description: string;
  meal_type_id: number;
  image: string;
  rating: number;
}

export interface RecipeDetail extends Recipe {
  img_url: string;
  creator_name: string;
  steps: string[];
}

export interface Comment {
  id: number;
  user_name: string;
  recipe_id: string;
  comment: string;
  rating: number;
}

export interface MealType {
  id: number;
  name: string;
}
