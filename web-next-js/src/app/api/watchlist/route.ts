import { APIBaseURL } from "@/constants/api";

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
