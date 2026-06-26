import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";
import { NextRequest } from "next/server";

export async function GET(
  _req: NextRequest,
  { params }: { params: { id: string } },
) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const url = `${APIBaseURL}/movies/${params.id}`;

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = `${tokenType} ${token}`;
  }

  const res = await fetch(url, { headers });
  return res;
}
