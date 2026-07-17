import { ConditionalNav } from "@/components/ConditionalNav";
import "bootstrap/dist/css/bootstrap.min.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Filmogophery",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <head>
        <link
          rel="stylesheet"
          href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.0/font/bootstrap-icons.css"
        />
      </head>
      <body className="bg-dark text-light">
        <ConditionalNav />
        {children}
      </body>
    </html>
  );
}
