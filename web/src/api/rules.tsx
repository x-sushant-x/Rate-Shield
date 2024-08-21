export interface rule {
    type: string;
    endpoint: string;
    http_method: string;
    bucket_capacity: string;
    token_add_rate: string;
}

interface getAllRuleResponse {
    data: rule[];
    status: string
}


export async function getAllRules(): Promise<rule[]> {
    const url = "http://127.0.0.1:8080/rules/all";

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
