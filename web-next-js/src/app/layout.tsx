"use client";

import { Inter } from "next/font/google";
import "bootstrap/dist/css/bootstrap.min.css";
import { NavLinks } from "@/app/components/nav-links";
import { usePathname } from "next/navigation";

const inter = Inter({ subsets: ["latin"] });

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname();
  const hideNav = pathname === "/login";

  return (
    <html lang="ja">
      <body className="bg-dark text-light">
        {!hideNav && <NavLinks />}
        {children}
      </body>
    </html>
  );
}
