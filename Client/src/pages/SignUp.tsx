import { useState } from "react"

const SignUp: React.FC = () => {
  const [form, setForm] = useState({
    user_name: "",
    email: "",
    password: "",
  })
  const [previewImage, setPreviewImage] = useState<string | ArrayBuffer | null>(null)
  const [imageFile, setImageFile] = useState<string | Blob | null>(null)
  const [info, setInfo] = useState({
    isSubmiting: false,
    error: "",
  })

  const handleChange = (e: any) => {
    const { name, value } = e.target;

    if (name === 'user_name') {
      if (/[0-9]/.test(value)) {
        setInfo({...info, error: "Numbers not allowed"})
        return;
      } else if (/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/.test(value)) {
        setInfo({...info, error: "No special characters allowed"})
        return;
      }
    }
    
    setForm({ ...form, [name]: value });
  };

  const handleImageChange = (e: any) => {
    const file = e.target.files[0];
    if (!file) return;


    if (!file.type.match('image.*')) {
      setInfo({...info, error: "Only PNG or JPEG images"})
      return;
    }

    if (file.size > 2 * 1024 * 1024) {
      setInfo({...info, error: "Image should be less than 2MB"})
      return;
    }

    setImageFile(file);
  
    const reader = new FileReader();
    reader.onloadend = () => {
      setPreviewImage(reader.result);
    };
    reader.readAsDataURL(file);
  };

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    if (info.isSubmiting) return

    setInfo({...info, isSubmiting: true})

    const formData = new FormData();
    formData.append("name", `${form.user_name}`);
    formData.append("email", form.email);
    formData.append("password", form.password);
    if (imageFile) {
      formData.append("image", imageFile);
    }

    try {
      const response = await fetch("http://localhost:8080/api/auth/signup", {
        method: "POST",
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error("Error en el registro");
      }

      cleanForm();
      
      if (data.token) {
        localStorage.setItem("token", data.token);
        //setTimeout(() => navigate("/profile"), 1500);
      }
    } catch (error) {
      setInfo({...info, error: "Error while signing up, please try again"})
    } finally {
      setInfo({...info, isSubmiting: false})
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
    setInfo({
      isSubmiting: false,
      error: "",
    })
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
                    <img src={previewImage} alt="Preview" className="w-full h-full object-cover" />
                  ) : (
                    <div className="text-center">
                      {/* <FiCamera className="mx-auto text-[#EFB8C8] text-2xl" />   */} 
                      <span className="text-xs text-pink-200 block mt-1">Subir foto</span>
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
                  {/* <FiCamera className="text-white text-xl" /> */}
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

            {info.error != "" ? (
              <label className="block text-sm/6 font-medium text-gray-900">
                {info.error}
              </label>
            ): null}

            <div>
              <button
                type="submit"
                className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Sign up
              </button>
            </div>
          </form>

          <p className="mt-10 text-center text-sm/6 text-gray-500">
            Already have an account? 
            <a href="#" className="font-semibold text-indigo-600 hover:text-indigo-500">
              Log In
            </a>
          </p>
        </div>
      </div>
    </>
  )
}

export default SignUp
