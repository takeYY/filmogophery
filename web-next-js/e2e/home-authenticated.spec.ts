import type { Page, Route } from "@playwright/test";
import { expect, test } from "@playwright/test";

/**
 * Middleware が access_token Cookie を見てリダイレクトするため、
 * テスト用のダミー Cookie をセットして Middleware を通過させる。
 */
const setAuthCookie = async (page: Page) => {
  await page.context().addCookies([
    {
      name: "access_token",
      value: "test-token",
      domain: "localhost",
      path: "/",
      httpOnly: true,
      sameSite: "Lax",
    },
  ]);
};

/** セッションAPIをモック（NavLinks / useAuth 向け） */
const mockAuthSession = async (page: Page) => {
  await page.route("**/api/auth/session", (route: Route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        authenticated: true,
        user: { id: 1, username: "testuser", email: "test@example.com" },
      }),
    }),
  );
};

/** 映画データなしのホームAPIをモック */
const mockEmptyHomeAPIs = async (page: Page) => {
  await page.route("**/api/movies**", (route: Route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: "[]" }),
  );
  await page.route("**/api/trending/movies", (route: Route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: "[]" }),
  );
};

/**
 * ホーム画面へ遷移してレンダリング完了を待つ。
 * useAuth の checked チェーンが不要になったため、
 * networkidle でのシンプルな待機で対応できる。
 */
const gotoHome = async (page: Page) => {
  await page.goto("/");
  await page.waitForLoadState("networkidle");
  await expect(page.getByRole("heading", { name: "Home" })).toBeVisible({
    timeout: 15000,
  });
};

/** レビュー済み映画データのモック */
const MOCK_MOVIES = [
  {
    id: 1,
    title: "テスト映画1",
    overview: "テスト映画1の概要",
    releaseDate: "2024-01-01",
    runtimeMinutes: 120,
    posterURL: null,
    tmdbID: 1001,
    genres: [{ code: "action", name: "アクション" }],
  },
  {
    id: 2,
    title: "テスト映画2",
    overview: "テスト映画2の概要",
    releaseDate: "2024-01-02",
    runtimeMinutes: 90,
    posterURL: null,
    tmdbID: 1002,
    genres: [],
  },
];

/** トレンド映画データのモック */
const MOCK_TRENDING = [
  { id: 10, title: "トレンド映画1", posterURL: null, hasReview: false },
  { id: 11, title: "トレンド映画2", posterURL: null, hasReview: true },
];

test.describe("ホーム画面（ログイン後）", () => {
  test.beforeEach(async ({ page }) => {
    await setAuthCookie(page);
    await mockAuthSession(page);
    await mockEmptyHomeAPIs(page);
  });

  test("認証済みでホームにアクセスするとリダイレクトされない", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(page).toHaveURL("/");
    await expect(page).not.toHaveURL(/\/login/);
  });

  test("ホーム画面の見出しが表示される", async ({ page }) => {
    await gotoHome(page);
    // gotoHome 内で既に heading を確認済み
  });

  test("ナビゲーションバーが表示される", async ({ page }) => {
    await gotoHome(page);
    await expect(page.getByRole("navigation")).toBeVisible();
    await expect(
      page.getByRole("link", { name: "Filmogophery" }),
    ).toBeVisible();
  });

  test("レビュー済み映画がない場合に空状態のメッセージが表示される", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(
      page.getByText("まだレビューした映画がありません。"),
    ).toBeVisible();
    await expect(
      page.getByRole("button", { name: "映画を探す" }),
    ).toBeVisible();
  });

  test("「映画を探す」ボタンをクリックすると検索画面へ遷移する", async ({
    page,
  }) => {
    await gotoHome(page);
    await page.getByRole("button", { name: "映画を探す" }).click();
    await expect(page).toHaveURL(/\/search/);
  });

  test("データ取得中はローディングスピナーが表示される", async ({ page }) => {
    // movies API だけ遅延させてローディング状態を確認
    await page.unroute("**/api/movies**");
    await page.unroute("**/api/trending/movies");
    await page.route("**/api/movies**", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 3000));
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: "[]",
      });
    });
    await page.route("**/api/trending/movies", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 3000));
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: "[]",
      });
    });

    await page.goto("/");
    // checked が true になりデータフェッチが始まった直後にスピナーが見える
    await expect(page.getByRole("status")).toBeVisible({ timeout: 10000 });
  });

  test("認証済みユーザーが /login に直接アクセスするとログイン画面がそのまま表示される", async ({
    page,
  }) => {
    await page.goto("/login");
    await expect(page).toHaveURL(/\/login/);
    await expect(page.getByRole("heading", { name: "ログイン" })).toBeVisible();
  });

  test("visual snapshot（レビューなし）", async ({ page }) => {
    await gotoHome(page);
    await expect(page).toHaveScreenshot("home-authenticated-empty.png");
  });
});

