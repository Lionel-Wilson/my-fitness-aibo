"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const items = [
  { href: "/plans", label: "Plans", icon: "🏋️" },
  { href: "/dashboard", label: "Progress", icon: "📈" },
  { href: "/settings", label: "Settings", icon: "⚙️" },
];

// Fixed bottom tab bar (mobile-first navigation).
export function BottomNav() {
  const pathname = usePathname();
  return (
    <nav className="safe-bottom sticky bottom-0 z-10 bg-surface/95 backdrop-blur border-t border-border">
      <div className="flex">
        {items.map((it) => {
          const active = pathname === it.href || pathname.startsWith(it.href + "/");
          return (
            <Link
              key={it.href}
              href={it.href}
              className={`flex-1 flex flex-col items-center gap-0.5 py-2.5 text-xs ${
                active ? "text-accent" : "text-muted"
              }`}
            >
              <span className="text-lg">{it.icon}</span>
              {it.label}
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
