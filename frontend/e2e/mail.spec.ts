import type { Page } from '@playwright/test';
import { expect, submitTurn, test } from './setup';

function mailDialog(page: Page) {
	return page.locator('[data-type="mail-dialog"]');
}

test('writes mail that is delivered with the next turn', async ({ newGamePage }) => {
	const { page, id } = newGamePage;
	const items = mailDialog(page).locator('[data-type="mail-item"]');

	await test.step('write mail', async () => {
		// mail is written from the messages pane
		await page.locator('a[href="#messages"]:visible').click();
		await page.locator('[data-type="write-mail"]:visible').click();
		// mail to yourself is delivered like any other, as a note
		await mailDialog(page).locator('[data-type="mail-to"]').selectOption('1');
		await mailDialog(page).locator('textarea').fill('remember the trcue');
		await mailDialog(page).getByRole('button', { name: 'Send' }).click();

		await expect(items).toHaveCount(1);
		await expect(items).toContainText('Sends next turn');
		await expect(mailDialog(page).locator('textarea')).toHaveValue('');
	});

	await test.step('change it before it is delivered', async () => {
		await items.getByRole('button', { name: 'Edit' }).click();
		await expect(mailDialog(page).locator('textarea')).toHaveValue('remember the trcue');
		await mailDialog(page).locator('textarea').fill('remember the truce');
		await mailDialog(page).getByRole('button', { name: 'Save' }).click();

		await expect(items).toHaveCount(1);
		await expect(items).toContainText('remember the truce');
		await mailDialog(page).getByRole('button', { name: 'Close' }).click();
		await expect(mailDialog(page).locator('textarea')).toBeHidden();
	});

	await test.step('the turn delivers it as a message', async () => {
		await submitTurn(page);
		await page.goto(`/games/${id}/messages`);
		await expect(page.getByRole('main')).toContainText('remember the truce');
		await expect(page.getByRole('main')).toContainText('Note:');
	});

	await test.step('reply to it', async () => {
		await page.getByRole('button', { name: 'Reply' }).click();
		await expect(mailDialog(page).locator('[data-type="mail-to"]')).toHaveValue('1');
		// delivered mail can't be changed
		await expect(items).toHaveCount(1);
		await expect(items).not.toContainText('Sends next turn');
		await expect(items.getByRole('button', { name: 'Edit' })).toBeHidden();
	});
});
