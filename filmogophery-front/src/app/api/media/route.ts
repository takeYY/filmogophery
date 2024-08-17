import { APIBaseURL } from "@/constants/api";

export async function GET() {
  const url = `${APIBaseURL}/media`;

  console.log("app api から情報を取得しました。");

  const res = await fetch(url);
  return res;
}
