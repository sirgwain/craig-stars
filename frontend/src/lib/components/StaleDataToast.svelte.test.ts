import { page } from 'vitest/browser';
import { expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { staleData } from '#lib/services/Stale.js';
import StaleDataToast from './StaleDataToast.svelte';

it('prompts for a reload once the server reports stale data', async () => {
	render(StaleDataToast);
	await expect.element(page.getByRole('alert')).not.toBeInTheDocument();
	staleData.set(true);
	await expect
		.element(page.getByRole('alert'))
		.toHaveTextContent('newer changes from another device');
	await expect.element(page.getByRole('button', { name: 'Reload' })).toBeVisible();
});
