import { useState } from "react";
import APIConfigurationHeader from "../components/APIConfigurationHeader";
import RulesTable from "../components/RulesTable";
import AddNewRule from "../components/AddOrUpdateRule";
import { rule } from "../api/rules";

export default function APIConfiguration() {
    const [isAddNewRuleDialogOpen, setIsAddRuleDialogOpen] = useState(false)
    const [selectedRule, setSelectedRule] = useState<rule | null>(null)

    const openAddNewRuleDialog = (rule: rule | null) => {
        setSelectedRule(rule)
        setIsAddRuleDialogOpen(true)
    }

    const closeAddNewRuleDialog = () => {
        setIsAddRuleDialogOpen(false)
        setSelectedRule(null)
    }

    return (
        <div className="h-screen bg-white rounded-xl shadow-lg">
            <APIConfigurationHeader openAddNewRuleDialog={openAddNewRuleDialog} />

            {
                isAddNewRuleDialogOpen ?
                    <AddNewRule
                        closeAddNewRule={closeAddNewRuleDialog}
                        action={selectedRule ? "UPDATE" : "ADD"}
                        endpoint={selectedRule?.endpoint}
                        httpMethod={selectedRule?.http_method}
                        bucketCapacity={selectedRule?.bucket_capacity}
                        tokenAddRate={selectedRule?.token_add_rate}
                    />
                    :
                    <RulesTable openAddOrUpdateRuleDialog={openAddNewRuleDialog} />
            }
        </div>
    )
}
