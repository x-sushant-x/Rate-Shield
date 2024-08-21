import APIConfigurationHeader from "../components/APIConfigurationHeader";
import RulesTable from "../components/RulesTable";


export default function APIConfiguration() {
    return (
        <div className="h-screen bg-white rounded-xl shadow-lg">
            <APIConfigurationHeader />

            <RulesTable />

        </div>
    )
}
