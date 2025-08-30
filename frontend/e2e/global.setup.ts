/// <reference types="node" />
import { chromium, FullConfig } from '@playwright/test';
import path from 'path';
import { fileURLToPath } from 'url';

async function globalSetup(config: FullConfig) {
	const { baseURL } = config.projects[0].use;
	const __filename = fileURLToPath(import.meta.url);
	const __dirname = path.dirname(__filename);
	const authFile = path.join(__dirname, '../.auth/user.json');

	// Launch browser
	const browser = await chromium.launch();
	const page = await browser.newPage();

	// Navigate to login page
	await page.goto(baseURL!);
	await page.getByRole('button', { name: "I'm an admin" }).click();
	await page.getByRole('textbox', { name: 'Username' }).click();
	await page.getByRole('textbox', { name: 'Username' }).fill('admin');
	await page.getByRole('textbox', { name: 'Username' }).press('Tab');
	await page.getByRole('textbox', { name: 'Password' }).fill('admin');
	await page.getByRole('button', { name: 'Submit' }).click();

	// Wait for the homepage to be visible
	await page.getByRole('link', { name: 'Single Player' }).waitFor();

	// Save signed-in state to 'authFile'
	await page.context().storageState({ path: authFile });
	await browser.close();
}

export default globalSetup;