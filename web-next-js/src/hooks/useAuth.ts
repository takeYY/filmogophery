import { User } from "@/interface";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

/**
 * 認証ガードフック
 * HttpOnly Cookie の有効性を /api/auth/session で確認する。
 * 未認証の場合は /login にリダイレクトする。
 */
export const useAuth = () => {
  const router = useRouter();
  const pathname = usePathname();
  const [user, setUser] = useState<User | null>(null);
  const [checked, setChecked] = useState(false);

  useEffect(() => {
    const checkSession = async () => {
      try {
        const res = await fetch("/api/auth/session");
        if (!res.ok) {
          router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
          return;
        }
        const data = await res.json();
        setUser(data.user ?? null);
      } catch {
        router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
      } finally {
        setChecked(true);
      }
    };

    checkSession();
  }, [router, pathname]);

  return { user, checked };
};
