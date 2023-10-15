import React from "react";
import { ChevronRightIcon } from "@heroicons/react/20/solid";
import { CircleStackIcon, CpuChipIcon } from "@heroicons/react/24/outline";

import { classNames } from "../utils";
import { Link } from "react-router-dom";

const items = [
	{
		name: "In-Memory",
		description: "An in-memory location for testing Dripfile.",
		href: "/locations/create/in-memory",
		iconColor: "bg-pink-500",
		icon: CpuChipIcon,
	},
	{
		name: "S3 Bucket",
		description: "An Amazon S3 (or compatible) object storage bucket.",
		href: "/locations/create/s3",
		iconColor: "bg-purple-500",
		icon: CircleStackIcon,
	},
];

export default function LocationCreate() {
	return (
		<div className="mt-24 mx-auto max-w-lg">
			<h2 className="text-base font-semibold leading-6 text-gray-900">Add your first location</h2>
			<p className="mt-1 text-sm text-gray-500">A location is a place where your data lives.</p>
			<ul role="list" className="mt-6 divide-y divide-gray-200 border-b border-t border-gray-200">
				{items.map((item, itemIdx) => (
					<li key={itemIdx}>
						<div className="group relative flex items-start space-x-3 py-4">
							<div className="flex-shrink-0">
								<span
									className={classNames(item.iconColor, "inline-flex h-10 w-10 items-center justify-center rounded-lg")}
								>
									<item.icon className="h-6 w-6 text-white" aria-hidden="true" />
								</span>
							</div>
							<div className="min-w-0 flex-1">
								<div className="text-sm font-medium text-gray-900">
									<Link to={item.href}>
										<span className="absolute inset-0" aria-hidden="true" />
										{item.name}
									</Link>
								</div>
								<p className="text-sm text-gray-500">{item.description}</p>
							</div>
							<div className="flex-shrink-0 self-center">
								<ChevronRightIcon className="h-5 w-5 text-gray-400 group-hover:text-gray-500" aria-hidden="true" />
							</div>
						</div>
					</li>
				))}
			</ul>
		</div>
	);
}
