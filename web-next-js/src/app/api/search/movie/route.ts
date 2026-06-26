import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const { searchParams } = new URL(req.url);
  const query = searchParams.get("query");

  const url = `${APIBaseURL}/search/movies?title=${query}`;
  console.log(`app apiから情報取得: ${query}`);

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = `${tokenType} ${token}`;
  }

  const res = await fetch(url, { headers });
  console.log("app apiから情報取得: 完了");
  return res;
}
