import { faSearch } from "@fortawesome/free-solid-svg-icons/faSearch";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { rule } from "../api/rules";
import { useState } from "react";

interface Props {
    openAddOrUpdateRuleDialog: (rule: rule | null) => void
    setSearchRuleText: (searchText: string) => void
}

const APIConfigurationHeader: React.FC<Props> = ({ openAddOrUpdateRuleDialog, setSearchRuleText }) => {
    const [searchedText, setSearchedText] = useState('')

    return (
        <div className="p-8 flex bg-white justify-between rounded-none">
            <div className="block lg:hidden py-5">Appbar</div>
            <p className="text-[1.375rem] font-poppins font-medium text-slate-900">API</p>
            <div className="flex space-x-4">
                <div className="relative">
                    <input
                        className="bg-slate-200 pl-4 pr-4 py-2 rounded-md w-full focus:outline-none"
                        placeholder="Search Rules"
                        onChange={(e) => {
                            setSearchedText(e.target.value)
                        }}
                    />
                </div>
                <button className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center"
                    onClick={() => {
                        setSearchRuleText(searchedText)
                    }}>
                    <FontAwesomeIcon icon={faSearch} className="text-white" />
                </button>

                <button className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center" onClick={() => {
                    openAddOrUpdateRuleDialog(null)}}>
                    <span className="mr-2">+</span>
                    Add New
                </button>
            </div>  
        </div>
    )
}

export default APIConfigurationHeader