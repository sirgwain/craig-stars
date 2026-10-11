import { TechService } from '#lib/services/TechService.js';
import { describe, expect, it } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import TechDescription from './TechDescription.svelte';

const techs = new TechService();

describe('TechDescription', () => {
	it('describes the jump gate', async () => {
		render(TechDescription, { tech: techs.getHullComponent('Jump Gate')! });

		await expect
			.element(page.getByText('This device allows a ship to jump to planetary stargates.'))
			.toBeInTheDocument();
	});

	it.each([
		['Enigma Pulsar', '1/8 square'],
		['Maneuvering Jet', '1/4 square'],
		['Overthruster', '1/2 square']
	])('%s increases battle speed by %s', async (name, bonus) => {
		render(TechDescription, { tech: techs.getHullComponent(name)! });

		await expect
			.element(page.getByText(`Increases speed in battle by ${bonus} of movement.`))
			.toBeInTheDocument();
	});

	it('renders nothing for a tech with no description', async () => {
		const { container } = render(TechDescription, {
			tech: techs.getHullComponent('Long Hump 6')!
		});

		expect(container.textContent?.trim()).toBe('');
	});
});
