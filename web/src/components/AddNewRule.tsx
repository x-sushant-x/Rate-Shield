import BackArrow from '../assets/BackArrow.png'

interface Props {
    closeAddNewRule: () => void
}

const AddNewRule: React.FC<Props> = ({ closeAddNewRule }) => {
    return (
        <div className="px-8">
            <div className='flex items-center mb-12'>
                <img src={BackArrow} className='cursor-pointer' width={25} onClick={
                    () => {
                        closeAddNewRule()
                    }
                } />
                <p className="text-xl ml-4">Add Rule</p>
            </div>

            <p className='mb-2'>API Endpoint</p>
            <input
                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                placeholder="Ex: - /api/v1/create"
            />

            <br></br>

            <p className='mb-2 mt-6'>Strategy</p>
            <select
                className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto appearance-none">
                <option value="token-bucket">Token Bucket</option>
            </select>

            <p className='mb-2 mt-6'>HTTP Method</p>
            <select
                className="bg-slate-200 px-4 py-2 rounded-md focus:outline-none w-auto appearance-none">
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
            />

            <p className='mb-2 mt-6'>Token Add Rate</p>
            <input
                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md  focus:outline-none w-auto"
                placeholder="Ex: - 100"
            />

            <button className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center mt-8" onClick={() => { }}>
                Add Rule
            </button>

        </div>
    )
}

export default AddNewRule