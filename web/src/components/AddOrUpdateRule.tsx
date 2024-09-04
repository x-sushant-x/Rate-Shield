import { useState } from 'react';
import toast, { Toaster } from 'react-hot-toast';

import BackArrow from '../assets/BackArrow.png'
import { createNewRule, deleteRule, fixedWindowCounterRule, rule, tokenBucketRule } from '../api/rules';
import { customToastStyle } from '../utils/toast_styles';

interface Props {
    closeAddNewRule: () => void
    strategy: string
    action: string
    endpoint?: string
    httpMethod?: string
    fixed_window_counter_rule: fixedWindowCounterRule | null
    token_bucket_rule: tokenBucketRule | null
}



const AddOrUpdateRule: React.FC<Props> = ({ closeAddNewRule, action, strategy, endpoint, httpMethod, token_bucket_rule, fixed_window_counter_rule }) => {
    const [apiEndpoint, setApiEndpoint] = useState(endpoint || '');
    const [limitStrategy, setLimitStrategy] = useState(strategy);
    const [method, setHttpMethod] = useState(httpMethod || '');
    const [tokenBucket, setTokenBucketRule] = useState(token_bucket_rule)
    const [fixedWindowCounter, setFixedWindowCounterRule] = useState(fixed_window_counter_rule)


    const addOrUpdateRule = async () => {
        const newRule : rule = {
            endpoint: apiEndpoint,
            http_method: method,
            strategy: limitStrategy,
            fixed_window_counter_rule: fixedWindowCounter,
            token_bucket_rule: tokenBucket
        }

        try {
            await createNewRule(newRule)
            closeAddNewRule()
        } catch(error) {
            toast.error("Unable to save rule: " + error, {
                style: customToastStyle
            })
        }
    }



    async function deleteExistingRule() {
        try {
            await deleteRule(apiEndpoint)
            closeAddNewRule()
        } catch (error) {
            toast.error("Unable to add rule: " + error)
        }
    }

    return (
        <div className="px-8">
            <div className='flex items-center mb-12'>
                <img src={BackArrow} className='cursor-pointer' width={25} onClick={
                    () => {
                        closeAddNewRule()
                    }
                } />
                <p className="text-xl ml-4">{action === "ADD" ? "Add Rule" : "Update Rule"}</p>
            </div>

            <p className='mb-2'>API Endpoint</p>
            <input
                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                required
                placeholder="Ex: - /api/v1/create"
                value={apiEndpoint}
                onChange={(e) => {
                    setApiEndpoint(e.target.value)
                }}
            />

            <br></br>

            {
                action === "UPDATE" ?
                    <div><p className='mb-2 mt-6'>Strategy</p>
                        <select
                            className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto appearance-none">
                            <option value={limitStrategy}>{limitStrategy}</option>
                        </select>
                    </div>
                    :
                    <div>
                        <div><p className='mb-2 mt-6'>Strategy</p>
                            <select
                                className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto min-w-[50px] inline-block appearance-none" onChange={(e) => {
                                    setLimitStrategy(e.target.value)
                                }}>
                                <option value=''>Select Limiting Strategy</option>
                                <option value='TOKEN BUCKET'>TOKEN BUCKET</option>
                                <option value='FIXED WINDOW COUNTER'>FIXED WINDOW COUNTER</option>
                            </select>
                        </div>
                    </div>
            }

            <p className='mb-2 mt-6'>HTTP Method</p>
            <select
                className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto appearance-none"
                value={method}
                onChange={(e) => {
                    setHttpMethod(e.target.value)
                }}>

                <option value="GET">GET</option>
                <option value="POST">POST</option>
                <option value="DELETE">DELETE</option>
                <option value="PUT">PUT</option>
                <option value="PATCH">PATCH</option>
            </select>


            {
                limitStrategy === "TOKEN BUCKET" ?
                    <div>
                        <p className='mb-2 mt-6'>Bucket Capacity</p>
                        <input
                            className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                            placeholder="Ex: - 10000"
                            value={tokenBucket?.bucket_capacity}
                            onChange={(e) => setTokenBucketRule({
                                bucket_capacity: Number.parseInt(e.target.value) || 0,
                                token_add_rate: tokenBucket?.token_add_rate || 0
                            })}
                        />

                        <p className='mb-2 mt-6'>Token Add Rate (per minute)</p>
                        <input
                            className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                            placeholder="Ex: - 100"
                            value={tokenBucket?.token_add_rate}
                            onChange={(e) => {
                                setTokenBucketRule({
                                    token_add_rate: Number.parseInt(e.target.value) || 0,
                                    bucket_capacity: tokenBucket?.bucket_capacity || 0
                                })
                            }}
                        />
                    </div>

                    : limitStrategy === "FIXED WINDOW COUNTER" ?
                        <div>
                            <p className='mb-2 mt-6'>Maximum Requests</p>
                            <input
                                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                                placeholder="Ex: - 10000"
                                value={fixedWindowCounter?.max_requests}
                                onChange={(e) => {
                                    setFixedWindowCounterRule({
                                        max_requests: Number.parseInt(e.target.value),
                                        window: fixedWindowCounter?.window || 0
                                    })
                                }}
                            />

                            <p className='mb-2 mt-6'>Window Time (in seconds)</p>
                            <input
                                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                                placeholder="Ex: - 100"
                                value={fixedWindowCounter?.window}
                                onChange={(e) => {
                                    setFixedWindowCounterRule({
                                        max_requests: fixedWindowCounter?.max_requests || 0,
                                        window: Number.parseInt(e.target.value) || 0
                                    })
                                }}
                            />
                        </div>

                        : <div></div>
            }

            <div className='flex'>
                <button className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center mt-8" onClick={() => {
                    addOrUpdateRule()
                }}>
                    {
                        action === "ADD" ? "Add" : "Update"
                    }
                </button>

                {
                    action === "UPDATE" ?
                        <button className="bg-[#bb2124] text-slate-200 py-2 px-4 rounded-md flex items-center mt-8 ml-4" onClick={() => {
                            deleteExistingRule()
                        }}>
                            Delete
                        </button> : <div></div>
                }
            </div>

            <Toaster position="bottom-right" reverseOrder={false} />
        </div>
    )
}


export default AddOrUpdateRule