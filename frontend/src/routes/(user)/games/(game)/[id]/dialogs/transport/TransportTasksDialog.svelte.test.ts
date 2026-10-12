import { create } from '@bufbuild/protobuf';
import type { TransportTasksDialogEvent } from '#lib/services/Events.js';
import { CommandedFleet } from '#lib/types/Fleet.js';
import {
	WaypointSchema,
	WaypointTaskTransportAction as Action,
	WaypointTransportTasksSchema,
	type WaypointTransportTasks
} from '#lib/types/cs-proto.js';
import { expect, it, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import TransportTasksDialog from './TransportTasksDialog.svelte';

function orders(transportTasks?: WaypointTransportTasks): TransportTasksDialogEvent {
	return {
		fleet: new CommandedFleet(),
		waypoint: create(WaypointSchema, { transportTasks }),
		waypointIndex: 0
	};
}

function action(container: HTMLElement, resource: string) {
	return page.elementLocator(
		container.querySelector<HTMLSelectElement>(`[name="action${resource}"]`)!
	);
}

function amount(container: HTMLElement, resource: string) {
	return page.elementLocator(
		container.querySelector<HTMLInputElement>(`[name="amount${resource}"]`)!
	);
}

it('opens with the current orders and preserves every resource when saving unchanged', async () => {
	const transportTasks = create(WaypointTransportTasksSchema, {
		fuel: { action: Action.LOAD_OPTIMAL },
		ironium: { action: Action.LOAD_AMOUNT, amount: 12 },
		boranium: { action: Action.UNLOAD_AMOUNT, amount: 34 },
		germanium: { action: Action.FILL_PERCENT, amount: 56 },
		colonists: { action: Action.LOAD_ALL }
	});
	const props = orders(transportTasks);
	const onOk = vi.fn();
	const view = await render(TransportTasksDialog, {
		props: { show: false, props, onOk, onCancel: vi.fn() }
	});
	await view.rerender({ show: true });

	for (const [resource, task] of Object.entries({
		Fuel: transportTasks.fuel!,
		Ironium: transportTasks.ironium!,
		Boranium: transportTasks.boranium!,
		Germanium: transportTasks.germanium!,
		Colonists: transportTasks.colonists!
	})) {
		await expect.element(action(view.container, resource)).toHaveValue(String(task.action));
		await expect.element(amount(view.container, resource)).toHaveValue(task.amount);
	}
	await page.getByRole('button', { name: 'OK', exact: true }).click();

	expect(onOk).toHaveBeenCalledOnce();
	expect(onOk.mock.calls[0][0]).toEqual({ ...props, transportTasks });
	expect(onOk.mock.calls[0][0].transportTasks).not.toBe(transportTasks);
});

it('discards cancelled edits and starts each opening with the selected waypoint orders', async () => {
	const props = orders(
		create(WaypointTransportTasksSchema, {
			colonists: { action: Action.LOAD_AMOUNT, amount: 10 }
		})
	);
	const onOk = vi.fn();
	const onCancel = vi.fn();
	const view = await render(TransportTasksDialog, { props: { show: true, props, onOk, onCancel } });

	await action(view.container, 'Colonists').selectOptions(String(Action.UNLOAD_AMOUNT));
	await amount(view.container, 'Colonists').fill('25');
	await expect
		.element(page.getByText('Unloads up to 2,500 colonists to the target', { exact: true }))
		.toBeVisible();
	await page.getByRole('button', { name: 'Cancel', exact: true }).click();
	expect(onCancel).toHaveBeenCalledOnce();
	expect(onOk).not.toHaveBeenCalled();
	expect(props.waypoint.transportTasks?.colonists).toMatchObject({
		action: Action.LOAD_AMOUNT,
		amount: 10
	});

	await view.rerender({ show: false });
	await view.rerender({ show: true });
	await expect.element(action(view.container, 'Colonists')).toHaveValue(String(Action.LOAD_AMOUNT));
	await expect.element(amount(view.container, 'Colonists')).toHaveValue(10);

	await view.rerender({ show: false });
	props.waypoint = create(WaypointSchema, {
		transportTasks: { colonists: { action: Action.UNLOAD_AMOUNT, amount: 42 } }
	});
	await view.rerender({ show: true });
	await expect
		.element(action(view.container, 'Colonists'))
		.toHaveValue(String(Action.UNLOAD_AMOUNT));
	await expect.element(amount(view.container, 'Colonists')).toHaveValue(42);
});

it('creates missing resource orders when editing a sparse waypoint', async () => {
	const props = orders(create(WaypointTransportTasksSchema));
	const onOk = vi.fn();
	const view = await render(TransportTasksDialog, {
		props: { show: true, props, onOk, onCancel: vi.fn() }
	});

	await expect.element(action(view.container, 'Colonists')).toHaveValue(String(Action.UNSPECIFIED));
	await expect.element(amount(view.container, 'Colonists')).toBeDisabled();
	await action(view.container, 'Colonists').selectOptions(String(Action.LOAD_AMOUNT));
	await amount(view.container, 'Colonists').fill('15');
	await page.getByRole('button', { name: 'OK', exact: true }).click();

	expect(onOk.mock.calls[0][0].transportTasks.colonists).toMatchObject({
		action: Action.LOAD_AMOUNT,
		amount: 15
	});
	expect(props.waypoint.transportTasks?.colonists).toBeUndefined();
});
