import { faSearch } from "@fortawesome/free-solid-svg-icons/faSearch";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { rule } from "../api/rules";

interface Props {
    openAddOrUpdateRuleDialog: (rule: rule | null) => void
}

const APIConfigurationHeader: React.FC<Props> = ({ openAddOrUpdateRuleDialog }) => {
    return (
        <div className="px-8 py-8 flex justify-between">
            <p className="text-[1.375rem] font-poppins font-medium text-slate-900">APIs Configurations</p>

            <div className="flex space-x-4">
                <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <FontAwesomeIcon icon={faSearch} className="text-gray-400" />
                    </div>
                    <input
                        className="bg-slate-200 pl-10 pr-4 py-2 rounded-md w-full focus:outline-none"
                        placeholder="Search"
                    />
                </div>
                <button className="bg-sidebar-bg text-slate-200 py-2 px-4 rounded-md flex items-center" onClick={() => {
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