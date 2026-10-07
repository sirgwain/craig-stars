import { page, userEvent } from 'vitest/browser';
import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import CargoBar from './CargoBar.svelte';
import { CargoSchema } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';

describe('CargoBar', () => {
	it('Should render a cargo bar with 30kT ironium', async () => {
		render(CargoBar, { value: create(CargoSchema, { ironium: 30 }), capacity: 50 });
		const item = page.getByText('30 of 50kT');
		await expect.element(item).toBeInTheDocument();
	});

	it('opens cargo transfer from the keyboard at the focused bar position', async () => {
		const onPointerDown = vi.fn<(event: PointerEvent) => void>();
		render(CargoBar, {
			value: create(CargoSchema, { ironium: 30 }),
			capacity: 50,
			canTransferCargo: true,
			onPointerDown
		});
		const bar = page.getByRole('button', { name: 'Transfer cargo' }).element() as HTMLElement;
		bar.focus();
		await userEvent.keyboard('{Enter}');
		expect(onPointerDown).toHaveBeenCalledOnce();
		const event = onPointerDown.mock.calls[0][0];
		const rect = bar.getBoundingClientRect();
		expect(event.clientX).toBeCloseTo(rect.left + rect.width / 2);
		expect(event.clientY).toBeCloseTo(rect.top + rect.height / 2);
		const repeatedKey = new KeyboardEvent('keydown', {
			key: 'Enter',
			repeat: true,
			bubbles: true,
			cancelable: true
		});
		bar.dispatchEvent(repeatedKey);
		expect(repeatedKey.defaultPrevented).toBe(true);
		expect(onPointerDown).toHaveBeenCalledOnce();
	});
});
