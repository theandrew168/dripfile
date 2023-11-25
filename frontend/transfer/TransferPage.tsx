import React from "react";
import { useQuery } from "@tanstack/react-query";

import type { TransferListResponse } from "../types";
import TransferEmpty from "./TransferEmpty";
import TransferList from "./TransferList";

export default function TransferPage() {
	const { isPending, isError, error, data } = useQuery({
		queryKey: ["transfer"],
		queryFn: async () => {
			const response = await fetch("/api/v1/transfer");
			if (!response.ok) {
				throw new Error("Network response was not OK");
			}

			const data: TransferListResponse = await response.json();
			return data;
		},
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	const transfers = data.transfers;
	const hasTransfers = transfers.length > 0;
	return hasTransfers ? <TransferList transfers={transfers} /> : <TransferEmpty />;
}
