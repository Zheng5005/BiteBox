import { useState, useRef } from "react";
import { Link } from "react-router";
import type { AiRecipe, AiRecipeSummary, AiConstraints } from "../types";
import ChefInput from "../components/ChefInput";
import ChefLoading from "../components/ChefLoading";
import ConstraintFilters, {
	type ConstraintValues,
} from "../components/ConstraintFilters";
import AIBadge from "../components/AIBadge";
import { streamAiChefRequest, submitAiFeedback } from "../api/aiChef";

type SearchStatus = "idle" | "searching" | "complete" | "error";

const AIChef: React.FC = () => {
	const [ingredients, setIngredients] = useState<string[]>([]);
	const [constraints, setConstraints] = useState<ConstraintValues>({});
	const [status, setStatus] = useState<SearchStatus>("idle");
	const [errorMsg, setErrorMsg] = useState("");
	const [catalogResults, setCatalogResults] = useState<AiRecipeSummary[]>([]);
	const [aiRecipe, setAiRecipe] = useState<AiRecipe | null>(null);
	const [feedbackDone, setFeedbackDone] = useState(false);
	const abortRef = useRef(false);

	const handleCook = async () => {
		if (ingredients.length === 0) return;

		setStatus("searching");
		setErrorMsg("");
		setCatalogResults([]);
		setAiRecipe(null);
		setFeedbackDone(false);
		abortRef.current = false;

		const aiConstraints: AiConstraints = {};
		if (constraints.max_time) aiConstraints.max_time = constraints.max_time;
		if (constraints.difficulty)
			aiConstraints.difficulty = constraints.difficulty;
		if (constraints.cuisine) aiConstraints.cuisine = constraints.cuisine;

		try {
			for await (const event of streamAiChefRequest(
				ingredients,
				aiConstraints,
			)) {
				if (abortRef.current) break;

				switch (event.source) {
					case "catalog":
						setCatalogResults(event.recipes);
						setStatus("complete");
						break;

					case "ai":
						if (event.status === "complete") {
							setAiRecipe(event.recipe);
							setStatus("complete");
						}
						// partial/starting events are handled by the loading state
						break;

					case "error":
						setErrorMsg(event.error);
						setStatus("error");
						break;
				}
			}
		} catch (err) {
			if (!abortRef.current) {
				setErrorMsg(
					err instanceof Error ? err.message : "Something went wrong",
				);
				setStatus("error");
			}
		}
	};

	const handleCancel = () => {
		abortRef.current = true;
		setStatus("idle");
	};

	const handleFeedback = async (action: "save" | "discard") => {
		if (!aiRecipe) return;
		try {
			await submitAiFeedback(aiRecipe.id, action);
			setFeedbackDone(true);
		} catch {
			// Silently fail — not critical
		}
	};

	const handleReset = () => {
		setIngredients([]);
		setConstraints({});
		setStatus("idle");
		setErrorMsg("");
		setCatalogResults([]);
		setAiRecipe(null);
		setFeedbackDone(false);
	};

	return (
		<div className="max-w-3xl mx-auto px-4 py-8">
			{/* Header */}
			<div className="text-center mb-8">
				<h1 className="text-3xl font-bold text-gray-800 mb-2">
					<span className="text-violet-500">&#10024;</span> AI Chef
				</h1>
				<p className="text-gray-500">
					Tell us what you have, we&apos;ll find you a recipe.
				</p>
			</div>

			{/* Input Section */}
			<div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mb-6">
				<label className="block text-sm font-semibold text-gray-700 mb-2">
					Your Ingredients
				</label>
				<ChefInput
					value={ingredients}
					onChange={setIngredients}
					disabled={status === "searching"}
					placeholder="chicken, rice, garlic, onion..."
				/>

				{/* Constraints */}
				<div className="mt-4 pt-4 border-t border-gray-100">
					<ConstraintFilters value={constraints} onChange={setConstraints} />
				</div>

				{/* Actions */}
				<div className="flex gap-3 mt-6">
					<button
						onClick={handleCook}
						disabled={ingredients.length === 0 || status === "searching"}
						className="flex-1 px-6 py-3 bg-gradient-to-r from-violet-500 to-purple-600 text-white font-semibold rounded-xl shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed transition"
					>
						{status === "searching" ? "Searching..." : "Let's Cook!"}
					</button>
					{status === "searching" && (
						<button
							onClick={handleCancel}
							className="px-4 py-3 border border-gray-200 text-gray-600 rounded-xl hover:bg-gray-50 transition"
						>
							Cancel
						</button>
					)}
					{status === "complete" && (
						<button
							onClick={handleReset}
							className="px-4 py-3 border border-gray-200 text-gray-600 rounded-xl hover:bg-gray-50 transition"
						>
							New Search
						</button>
					)}
				</div>
			</div>

			{/* Error */}
			{status === "error" && (
				<div className="bg-red-50 border border-red-200 rounded-xl p-4 mb-6 text-center">
					<p className="text-red-600 font-medium">{errorMsg}</p>
					{errorMsg.includes("limit") && (
						<p className="text-red-500 text-sm mt-1">
							Try again tomorrow or browse existing recipes.
						</p>
					)}
					<button
						onClick={handleReset}
						className="mt-3 px-4 py-2 bg-red-100 text-red-700 rounded-lg text-sm hover:bg-red-200 transition"
					>
						Try Again
					</button>
				</div>
			)}

			{/* Loading */}
			{status === "searching" && (
				<div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mb-6">
					<ChefLoading
						message={
							catalogResults.length > 0
								? "Checking the catalog..."
								: "Consulting the AI Chef..."
						}
					/>
				</div>
			)}

			{/* Catalog Results */}
			{status === "complete" && catalogResults.length > 0 && !aiRecipe && (
				<div>
					<h2 className="text-lg font-semibold text-gray-700 mb-3">
						Found {catalogResults.length} matching recipe
						{catalogResults.length > 1 ? "s" : ""} in our catalog
					</h2>
					<div className="space-y-4">
						{catalogResults.map((recipe) => (
							<Link
								key={recipe.id}
								to={`/details/${recipe.id}`}
								className="block bg-white rounded-xl shadow-sm border border-gray-100 p-4 hover:shadow-md transition"
							>
								<div className="flex items-start gap-4">
									{recipe.img_url && (
										<img
											src={recipe.img_url}
											alt={recipe.name_recipe}
											className="w-20 h-20 rounded-lg object-cover flex-shrink-0"
										/>
									)}
									<div className="flex-1 min-w-0">
										<h3 className="font-semibold text-gray-800 truncate">
											{recipe.name_recipe}
										</h3>
										<p className="text-sm text-gray-500 mt-1 line-clamp-2">
											{recipe.description}
										</p>
										<div className="flex items-center gap-3 mt-2 text-sm">
											<span className="text-yellow-500">
												{recipe.rating !== "0"
													? `⭐ ${recipe.rating}`
													: "No ratings yet"}
											</span>
											<span className="text-red-400">
												❤ {recipe.likes ?? 0}
											</span>
											<span className="text-violet-500 font-medium">
												{Math.round(recipe.match_score * 100)}% match
											</span>
										</div>
									</div>
								</div>
							</Link>
						))}
					</div>
				</div>
			)}

			{/* AI Generated Recipe */}
			{status === "complete" && aiRecipe && (
				<div>
					<div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden mb-6">
						{/* Header */}
						<div className="p-6 border-b border-gray-100">
							<div className="flex items-center gap-2 mb-2">
								<AIBadge />
								<h2 className="text-xl font-bold text-gray-800">
									{aiRecipe.name_recipe}
								</h2>
							</div>
							<p className="text-gray-500 text-sm">{aiRecipe.description}</p>
							<div className="flex gap-3 mt-3 text-sm text-gray-500">
								{aiRecipe.estimated_minutes > 0 && (
									<span>⏱ {aiRecipe.estimated_minutes} min</span>
								)}
								{aiRecipe.difficulty && (
									<span className="capitalize">📊 {aiRecipe.difficulty}</span>
								)}
								{aiRecipe.cuisine && (
									<span className="capitalize">🌍 {aiRecipe.cuisine}</span>
								)}
							</div>
						</div>

						{/* Ingredients */}
						{aiRecipe.ingredients.length > 0 && (
							<div className="p-6 border-b border-gray-100">
								<h3 className="font-semibold text-gray-700 mb-3">
									Ingredients
								</h3>
								<ul className="grid grid-cols-1 sm:grid-cols-2 gap-2">
									{aiRecipe.ingredients.map((ing, i) => (
										<li
											key={i}
											className="flex items-center gap-2 text-sm text-gray-600"
										>
											<span className="w-1.5 h-1.5 rounded-full bg-violet-400 flex-shrink-0" />
											<span className="font-medium text-gray-800">
												{ing.quantity > 0
													? `${ing.quantity} ${ing.unit}`
													: ing.unit}
											</span>
											{ing.name}
										</li>
									))}
								</ul>
							</div>
						)}

						{/* Steps */}
						{aiRecipe.steps.length > 0 && (
							<div className="p-6">
								<h3 className="font-semibold text-gray-700 mb-3">Steps</h3>
								<ol className="space-y-4">
									{aiRecipe.steps
										.sort((a, b) => a.order - b.order)
										.map((step) => (
											<li key={step.order} className="flex gap-3">
												<span className="flex-shrink-0 w-7 h-7 rounded-full bg-violet-100 text-violet-600 text-sm font-semibold flex items-center justify-center">
													{step.order}
												</span>
												<p className="text-sm text-gray-600 leading-relaxed pt-1">
													{step.instruction}
												</p>
											</li>
										))}
								</ol>
							</div>
						)}
					</div>

					{/* Feedback */}
					{!feedbackDone && (
						<div className="flex gap-3 justify-center">
							<button
								onClick={() => handleFeedback("save")}
								className="px-6 py-2.5 bg-green-500 text-white rounded-xl font-medium hover:bg-green-600 transition"
							>
								❤ Save Recipe
							</button>
							<button
								onClick={() => handleFeedback("discard")}
								className="px-6 py-2.5 border border-gray-200 text-gray-600 rounded-xl font-medium hover:bg-gray-50 transition"
							>
								Not what I needed
							</button>
						</div>
					)}
					{feedbackDone && (
						<p className="text-center text-green-600 text-sm font-medium">
							Thanks for the feedback! &#10003;
						</p>
					)}
				</div>
			)}

			{/* Empty state */}
			{status === "idle" && ingredients.length === 0 && (
				<div className="text-center py-8">
					<p className="text-gray-400 text-sm">
						Add your ingredients above and hit &quot;Let&apos;s Cook!&quot; to
						get started.
					</p>
				</div>
			)}
		</div>
	);
};

export default AIChef;
