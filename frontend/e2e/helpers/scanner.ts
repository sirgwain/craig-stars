import type { Locator, Page } from '@playwright/test';
import type { MapObjectLike } from '../../src/lib/types/MapObject';

type ClickOptions = Parameters<Locator['click']>[0];

/**
 * The scanner is drawn on a canvas, so map objects aren't elements we can click. Instead, read
 * the scanner's current view (pixels per light year and zoom transform) and click the overlay
 * at the map object's screen position.
 */
async function mapObjectPosition(page: Page, mo: MapObjectLike | undefined) {
	if (!mo?.mapObject) {
		throw new Error('map object not found');
	}
	const canvas = page.locator('[data-type="scanner-canvas"]');
	await canvas.and(page.locator('[data-view]')).waitFor();
	// the view is updated when the canvas draws, so let any pending zoom/pan draw first
	await page.evaluate(
		() => new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)))
	);
	const view = await canvas.getAttribute('data-view');
	const [sx, sy, k, tx, ty] = (view ?? '').split(',').map(Number);
	return {
		x: Number(mo.mapObject.position?.x ?? 0) * sx * k + tx,
		y: Number(mo.mapObject.position?.y ?? 0) * sy * k + ty
	};
}

export function scannerOverlay(page: Page): Locator {
	return page.locator('[data-type="scanner-overlay"]');
}

export async function clickMapObject(
	page: Page,
	mo: MapObjectLike | undefined,
	options: ClickOptions = {}
): Promise<void> {
	const position = await mapObjectPosition(page, mo);
	await scannerOverlay(page).click({ ...options, position, force: true });
}

export async function dblclickMapObject(
	page: Page,
	mo: MapObjectLike | undefined,
	options: Parameters<Locator['dblclick']>[0] = {}
): Promise<void> {
	const position = await mapObjectPosition(page, mo);
	await scannerOverlay(page).dblclick({ ...options, position, force: true });
}

// tap a map object, like a finger on a phone. The page must have touch enabled.
export async function tapMapObject(
	page: Page,
	mo: MapObjectLike | undefined,
	options: Parameters<Locator['tap']>[0] = {}
): Promise<void> {
	const position = await mapObjectPosition(page, mo);
	await scannerOverlay(page).tap({ ...options, position, force: true });
}
