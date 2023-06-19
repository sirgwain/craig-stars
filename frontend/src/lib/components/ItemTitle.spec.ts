import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import ItemTitle from './ItemTitle.svelte';

describe('ItemTitle', () => {
	it('Should show an h1 element', () => {
		render(ItemTitle, { props: {} });
		expect(screen.getByRole('heading', { level: 3 })).toHaveClass("text-2xl");
	});
});
