import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function GET(
  req: NextRequest,
  { params }: { params: { id: string } }
) {
  const token = req.headers.get("authorization");
  const url = `${APIBaseURL}/movies/${params.id}`;

  const headers: HeadersInit = {};
  if (token) {
    headers.Authorization = token;
  }

  const res = await fetch(url, { headers });
  return res;
}
