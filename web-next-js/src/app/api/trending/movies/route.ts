import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const url = `${APIBaseURL}/trending/movies`;

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = `${tokenType} ${token}`;
  }

  const res = await fetch(url, { headers });
  console.log("app apiから情報取得: 完了");
  return res;
}
