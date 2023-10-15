import React, { useEffect, useState } from "react";

import type { Location, LocationListResponse } from "./types";
import LocationEmpty from "./LocationEmpty";
import LocationList from "./LocationList";

export default function LocationPage() {
	const [locations, setLocations] = useState<Location[]>([]);
	const hasLocations = locations.length > 0;

	useEffect(() => {
		const fetchLocations = async () => {
			const response = await fetch("/api/v1/locations");
			const data: LocationListResponse = await response.json();
			setLocations(data.locations);
		};
		fetchLocations();
	}, []);

	return hasLocations ? <LocationList locations={locations} /> : <LocationEmpty />;
}
