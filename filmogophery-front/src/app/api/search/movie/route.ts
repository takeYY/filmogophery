import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const query = searchParams.get("query");

  const url = `${APIBaseURL}/tmdb/search/movies?query=${query}`;
  console.log(`app apiから情報取得: ${query}`);

  const res = await fetch(url);
  console.log("app apiから情報取得: 完了");
  return res;
}
