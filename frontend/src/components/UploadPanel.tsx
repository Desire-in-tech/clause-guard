interface UploadPanelProps {
  value: string;
  disabled: boolean;
  onChange: (value: string) => void;
}

export default function UploadPanel({
  value,
  disabled,
  onChange,
}: UploadPanelProps) {
  return (
    <>
      <div className="flex flex-col gap-3 border-b border-black/15 px-5 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <label
          htmlFor="contract-text"
          className="text-sm font-semibold text-black"
        >
          Contract text
        </label>
        <button
          type="button"
          disabled
          aria-disabled="true"
          className="w-fit cursor-not-allowed border border-dashed border-black/25 px-3 py-2 text-xs font-medium text-black/45"
          data-testid="button-upload-file"
        >
          Upload file (coming soon)
        </button>
      </div>
      <div className="p-5 sm:p-6">
        <textarea
          id="contract-text"
          value={value}
          onChange={(event) => onChange(event.target.value)}
          placeholder="Paste your contract here…"
          rows={12}
          maxLength={100000}
          className="min-h-64 w-full resize-y border border-black/20 bg-white p-4 text-sm leading-7 text-black placeholder:text-black/40 focus:border-primary focus:outline-none"
          disabled={disabled}
          required
          data-testid="input-contract-text"
        />
      </div>
    </>
  );
}