test.describe("ホーム画面（ログイン後・映画データあり）", () => {
  test.beforeEach(async ({ page }) => {
    await setAuthCookie(page);
    await mockAuthSession(page);
    await page.route("**/api/movies**", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_MOVIES),
      }),
    );
    await page.route("**/api/trending/movies", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(MOCK_TRENDING),
      }),
    );
  });

  test("レビュー済み映画セクションが表示される", async ({ page }) => {
    await gotoHome(page);
    await expect(page.getByText("レビュー済み映画")).toBeVisible();
  });

  test("トレンド（最近の映画）セクションが表示される", async ({ page }) => {
    await gotoHome(page);
    await expect(page.getByText("最近の映画")).toBeVisible();
  });

  test("映画カードをクリックすると詳細画面へ遷移する", async ({ page }) => {
    await gotoHome(page);
    await expect(page.getByText("レビュー済み映画")).toBeVisible();
    await page.locator(".row.row-cols-md-3 .col").first().click();
    await expect(page).toHaveURL(/\/movie\/1/);
  });

  test("hasReview:false のトレンド映画をクリックするとレビュー作成画面へ遷移する", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(page.getByText("最近の映画")).toBeVisible();
    const posterImages = page.getByAltText("ポスター画像");
    await posterImages.first().click();
    await expect(page).toHaveURL(/\/movie\/10\/review\/create/);
  });

  test("hasReview:true のトレンド映画をクリックすると詳細画面へ遷移する", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(page.getByText("最近の映画")).toBeVisible();
    const posterImages = page.getByAltText("ポスター画像");
    await posterImages.nth(1).click();
    await expect(page).toHaveURL(/\/movie\/11/);
  });

  test("ナビバーの検索フォームにキーワードを入力してsubmitすると検索結果画面へ遷移する", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(page.getByRole("navigation")).toBeVisible();

    const searchBox = page.getByRole("searchbox");
    await searchBox.fill("インセプション");
    // フォームを Enter キーで submit
    await searchBox.press("Enter");

    await expect(page).toHaveURL(
      /\/search\/movie\?query=%E3%82%A4%E3%83%B3%E3%82%BB%E3%83%97%E3%82%B7%E3%83%A7%E3%83%B3/,
    );
  });

  test("検索ワードが空のときは検索ボタンが無効化されている", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(page.getByRole("navigation")).toBeVisible();
    // form[role="search"] 内の submit ボタンを取得
    const searchButton = page.locator(
      'form[role="search"] button[type="submit"]',
    );
    await expect(searchButton).toBeDisabled();
  });

  test("ナビバーのユーザーアイコンをクリックするとドロップダウンが開く", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(page.getByRole("navigation")).toBeVisible();

    await page.locator(".bi-person-circle").click();
    await expect(page.getByText("testuser")).toBeVisible();
    await expect(page.getByText("test@example.com")).toBeVisible();
    await expect(
      page.getByRole("button", { name: "ログアウト" }),
    ).toBeVisible();
  });

  test("ドロップダウン外をクリックするとドロップダウンが閉じる", async ({
    page,
  }) => {
    await gotoHome(page);
    await page.locator(".bi-person-circle").click();
    await expect(
      page.getByRole("button", { name: "ログアウト" }),
    ).toBeVisible();

    await page.mouse.click(640, 400);
    await expect(
      page.getByRole("button", { name: "ログアウト" }),
    ).not.toBeVisible();
  });

  test("ログアウトボタンをクリックするとログイン画面へ遷移する", async ({
    page,
  }) => {
    await page.route("**/api/auth/logout", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ ok: true }),
      }),
    );

    await gotoHome(page);
    await page.locator(".bi-person-circle").click();
    await page.getByRole("button", { name: "ログアウト" }).click();

    await page.waitForURL(/\/login/);
    await expect(page).toHaveURL(/\/login/);
  });

  test("ハンバーガーボタンをクリックするとサイドバーが開く", async ({
    page,
  }) => {
    await gotoHome(page);
    await expect(page.getByRole("navigation")).toBeVisible();

    await page.locator(".bi-list").click();
    await expect(page.getByRole("heading", { name: "Menu" })).toBeVisible();
    await expect(
      page.getByRole("link", { name: /Home/ }).first(),
    ).toBeVisible();
    await expect(page.getByRole("link", { name: /Watch List/ })).toBeVisible();
  });

  test("サイドバーのオーバーレイをクリックするとサイドバーが閉じる", async ({
    page,
  }) => {
    await gotoHome(page);
    await page.locator(".bi-list").click();
    await expect(page.getByRole("heading", { name: "Menu" })).toBeVisible();

    await page.locator(".bg-opacity-50").click();
    await expect(page.locator(".bg-opacity-50")).not.toBeVisible();
  });

  test("visual snapshot（映画データあり）", async ({ page }) => {
    await gotoHome(page);
    await expect(page.getByText("レビュー済み映画")).toBeVisible();
    await expect(page).toHaveScreenshot("home-authenticated-with-movies.png");
  });
});

