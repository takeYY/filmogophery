import { APIBaseURL } from "@/constants/api";
import { NextResponse } from "next/server";

export async function POST(
  req: Request,
  { params }: { params: { _: string; reviewId: string } }
) {
  const token = req.headers.get("authorization");
  const url = `${APIBaseURL}/reviews/${params.reviewId}/history`;

  try {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
      "Access-Control-Allow-Origin": "*",
    };
    if (token) {
      headers.Authorization = token;
    }

    const requestData = await req.json();
    const processedData = {
      ...requestData,
      platformId: parseInt(requestData.platformId), // int にする
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
      { status: 500 }
    );
  }
}

export async function PUT(
  req: Request,
  { params }: { params: { id: string; reviewId: string } }
) {
  const token = req.headers.get("authorization");
  const url = `${APIBaseURL}/reviews/${params.reviewId}`;

  try {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
      "Access-Control-Allow-Origin": "*",
    };
    if (token) {
      headers.Authorization = token;
    }

    const requestData = await req.json();
    const processedData = {
      ...requestData,
      rating: parseFloat(requestData.rating), // float にする
    };
    console.log("リクエストデータ:", processedData);

    const res = await fetch(url, {
      method: "PUT",
      headers: headers,
      body: JSON.stringify(processedData),
    });
    console.log(`APIレスポンスステータス: ${res.status}`);

    return res;
  } catch (error) {
    console.error("Error in API route:", error);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
