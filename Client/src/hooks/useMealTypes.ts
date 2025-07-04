import { useEffect, useState } from "react";
import type { MealType } from "../components/MealTypeFilter";

export function useMealTypes() {
  const [mealTypes, setMealTypes] = useState<MealType[]>([])

  useEffect(() => {
    fetch("http://localhost:8080/api/mealtypes")
      .then((res) => res.json())
      .then((data) => setMealTypes(data))
      .catch(() => setMealTypes([]))
  }, [])

  return mealTypes
}
