import { flushSync, mount, unmount } from 'svelte';
import { expect, it } from 'vitest';
import FleetFilterMenu from './FleetFilterMenu.svelte';

it('chooses what a filter shows and turns the filter on', () => {
	const props = $state({
		title: 'Design Filter Options',
		items: [
			{ value: 1, label: 'Long Range Scout' },
			{ value: 2, label: 'Santa Maria' },
			{ value: 3, label: 'Teamster' }
		],
		selected: [] as number[],
		enabled: false
	});
	const target = document.createElement('div');
	document.body.appendChild(target);
	const menu = mount(FleetFilterMenu<number>, {
		target,
		props: {
			get title() {
				return props.title;
			},
			get items() {
				return props.items;
			},
			get selected() {
				return props.selected;
			},
			set selected(value) {
				props.selected = value;
			},
			get enabled() {
				return props.enabled;
			},
			set enabled(value) {
				props.enabled = value;
			}
		}
	});
	const click = (label: string) => {
		[...target.querySelectorAll('button')].find((b) => b.textContent?.includes(label))!.click();
		flushSync();
	};

	// clearing an empty filter leaves it off
	click('No Designs');
	expect(props.enabled).toBe(false);

	click('Santa Maria');
	expect(props.selected).toEqual([2]);
	expect(props.enabled).toBe(true);

	click('Invert Filter');
	expect(props.selected).toEqual([1, 3]);

	// unchecking doesn't turn the filter on
	props.enabled = false;
	click('Teamster');
	expect(props.selected).toEqual([1]);
	expect(props.enabled).toBe(false);

	click('All Designs');
	expect(props.selected).toEqual([1, 2, 3]);
	expect(props.enabled).toBe(true);

	click('No Designs');
	expect(props.selected).toEqual([]);
	expect(props.enabled).toBe(true);

	unmount(menu);
	target.remove();
});
