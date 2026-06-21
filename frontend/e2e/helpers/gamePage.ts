import { expect, type Locator, type Page } from '@playwright/test';

export type CargoName = 'fuel' | 'ironium' | 'boranium' | 'germanium' | 'colonists';
export type CargoAmounts = Partial<Record<CargoName, number>>;
export type CargoText = Partial<Record<CargoName, string>>;
export type TransferDirection = 'source' | 'dest';

const cargoLabels: Record<CargoName, string> = {
	fuel: 'Fuel',
	ironium: 'Ironium',
	boranium: 'Boranium',
	germanium: 'Germanium',
	colonists: 'Colonists'
};

export class GamePage {
	constructor(readonly page: Page) {}

	tile(name: string): Locator {
		return this.page.locator(`[data-type="command-tile"][data-id="${name}"]`).first();
	}

	summary(): Locator {
		return this.page.locator('[data-type="map-object-summary"]').first();
	}

	transferButton(cargo: CargoName, direction: TransferDirection): Locator {
		return this.page.locator(`[data-id="${cargo}"][data-type="transfer-to-${direction}-button"]`);
	}

	async clickTileButton(tileName: string, buttonName: string | RegExp): Promise<void> {
		await this.tile(tileName).getByRole('button', { name: buttonName }).first().click();
	}

	async openTransferFromTile(tileName: string): Promise<void> {
		await this.clickTileButton(tileName, /Jettison|Transfer/);
	}

	async transferCargo(amounts: CargoAmounts, direction: TransferDirection): Promise<void> {
		for (const [cargo, amount] of Object.entries(amounts) as [CargoName, number][]) {
			if (!amount) {
				continue;
			}
			await this.transferButton(cargo, direction).click({ clickCount: amount });
		}
	}

	async confirmDialog(buttonName = 'Ok'): Promise<void> {
		await this.page.getByRole('button', { name: buttonName }).click();
	}

	async expectTileCargo(tileName: string, cargo: CargoText): Promise<void> {
		const tile = this.tile(tileName);
		for (const [cargoName, value] of Object.entries(cargo) as [CargoName, string][]) {
			await expect(tile.getByText(`${cargoLabels[cargoName]} ${value}`).first()).toBeVisible();
		}
	}

	async clickNextSelectedObject(): Promise<void> {
		await this.page.getByRole('button', { name: 'Next', exact: true }).first().click();
	}

	async cycleSummarySelection(): Promise<void> {
		await this.summary().locator('[data-type="cycle-selected-map-object-button"]').first().click();
	}

	async expectSummarySelection(name: RegExp): Promise<void> {
		await expect(this.summary().locator('div').filter({ hasText: name }).first()).toBeVisible();
	}

	async scrollSummaryText(text: string): Promise<void> {
		await this.summary().getByText(text).first().scrollIntoViewIfNeeded();
	}

	async expectSummaryCargo(cargo: CargoText): Promise<void> {
		const summary = this.summary();
		for (const [cargoName, value] of Object.entries(cargo) as [CargoName, string][]) {
			await expect(summary.getByText(`${cargoLabels[cargoName]} ${value}`).first()).toBeVisible();
		}
	}
}
