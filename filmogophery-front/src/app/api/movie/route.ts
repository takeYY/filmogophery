import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const movieID = searchParams.get("id");

  const url = `${APIBaseURL}/movies/${movieID}`;
  console.log("app apiから情報を取得中...");

  const res = await fetch(url);
  console.log("app apiから情報取得: 完了");
  return res;
}

export async function POST(req: NextRequest) {
  const url = `${APIBaseURL}/movie`;
  console.log("app apiからリクエスト中...");

  const res = await fetch(url, {
    method: "POST",
    body: JSON.stringify(req.json()),
  });
  console.log("app apiからリクエスト: 完了");

  return res;
}
