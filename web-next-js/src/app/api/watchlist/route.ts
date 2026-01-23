import { APIBaseURL } from "@/constants/api";

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url);
  const offset = searchParams.get("offset") || "0";
  const token = req.headers.get("authorization");

  const url = `${APIBaseURL}/watchlist?offset=${offset}`;
  console.log("app apiから情報を取得中...");

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, { headers });
  return res;
}

export async function POST(req: Request) {
  const token = req.headers.get("authorization");
  const body = await req.json();

  const url = `${APIBaseURL}/watchlist`;

  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(body),
  });

  return res;
}
