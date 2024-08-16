import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const query = searchParams.get("query");

  const url = `http://127.0.0.1:8000/tmdb/search/movies?query=${query}`;
  console.log(`app apiから情報取得: ${query}`);

  const res = await fetch(url);
  console.log("app apiから情報取得: 完了");
  return res;
}
