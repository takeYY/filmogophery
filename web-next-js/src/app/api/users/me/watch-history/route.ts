import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";

export async function GET(request: Request) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const { searchParams } = new URL(request.url);
  const limit = searchParams.get("limit") || "12";
  const offset = searchParams.get("offset") || "0";

  const res = await fetch(
    `${APIBaseURL}/users/me/watch-history?limit=${limit}&offset=${offset}`,
    {
      headers: {
        Authorization: token ? `${tokenType} ${token}` : "",
      },
    },
  );

  return res;
}
