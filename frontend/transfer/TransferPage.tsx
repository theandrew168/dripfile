import React from "react";
import { useQuery } from "@tanstack/react-query";

import TransferEmpty from "./TransferEmpty";
import TransferList from "./TransferList";
import { listTransfers } from "../fetch";

export default function TransferPage() {
	const {
		isPending,
		isError,
		error,
		data: transfers,
	} = useQuery({
		queryKey: ["transfer"],
		queryFn: async () => listTransfers(),
	});

	// TODO: build a generic loading component
	if (isPending) {
		return <div>Loading...</div>;
	}

	// TODO: build a generic error component
	if (isError) {
		return <div>Error: {error.message}</div>;
	}

	const hasTransfers = transfers.length > 0;
	return hasTransfers ? <TransferList transfers={transfers} /> : <TransferEmpty />;
}
