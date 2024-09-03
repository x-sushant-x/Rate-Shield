import { useState } from "react";
import APIConfigurationHeader from "../components/APIConfigurationHeader";
import RulesTable from "../components/RulesTable";
import AddOrUpdateRule from "../components/AddOrUpdateRule";
import { rule } from "../api/rules";

export default function APIConfiguration() {
    const [isAddNewRuleDialogOpen, setIsAddRuleDialogOpen] = useState(false)
    const [selectedRule, setSelectedRule] = useState<rule | null>(null)

    const openAddOrUpdateRuleDialog = (rule: rule | null) => {
        setSelectedRule(rule)
        setIsAddRuleDialogOpen(true)
    }

    const closeAddNewRuleDialog = () => {
        setIsAddRuleDialogOpen(false)
        setSelectedRule(null)
    }

    return (
        <div className="h-screen bg-white rounded-xl shadow-lg">
            <APIConfigurationHeader openAddOrUpdateRuleDialog={openAddOrUpdateRuleDialog} />

            {
                isAddNewRuleDialogOpen ?
                    <AddOrUpdateRule
                        closeAddNewRule={closeAddNewRuleDialog}
                        strategy={selectedRule?.strategy ?? 'UNDEFINED'}
                        action={selectedRule ? "UPDATE" : "ADD"}
                        endpoint={selectedRule?.endpoint}
                        httpMethod={selectedRule?.http_method}
                        fixed_window_counter_rule={selectedRule?.fixed_window_counter_rule || null}
                        token_bucket_rule={selectedRule?.token_bucket_rule || null}
                    />
                    :
                    <RulesTable openAddOrUpdateRuleDialog={openAddOrUpdateRuleDialog} />
            }
        </div>
    )
}
