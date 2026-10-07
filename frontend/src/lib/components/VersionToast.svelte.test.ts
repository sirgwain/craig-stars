import { page } from 'vitest/browser';
import { expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import VersionToast from './VersionToast.svelte';

vi.mock('#lib/services/Version.js', async (importOriginal) => {
	const original = await importOriginal<typeof import('#lib/services/Version.js')>();
	return { ...original, versionUpdate: original.createVersionStore('v1.42.0') };
});

it('shows the prompt on a frontend update and hides it on Later', async () => {
	render(VersionToast, { frontendUpdated: true });
	await expect.element(page.getByRole('status')).toHaveTextContent('A new version of CraigStars!');
	await page.getByRole('button', { name: 'Later' }).click();
	await expect.element(page.getByRole('status')).not.toBeInTheDocument();
});
