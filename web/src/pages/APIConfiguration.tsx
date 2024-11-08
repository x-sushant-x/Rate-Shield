import { useEffect, useRef, useState } from "react";
import APIConfigurationHeader from "../components/APIConfigurationHeader";
import RulesTable from "../components/RulesTable";
import AddOrUpdateRule from "../components/AddOrUpdateRule";
import { getPaginatedRules, rule, searchRulesViaEndpoint } from "../api/rules";
import toast from "react-hot-toast";
import { customToastStyle } from "../utils/toast_styles";

export default function APIConfiguration() {
    const [rulesData, setRulesData] = useState<rule[]>();
    const [searchRuleText, setSearchRuleText] = useState("");
    const errorShown = useRef(false);
    const [pageNumber, setPageNumber] = useState(1);
    const [hasNextPage, setHasNextPage] = useState(false);

    const fetchRules = async () => {
        try {
            const rules = await getPaginatedRules(pageNumber);
            setRulesData(rules.data.rules);
            setHasNextPage(rules.data.has_next_page);

            errorShown.current = false;
        } catch (error) {
            console.error("Failed to fetch rules:", error);
            if (errorShown.current === false) {
                toast.error("Error: " + error, {
                    style: customToastStyle,
                });
                errorShown.current = true;
            }
        }
    };

    useEffect(() => {
        fetchRules();
        console.log("Page Number: " + pageNumber);
    }, [pageNumber]);

    useEffect(() => {
        const searchRules = async () => {
            if (searchRuleText) {
                try {
                    const rules = await searchRulesViaEndpoint(searchRuleText);
                    setRulesData(rules);
                    errorShown.current = false;
                } catch (error) {
                    console.error("Failed to fetch rules:", error);
                    if (errorShown.current === false) {
                        toast.error("Error: " + error, {
                            style: customToastStyle,
                        });
                        errorShown.current = true;
                    }
                }
            }
        };
        searchRules();
    }, [searchRuleText]);

    const [isAddNewRuleDialogOpen, setIsAddRuleDialogOpen] = useState(false);
    const [selectedRule, setSelectedRule] = useState<rule | null>(null);

    const openAddOrUpdateRuleDialog = (rule: rule | null) => {
        setSelectedRule(rule);
        setIsAddRuleDialogOpen(true);
    };

    const closeAddNewRuleDialog = () => {
        setIsAddRuleDialogOpen(false);
        setSelectedRule(null);
    };

    return (
        <div className="bg-white rounded-xl h-full">
            <APIConfigurationHeader
                openAddOrUpdateRuleDialog={openAddOrUpdateRuleDialog}
                setSearchRuleText={setSearchRuleText}
            />

            {isAddNewRuleDialogOpen ? (
                <AddOrUpdateRule
                    closeAddNewRule={closeAddNewRuleDialog}
                    strategy={selectedRule?.strategy ?? "UNDEFINED"}
                    action={selectedRule ? "UPDATE" : "ADD"}
                    endpoint={selectedRule?.endpoint}
                    httpMethod={selectedRule?.http_method}
                    fixed_window_counter_rule={
                        selectedRule?.fixed_window_counter_rule || null
                    }
                    token_bucket_rule={selectedRule?.token_bucket_rule || null}
                    allow_on_error={selectedRule?.allow_on_error || false}
                />
            ) : (
                <RulesTable
                    openAddOrUpdateRuleDialog={openAddOrUpdateRuleDialog}
                    rulesData={rulesData}
                    hasNextPage={hasNextPage}
                    currentPageNumber={pageNumber}
                    setPageNumber={setPageNumber}
                />
            )}
        </div>
    );
}
