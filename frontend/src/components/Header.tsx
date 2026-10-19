import { Link } from "react-router-dom";

export default function Header() {
  return (
    <header className="border-b border-black/15 bg-white">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-5 py-4 sm:px-8">
        <Link
          to="/"
          className="flex items-center gap-3 font-semibold tracking-tight"
          aria-label="ClauseGuard home"
          data-testid="link-home"
        >
          <span className="grid size-9 place-items-center bg-primary text-sm font-bold text-white">
            CG
          </span>
          <span className="text-lg text-black">ClauseGuard</span>
        </Link>
        <span className="text-xs font-medium uppercase tracking-[0.14em] text-black/55">
          Contract review
        </span>
      </div>
    </header>
  );
}
