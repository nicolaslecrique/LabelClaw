import { expect, test } from '@playwright/test';

test('loads the app and shows both tabs', async ({ page }) => {
  await page.goto('/');

  const tabs = page.getByRole('navigation', { name: 'Main tabs' });

  await expect(tabs.getByRole('button', { name: 'Configuration', exact: true })).toBeVisible();
  await expect(tabs.getByRole('button', { name: 'Labelling', exact: true })).toBeVisible();
});

test('generates a panel, saves it, and restores it after reload', async ({
  page,
}) => {
  await page.goto('/');

  await page.getByRole('button', { name: 'Generate labelling panel' }).click();
  await expect(page.getByRole('heading', { name: 'Example article' })).toBeVisible();

  await page.getByLabel('Sentiment').selectOption('positive');
  await page.getByLabel('Notes').fill('Clear rationale');

  await expect(page.getByText('Output matches the label schema.')).toBeVisible();
  await page.getByRole('button', { name: 'Save labelling panel' }).click();
  await expect(page.getByText('Labelling panel saved.')).toBeVisible();

  await page.reload();
  await expect(page.getByText('Loaded saved configuration.')).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Example article' })).toBeVisible();
});

test('blocks generation when a schema is invalid', async ({ page }) => {
  await page.goto('/');

  await page.getByLabel('Sample JSON Schema').fill('{');
  await page.getByRole('button', { name: 'Generate labelling panel' }).click();

  await expect(
    page.getByText('Both schemas must be valid JSON before generation.'),
  ).toBeVisible();
});

test('keeps the labelling tab empty for now', async ({ page }) => {
  await page.goto('/');
  await page
    .getByRole('navigation', { name: 'Main tabs' })
    .getByRole('button', { name: 'Labelling', exact: true })
    .click();
  await expect(page.getByText('To be defined later')).toBeVisible();
});
