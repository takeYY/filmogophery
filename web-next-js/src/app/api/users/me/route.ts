import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";

export async function GET(request: Request) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const res = await fetch(`${APIBaseURL}/users/me`, {
    headers: {
      Authorization: token ? `${tokenType} ${token}` : "",
    },
  });

  return res;
}
