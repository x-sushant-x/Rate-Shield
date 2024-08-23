import { useState } from 'react';
import toast, { Toaster } from 'react-hot-toast';

import BackArrow from '../assets/BackArrow.png'
import { createNewRule, rule } from '../api/rules';
import { customToastStyle } from '../utils/toast_styles';

interface Props {
    closeAddNewRule: () => void
    action: string
    endpoint?: string
    httpMethod?: string
    bucketCapacity?: number
    tokenAddRate?: number
}



const AddOrUpdateRule: React.FC<Props> = ({ closeAddNewRule, action, bucketCapacity, endpoint, httpMethod, tokenAddRate }) => {
    const [apiEndpoint, setApiEndpoint] = useState(endpoint || '');
    const [strategy,] = useState('TOKEN BUCKET');
    const [method, setHttpMethod] = useState('GET');
    const [capacity, setBucketCapacity] = useState(bucketCapacity || 0);
    const [addRate, setTokenAddRate] = useState(tokenAddRate || 0);


    async function addRule() {
        const newRule: rule = {
            bucket_capacity: capacity,
            endpoint: apiEndpoint,
            http_method: method,
            token_add_rate: addRate,
            strategy: strategy
        }

        if (capacity === 0 || apiEndpoint === "" || httpMethod === "" || addRate === 0 || strategy === "") {
            toast.error("Please ensure entered data is valid.", {
                style: customToastStyle
            })
            return
        }

        try {
            await createNewRule(newRule)
            closeAddNewRule()
        } catch (e) {
            toast.error("Unable to add rule: " + e)
        }
    }

    // useEffect(() => {
    //     if (action === "UPDATE") {
    //         setApiEndpoint(endpoint || '');
    //         setHttpMethod(httpMethod || 'GET');
    //         setBucketCapacity(bucketCapacity ? bucketCapacity.toString() : '');
    //         setTokenAddRate(tokenAddRate ? tokenAddRate.toString() : '');
    //     }
    // }, [action, endpoint, httpMethod, bucketCapacity, tokenAddRate])

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

            <p className='mb-2 mt-6'>Strategy</p>
            <select
                className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto appearance-none">
                <option value="token-bucket">Token Bucket</option>
            </select>

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


            <p className='mb-2 mt-6'>Bucket Capacity</p>
            <input
                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                placeholder="Ex: - 10000"
                value={capacity}
                onChange={(e) => setBucketCapacity(Number.parseInt(e.target.value) || 0)}
            />

            <p className='mb-2 mt-6'>Token Add Rate (per minute)</p>
            <input
                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                placeholder="Ex: - 100"
                value={addRate}
                onChange={(e) => {
                    console.log('New Value: ', e.target.value)
                    setTokenAddRate(Number.parseInt(e.target.value) || 0)
                }}
            />

            <button className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center mt-8" onClick={() => {
                addRule()
            }}>
                {
                    action === "ADD" ? "Add" : "Update"
                }
            </button>

            <Toaster position="bottom-right" reverseOrder={false} />
        </div>
    )
}


export default AddOrUpdateRule