test.describe("ホーム画面（ユーザー登録後）", () => {
  const fillRegisterForm = async (page: Page) => {
    await page.getByLabel("ユーザー名").fill("newuser");
    await page.getByLabel("メールアドレス").fill("new@example.com");
    await page.locator("#password").fill("password123");
    await page.locator("#confirmPassword").fill("password123");
    // 登録後の router.push("/") が Middleware を通過するために
    // click() より前に Cookie をセットしておく。
    await setAuthCookie(page);
    await page.getByRole("button", { name: "登録" }).click();
  };

  test.beforeEach(async ({ page }) => {
    await page.route("**/api/auth/register", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ ok: true }),
      }),
    );
    await page.route("**/api/auth/session", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          authenticated: true,
          user: { id: 2, username: "newuser", email: "new@example.com" },
        }),
      }),
    );
    await page.route("**/api/movies**", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: "[]",
      }),
    );
    await page.route("**/api/trending/movies", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: "[]",
      }),
    );
  });

  test("登録直後のセッションでホーム画面が表示される", async ({ page }) => {
    await page.goto("/register");
    await fillRegisterForm(page);

    await page.waitForURL("/");
    await page.waitForLoadState("networkidle");
    await expect(page.getByRole("heading", { name: "Home" })).toBeVisible({
      timeout: 15000,
    });
    await expect(page).toHaveURL("/");
  });

  test("登録後のホームは空のレビュー状態を表示する", async ({ page }) => {
    await page.goto("/register");
    await fillRegisterForm(page);

    await page.waitForURL("/");
    await page.waitForLoadState("networkidle");
    await expect(page.getByRole("heading", { name: "Home" })).toBeVisible({
      timeout: 15000,
    });
    await expect(
      page.getByText("まだレビューした映画がありません。"),
    ).toBeVisible();
  });

  test("visual snapshot（登録後ホーム）", async ({ page }) => {
    await page.goto("/register");
    await fillRegisterForm(page);

    await page.waitForURL("/");
    await page.waitForLoadState("networkidle");
    await expect(page.getByRole("heading", { name: "Home" })).toBeVisible({
      timeout: 15000,
    });
    await expect(page).toHaveScreenshot("home-after-register.png");
  });
});
