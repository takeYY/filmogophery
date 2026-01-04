// app/login/page.tsx
/**
 * ログイン
 * パス: /login
 */

"use client";

import { saveToken } from "@/utils/auth";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import React, { useState } from "react";

export default function Login() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const redirect = searchParams.get("redirect") || "/";

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      const response = await fetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data = await response.json();
        saveToken(data);
        router.push(redirect); // 元のページに戻る
      } else {
        const data = await response.json();
        setError(data.message || "ログインに失敗しました");
      }
    } catch (err) {
      setError("ログインに失敗しました");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main>
      <div className="container py-5">
        <div className="row justify-content-center">
          <div className="col-md-6">
            <h3 className="text-center mb-4">ログイン</h3>
            <form onSubmit={handleSubmit}>
              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}
              <div className="mb-3">
                <label htmlFor="email" className="form-label">
                  メールアドレス
                </label>
                <input
                  type="email"
                  className="form-control"
                  id="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>
              <div className="mb-3">
                <label htmlFor="password" className="form-label">
                  パスワード
                </label>
                <input
                  type="password"
                  className="form-control"
                  id="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
              <button
                type="submit"
                className="btn btn-info w-100"
                disabled={isLoading}
              >
                {isLoading ? "ログイン中..." : "ログイン"}
              </button>
            </form>
            <div className="text-center mt-3">
              <Link href="/register" className="text-decoration-none">
                アカウントをお持ちでない方はこちら
              </Link>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
