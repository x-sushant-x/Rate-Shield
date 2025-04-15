import { useState } from "react";
import toast, { Toaster } from "react-hot-toast";

import BackArrow from "../assets/BackArrow.png";
import {
    createNewRule,
    deleteRule,
    fixedWindowCounterRule,
    rule,
    slidingWindowCounterRule,
    tokenBucketRule,
} from "../api/rules";
import { customToastStyle } from "../utils/toast_styles";
import { validateNewFixedWindowCounterRule, validateNewRule, validateNewSlidingWindowCounterRule, validateNewTokenBucketRule } from "../utils/validators";

interface Props {
    closeAddNewRule: () => void;
    strategy: string;
    action: string;
    endpoint?: string;
    httpMethod?: string;
    fixed_window_counter_rule: fixedWindowCounterRule | null;
    token_bucket_rule: tokenBucketRule | null;
    sliding_window_counter_rule: slidingWindowCounterRule | null;
    allow_on_error: boolean;
}

const AddOrUpdateRule: React.FC<Props> = ({
    closeAddNewRule,
    action,
    strategy,
    endpoint,
    httpMethod,
    token_bucket_rule,
    fixed_window_counter_rule,
    sliding_window_counter_rule,
    allow_on_error,
}) => {
    const [apiEndpoint, setApiEndpoint] = useState(endpoint || "");
    const [limitStrategy, setLimitStrategy] = useState(strategy);
    const [method, setHttpMethod] = useState(httpMethod || "GET");
    const [tokenBucket, setTokenBucketRule] = useState(token_bucket_rule);
    const [fixedWindowCounter, setFixedWindowCounterRule] = useState(fixed_window_counter_rule);
    const [slidingWindowCounter, setSlidingWindowCounterRule] = useState(sliding_window_counter_rule)
    const [allowOnError, setAllowOnError] = useState(allow_on_error || false);

    const addOrUpdateRule = async () => {
        const newRule: rule = {
            endpoint: apiEndpoint,
            http_method: method,
            strategy: limitStrategy,
            fixed_window_counter_rule: fixedWindowCounter,
            token_bucket_rule: tokenBucket,
            sliding_window_counter_rule: slidingWindowCounter,
            allow_on_error: allowOnError,
        };
        

        if(!validateNewRule(newRule)) {
            console.log("validateNewRule")
            return;
        }

        if(!validateNewTokenBucketRule(newRule)) {
            console.log("validateNewTokenBucketRule")
            return;
        }

        if(!validateNewFixedWindowCounterRule(newRule)) {
            console.log("validateNewFixedWindowCounterRule")
            return;
        }
        
        if(!validateNewSlidingWindowCounterRule(newRule)) {
            console.log("validateNewSlidingWindowCounterRule")
            return;
        }

        try {
            await createNewRule(newRule);
            closeAddNewRule();
        } catch (error) {
            toast.error("Unable to save rule: " + error, {
                style: customToastStyle,
            });
        }
    };

    const handleAllowOnErrorCheckbox = () => {
        setAllowOnError(!allowOnError);
    };

    async function deleteExistingRule() {
        try {
            await deleteRule(apiEndpoint);
            closeAddNewRule();
        } catch (error) {
            toast.error("Unable to add rule: " + error);
        }
    }

    return (
        <div className="px-8">
            <div className="flex items-center mb-12">
                <img
                    src={BackArrow}
                    className="cursor-pointer"
                    width={25}
                    onClick={() => {
                        closeAddNewRule();
                    }}
                />
                <p className="text-xl ml-4">
                    {action === "ADD" ? "Add Rule" : "Update Rule"}
                </p>
            </div>

            <p className="mb-2">API Endpoint</p>
            <input
                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                required
                placeholder="Ex: - /api/v1/create"
                value={apiEndpoint}
                readOnly={action === "UPDATE" ? true : false}
                onChange={(e) => {
                    setApiEndpoint(e.target.value);
                }}
            />

            <br></br>

            {action === "UPDATE" ? (
                <div>
                    <p className="mb-2 mt-6">Strategy</p>
                    <select className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto appearance-none">
                        <option value={limitStrategy}>{limitStrategy}</option>
                    </select>
                </div>
            ) : (
                <div>
                    <div>
                        <p className="mb-2 mt-6">Strategy</p>
                        <select
                            className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto min-w-[50px] inline-block appearance-none"
                            onChange={(e) => {
                                setLimitStrategy(e.target.value);
                            }}
                        >
                            <option value="">Select Limiting Strategy</option>
                            <option value="TOKEN BUCKET">TOKEN BUCKET</option>
                            <option value="FIXED WINDOW COUNTER">
                                FIXED WINDOW COUNTER
                            </option>
                            <option value="SLIDING WINDOW COUNTER">SLIDING WINDOW COUNTER</option>
                        </select>
                    </div>
                </div>
            )}

            <p className="mb-2 mt-6">HTTP Method</p>
            <select
                className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto appearance-none"
                value={method}
                onChange={(e) => {
                    setHttpMethod(e.target.value);
                }}
            >
                <option value="GET">GET</option>
                <option value="POST">POST</option>
                <option value="DELETE">DELETE</option>
                <option value="PUT">PUT</option>
                <option value="PATCH">PATCH</option>
            </select>

            {limitStrategy === "TOKEN BUCKET" ? (
                <div>
                    <p className="mb-2 mt-6">Bucket Capacity</p>
                    <input
                        className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                        placeholder="Ex: - 10000"
                        value={tokenBucket?.bucket_capacity}
                        onChange={(e) =>
                            setTokenBucketRule({
                                bucket_capacity:
                                    Number.parseInt(e.target.value) || 0,
                                token_add_rate:
                                    tokenBucket?.token_add_rate || 0,
                                retention_time: 
                                    tokenBucket?.retention_time || 0
                            })
                        }
                    />

                    <p className="mb-2 mt-6">Token Add Rate (per minute)</p>
                    <input
                        className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                        placeholder="Ex: - 100"
                        value={tokenBucket?.token_add_rate}
                        onChange={(e) => {
                            setTokenBucketRule({
                                token_add_rate:
                                    Number.parseInt(e.target.value) || 0,
                                bucket_capacity:
                                    tokenBucket?.bucket_capacity || 0,
                                retention_time: 
                                    tokenBucket?.retention_time || 0
                            });
                        }}
                    />

                    <p className="mb-2 mt-6">Retention Time (in seconds)</p>
                    <input
                        className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-1/4"
                        placeholder="Time to keep inactive bucket. Default 60"
                        value={tokenBucket?.retention_time}
                        onChange={(e) => {
                            setTokenBucketRule({
                                token_add_rate:
                                    tokenBucket?.bucket_capacity || 0,
                                bucket_capacity:
                                    tokenBucket?.bucket_capacity || 0,
                                retention_time: 
                                    Number.parseInt(e.target.value) || 0
                            });
                        }}
                    />
                </div>
            ) : limitStrategy === "FIXED WINDOW COUNTER" || limitStrategy === "SLIDING WINDOW COUNTER" ? (
                <div>
                    <p className="mb-2 mt-6">Maximum Requests</p>
                    <input
                        className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                        placeholder="Ex: - 10000"
                        value={fixedWindowCounter?.max_requests}
                        onChange={(e) => {
                            if(limitStrategy === "FIXED WINDOW COUNTER") {
                                setFixedWindowCounterRule({
                                    max_requests: Number.parseInt(e.target.value),
                                    window: fixedWindowCounter?.window || 0,
                                });
                            } else if(limitStrategy === "SLIDING WINDOW COUNTER") {
                                setSlidingWindowCounterRule({
                                    max_requests: Number.parseInt(e.target.value),
                                    window: fixedWindowCounter?.window || 0,
                                });
                            }
                        }}
                    />

                    <p className="mb-2 mt-6">Window Time (in seconds)</p>
                    <input
                        className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                        placeholder="Ex: - 100"
                        value={fixedWindowCounter?.window}
                        onChange={(e) => {
                            if(limitStrategy === "FIXED WINDOW COUNTER") {
                                setFixedWindowCounterRule({
                                    max_requests:
                                        fixedWindowCounter?.max_requests || 0,
                                    window: Number.parseInt(e.target.value) || 0,
                                });
                            } else if(limitStrategy === "SLIDING WINDOW COUNTER") {
                                setSlidingWindowCounterRule({
                                    max_requests:
                                        fixedWindowCounter?.max_requests || 0,
                                    window: Number.parseInt(e.target.value) || 0,
                                });
                            }
                        }}
                    />
                </div>
            ) : (
                <div></div>
            )}

            <br></br>

            <label className="flex items-center space-x-3">
                <input
                    type="checkbox"
                    checked={allowOnError}
                    onChange={handleAllowOnErrorCheckbox}
                    className="appearance-none h-4 w-4 border border-gray-700 rounded-md checked:bg-blue-600 checked:border-transparent focus:outline-none transition duration-300 ease-in-out"
                />
                <span className="text-gray-700">Allow on error?</span>
            </label>

            <div className="flex">
                <button
                    className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center mt-8"
                    onClick={() => {
                        addOrUpdateRule();
                    }}
                >
                    {action === "ADD" ? "Add" : "Update"}
                </button>

                {action === "UPDATE" ? (
                    <button
                        className="bg-[#bb2124] text-slate-200 py-2 px-4 rounded-md flex items-center mt-8 ml-4"
                        onClick={() => {
                            deleteExistingRule();
                        }}
                    >
                        Delete
                    </button>
                ) : (
                    <div></div>
                )}
            </div>

            <Toaster position="bottom-right" reverseOrder={false} />
        </div>
    );
};

export default AddOrUpdateRule;
