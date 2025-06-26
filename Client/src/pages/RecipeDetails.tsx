import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import { useAuth } from '../context/AuthContext';

interface Recipe {
  id: number;
  recipe_name: string;
  description: string;
  meal_type_id: string;
  image: string;
  creator_name: string;
  rating: number;
  steps: string[];
}

interface Comment {
  id: number;
  user_name: string;
  recipe_id: string;
  comment: string;
  rating: number;
}

const RecipeDetails: React.FC = () => {
  const [recipe, setRecipe] = useState<Recipe | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState({
    comment: "",
    rating: 1
  });
  const { id } = useParams()
  const { user } = useAuth();
  const token = localStorage.getItem("token")

  useEffect(() => {
    async function fetchData() {
      const res = await fetch(`http://localhost:8080/api/recipes/${id}`);
      const recipeData = await res.json();
      setRecipe(recipeData);

      const commentsRes = await fetch(`http://localhost:8080/api/comments/${id}`);
      const commentsData = await commentsRes.json();

      if (Array.isArray(commentsData)) {
        setComments(commentsData);
      } else {
        setComments([])
      }
    }

    fetchData();
  }, [id, comments]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setNewComment(prev => ({
      ...prev,
      [name]: name === "rating" ? Number(value) : value
    }))

  };

  const handleCommentSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.comment.trim() || newComment.rating == 0) return;

    await fetch(`http://localhost:8080/api/comments/post/${id}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ comment: newComment.comment, rating: newComment.rating }),
    });

    setNewComment({
      comment: "",
      rating: 0
    });

    // Refresh comments
    const res = await fetch(`http://localhost:8080/api/comments/${id}`);
    const resData = res.json()
    if (Array.isArray(resData)) {
      setComments(resData);
    } else {
      setComments([])
    }
  };

  if (!recipe) return <p className="text-center">Loading recipe...</p>;

  return (
    <div className="max-w-4xl mx-auto p-6 font-sans space-y-8">
      {/* Recipe info */}
      <div className="space-y-4">
        <h1 className="text-3xl font-bold">{recipe.recipe_name}</h1>
        <p className="text-sm text-gray-600">
          By: {recipe.creator_name ?? 'Anonymous'}
        </p>
        <img
          src="https://i.pravatar.cc/150?u=a042581f4e29026024d"
          //src={recipe.image}
          alt={recipe.recipe_name}
          className="w-full max-h-96 object-cover rounded-lg shadow"
        />
        {recipe.rating == 0 ? (
          <div className="text-yellow-500 text-lg">⭐ No rating yet</div>
        ) : (
          <div className="text-yellow-500 text-lg">⭐ {recipe.rating}</div>
        )}
        <p className="text-gray-700">{recipe.description}</p>
        <div>
          <h2 className="text-xl font-semibold mt-4 mb-2">Steps:</h2>
          <p className="text-lg mt-4 mb-2">{recipe.steps}</p>
        </div>
      </div>

      {/* Comments */}
      <div className="space-y-4">
        <h2 className="text-2xl font-bold">Comments</h2>
        {comments.length === 0 ? (
          <p className="text-gray-500">No comments yet</p>
        ) : (
          <ul className="space-y-3">
            {comments.map((comment) => (
              <li key={comment.id} className="border-b pb-2">
                <p className="font-semibold">{comment.user_name}</p>
                <p className="text-gray-700">{comment.comment}</p>
                <span className="text-yellow-500 text-lg">⭐ {comment.rating}</span>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Comment Form */}
      { user && (
        <form onSubmit={handleCommentSubmit} className="space-y-4">
          <h3 className="text-xl font-semibold">Leave a Comment</h3>
          <textarea
            id="comment"
            name="comment"
            required
            value={newComment.comment}
            onChange={handleChange}
            placeholder="Write something nice..."
            className="w-full p-2 border border-gray-300 rounded-md focus:ring"
            rows={4}
          />

          <div>
            <label htmlFor="rating" className="block text-sm/6 font-medium text-gray-900">
              Rating
            </label>
            <div className="mt-2">
              <input
                id="rating"
                name="rating"
                type="number"
                required
                min="1"
                max="5"
                step="0.5"
                className="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                onChange={handleChange}
                value={newComment.rating}
              />
            </div>
          </div>

          <button
            type="submit"
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition"
          >
            Submit
          </button>
        </form>
      )}
    </div>
  );
};

export default RecipeDetails;
