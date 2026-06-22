import type { Position } from '$lib/types/MapObject';
import ScannerContextPopup from './ScannerContextPopup.svelte';
import { showPopup, type PopupPropsBase } from './Popup';

export type ScannerContextPopupProps = {
	position: Position;
} & PopupPropsBase;

export function onScannerContextPopup(e: PointerEvent | MouseEvent, position: Position) {
	showPopup<ScannerContextPopupProps>(e.x, e.y, ScannerContextPopup, { position });
}
