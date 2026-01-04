"use client";

import { NavLinks } from "@/components/nav-links";
import "bootstrap/dist/css/bootstrap.min.css";
import { usePathname } from "next/navigation";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname();
  const hideNav = pathname === "/login" || pathname === "/register";

  return (
    <html lang="ja">
      <head>
        <link
          rel="stylesheet"
          href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.0/font/bootstrap-icons.css"
        />
      </head>
      <body className="bg-dark text-light">
        {!hideNav && <NavLinks />}
        {children}
      </body>
    </html>
  );
}
