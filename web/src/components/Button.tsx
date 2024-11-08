export function Button({ onClick, text }) {
    return (
        <button
            className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center"
            onClick={onClick}
        >
            {text}
        </button>
    );
}
