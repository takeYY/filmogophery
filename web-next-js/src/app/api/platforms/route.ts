import { APIBaseURL } from "@/constants/api";

export async function GET(req: Request) {
  const token = req.headers.get("authorization");
  const url = `${APIBaseURL}/platforms`;
  console.log("app apiから情報を取得中...");

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, { headers });
  return res;
}
