import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import CargoBar from './CargoBar.svelte';

describe('CargoBar', () => {
	it('Should render a cargo bar with 30kT ironium', () => {
		render(CargoBar, { props: { value: { ironium: 30 }, capacity: 50 } });
		screen.debug
		expect(screen.getByText('30 of 50kT')).toBeInTheDocument();
	});
});
