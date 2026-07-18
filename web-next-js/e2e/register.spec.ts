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

test.describe("ユーザー登録画面", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/register");
  });

  test("ページタイトルと各フォーム要素が表示される", async ({ page }) => {
    await expect(
      page.getByRole("heading", { name: "ユーザー登録" }),
    ).toBeVisible();
    await expect(page.getByLabel("ユーザー名")).toBeVisible();
    await expect(page.getByLabel("メールアドレス")).toBeVisible();
    // PasswordInput は toggle ボタンも aria-label を持つため id で特定
    await expect(page.locator("#password")).toBeVisible();
    await expect(page.locator("#confirmPassword")).toBeVisible();
    await expect(page.getByRole("button", { name: "登録" })).toBeVisible();
    await expect(
      page.getByRole("link", { name: "既にアカウントをお持ちの方はこちら" }),
    ).toBeVisible();
  });

  test("ログイン画面へのリンクが /login に遷移する", async ({ page }) => {
    await page
      .getByRole("link", { name: "既にアカウントをお持ちの方はこちら" })
      .click();
    await expect(page).toHaveURL(/\/login/);
  });

  test("必須項目が空の場合はsubmitできない", async ({ page }) => {
    await page.getByRole("button", { name: "登録" }).click();
    await expect(page).toHaveURL(/\/register/);
  });

  test("パスワードと確認用パスワードが一致しない場合はエラーが表示される", async ({
    page,
  }) => {
    await page.getByLabel("ユーザー名").fill("testuser");
    await page.getByLabel("メールアドレス").fill("test@example.com");
    await page.locator("#password").fill("password123");
    await page.locator("#confirmPassword").fill("differentpassword");
    await page.getByRole("button", { name: "登録" }).click();

    await expect(page.getByText("パスワードが一致しません")).toBeVisible();
    await expect(page).toHaveURL(/\/register/);
  });

  test("登録失敗時にエラーメッセージが表示される", async ({ page }) => {
    await page.route("**/api/auth/register", (route) =>
      route.fulfill({
        status: 409,
        contentType: "application/json",
        body: JSON.stringify({
          message: "このメールアドレスは既に使用されています",
        }),
      }),
    );

    await page.getByLabel("ユーザー名").fill("testuser");
    await page.getByLabel("メールアドレス").fill("existing@example.com");
    await page.locator("#password").fill("password123");
    await page.locator("#confirmPassword").fill("password123");
    await page.getByRole("button", { name: "登録" }).click();

    await expect(
      page.getByText("このメールアドレスは既に使用されています"),
    ).toBeVisible();
    await expect(page).toHaveURL(/\/register/);
  });

  test("登録成功後にホーム画面へ遷移する", async ({ page }) => {
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
          user: { id: 1, username: "newuser", email: "new@example.com" },
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

    await page.getByLabel("ユーザー名").fill("newuser");
    await page.getByLabel("メールアドレス").fill("new@example.com");
    await page.locator("#password").fill("password123");
    await page.locator("#confirmPassword").fill("password123");
    // click() より前に Cookie をセットしておく。
    await setAuthCookie(page);
    await page.getByRole("button", { name: "登録" }).click();

    await page.waitForURL("/");
    await expect(page).toHaveURL("/");
  });

  test("登録中はボタンが無効化されローディング表示になる", async ({ page }) => {
    await page.route("**/api/auth/register", async (route) => {
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
          user: { id: 1, username: "newuser", email: "new@example.com" },
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

    await page.getByLabel("ユーザー名").fill("newuser");
    await page.getByLabel("メールアドレス").fill("new@example.com");
    await page.locator("#password").fill("password123");
    await page.locator("#confirmPassword").fill("password123");
    await page.getByRole("button", { name: "登録" }).click();

    await expect(
      page.getByRole("button", { name: "登録中..." }),
    ).toBeDisabled();
  });

  test("ネットワークエラー時にフォールバックエラーメッセージが表示される", async ({
    page,
  }) => {
    await page.route("**/api/auth/register", (route) => route.abort("failed"));

    await page.getByLabel("ユーザー名").fill("testuser");
    await page.getByLabel("メールアドレス").fill("test@example.com");
    await page.locator("#password").fill("password123");
    await page.locator("#confirmPassword").fill("password123");
    await page.getByRole("button", { name: "登録" }).click();

    await expect(page.getByText("ユーザー登録に失敗しました")).toBeVisible();
    await expect(page).toHaveURL(/\/register/);
  });

  test("不正なメールアドレス形式はHTML5バリデーションで弾かれる", async ({
    page,
  }) => {
    await page.getByLabel("ユーザー名").fill("testuser");
    await page.getByLabel("メールアドレス").fill("not-an-email");
    await page.locator("#password").fill("password123");
    await page.locator("#confirmPassword").fill("password123");
    await page.getByRole("button", { name: "登録" }).click();
    // HTML5 バリデーションにより遷移しない
    await expect(page).toHaveURL(/\/register/);
  });

  test("visual snapshot", async ({ page }) => {
    await expect(page).toHaveScreenshot("register.png");
  });
});
