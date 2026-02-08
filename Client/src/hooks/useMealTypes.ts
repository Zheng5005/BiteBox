import { useEffect, useState } from "react";
import type { MealType } from "../types";
import { getMealTypes } from "../api/meals";

export function useMealTypes() {
  const [mealTypes, setMealTypes] = useState<MealType[]>([])

  useEffect(() => {
    getMealTypes()
      .then((res) => setMealTypes(res.data))
      .catch(() => setMealTypes([]))
  }, [])

  return mealTypes
}
