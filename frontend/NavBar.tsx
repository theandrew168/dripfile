import React, { Fragment } from "react";
import { Disclosure, Menu, Transition } from "@headlessui/react";
import { Bars3Icon, BellIcon, XMarkIcon } from "@heroicons/react/24/outline";
import { Outlet, useLocation } from "react-router";
import { Link } from "react-router-dom";

import { classNames } from "./utils";

export default function NavBar() {
	const location = useLocation();

	const user = {
		name: "Andrew Dailey",
		email: "andrew@shallowbrooksoftware.com",
		imageUrl: "https://avatars.githubusercontent.com/u/13205759",
	};
	const navigation = [
		{ name: "Dashboard", href: "/", current: location.pathname === "/" },
		{ name: "Locations", href: "/location", current: location.pathname.startsWith("/location") },
		{ name: "Itineraries", href: "/itinerary", current: location.pathname.startsWith("/itinerary") },
		{ name: "Transfers", href: "/transfer", current: location.pathname.startsWith("/transfer") },
	];
	const userNavigation = [
		{ name: "Your Profile", href: "#" },
		{ name: "Settings", href: "#" },
		{ name: "Sign out", href: "#" },
	];

	// https://tailwindui.com/components/application-ui/application-shells/stacked#component-e1d4a342c5903e8ade21c57c20448d69
	return (
		<div className="min-h-full">
			<Disclosure as="nav" className="bg-indigo-600">
				{({ open }) => (
					<>
						<div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
							<div className="flex h-16 items-center justify-between">
								<div className="flex items-center">
									<div className="flex-shrink-0">
										<img className="h-8 w-8" src="/static/img/logo-white.svg" alt="Dripfile" />
									</div>
									<div className="hidden md:block">
										<div className="ml-10 flex items-baseline space-x-4">
											{navigation.map((item) => (
												<Link
													key={item.name}
													to={item.href}
													className={classNames(
														item.current
															? "bg-indigo-700 text-white"
															: "text-white hover:bg-indigo-500 hover:bg-opacity-75",
														"rounded-md px-3 py-2 text-sm font-medium",
													)}
													aria-current={item.current ? "page" : undefined}
												>
													{item.name}
												</Link>
											))}
										</div>
									</div>
								</div>
								<div className="hidden md:block">
									<div className="ml-4 flex items-center md:ml-6">
										<button
											type="button"
											className="relative rounded-full bg-indigo-600 p-1 text-indigo-200 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600"
										>
											<span className="absolute -inset-1.5" />
											<span className="sr-only">View notifications</span>
											<BellIcon className="h-6 w-6" aria-hidden="true" />
										</button>

										{/* Profile dropdown */}
										<Menu as="div" className="relative ml-3">
											<div>
												<Menu.Button className="relative flex max-w-xs items-center rounded-full bg-indigo-600 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600">
													<span className="absolute -inset-1.5" />
													<span className="sr-only">Open user menu</span>
													<img className="h-8 w-8 rounded-full" src={user.imageUrl} alt="" />
												</Menu.Button>
											</div>
											<Transition
												as={Fragment}
												enter="transition ease-out duration-100"
												enterFrom="transform opacity-0 scale-95"
												enterTo="transform opacity-100 scale-100"
												leave="transition ease-in duration-75"
												leaveFrom="transform opacity-100 scale-100"
												leaveTo="transform opacity-0 scale-95"
											>
												<Menu.Items className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
													{userNavigation.map((item) => (
														<Menu.Item key={item.name}>
															{({ active }) => (
																<Link
																	to={item.href}
																	className={classNames(
																		active ? "bg-gray-100" : "",
																		"block px-4 py-2 text-sm text-gray-700",
																	)}
																>
																	{item.name}
																</Link>
															)}
														</Menu.Item>
													))}
												</Menu.Items>
											</Transition>
										</Menu>
									</div>
								</div>
								<div className="-mr-2 flex md:hidden">
									{/* Mobile menu button */}
									<Disclosure.Button className="relative inline-flex items-center justify-center rounded-md bg-indigo-600 p-2 text-indigo-200 hover:bg-indigo-500 hover:bg-opacity-75 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600">
										<span className="absolute -inset-0.5" />
										<span className="sr-only">Open main menu</span>
										{open ? (
											<XMarkIcon className="block h-6 w-6" aria-hidden="true" />
										) : (
											<Bars3Icon className="block h-6 w-6" aria-hidden="true" />
										)}
									</Disclosure.Button>
								</div>
							</div>
						</div>

						<Disclosure.Panel className="md:hidden">
							<div className="space-y-1 px-2 pb-3 pt-2 sm:px-3">
								{navigation.map((item) => (
									<Disclosure.Button
										key={item.name}
										as="a"
										href={item.href}
										className={classNames(
											item.current ? "bg-indigo-700 text-white" : "text-white hover:bg-indigo-500 hover:bg-opacity-75",
											"block rounded-md px-3 py-2 text-base font-medium",
										)}
										aria-current={item.current ? "page" : undefined}
									>
										{item.name}
									</Disclosure.Button>
								))}
							</div>
							<div className="border-t border-indigo-700 pb-3 pt-4">
								<div className="flex items-center px-5">
									<div className="flex-shrink-0">
										<img className="h-10 w-10 rounded-full" src={user.imageUrl} alt="" />
									</div>
									<div className="ml-3">
										<div className="text-base font-medium text-white">{user.name}</div>
										<div className="text-sm font-medium text-indigo-300">{user.email}</div>
									</div>
									<button
										type="button"
										className="relative ml-auto flex-shrink-0 rounded-full bg-indigo-600 p-1 text-indigo-200 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600"
									>
										<span className="absolute -inset-1.5" />
										<span className="sr-only">View notifications</span>
										<BellIcon className="h-6 w-6" aria-hidden="true" />
									</button>
								</div>
								<div className="mt-3 space-y-1 px-2">
									{userNavigation.map((item) => (
										<Disclosure.Button
											key={item.name}
											as="a"
											href={item.href}
											className="block rounded-md px-3 py-2 text-base font-medium text-white hover:bg-indigo-500 hover:bg-opacity-75"
										>
											{item.name}
										</Disclosure.Button>
									))}
								</div>
							</div>
						</Disclosure.Panel>
					</>
				)}
			</Disclosure>

			<header className="bg-white shadow-sm">
				<div className="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
					<h1 className="text-lg font-semibold leading-6 text-gray-900">Dashboard</h1>
				</div>
			</header>
			<main>
				<div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
					<Outlet />
				</div>
			</main>
		</div>
	);
}
