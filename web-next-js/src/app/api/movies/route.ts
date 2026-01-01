import { APIBaseURL } from "@/constants/api";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const offset = searchParams.get("offset") || "0";
  const token = req.headers.get("authorization");

  const url = `${APIBaseURL}/movies?offset=${offset}`;
  console.log("app apiから情報を取得中...");

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, { headers });
  return res;
}
