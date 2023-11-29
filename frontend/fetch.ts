import {
	isErrorResponse,
	type Itinerary,
	type ItineraryListResponse,
	type ItineraryReadResponse,
	type Location,
	type LocationListResponse,
	type LocationReadResponse,
	type NewInMemoryLocation,
	type NewS3Location,
	type Transfer,
	type TransferListResponse,
	type TransferReadResponse,
} from "./types";

export async function createInMemoryLocation(_location: NewInMemoryLocation): Promise<void> {
	const response = await fetch("/api/v1/location", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ kind: "memory" }),
	});
	if (!response.ok) {
		const error = await response.json();
		if (isErrorResponse(error)) {
			throw new Error(JSON.stringify(error.error));
		} else {
			throw new Error("Network response was not OK");
		}
	}

	return response.json();
}

export async function createS3Location(location: NewS3Location): Promise<void> {
	const response = await fetch("/api/v1/location", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(location),
	});
	if (!response.ok) {
		const error = await response.json();
		if (isErrorResponse(error)) {
			throw new Error(JSON.stringify(error.error));
		} else {
			throw new Error("Network response was not OK");
		}
	}

	return response.json();
}

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
