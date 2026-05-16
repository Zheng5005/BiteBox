import { useState, useCallback } from "react";

interface ChefInputProps {
	value: string[];
	onChange: (tags: string[]) => void;
	placeholder?: string;
	disabled?: boolean;
}

const ChefInput: React.FC<ChefInputProps> = ({
	value,
	onChange,
	placeholder = "Type ingredients, press Enter or comma...",
	disabled = false,
}) => {
	const [input, setInput] = useState("");

	const addTag = useCallback(() => {
		const trimmed = input.trim();
		if (trimmed && !value.includes(trimmed.toLowerCase()) && !disabled) {
			onChange([...value, trimmed.toLowerCase()]);
		}
		setInput("");
	}, [input, value, onChange, disabled]);

	const removeTag = (index: number) => {
		if (!disabled) {
			onChange(value.filter((_, i) => i !== index));
		}
	};

	const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
		if (e.key === "Enter" || e.key === ",") {
			e.preventDefault();
			addTag();
		} else if (
			e.key === "Backspace" &&
			input === "" &&
			value.length > 0 &&
			!disabled
		) {
			removeTag(value.length - 1);
		}
	};

	return (
		<div
			className={`flex flex-wrap items-center gap-2 p-3 min-h-[48px] border-2 rounded-xl transition-colors ${
				disabled
					? "bg-gray-50 border-gray-200"
					: "bg-white border-gray-200 focus-within:border-violet-400"
			}`}
		>
			{value.map((tag, i) => (
				<span
					key={i}
					className="inline-flex items-center gap-1 px-3 py-1 bg-violet-100 text-violet-700 text-sm font-medium rounded-lg"
				>
					{tag}
					{!disabled && (
						<button
							type="button"
							onClick={() => removeTag(i)}
							className="ml-1 text-violet-400 hover:text-violet-600 transition"
							aria-label={`Remove ${tag}`}
						>
							&times;
						</button>
					)}
				</span>
			))}
			<input
				type="text"
				value={input}
				onChange={(e) => setInput(e.target.value)}
				onKeyDown={handleKeyDown}
				onBlur={addTag}
				placeholder={value.length === 0 ? placeholder : ""}
				disabled={disabled}
				className="flex-1 min-w-[120px] bg-transparent outline-none text-sm text-gray-700 placeholder-gray-400"
				aria-label="Add ingredient"
			/>
		</div>
	);
};

export default ChefInput;
