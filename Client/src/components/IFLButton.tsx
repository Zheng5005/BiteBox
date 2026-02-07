import { useNavigate } from "react-router"
import type { Recipe } from '../types'

interface IFLButtonProps {
  recipes: Recipe[]
}

const IFLButton: React.FC<IFLButtonProps> = ({recipes}) => {
  const navigate = useNavigate()

  const handleClick = () => {
    if (recipes.length === 0) return

    const randomIndex = Math.floor(Math.random() * recipes.length)
    const randomId = recipes[randomIndex].id;
    navigate(`/details/${randomId}`)
  }
  return(
    <button
      onClick={handleClick}
      className="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-100 transition"
    >
      I feel lucky
    </button>
  )

}

export default IFLButton
