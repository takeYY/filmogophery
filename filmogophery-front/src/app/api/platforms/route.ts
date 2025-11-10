import { APIBaseURL } from "@/constants/api";

export async function GET() {
  const url = `${APIBaseURL}/platforms`;
  console.log("app apiから情報を取得中...");

  const res = await fetch(url);
  return res;
}
