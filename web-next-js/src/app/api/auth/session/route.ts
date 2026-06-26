import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

/**
 * セッション確認用エンドポイント
 * HttpOnly Cookie が有効かつバックエンドでも認証できる場合に 200 を返す
 */
export async function GET() {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  if (!token) {
    return NextResponse.json({ authenticated: false }, { status: 401 });
  }

  const res = await fetch(`${APIBaseURL}/users/me`, {
    headers: {
      Authorization: `${tokenType} ${token}`,
    },
  });

  if (!res.ok) {
    // トークンが無効なら Cookie も削除
    cookieStore.delete("access_token");
    cookieStore.delete("token_type");
    return NextResponse.json({ authenticated: false }, { status: 401 });
  }

  const user = await res.json();
  return NextResponse.json({ authenticated: true, user });
}
