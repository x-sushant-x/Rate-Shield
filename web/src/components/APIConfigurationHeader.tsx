import { faSearch } from "@fortawesome/free-solid-svg-icons/faSearch";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { rule } from "../api/rules";
import { useState } from "react";
import logo from './../assets/logo.svg'


interface Props {
    openAddOrUpdateRuleDialog: (rule: rule | null) => void
    setSearchRuleText: (searchText: string) => void
}

const APIConfigurationHeader: React.FC<Props> = ({ openAddOrUpdateRuleDialog, setSearchRuleText }) => {
    const [searchedText, setSearchedText] = useState('')
 return (
 <div className="lg:p-8 p-4 flex flex-col lg:flex-row lg:items-center justify-between rounded-2xl">
    <div className="sm:block lg:hidden py-3 bg-black flex items-center justify-center rounded-2xl">
        <img src={logo} alt="logo" />
    </div>

    <p className="hidden lg:block text-[1.375rem] font-poppins font-medium text-slate-900">
        API Configurations
    </p>

    <div className="flex flex-col lg:flex-row w-full lg:w-auto items-center space-y-2 lg:space-y-0 lg:space-x-4 mt-4 lg:mt-0">
        {/* Search bar */}
        <div className="relative w-full lg:w-auto flex-grow">
            <input
                className="bg-slate-200 pl-4 pr-4 py-2 rounded-md w-full focus:outline-none"
                placeholder="Search Rules"
                onChange={(e) => {
                    setSearchedText(e.target.value)
                }}
            />
        </div>

    <div className="flex flex-row lg:flex-row w-full lg:w-auto gap-2 lg:gap-5 items-center mt-2 py-2 lg:mt-0">
    <button className="bg-sidebar-bg text-slate-200 py-3 px-5 rounded-md flex items-center justify-center w-full lg:flex-grow"
        onClick={() => {
            setSearchRuleText(searchedText)
        }}>
        <FontAwesomeIcon icon={faSearch} className="text-white" />
    </button>

    <button className="lg:hidden bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center justify-center w-full lg:w-auto"
        onClick={() => {
            openAddOrUpdateRuleDialog(null)
        }}>
        Add New
    </button>
   </div>
<button className="hidden lg:block bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center justify-center w-full lg:w-auto"
                onClick={() => {
                    openAddOrUpdateRuleDialog(null)
        }}>
            <span className="mr-2">+</span>
    Add New
</button>

    </div>
</div>

    )
}

export default APIConfigurationHeader