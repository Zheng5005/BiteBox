import type { MealType } from '../types';

export interface MealTypeFilterProps {
  mealTypes: MealType[];
  selected: string;
  onChange: (value: string) => void;
}

const MealTypeFilter: React.FC<MealTypeFilterProps> = ({ mealTypes, selected, onChange}) => {
  return(
    <div className="relative w-full md:w-auto">
      <select
        value={selected}
        onChange={(e) => onChange(e.target.value)}
        className="w-full appearance-none pl-4 pr-10 py-3 bg-white/95 backdrop-blur border-0 rounded-xl focus:ring-4 focus:ring-blue-300 text-gray-900 transition-all shadow-inner cursor-pointer font-medium"
      >
        <option value="">All meals</option>
        {mealTypes.map((mt) => (
          <option key={mt.id} value={mt.id}>
            {mt.name}
          </option>
        ))}
      </select>
      <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400">
        <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </span>
    </div>
  )
}

export default MealTypeFilter;
