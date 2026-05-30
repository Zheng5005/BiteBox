const DIFFICULTY_OPTIONS = [
	{ value: "easy", label: "Easy" },
	{ value: "medium", label: "Medium" },
	{ value: "hard", label: "Hard" },
];

const CUISINE_OPTIONS = [
	{ value: "italian", label: "🇮🇹 Italian" },
	{ value: "mexican", label: "🇲🇽 Mexican" },
	{ value: "japanese", label: "🇯🇵 Japanese" },
	{ value: "indian", label: "🇮🇳 Indian" },
	{ value: "chinese", label: "🇨🇳 Chinese" },
	{ value: "american", label: "🇺🇸 American" },
];

const TIME_OPTIONS = [
	{ value: 15, label: "15 min" },
	{ value: 30, label: "30 min" },
	{ value: 45, label: "45 min" },
	{ value: 60, label: "60 min" },
];

export interface ConstraintValues {
	max_time?: number;
	difficulty?: string;
	cuisine?: string;
}

interface ConstraintFiltersProps {
	value: ConstraintValues;
	onChange: (value: ConstraintValues) => void;
}

const ConstraintFilters: React.FC<ConstraintFiltersProps> = ({
	value,
	onChange,
}) => {
	const toggleDifficulty = (diff: string) => {
		onChange({
			...value,
			difficulty: value.difficulty === diff ? undefined : diff,
		});
	};

	const toggleCuisine = (cuisine: string) => {
		onChange({
			...value,
			cuisine: value.cuisine === cuisine ? undefined : cuisine,
		});
	};

	const setTime = (time?: number) => {
		onChange({
			...value,
			max_time: time,
		});
	};

	const hasFilters = value.max_time || value.difficulty || value.cuisine;

	return (
		<div className="space-y-3">
			{/* Max Time */}
			<div>
				<label className="text-xs font-semibold text-gray-500 uppercase tracking-wide">
					Max Time
				</label>
				<div className="flex gap-2 mt-1 flex-wrap">
					{TIME_OPTIONS.map((opt) => (
						<button
							key={opt.value}
							type="button"
							onClick={() =>
								setTime(value.max_time === opt.value ? undefined : opt.value)
							}
							className={`px-3 py-1.5 text-sm rounded-lg border transition ${
								value.max_time === opt.value
									? "bg-violet-100 border-violet-300 text-violet-700 font-medium"
									: "bg-white border-gray-200 text-gray-600 hover:border-gray-300"
							}`}
						>
							{opt.label}
						</button>
					))}
				</div>
			</div>

			{/* Difficulty */}
			<div>
				<label className="text-xs font-semibold text-gray-500 uppercase tracking-wide">
					Difficulty
				</label>
				<div className="flex gap-2 mt-1 flex-wrap">
					{DIFFICULTY_OPTIONS.map((opt) => (
						<button
							key={opt.value}
							type="button"
							onClick={() => toggleDifficulty(opt.value)}
							className={`px-3 py-1.5 text-sm rounded-lg border transition ${
								value.difficulty === opt.value
									? "bg-violet-100 border-violet-300 text-violet-700 font-medium"
									: "bg-white border-gray-200 text-gray-600 hover:border-gray-300"
							}`}
						>
							{opt.label}
						</button>
					))}
				</div>
			</div>

			{/* Cuisine */}
			<div>
				<label className="text-xs font-semibold text-gray-500 uppercase tracking-wide">
					Cuisine
				</label>
				<div className="flex gap-2 mt-1 flex-wrap">
					{CUISINE_OPTIONS.map((opt) => (
						<button
							key={opt.value}
							type="button"
							onClick={() => toggleCuisine(opt.value)}
							className={`px-3 py-1.5 text-sm rounded-lg border transition ${
								value.cuisine === opt.value
									? "bg-violet-100 border-violet-300 text-violet-700 font-medium"
									: "bg-white border-gray-200 text-gray-600 hover:border-gray-300"
							}`}
						>
							{opt.label}
						</button>
					))}
				</div>
			</div>

			{/* Clear all */}
			{hasFilters && (
				<button
					type="button"
					onClick={() => onChange({})}
					className="text-xs text-violet-500 hover:text-violet-700 transition"
				>
					Clear all filters
				</button>
			)}
		</div>
	);
};

export default ConstraintFilters;
