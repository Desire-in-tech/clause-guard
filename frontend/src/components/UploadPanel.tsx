import { useRef, useState, type ChangeEvent } from "react";
import { extractFile } from "../api/client";

interface UploadPanelProps {
  value: string;
  disabled: boolean;
  onChange: (value: string) => void;
  onLoadingChange: (loading: boolean) => void;
}

export default function UploadPanel({
  value,
  disabled,
  onChange,
  onLoadingChange,
}: UploadPanelProps) {
  const fileInput = useRef<HTMLInputElement>(null);
  const [filename, setFilename] = useState("");
  const [isExtracting, setIsExtracting] = useState(false);
  const [uploadError, setUploadError] = useState("");

  async function handleFileChange(event: ChangeEvent<HTMLInputElement>) {
    const input = event.currentTarget;
    const file = input.files?.[0];
    if (!file) return;

    setFilename(file.name);
    setUploadError("");
    setIsExtracting(true);
    onLoadingChange(true);

    try {
      const extracted = await extractFile(file);
      setFilename(extracted.filename);
      onChange(extracted.text);
    } catch (error) {
      setUploadError(
        error instanceof Error && error.message !== "Failed to fetch"
          ? error.message
          : "Could not upload this file. Check your connection and try again.",
      );
    } finally {
      setIsExtracting(false);
      onLoadingChange(false);
      input.value = "";
    }
  }

  return (
    <>
      <div className="flex flex-col gap-3 border-b border-black/15 px-5 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <label
          htmlFor="contract-text"
          className="text-sm font-semibold text-black"
        >
          Contract text
        </label>
        <div className="flex flex-wrap items-center gap-3">
          <input
            ref={fileInput}
            id="contract-file"
            type="file"
            accept=".txt,.docx,.pdf"
            onChange={handleFileChange}
            disabled={disabled || isExtracting}
            className="sr-only"
            data-testid="input-contract-file"
          />
          <label
            htmlFor="contract-file"
            aria-disabled={disabled || isExtracting}
            className={`w-fit border border-dashed border-black/35 px-3 py-2 text-xs font-medium ${
              disabled || isExtracting
                ? "cursor-not-allowed text-black/45"
                : "cursor-pointer text-black hover:border-primary"
            }`}
            data-testid="button-upload-file"
          >
            {isExtracting ? "Extracting…" : "Upload file"}
          </label>
          <span
            className="max-w-56 truncate text-xs text-black/60"
            role="status"
            aria-live="polite"
            title={filename}
          >
            {isExtracting
              ? `Reading ${filename}…`
              : filename || "TXT, DOCX, PDF · max 5 MB"}
          </span>
        </div>
      </div>
      <div className="p-5 sm:p-6">
        {uploadError && (
          <p
            className="mb-4 border-l-2 border-[#bb2634] bg-[#fff0f1] px-4 py-3 text-sm text-[#8d1a25]"
            role="alert"
            data-testid="error-upload"
          >
            {uploadError}
          </p>
        )}
        <textarea
          id="contract-text"
          value={value}
          onChange={(event) => onChange(event.target.value)}
          placeholder="Paste your contract here…"
          rows={12}
          maxLength={100000}
          className="min-h-64 w-full resize-y border border-black/20 bg-white p-4 text-sm leading-7 text-black placeholder:text-black/40 focus:border-primary focus:outline-none"
          disabled={disabled || isExtracting}
          required
          data-testid="input-contract-text"
        />
      </div>
    </>
  );
}
