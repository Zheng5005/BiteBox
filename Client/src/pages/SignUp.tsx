import { useState } from "react"
import axios from "axios"
import { signup } from "../api/auth"
import { Link } from "react-router"

const SignUp: React.FC = () => {
  const [form, setForm] = useState({
    user_name: "",
    email: "",
    password: "",
  })
  const [previewImage, setPreviewImage] = useState<string | ArrayBuffer | null>(null)
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [info, setInfo] = useState({
    isSubmiting: false,
    error: "",
    success: false,
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setInfo(prev => ({...prev, error: ""}));

    if (name === 'user_name') {
      if (/[0-9]/.test(value)) {
        setInfo(prev => ({...prev, error: "Numbers not allowed in username."}))
        return;
      } else if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/.test(value)) {
        setInfo(prev => ({...prev, error: "No special characters allowed in username."}))
        return;
      }
    }
    
    setForm({ ...form, [name]: value });
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files ? e.target.files[0] : null;
    if (!file) return;

    setInfo(prev => ({...prev, error: ""}));

    if (!file.type.match('image.*')) {
      setInfo(prev => ({...prev, error: "Only PNG or JPEG images are allowed."}))
      return;
    }

    if (file.size > 2 * 1024 * 1024) {
      setInfo(prev => ({...prev, error: "Image should be less than 2MB."}))
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

    setInfo(prev => ({...prev, isSubmiting: true, error: "", success: false}))

    const formData = new FormData();
    formData.append("name", form.user_name);
    formData.append("email", form.email);
    formData.append("password", form.password);
    if (imageFile) {
      formData.append("image", imageFile);
    }

    try {
      const response = await signup(formData);

      if (response.status == 201) {
        setInfo(prev => ({...prev, success: true}));
        setTimeout(() => {
          window.location.href = "/login"
        }, 1500)
      } else {
        setInfo(prev => ({...prev, error: response.data.message || "Sign up failed."}));
      }
      cleanForm();
    } catch (err) {
      console.error("Error during sign up:", err);
      if (axios.isAxiosError(err) && err.response) {
        setInfo(prev => ({...prev, error: err?.response?.data.message || "Error while signing up."}));
      } else {
        setInfo(prev => ({...prev, error: "An unexpected error occurred during sign up."}));
      }
    } finally {
      setInfo(prev => ({...prev, isSubmiting: false}))
    }
  };

  const cleanForm = () => {
    setForm({
      user_name: "",
      email: "",
      password: "",
    });

    setPreviewImage(null);
    setImageFile(null);
  };

  return (
    <>
      <div className="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-sm">
          <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">
            Sign Up
          </h2>
        </div>

        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="flex justify-center">
              <label className="relative group cursor-pointer">
                <div className="w-24 h-24 rounded-full border-2 border-dashed flex items-center justify-center overflow-hidden">
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
                <div className="absolute inset-0 bg-black bg-opacity-30 opacity-0 group-hover:opacity-100 rounded-full flex items-center justify-center transition">
                </div>
              </label>
            </div>


            <div>
              <label htmlFor="user_name" className="block text-sm/6 font-medium text-gray-900">
                Username
              </label>
              <div className="mt-2">
                <input
                  id="user_name"
                  name="user_name"
                  type="string"
                  required
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                  onChange={handleChange}
                  value={form.user_name}
                />
              </div>
            </div>

            <div>
              <label htmlFor="email" className="block text-sm/6 font-medium text-gray-900">
                Email address
              </label>
              <div className="mt-2">
                <input
                  id="email"
                  name="email"
                  type="email"
                  required
                  autoComplete="email"
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                  onChange={handleChange}
                  value={form.email}
                />
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between">
                <label htmlFor="password" className="block text-sm/6 font-medium text-gray-900">
                  Password
                </label>
              </div>
              <div className="mt-2">
                <input
                  id="password"
                  name="password"
                  type="password"
                  required
                  autoComplete="current-password"
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                  onChange={handleChange}
                  value={form.password}
                />
              </div>
            </div>

            {info.error && (
              <p className="text-red-500 text-sm text-center">{info.error}</p>
            )}
            {info.success && (
              <p className="text-green-500 text-sm text-center">Sign Up successful! Redirecting...</p>
            )}

            <div>
              <button
                type="submit"
                className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                disabled={info.isSubmiting} // Disable button when submitting
              >
                {info.isSubmiting ? "Signing Up..." : "Sign up"}
              </button>
            </div>
          </form>

          <p className="mt-10 text-center text-sm/6 text-gray-500">
            Already have an account? 
            <Link to="/login" className="font-semibold text-indigo-600 hover:text-indigo-500">
              Log In
            </Link>
          </p>
        </div>
      </div>
    </>
  )
}

export default SignUp
