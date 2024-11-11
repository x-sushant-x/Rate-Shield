import toast from "react-hot-toast";
import { rule } from "../api/rules";
import { customToastStyle } from "./toast_styles";

export function validateNewRule(newRule: rule) {
    if (newRule.endpoint === "" || newRule.endpoint === undefined) {
        toast.error("API Endpoint can't be null.", {
            style: customToastStyle,
        });
        return false;
    }

    if (
        newRule.strategy === "" ||
        newRule.strategy === undefined ||
        newRule.strategy === "UNDEFINED"
    ) {
        toast.error("API limit strategy can't be null.", {
            style: customToastStyle,
        });
        return false;
    }

    if (newRule.http_method === "" || newRule.http_method === undefined) {
        toast.error("API HTTP Method can't be null.", {
            style: customToastStyle,
        });
        return false;
    }
    return true
}

export function validateNewTokenBucketRule(newRule: rule) {
    if (newRule.strategy === "TOKEN BUCKET") {
        if (
            newRule.token_bucket_rule?.bucket_capacity === 0 ||
            !newRule.token_bucket_rule?.bucket_capacity ||
            newRule.token_bucket_rule.bucket_capacity <= 0
        ) {
            toast.error("Invalid value for bucket capacity.", {
                style: customToastStyle,
            });
            return false;
        }

        if (
            newRule.token_bucket_rule?.token_add_rate === 0 ||
            !newRule.token_bucket_rule?.token_add_rate ||
            newRule.token_bucket_rule.token_add_rate <= 0
        ) {
            toast.error("Invalid value for bucket capacity.", {
                style: customToastStyle,
            });
            return false;
        }

        if (
            newRule.token_bucket_rule?.token_add_rate >
            newRule.token_bucket_rule.bucket_capacity
        ) {
            toast.error(
                "Token add rate should not be more than bucket capacity.",
                {
                    style: customToastStyle,
                },
            );
            return false;
        }
    }
    return true
}

export function validateNewFixedWindowCounterRule(newRule: rule) {
    if (newRule.strategy === "FIXED WINDOW COUNTER") {
        if (
            newRule.fixed_window_counter_rule?.max_requests === 0 ||
            !newRule.fixed_window_counter_rule?.max_requests ||
            newRule.fixed_window_counter_rule?.max_requests <= 0
        ) {
            toast.error("Invalid value for maximum requests.", {
                style: customToastStyle,
            });
            return false;
        }

        if (
            newRule.fixed_window_counter_rule?.window === 0 ||
            !newRule.fixed_window_counter_rule?.window ||
            newRule.fixed_window_counter_rule?.window <= 0
        ) {
            toast.error("Invalid value for window time.", {
                style: customToastStyle,
            });
            return false;
        }
    }
    return true
}

export function validateNewSlidingWindowCounterRule(newRule: rule) {
    if (newRule.strategy === "FIXED WINDOW COUNTER") {
        if (
            newRule.sliding_window_counter_rule?.max_requests === 0 ||
            !newRule.sliding_window_counter_rule?.max_requests ||
            newRule.sliding_window_counter_rule?.max_requests <= 0
        ) {
            toast.error("Invalid value for maximum requests.", {
                style: customToastStyle,
            });
            return false;
        }

        if (
            newRule.sliding_window_counter_rule?.window === 0 ||
            !newRule.sliding_window_counter_rule?.window ||
            newRule.sliding_window_counter_rule?.window <= 0
        ) {
            toast.error("Invalid value for window time.", {
                style: customToastStyle,
            });
            return false;
        }
    }
    return true
}