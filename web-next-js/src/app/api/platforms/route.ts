import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";

export async function GET(req: Request) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const url = `${APIBaseURL}/platforms`;
  console.log("app apiから情報を取得中...");

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = `${tokenType} ${token}`;
  }

  const res = await fetch(url, { headers });
  return res;
}
