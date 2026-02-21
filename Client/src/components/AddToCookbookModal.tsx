import { useEffect, useState } from 'react';
import type { Cookbook } from '../types';
import { getCookbooks, addRecipeToCookbook } from '../api/cookbooks';

interface AddToCookbookModalProps {
  recipeId: number;
  onClose: () => void;
}

const AddToCookbookModal: React.FC<AddToCookbookModalProps> = ({ recipeId, onClose }) => {
  const [cookbooks, setCookbooks] = useState<Cookbook[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    async function fetchCookbooks() {
      try {
        setCookbooks(await getCookbooks());
      } catch {
        setError('Failed to load cookbooks.');
      } finally {
        setLoading(false);
      }
    }
    fetchCookbooks();
  }, []);

  const handleSelect = async (cookbook: Cookbook) => {
    if (saving) return;
    setSaving(true);
    setError(null);
    setSuccess(null);
    try {
      await addRecipeToCookbook(cookbook.id, recipeId);
      setSuccess(`Added to "${cookbook.name}"!`);
      setTimeout(onClose, 1200);
    } catch {
      setError('Failed to add recipe to cookbook.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div
        className="bg-white rounded-2xl shadow-lg w-full max-w-md mx-4 p-6"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold">Save to Cookbook</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-2xl leading-none">
            &times;
          </button>
        </div>

        {error && <p className="text-red-500 text-sm mb-3">{error}</p>}
        {success && <p className="text-green-600 text-sm mb-3">{success}</p>}

        {loading ? (
          <p className="text-center text-gray-500 py-4">Loading cookbooks...</p>
        ) : cookbooks.length === 0 ? (
          <p className="text-center text-gray-500 py-4">You don't have any cookbooks yet.</p>
        ) : (
          <ul className="space-y-2 max-h-64 overflow-y-auto">
            {cookbooks.map((cookbook) => (
              <li key={cookbook.id}>
                <button
                  onClick={() => handleSelect(cookbook)}
                  disabled={saving}
                  className="w-full text-left px-4 py-3 rounded-lg border border-gray-200 hover:bg-green-50 hover:border-green-400 transition disabled:opacity-50"
                >
                  <p className="font-semibold">{cookbook.name}</p>
                  {cookbook.description && (
                    <p className="text-gray-500 text-sm truncate">{cookbook.description}</p>
                  )}
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default AddToCookbookModal;
