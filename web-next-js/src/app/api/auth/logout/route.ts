import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function POST() {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  if (token) {
    await fetch(`${APIBaseURL}/auth/logout`, {
      method: "POST",
      headers: {
        Authorization: `${tokenType} ${token}`,
      },
    });
  }

  // Cookieを削除
  cookieStore.delete("access_token");
  cookieStore.delete("token_type");

  return NextResponse.json({ ok: true });
}
