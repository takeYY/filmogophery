import { User } from "@/interface";
import { useEffect, useState } from "react";

/**
 * ログイン中ユーザーの情報を取得するフック。
 * 認証ガードは Middleware に委譲しているため、このフックはリダイレクトしない。
 * NavLinks などユーザー情報を表示したいコンポーネントで使用する。
 */
export const useUser = () => {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await fetch("/api/auth/session");
        if (res.ok) {
          const data = await res.json();
          setUser(data.user ?? null);
        } else {
          setUser(null);
        }
      } catch {
        setUser(null);
      }
    };

    fetchUser();
  }, []);

  return { user };
};
