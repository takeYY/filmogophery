import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function POST(req: NextRequest) {
  const url = `${APIBaseURL}/movie/record`;

  console.log("app api から情報を取得しました。");

  const res = await fetch(url, { method: "POST" });
  return res;
}
