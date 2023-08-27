import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";

type Location = {
	id: string;
	kind: string;
};

type LocationListReponse = {
	locations: Location[];
};

export default function LocationList() {
	const [locations, setLocations] = useState<Location[]>([]);

	useEffect(() => {
		const fetchLocations = async () => {
			const response = await fetch("/api/v1/locations");
			const data: LocationListReponse = await response.json();
			setLocations(data.locations);
		};
		fetchLocations();
	}, []);

	return (
		<>
			{locations && (
				<ul>
					{locations.map((location) => (
						<li key={location.id}>
							<Link to={`/locations/${location.id}`}>{location.id}</Link> {location.kind}
						</li>
					))}
				</ul>
			)}
		</>
	);
}
