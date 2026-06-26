import { NextRequest, NextResponse } from "next/server";

/**
 * 保護ルート: access_token Cookie がなければ /login にリダイレクト
 */
const PROTECTED_PATHS = [
  "/",
  "/watchlist",
  "/mypage",
  "/movie",
  "/watch-calendar",
];

export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  const isProtected = PROTECTED_PATHS.some(
    (p) => pathname === p || pathname.startsWith(p + "/"),
  );

  if (!isProtected) {
    return NextResponse.next();
  }

  const token = req.cookies.get("access_token")?.value;
  if (!token) {
    const loginUrl = req.nextUrl.clone();
    loginUrl.pathname = "/login";
    loginUrl.searchParams.set("redirect", pathname);
    return NextResponse.redirect(loginUrl);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * /api, /_next, /favicon.ico, /login, /register は除外
     */
    "/((?!api|_next/static|_next/image|favicon.ico|login|register).*)",
  ],
};
