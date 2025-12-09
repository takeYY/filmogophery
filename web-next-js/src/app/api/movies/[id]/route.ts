import { APIBaseURL } from "@/constants/api";
import { NextRequest } from "next/server";

export async function GET(
  req: NextRequest,
  { params }: { params: { id: string } }
) {
  const url = `${APIBaseURL}/movies/${params.id}`;
  const res = await fetch(url);
  return res;
}
