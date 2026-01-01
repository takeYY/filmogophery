import { APIBaseURL } from "@/constants/api";

// 視聴履歴一覧
export async function GET(
  req: Request,
  { params }: { params: { id: number } }
) {
  const reviewID = params.id;

  const token = req.headers.get("authorization");
  const url = `${APIBaseURL}/reviews/${reviewID}/history`;
  console.log("app apiから情報を取得中...");

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, { headers });
  console.log("app apiから情報取得: 完了");
  return res;
}
