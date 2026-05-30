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
	likes?: number;
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

export interface CookbookRecipe extends Recipe {
	notes: string;
}

export interface Cookbook {
	id: number;
	user_id: number;
	name: string;
	description: string;
	is_public: boolean;
	created_at: string;
	updated_at: string;
}

// --- AI Chef types ---

export interface AiIngredient {
	name: string;
	quantity: number;
	unit: string;
}

export interface AiStep {
	order: number;
	instruction: string;
}

export interface AiRecipe {
	id: string;
	name_recipe: string;
	description: string;
	img_url: string;
	ingredients: AiIngredient[];
	steps: AiStep[];
	tags: string[];
	estimated_minutes: number;
	difficulty: string;
	cuisine: string;
	ai_generated: boolean;
	created_at: string;
}

export interface AiRecipeSummary {
	id: string;
	name_recipe: string;
	description: string;
	img_url: string;
	rating: string;
	likes: number;
	match_score: number;
}

export interface AiConstraints {
	max_time?: number;
	difficulty?: string;
	cuisine?: string;
}

export type AiEvent =
	| { source: "catalog"; recipes: AiRecipeSummary[]; match_count: number }
	| { source: "ai"; status: "starting" }
	| { source: "ai"; status: "partial"; token: string }
	| { source: "ai"; status: "complete"; recipe: AiRecipe }
	| { source: "error"; error: string; code: string };
