interface ChefLoadingProps {
	message?: string;
}

const ChefLoading: React.FC<ChefLoadingProps> = ({
	message = "Chef is cooking...",
}) => {
	return (
		<div className="flex flex-col items-center justify-center py-12 gap-4">
			{/* Animated dots */}
			<div className="flex gap-2">
				{[0, 1, 2].map((i) => (
					<div
						key={i}
						className="w-3 h-3 rounded-full bg-violet-500 animate-bounce"
						style={{ animationDelay: `${i * 0.15}s` }}
					/>
				))}
			</div>
			<p className="text-gray-500 text-sm font-medium animate-pulse">
				{message}
			</p>
			<div className="w-48 h-1.5 bg-gray-100 rounded-full overflow-hidden">
				<div className="h-full bg-gradient-to-r from-violet-400 to-purple-500 rounded-full animate-progress" />
			</div>
		</div>
	);
};

export default ChefLoading;
