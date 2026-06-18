import type { Metadata, Viewport } from "next";
import "./globals.css";
import { AuthProvider } from "@/lib/auth";
import { SettingsProvider } from "@/lib/settings";
import { ServiceWorker } from "@/components/ServiceWorker";

export const metadata: Metadata = {
  title: "My Fitness Aibo",
  description: "Track your training plans, cycles and progress.",
  manifest: "/manifest.webmanifest",
  icons: {
    icon: "/icon-192.png",
    apple: "/icon-192.png",
  },
  appleWebApp: {
    capable: true,
    statusBarStyle: "black-translucent",
    title: "My Fitness Aibo",
  },
};

export const viewport: Viewport = {
  themeColor: "#0b0b0f",
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
  viewportFit: "cover",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-bg">
        <AuthProvider>
          <SettingsProvider>
            <div className="mx-auto w-full max-w-md min-h-screen flex flex-col">{children}</div>
          </SettingsProvider>
        </AuthProvider>
        <ServiceWorker />
      </body>
    </html>
  );
}
