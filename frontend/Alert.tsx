import React, { useState } from "react";
import { CheckCircleIcon, XCircleIcon, XMarkIcon } from "@heroicons/react/20/solid";
import { classNames } from "./utils";

export type AlertType = "success" | "failure";

type Props = {
	type: AlertType;
	message: string;
};

const backgroundColorByType: Record<AlertType, string> = {
	success: "bg-green-100",
	failure: "bg-red-100",
};

const backgroundHoverColorByType: Record<AlertType, string> = {
	success: "hover:bg-green-200",
	failure: "hover:bg-red-200",
};

const iconColorByType: Record<AlertType, string> = {
	success: "text-green-400",
	failure: "text-red-400",
};

const textColorByType: Record<AlertType, string> = {
	success: "text-green-800",
	failure: "text-red-800",
};

const ringColorByType: Record<AlertType, string> = {
	success: "focus:ring-green-600 focus:ring-offset-green-50",
	failure: "focus:ring-red-600 focus:ring-offset-red-50",
};

export default function Alert({ type, message }: Props) {
	const backgroundColor = backgroundColorByType[type];
	const backgroundHoverColor = backgroundHoverColorByType[type];
	const iconColor = iconColorByType[type];
	const textColor = textColorByType[type];
	const ringColor = ringColorByType[type];

	const [isHidden, setIsHidden] = useState(false);

	// https://tailwindui.com/components/application-ui/feedback/alerts#component-aa7cc38968c95d870db6ba62e76b8e0f
	return (
		<div className={classNames(backgroundColor, isHidden ? "hidden" : "", "rounded-md p-4 mb-6")}>
			<div className="flex">
				<div className="flex-shrink-0">
					{type === "success" && <CheckCircleIcon className={classNames(iconColor, "h-5 w-5")} aria-hidden="true" />}
					{type === "failure" && <XCircleIcon className={classNames(iconColor, "h-5 w-5")} aria-hidden="true" />}
				</div>
				<div className="ml-3">
					<p className={classNames(textColor, "text-sm font-medium")}>{message}</p>
				</div>
				<div className="ml-auto pl-3">
					<div className="-mx-1.5 -my-1.5" onClick={() => setIsHidden(true)}>
						<button
							type="button"
							className={classNames(
								backgroundColor,
								backgroundHoverColor,
								textColor,
								ringColor,
								"inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2",
							)}
						>
							<span className="sr-only">Dismiss</span>
							<XMarkIcon className="h-5 w-5" aria-hidden="true" />
						</button>
					</div>
				</div>
			</div>
		</div>
	);
}
