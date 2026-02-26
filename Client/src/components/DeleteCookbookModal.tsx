import { useState } from 'react';
import { deleteCookbook } from '../api/cookbooks';
import type { Cookbook } from '../types';

interface DeleteCookbookModalProps {
  cookbook: Cookbook;
  onClose: () => void;
  onDeleted: () => void;
}

const DeleteCookbookModal: React.FC<DeleteCookbookModalProps> = ({ cookbook, onClose, onDeleted }) => {
  const [deleting, setDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleDelete = async () => {
    setDeleting(true);
    setError(null);
    try {
      await deleteCookbook(cookbook.id);
      onDeleted();
    } catch {
      setError('Failed to delete cookbook.');
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div
        className="bg-white rounded-2xl shadow-lg w-full max-w-md mx-4 p-6"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold">Delete Cookbook</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-2xl leading-none">
            &times;
          </button>
        </div>

        <p className="text-gray-700 mb-2">
          Are you sure you want to delete <span className="font-semibold">{cookbook.name}</span>?
        </p>
        <p className="text-gray-500 text-sm mb-4">
          This action cannot be undone. All saved recipes in this cookbook will be removed.
        </p>

        {error && <p className="text-red-500 text-sm mb-3">{error}</p>}

        <div className="flex gap-3">
          <button
            onClick={onClose}
            disabled={deleting}
            className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="flex-1 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition disabled:opacity-50"
          >
            {deleting ? 'Deleting...' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default DeleteCookbookModal;
