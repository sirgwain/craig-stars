import './app.css';
import { afterEach, describe, expect, it } from 'vitest';
import { page } from 'vitest/browser';

// the classes the game view uses to switch between its desktop and phone layouts
function build() {
	const host = document.createElement('div');
	host.innerHTML = `
		<div id="row" class="flex flex-col desk:flex-row phone-landscape:flex-row">
			<div id="pane" class="hidden desk:flex desk:w-62 desk:lg:w-120 phone-landscape:flex phone-landscape:flex-none phone-landscape:w-62 h-10">
				<div id="paneSummary" class="hidden desk:lg:block"></div>
			</div>
			<div id="scanner" class="grow h-10">
				<div id="summary" class="hidden desk:block desk:lg:hidden"></div>
				<button id="summaryButton" class="hidden phone-landscape:block"></button>
			</div>
			<div id="drawer" class="desk:hidden phone-landscape:hidden"></div>
		</div>`;
	document.body.appendChild(host);
	const el = (id: string) => host.querySelector<HTMLElement>(`#${id}`)!;
	const shown = (id: string) => getComputedStyle(el(id)).display !== 'none';
	return {
		host,
		layout: () => ({
			row: getComputedStyle(el('row')).flexDirection,
			pane: shown('pane'),
			paneWidth: shown('pane') ? el('pane').getBoundingClientRect().width : 0,
			// the summary is in the pane, under the scanner, or a popup behind a button
			summary: shown('paneSummary')
				? 'pane'
				: shown('summary')
					? 'scanner'
					: shown('summaryButton')
						? 'popup'
						: 'drawer',
			drawer: shown('drawer'),
			scannerOnRight:
				shown('pane') &&
				el('pane').getBoundingClientRect().right <= el('scanner').getBoundingClientRect().left
		})
	};
}

describe('game view layout', () => {
	let host: HTMLElement | undefined;
	afterEach(() => host?.remove());

	async function layoutAt(width: number, height: number) {
		await page.viewport(width, height);
		const built = build();
		host = built.host;
		return built.layout();
	}

	it('puts the command pane and summary around the scanner on a tablet or desktop', async () => {
		const wide = {
			row: 'row',
			pane: true,
			paneWidth: 480,
			summary: 'pane',
			drawer: false,
			scannerOnRight: true
		};
		expect(await layoutAt(1280, 800)).toEqual(wide);
		expect(await layoutAt(1133, 744)).toEqual(wide);
		expect(await layoutAt(768, 1024)).toEqual({ ...wide, paneWidth: 248, summary: 'scanner' });
	});

	it('uses the drawer under the scanner on an upright phone', async () => {
		expect(await layoutAt(390, 844)).toEqual({
			row: 'column',
			pane: false,
			paneWidth: 0,
			summary: 'drawer',
			drawer: true,
			scannerOnRight: false
		});
	});

	it('shows the command pane left of the scanner, and pops up the summary, on a sideways phone', async () => {
		const sideways = {
			row: 'row',
			pane: true,
			paneWidth: 248,
			summary: 'popup',
			drawer: false,
			scannerOnRight: true
		};
		// wide enough to count as a tablet by width alone
		expect(await layoutAt(844, 390)).toEqual(sideways);
		expect(await layoutAt(932, 430)).toEqual(sideways);
		expect(await layoutAt(667, 375)).toEqual(sideways);
		// and as a desktop
		expect(await layoutAt(1100, 450)).toEqual(sideways);
	});
});
