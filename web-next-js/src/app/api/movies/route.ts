import { APIBaseURL } from "@/constants/api";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const offset = searchParams.get("offset") || "0";

  const url = `${APIBaseURL}/movies?offset=${offset}`;
  console.log("app apiから情報を取得中...");

  const res = await fetch(url);
  return res;
}
