import { useState } from "react";
import APIConfigurationHeader from "../components/APIConfigurationHeader";
import RulesTable from "../components/RulesTable";
import AddNewRule from "../components/AddNewRule";

export default function APIConfiguration() {
    const [isAddNewRuleDialogOpen, setIsAddRuleDialogOpen] = useState(false)

    const openAddNewRuleDialog = () => {
        setIsAddRuleDialogOpen(true)
    }

    const closeAddNewRuleDialog = () => {
        setIsAddRuleDialogOpen(false)
    }

    return (
        <div className="h-screen bg-white rounded-xl shadow-lg">
            <APIConfigurationHeader openAddNewRuleDialog={openAddNewRuleDialog} />

            {
                isAddNewRuleDialogOpen ?
                    <AddNewRule closeAddNewRule={closeAddNewRuleDialog} />
                    :
                    <RulesTable />
            }
        </div>
    )
}
