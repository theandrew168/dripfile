import React, { useEffect, useState } from "react";
import { useParams } from "react-router";

type Location = {
	id: string;
	kind: string;
};

type LocationReadReponse = {
	location: Location;
};

export default function LocationRead() {
	const { id } = useParams();
	const [location, setLocation] = useState<Location | null>(null);

	useEffect(() => {
		const fetchLocation = async () => {
			const response = await fetch(`/api/v1/locations/${id}`);
			const data: LocationReadReponse = await response.json();
			setLocation(data.location);
		};
		fetchLocation();
	}, []);

	return (
		location && (
			<div>
				<p>ID: {location.id}</p>
				<p>Kind: {location.kind}</p>
			</div>
		)
	);
}
