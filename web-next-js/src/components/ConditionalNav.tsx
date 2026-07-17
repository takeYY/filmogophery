"use client";

import { NavLinks } from "@/components/nav-links";
import { usePathname } from "next/navigation";

export function ConditionalNav() {
  const pathname = usePathname();
  if (pathname === "/login" || pathname === "/register") return null;
  return <NavLinks />;
}
