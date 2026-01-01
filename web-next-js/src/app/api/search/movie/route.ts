import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const token = req.headers.get("authorization");
  const { searchParams } = new URL(req.url);
  const query = searchParams.get("query");

  const url = `${APIBaseURL}/search/movies?title=${query}`;
  console.log(`app apiから情報取得: ${query}`);

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, { headers });
  console.log("app apiから情報取得: 完了");
  return res;
}
