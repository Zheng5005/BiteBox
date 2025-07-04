export interface MealType {
  id: number;
  name: string;
}

export interface MealTypeFilterProps {
  mealTypes: MealType[];
  selected: string;
  onChange: (value: string) => void;
}

const MealTypeFilter: React.FC<MealTypeFilterProps> = ({ mealTypes, selected, onChange}) => {
  return(
    <select
      value={selected}
      onChange={(e) => onChange(e.target.value)}
      className="w-full md:w-1/4 px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring focus:border-red-300"
    >
      <option value="">All meals</option>
      {mealTypes.map((mt) => (
        <option key={mt.id} value={mt.id}>
          {mt.name}
        </option>
      ))}
    </select>
  )
}

export default MealTypeFilter;
