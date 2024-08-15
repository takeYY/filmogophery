import { NextRequest } from "next/server";

export async function POST(req: NextRequest) {
  const url = `http://127.0.0.1:8000/movie/record`;

  console.log("app api から情報を取得しました。");

  const res = await fetch(url, { method: "POST" });
  return res;
}
