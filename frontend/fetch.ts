import type {
	Itinerary,
	ItineraryListResponse,
	ItineraryReadResponse,
	Location,
	LocationListResponse,
	LocationReadResponse,
	Transfer,
	TransferListResponse,
	TransferReadResponse,
} from "./types";

export async function listLocations(): Promise<Location[]> {
	const response = await fetch("/api/v1/location");
	if (!response.ok) {
		throw new Error("Network response was not OK");
	}

	const data: LocationListResponse = await response.json();
	return data.locations;
}

export async function readLocation(id: string): Promise<Location> {
	const response = await fetch(`/api/v1/location/${id}`);
	if (!response.ok) {
		throw new Error("Network response was not OK");
	}

	const data: LocationReadResponse = await response.json();
	return data.location;
}

export async function listItineraries(): Promise<Itinerary[]> {
	const response = await fetch("/api/v1/itinerary");
	if (!response.ok) {
		throw new Error("Network response was not OK");
	}

	const data: ItineraryListResponse = await response.json();
	return data.itineraries;
}

export async function readItinerary(id: string): Promise<Itinerary> {
	const response = await fetch(`/api/v1/itinerary/${id}`);
	if (!response.ok) {
		throw new Error("Network response was not OK");
	}

	const data: ItineraryReadResponse = await response.json();
	return data.itinerary;
}

export async function listTransfers(): Promise<Transfer[]> {
	const response = await fetch("/api/v1/transfer");
	if (!response.ok) {
		throw new Error("Network response was not OK");
	}

	const data: TransferListResponse = await response.json();
	return data.transfers;
}

export async function readTransfer(id: string): Promise<Transfer> {
	const response = await fetch(`/api/v1/transfer/${id}`);
	if (!response.ok) {
		throw new Error("Network response was not OK");
	}

	const data: TransferReadResponse = await response.json();
	return data.transfer;
}
