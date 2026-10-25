const { test, expect } = require("@playwright/test");

test.describe("WaqfWise Widget", () => {
  test("widget loads and displays campaigns", async ({ page }) => {
    await page.goto(
      "http://localhost:8085/widget?mosque_id=masjid-al-ikhlas&limit=3",
    );

    // Wait for campaigns to load (skeleton disappears)
    await page.waitForSelector(".waqf-card, .waqf-empty", { timeout: 5000 });

    const cards = await page.$$(".waqf-card");
    expect(cards.length).toBeGreaterThan(0);
  });

  test("clicking campaign redirects to donate page", async ({ page }) => {
    await page.goto("http://localhost:8085/widget?mosque_id=masjid-al-ikhlas");
    await page.waitForSelector(".waqf-card", { timeout: 5000 });

    const firstCard = await page.$(".waqf-card");
    const url = await firstCard.getAttribute("data-url");
    expect(url).toContain("/donate/");
  });

  test("theming applies primary color", async ({ page }) => {
    await page.goto(
      "http://localhost:8085/widget?mosque_id=masjid-al-ikhlas&primary_color=FF0000",
    );
    await page.waitForSelector(".waqf-card, .waqf-empty", { timeout: 5000 });

    const bgColor = await page.evaluate(() =>
      getComputedStyle(document.documentElement).getPropertyValue(
        "--waqf-primary",
      ),
    );
    expect(bgColor.trim()).toBe("#FF0000");
  });

  test("empty state shows when no campaigns", async ({ page }) => {
    await page.goto(
      "http://localhost:8085/widget?mosque_id=nonexistent-mosque",
    );
    // Should show error or empty state
    const visible = await page.waitForSelector(".waqf-error, .waqf-empty", {
      timeout: 5000,
    });
    expect(visible).toBeTruthy();
  });
});
