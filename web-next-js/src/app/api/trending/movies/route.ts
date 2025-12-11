import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const url = `${APIBaseURL}/trending/movies`;

  const res = await fetch(url);
  console.log("app apiから情報取得: 完了");
  return res;
}
