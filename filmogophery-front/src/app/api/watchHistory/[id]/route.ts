import { APIBaseURL } from "@/constants/api";

// 視聴履歴一覧
export async function GET(_: Request, { params }: { params: { id: number } }) {
  const reviewID = params.id;

  const url = `${APIBaseURL}/reviews/${reviewID}/history`;
  console.log("app apiから情報を取得中...");

  const res = await fetch(url);
  console.log("app apiから情報取得: 完了");
  return res;
}
