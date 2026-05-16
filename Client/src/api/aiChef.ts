import type { AiConstraints, AiEvent, AiRecipe } from "../types";

const API_BASE = "http://localhost:8080/api";

/**
 * Submit an AI Chef request and return an async iterable of SSE events.
 * Uses fetch + ReadableStream since SSE requires POST.
 */
export async function* streamAiChefRequest(
	ingredients: string[],
	constraints?: AiConstraints,
): AsyncIterableIterator<AiEvent> {
	const token = localStorage.getItem("token");

	const response = await fetch(`${API_BASE}/ai/recipes`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			...(token ? { Authorization: `Bearer ${token}` } : {}),
		},
		body: JSON.stringify({ ingredients, constraints }),
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		yield {
			source: "error",
			error: errorBody?.error ?? `HTTP ${response.status}`,
			code: errorBody?.code ?? "NETWORK_ERROR",
		};
		return;
	}

	const reader = response.body?.getReader();
	if (!reader) {
		yield {
			source: "error",
			error: "Streaming not supported",
			code: "INTERNAL_ERROR",
		};
		return;
	}

	const decoder = new TextDecoder();
	let buffer = "";

	try {
		while (true) {
			const { done, value } = await reader.read();
			if (done) break;

			buffer += decoder.decode(value, { stream: true });

			// Parse SSE events from buffer
			const lines = buffer.split("\n");
			buffer = lines.pop() ?? ""; // Keep incomplete line in buffer

			for (const line of lines) {
				if (line.startsWith("data: ")) {
					try {
						const data = JSON.parse(line.slice(6));
						yield data as AiEvent;
					} catch {
						// Skip malformed JSON
					}
				}
			}
		}

		// Process remaining buffer
		if (buffer.startsWith("data: ")) {
			try {
				const data = JSON.parse(buffer.slice(6));
				yield data as AiEvent;
			} catch {
				// ignore
			}
		}
	} finally {
		reader.releaseLock();
	}
}

/**
 * Submit feedback (save/discard) for an AI recipe.
 */
export async function submitAiFeedback(
	recipeId: string,
	action: "save" | "discard",
): Promise<void> {
	const token = localStorage.getItem("token");
	const response = await fetch(`${API_BASE}/ai/recipes/${recipeId}/feedback`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			...(token ? { Authorization: `Bearer ${token}` } : {}),
		},
		body: JSON.stringify({ action }),
	});

	if (!response.ok) {
		throw new Error(`Feedback failed: ${response.status}`);
	}
}

/**
 * Fetch detail for an AI-generated recipe.
 */
export async function getAiRecipeDetail(id: string): Promise<AiRecipe> {
	const response = await fetch(`${API_BASE}/ai/recipes/${id}`, {
		headers: {
			Accept: "application/json",
		},
	});

	if (!response.ok) {
		throw new Error(`Recipe not found: ${response.status}`);
	}

	return response.json();
}
