const baseUrl = "http://127.0.0.1:8080/"

export interface rule {
    strategy: string;
    endpoint: string;
    http_method: string;
    bucket_capacity: number;
    token_add_rate: number;
}

interface getAllRuleResponse {
    data: rule[];
    status: string
}


export async function getAllRules(): Promise<rule[]> {
    const url = `${baseUrl}rules/all`;

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
    const url = `${baseUrl}rules/add`;

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(rule)
        })

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }
    } catch (error) {
        console.error("Failed to add rule: ", error)
        throw error
    }
}