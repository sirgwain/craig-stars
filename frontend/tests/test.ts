// NOTE: jest-dom adds handy assertions to Jest and it is recommended, but not required.
import '@testing-library/jest-dom';

import { render, screen } from '@testing-library/svelte';

import ItemTitle from '$lib/components/ItemTitle.svelte';

test('shows proper heading when rendered', () => {
	render(ItemTitle, {});
	const heading = screen.getByText('');
	expect(heading).toBeInTheDocument();
});

