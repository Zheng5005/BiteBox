interface StickyNoteProps {
  note: string;
}

const StickyNote: React.FC<StickyNoteProps> = ({ note }) => {
  if (!note) return null;

  return (
    <div className="relative mt-2 mx-4 mb-1 max-w-xs rotate-1">
      <div className="bg-yellow-200 shadow-md px-4 py-3 rounded-sm border-l-4 border-yellow-400">
        <p className="text-sm text-yellow-900 italic whitespace-pre-line">📝 {note}</p>
      </div>
    </div>
  );
};

export default StickyNote;
