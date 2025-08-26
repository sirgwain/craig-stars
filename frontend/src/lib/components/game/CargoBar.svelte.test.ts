import { page } from '@vitest/browser/context';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import CargoBar from './CargoBar.svelte';
import { CargoSchema } from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';

describe('CargoBar', () => {
	it('Should render a cargo bar with 30kT ironium', async () => {
		render(CargoBar, { value: create(CargoSchema, { ironium: 30 }), capacity: 50 });
		const item = page.getByText('30 of 50kT');
		await expect.element(item).toBeInTheDocument();
	});
});
