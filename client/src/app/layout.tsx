import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

// Using Inter for a clean, modern dashboard look
const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Contest Rating Engine",
  description: "Production-grade competitive programming rating platform.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>{children}</body>
    </html>
  );
}