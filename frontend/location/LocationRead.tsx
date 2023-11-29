import React from "react";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";

import { readLocation } from "../fetch";

export default function LocationRead() {
	const { id } = useParams();
	if (!id) {
		return null;
	}

	const {
		isPending,
		isError,
		error,
		data: location,
	} = useQuery({
		queryKey: ["location", id],
		queryFn: async () => readLocation(id),
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	return (
		<div>
			<p>ID: {location.id}</p>
			<p>Kind: {location.kind}</p>
			<p>CreatedAt: {location.createdAt.toString()}</p>
			<p>UpdatedAt: {location.updatedAt.toString()}</p>
		</div>
	);
}
