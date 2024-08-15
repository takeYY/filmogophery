import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const movieID = searchParams.get("id");

  const url = `http://127.0.0.1:8000/movies/${movieID}`;
  console.log("app apiから情報を取得中...");

  const res = await fetch(url);
  console.log("app apiから情報取得: 完了");
  return res;
}
