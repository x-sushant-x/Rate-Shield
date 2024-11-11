const baseUrl = "http://localhost:8080/";

export interface rule {
    strategy: string;
    endpoint: string;
    http_method: string;
    fixed_window_counter_rule: fixedWindowCounterRule | null;
    sliding_window_counter_rule: slidingWindowCounterRule | null;
    token_bucket_rule: tokenBucketRule | null;
    allow_on_error: boolean;
}

export interface paginatedRules {
    page: number;
    total_items: number;
    has_next_page: boolean;
    rules: rule[];
    status: string;
}

export interface paginatedRulesResponse {
    data: paginatedRules;
    status: string;
}

export interface fixedWindowCounterRule {
    max_requests: number;
    window: number;
}

export interface slidingWindowCounterRule {
    max_requests: number;
    window: number;
}


export interface tokenBucketRule {
    bucket_capacity: number;
    token_add_rate: number;
}

interface getAllRuleResponse {
    data: rule[];
    status: string;
}

export async function getAllRules(): Promise<rule[]> {
    const url = `${baseUrl}rule/list`;

    try {
        const response = await fetch(url, {
            method: "GET",
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data: getAllRuleResponse = await response.json();
        console.log("Response", data);

        return data.data;
    } catch (error) {
        console.error("Failed to fetch rules:", error);
        throw error;
    }
}

export async function getPaginatedRules(
    pageNumber: number,
): Promise<paginatedRulesResponse> {
    const url = `${baseUrl}/rule/list?page=${pageNumber}&items=10`;

    try {
        const response = await fetch(url, {
            method: "GET",
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data: paginatedRulesResponse = await response.json();

        console.log("Paginated Rules Response: " + data.data.rules);

        return data;
    } catch (error) {
        console.error("Failed to fetch rules:", error);
        throw error;
    }
}

export async function searchRulesViaEndpoint(
    searchText: string,
): Promise<rule[]> {
    const url = `${baseUrl}rule/search?endpoint=${searchText}`;

    try {
        const response = await fetch(url, {
            method: "GET",
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data: getAllRuleResponse = await response.json();
        console.log("Response", data);

        return data.data;
    } catch (error) {
        console.error("Failed to fetch rules:", error);
        throw error;
    }
}

export async function createNewRule(rule: rule) {
    const url = `${baseUrl}rule/add`;

    try {
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(rule),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }
    } catch (error) {
        console.error("Failed to add rule: ", error);
        throw error;
    }
}

export async function deleteRule(ruleKey: string) {
    const url = `${baseUrl}rule/delete`;

    try {
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                rule_key: ruleKey,
            }),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }
    } catch (error) {
        console.error("Failed to delete rule: ", error);
        throw error;
    }
}
