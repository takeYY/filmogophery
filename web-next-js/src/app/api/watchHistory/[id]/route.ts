import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";

// 視聴履歴一覧
export async function GET(
  _req: Request,
  { params }: { params: { id: number } },
) {
  const reviewID = params.id;

  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const url = `${APIBaseURL}/reviews/${reviewID}/history`;
  console.log("app apiから情報を取得中...");

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = `${tokenType} ${token}`;
  }

  const res = await fetch(url, { headers });
  console.log("app apiから情報取得: 完了");
  return res;
}
