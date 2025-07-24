import { page } from '@vitest/browser/context';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import ItemTitle from './ItemTitle.svelte';

describe('ItemTitle', () => {
	it('Should show an h1 element', async () => {
		render(ItemTitle);

		const heading = page.getByRole('heading', { level: 3 });
		await expect.element(heading).toBeInTheDocument();
	});
});
