import { expect, Page, test } from "@playwright/test";

/** Middleware を通過させるための Cookie をセット */
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

test.describe("ログイン画面", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/login");
  });

  test("ページタイトルと各フォーム要素が表示される", async ({ page }) => {
    await expect(page.getByRole("heading", { name: "ログイン" })).toBeVisible();
    await expect(page.getByLabel("メールアドレス")).toBeVisible();
    // PasswordInput は toggle ボタンも aria-label を持つため id で特定
    await expect(page.locator("#password")).toBeVisible();
    await expect(page.getByRole("button", { name: "ログイン" })).toBeVisible();
    await expect(
      page.getByRole("link", { name: "アカウントをお持ちでない方はこちら" }),
    ).toBeVisible();
  });

  test("登録画面へのリンクが /register に遷移する", async ({ page }) => {
    await page
      .getByRole("link", { name: "アカウントをお持ちでない方はこちら" })
      .click();
    await expect(page).toHaveURL(/\/register/);
  });

  test("メールアドレス・パスワードが空の場合はsubmitできない", async ({
    page,
  }) => {
    await page.getByRole("button", { name: "ログイン" }).click();
    // HTML5 バリデーションにより遷移しない
    await expect(page).toHaveURL(/\/login/);
  });

  test("認証失敗時にエラーメッセージが表示される", async ({ page }) => {
    await page.route("**/api/auth/login", (route) =>
      route.fulfill({
        status: 401,
        contentType: "application/json",
        body: JSON.stringify({
          message: "メールアドレスまたはパスワードが間違っています",
        }),
      }),
    );

    await page.getByLabel("メールアドレス").fill("wrong@example.com");
    await page.locator("#password").fill("wrongpassword");
    await page.getByRole("button", { name: "ログイン" }).click();

    await expect(
      page.getByText("メールアドレスまたはパスワードが間違っています"),
    ).toBeVisible();
    await expect(page).toHaveURL(/\/login/);
  });

  test("ログイン成功後にホーム画面へ遷移する", async ({ page }) => {
    await page.route("**/api/auth/login", (route) =>
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
          user: { id: 1, username: "testuser", email: "test@example.com" },
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

    await page.getByLabel("メールアドレス").fill("test@example.com");
    await page.locator("#password").fill("password123");
    // click() より前に Cookie をセットしておく。
    // ログイン後の router.push が Middleware を通過するために必要。
    await setAuthCookie(page);
    await page.getByRole("button", { name: "ログイン" }).click();

    await page.waitForURL("/");
    await expect(page).toHaveURL("/");
  });

  test("ログイン中はボタンが無効化されローディング表示になる", async ({
    page,
  }) => {
    await page.route("**/api/auth/login", async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 500));
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ ok: true }),
      });
    });
    await page.route("**/api/auth/session", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          authenticated: true,
          user: { id: 1, username: "testuser", email: "test@example.com" },
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

    await page.getByLabel("メールアドレス").fill("test@example.com");
    await page.locator("#password").fill("password123");
    await page.getByRole("button", { name: "ログイン" }).click();

    await expect(
      page.getByRole("button", { name: "ログイン中..." }),
    ).toBeDisabled();
  });

  test("redirect パラメータ付きでログイン成功すると元のページに戻る", async ({
    page,
  }) => {
    await page.route("**/api/auth/login", (route) =>
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
          user: { id: 1, username: "testuser", email: "test@example.com" },
        }),
      }),
    );
    await page.route("**/api/watchlist**", (route) =>
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: "[]",
      }),
    );

    await page.goto("/login?redirect=%2Fwatchlist");
    await page.getByLabel("メールアドレス").fill("test@example.com");
    await page.locator("#password").fill("password123");
    // click() より前に Cookie をセットしておく。
    await setAuthCookie(page);
    await page.getByRole("button", { name: "ログイン" }).click();

    await page.waitForURL("/watchlist");
    await expect(page).toHaveURL("/watchlist");
  });

  test("ネットワークエラー時にフォールバックエラーメッセージが表示される", async ({
    page,
  }) => {
    await page.route("**/api/auth/login", (route) => route.abort("failed"));

    await page.getByLabel("メールアドレス").fill("test@example.com");
    await page.locator("#password").fill("password123");
    await page.getByRole("button", { name: "ログイン" }).click();

    await expect(page.getByText("ログインに失敗しました")).toBeVisible();
    await expect(page).toHaveURL(/\/login/);
  });

  test("不正なメールアドレス形式はHTML5バリデーションで弾かれる", async ({
    page,
  }) => {
    await page.getByLabel("メールアドレス").fill("not-an-email");
    await page.locator("#password").fill("password123");
    await page.getByRole("button", { name: "ログイン" }).click();
    // HTML5 バリデーションにより遷移しない
    await expect(page).toHaveURL(/\/login/);
  });

  test("visual snapshot", async ({ page }) => {
    await expect(page).toHaveScreenshot("login.png");
  });
});
