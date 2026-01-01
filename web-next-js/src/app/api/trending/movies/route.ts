import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const token = req.headers.get("authorization");
  const url = `${APIBaseURL}/trending/movies`;

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, { headers });
  console.log("app apiから情報取得: 完了");
  return res;
}
