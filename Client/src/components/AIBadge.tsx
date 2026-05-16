interface AIBadgeProps {
	className?: string;
}

const AIBadge: React.FC<AIBadgeProps> = ({ className = "" }) => {
	return (
		<span
			className={`inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full bg-gradient-to-r from-violet-500 to-purple-500 text-white shadow-sm ${className}`}
			title="AI-generated recipe"
		>
			<span>&#10024;</span> AI
		</span>
	);
};

export default AIBadge;
