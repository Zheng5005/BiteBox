import { useState } from "react"
import { useAuth } from "../context/AuthContext"
import { useMealTypes } from "../hooks/useMealTypes"
import MealTypeFilter from "../components/MealTypeFilter"
import axiosInstance from "../api/axiosInstance"
import axios from "axios" // Import axios for error type checking

const PostRecipe: React.FC = () => {
  const mealTypes = useMealTypes()
  const [mealTypeSelected, setMealTypeSelected] = useState<string>('')
  const [form, setForm] = useState({
    name_recipe: "",
    description: "",
    steps: "",
    guest_name: ""
  })
  const [info, setInfo] = useState({
    isSubmiting: false,
    error: "",
    success: false // Added success state
  })
  const [previewImage, setPreviewImage] = useState<string | ArrayBuffer | null>(null)
  const [imageFile, setImageFile] = useState<File | null>(null) // Changed type to File | null
  const { user } = useAuth();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files ? e.target.files[0] : null;
    if (!file) return;

    if (!file.type.match('image.*')) {
      setInfo(prev => ({...prev, error: "Only PNG or JPEG images"})) // Use functional update
      return;
    }

    if (file.size > 2 * 1024 * 1024) {
      setInfo(prev => ({...prev, error: "Image should be less than 2MB"})) // Use functional update
      return;
    }

    setImageFile(file);
  
    const reader = new FileReader();
    reader.onloadend = () => {
      setPreviewImage(reader.result);
    };
    reader.readAsDataURL(file);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (info.isSubmiting) return

    setInfo(prev => ({...prev, isSubmiting: true, error: "", success: false})) // Reset error and success

    const formData = new FormData();
    formData.append("name", form.name_recipe);
    formData.append("description", form.description);
    formData.append("steps", form.steps);
    formData.append("meal_type_id", mealTypeSelected);
    if (!user) {
      formData.append("guest_name", form.guest_name);
    }
    if (imageFile) {
      formData.append("image", imageFile);
    }

    try {
      const response = await axiosInstance.postForm("/recipes/post", formData);

      if (response.status === 201) { // Assuming 201 Created for successful post
        setInfo(prev => ({...prev, success: true}))
        cleanForm();
      } else {
        setInfo(prev => ({...prev, error: "Unexpected response from server."}))
      }
    } catch (err) {
      console.error("Error posting recipe:", err);
      if (axios.isAxiosError(err) && err.response) {
        setInfo(prev => ({...prev, error: err?.response?.data.message || "Error while posting recipe."}))
      } else {
        setInfo(prev => ({...prev, error: "An unexpected error occurred."}))
      }
    } finally {
      setInfo(prev => ({...prev, isSubmiting: false}))
    }
  }

  const cleanForm = () => {
    setForm({
      name_recipe: "",
      description: "",
      steps: "",
      guest_name: ""
    });

    setPreviewImage(null);
    setImageFile(null);
    setInfo(prev => ({
      ...prev,
      isSubmiting: false,
      error: "",
    }))
    setMealTypeSelected('')
  };

  return (
    <form className="max-w-4xl mx-auto p-6 font-sans space-y-8 justify-center" onSubmit={handleSubmit} >
      <div className="space-y-4">
        <label htmlFor="name_recipe" className="block text-sm/6 font-medium text-gray-900">
          Recipe name:
        </label>
        <input
          placeholder="Recipe name"
          id="name_recipe"
          name="name_recipe"
          type="string"
          required
          className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
          onChange={handleChange}
          value={form.name_recipe}
        />

        {user ? null : (
          <div>
            <label htmlFor="guest_name" className="block text-sm/6 font-medium text-gray-900">
              By:
            </label>
            <input
              placeholder="Author's name"
              id="guest_name"
              name="guest_name"
              type="string"
              required
              className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
              onChange={handleChange}
              value={form.guest_name}
              />
          </div>
        )}

        <div className="flex justify-center">
          <label className="relative group cursor-pointer">
            <div className="w-210 h-100 border-2 border-dashed flex items-center justify-center overflow-hidden">
              {previewImage ? (
                <img src={previewImage as string} alt="Preview" className="w-full h-full object-cover" />
              ) : (
                <div className="text-center">
                  <span className="text-xs text-pink-200 block mt-1">Upload photo</span>
                </div>
              )}
            </div>
            <input 
              type="file" 
              accept="image/*" 
              onChange={handleImageChange}
              className="hidden" 
            />
          </label>
        </div>

        <label htmlFor="description" className="block text-sm/6 font-medium text-gray-900">
          Small description
        </label>
        <textarea
          name="description"
          value={form.description}
          onChange={handleChange}
          className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
          rows="2"
          required
        />

        <label htmlFor="steps" className="block text-sm/6 font-medium text-gray-900">
          Steps
        </label>
        <textarea
          name="steps"
          value={form.steps}
          onChange={handleChange}
          className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
          rows="4"
          required
        />
        
        <label className="block text-sm/6 font-medium text-gray-900">
          Meal Type
        </label>
        <MealTypeFilter mealTypes={mealTypes} selected={mealTypeSelected} onChange={setMealTypeSelected}/>
      </div>

      {info.error && (
        <p className="text-red-500 text-sm text-center">{info.error}</p>
      )}
      {info.success && (
        <p className="text-green-500 text-sm text-center">Recipe posted successfully!</p>
      )}

      <div>
        <button
          type="submit"
          className="flex w-full justify-center rounded-md bg-green-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-green-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
          disabled={info.isSubmiting} // Disable button when submitting
        >
          {info.isSubmiting ? "Posting..." : "Post"}
        </button>
      </div>
    </form>
  );

}

export default PostRecipe
