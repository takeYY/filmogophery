import { APIBaseURL } from "@/constants/api";

export async function GET(request: Request) {
  const token = request.headers.get("Authorization");
  const { searchParams } = new URL(request.url);
  const limit = searchParams.get("limit") || "12";
  const offset = searchParams.get("offset") || "0";

  const res = await fetch(
    `${APIBaseURL}/users/me/watch-history?limit=${limit}&offset=${offset}`,
    {
      headers: {
        Authorization: token || "",
      },
    }
  );

  return res;
}
