import '../../../css/app.css';
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

	// A hold split three ways used to leave a sliver of the empty bar showing at the end.
	it('fills the whole bar when the hold is full', async () => {
		const { container } = render(CargoBar, {
			value: create(CargoSchema, { ironium: 100, boranium: 100, germanium: 100, colonists: 1 }),
			capacity: 301
		});
		const bar = container.querySelector<HTMLElement>('[role="button"]')!;
		bar.style.width = '317px';
		const fill = container.querySelector('[data-type="cargo-fill"]')!;
		const segments = [...fill.children].map((el) => el.getBoundingClientRect());
		// the segments run edge to edge, on one line, to the end of the bar
		const inner = fill.getBoundingClientRect();
		expect(segments[0].left).toBeCloseTo(inner.left, 0);
		expect(segments[3].right).toBeCloseTo(inner.right, 0);
		for (let i = 1; i < segments.length; i++) {
			expect(segments[i].left).toBeCloseTo(segments[i - 1].right, 0);
			expect(segments[i].top).toBeCloseTo(segments[0].top, 0);
		}
		// even the single kT of colonists gets drawn
		expect(segments[3].width).toBeGreaterThan(0);
	});
});
