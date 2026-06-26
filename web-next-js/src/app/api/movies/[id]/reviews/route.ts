import { APIBaseURL } from "@/constants/api";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function POST(
  req: Request,
  { params }: { params: { id: string } },
) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const tokenType = cookieStore.get("token_type")?.value ?? "Bearer";

  const url = `${APIBaseURL}/movies/${params.id}/reviews`;

  try {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
    };
    if (token) {
      headers.Authorization = `${tokenType} ${token}`;
    }

    const requestData = await req.json();
    const processedData = {
      ...requestData,
      rating: parseFloat(requestData.rating), // float にする
    };
    console.log("リクエストデータ:", processedData);

    const res = await fetch(url, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(processedData),
    });
    console.log(`APIレスポンスステータス: ${res.status}`);

    return res;
  } catch (error) {
    console.error("Error in API route:", error);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 },
    );
  }
}
