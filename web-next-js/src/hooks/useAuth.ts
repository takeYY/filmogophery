import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { getToken, clearToken } from "@/utils/auth";

export const useAuth = () => {
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
      return;
    }

    // 有効期限チェック
    const expiresAt = new Date(token.expiresAt);
    const now = new Date();

    if (expiresAt <= now) {
      clearToken();
      router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
    }
  }, [router, pathname]);

  return getToken();
};
