import Disclaimer from "./Disclaimer";

export default function Footer() {
  return (
    <footer className="mt-16 border-t border-black/15">
      <div className="mx-auto flex max-w-6xl flex-col gap-3 px-5 py-7 text-sm sm:flex-row sm:items-center sm:justify-between sm:px-8">
        <p className="font-medium text-black">ClauseGuard</p>
        <Disclaimer />
      </div>
    </footer>
  );
}
