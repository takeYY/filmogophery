import { clearToken, getToken } from "@/utils/auth";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

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
