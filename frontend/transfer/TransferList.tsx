import React from "react";
import { Link } from "react-router-dom";

import type { Transfer } from "../types";

type Props = {
	transfers: Transfer[];
};

export default function TransferList({ transfers }: Props) {
	// https://tailwindui.com/components/application-ui/lists/tables#component-4738eac883e67bf84a9f7db2446e838a
	return (
		<>
			<div className="sm:flex sm:items-center">
				<div className="sm:flex-auto">
					<h1 className="text-base font-semibold leading-6 text-gray-900">Transfers</h1>
					<p className="mt-2 text-sm text-gray-700">A list of all recent transfers.</p>
				</div>
			</div>
			<div className="mt-8 flow-root">
				<div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
					<div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
						<div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
							<table className="min-w-full divide-y divide-gray-300">
								<thead className="bg-gray-50">
									<tr>
										<th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">
											ID
										</th>
										<th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
											ItineraryID
										</th>
										<th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
											Status
										</th>
										<th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
											Progress
										</th>
										<th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
											Error
										</th>
										<th scope="col" className="relative py-3.5 pl-3 pr-4 sm:pr-6">
											<span className="sr-only">Edit</span>
										</th>
									</tr>
								</thead>
								<tbody className="divide-y divide-gray-200 bg-white">
									{transfers.map((transfer) => (
										<tr key={transfer.id}>
											<td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">
												<Link to={`/transfer/${transfer.id}`}>{transfer.id}</Link>
											</td>
											<td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
												<Link to={`/itinerary/${transfer.itineraryID}`}>{transfer.itineraryID}</Link>
											</td>
											<td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{transfer.status}</td>
											<td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{transfer.progress}</td>
											<td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{transfer.error}</td>
											<td className="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
												<a href="#" className="text-indigo-600 hover:text-indigo-900">
													Edit<span className="sr-only">, {transfer.id}</span>
												</a>
											</td>
										</tr>
									))}
								</tbody>
							</table>
						</div>
					</div>
				</div>
			</div>
		</>
	);
}
