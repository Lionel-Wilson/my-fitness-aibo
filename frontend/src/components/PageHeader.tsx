"use client";

import { useRouter } from "next/navigation";

interface Props {
  title: string;
  subtitle?: string;
  back?: boolean;
  right?: React.ReactNode;
}

// Sticky mobile top bar with optional back button and a right-side action slot.
export function PageHeader({ title, subtitle, back, right }: Props) {
  const router = useRouter();
  return (
    <header className="safe-top sticky top-0 z-10 bg-bg/90 backdrop-blur border-b border-border">
      <div className="flex items-center gap-2 px-4 h-14">
        {back && (
          <button
            onClick={() => router.back()}
            className="-ml-2 px-2 py-2 text-2xl leading-none text-muted active:text-white"
            aria-label="Back"
          >
            ‹
          </button>
        )}
        <div className="flex-1 min-w-0">
          <h1 className="text-lg font-semibold truncate">{title}</h1>
          {subtitle && <p className="text-xs text-muted truncate">{subtitle}</p>}
        </div>
        {right}
      </div>
    </header>
  );
}
