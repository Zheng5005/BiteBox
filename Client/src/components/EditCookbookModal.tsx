import { useState } from 'react';
import { updateCookbook } from '../api/cookbooks';
import type { Cookbook } from '../types';

interface EditCookbookModalProps {
  cookbook: Cookbook;
  onClose: () => void;
  onUpdated: (updated: Cookbook) => void;
}

const EditCookbookModal: React.FC<EditCookbookModalProps> = ({ cookbook, onClose, onUpdated }) => {
  const [name, setName] = useState(cookbook.name);
  const [description, setDescription] = useState(cookbook.description);
  const [isPublic, setIsPublic] = useState(cookbook.is_public);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setSaving(true);
    setError(null);
    try {
      await updateCookbook(cookbook.id, {
        name: name.trim(),
        description: description.trim(),
        is_public: isPublic,
      });
      onUpdated({ ...cookbook, name: name.trim(), description: description.trim(), is_public: isPublic });
    } catch {
      setError('Failed to update cookbook.');
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
          <h2 className="text-xl font-bold">Edit Cookbook</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-2xl leading-none">
            &times;
          </button>
        </div>

        {error && <p className="text-red-500 text-sm mb-3">{error}</p>}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="edit-cookbook-name" className="block text-sm font-medium text-gray-900">
              Name
            </label>
            <input
              id="edit-cookbook-name"
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="My Cookbook"
              className="mt-1 block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            />
          </div>

          <div>
            <label htmlFor="edit-cookbook-description" className="block text-sm font-medium text-gray-900">
              Description
            </label>
            <textarea
              id="edit-cookbook-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="What's this cookbook about?"
              rows={3}
              className="mt-1 block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              id="edit-cookbook-public"
              type="checkbox"
              checked={isPublic}
              onChange={(e) => setIsPublic(e.target.checked)}
              className="h-4 w-4 rounded border-gray-300 text-green-600 focus:ring-green-500"
            />
            <label htmlFor="edit-cookbook-public" className="text-sm text-gray-700">
              Make this cookbook public
            </label>
          </div>

          <button
            type="submit"
            disabled={saving || !name.trim()}
            className="w-full px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition disabled:opacity-50"
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default EditCookbookModal;
