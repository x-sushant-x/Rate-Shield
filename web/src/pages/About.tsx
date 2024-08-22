const About = () => {
    return (
        <div className="h-screen bg-white rounded-xl shadow-lg">
            <div className="px-8 py-8">
                <p className="text-[1.375rem] font-poppins font-medium text-slate-900">About</p>

                <div>
                    <p className="mb-6 mt-4">
                        A completely configurable rate limiter that can apply rate limiting on individual APIs with individual rules.
                    </p>


                    <h4 className="text-sm font-semibold mb-4">🎯 Why?</h4>
                    <p className="mb-6">
                        Why not? I've got some free time, so I decided to build something.
                    </p>

                    <h4 className="text-sm font-semibold mb-4">🌟 Features</h4>
                    <ul className="list-disc list-inside mb-6 space-y-2">
                        <li>
                            🛠 <b>Customizable Limiting</b>: Apply rate limiting to individual APIs with tailored rules.
                        </li>
                        <li>
                            🖥️ <b>Dashboard</b>: Manage all your API rules in one place, with a user-friendly dashboard.
                        </li>
                        <li>
                            ⚙️ <b>Plug-and-Play Middleware</b>: Seamless integration with various frameworks, just plug it in and go.
                        </li>
                    </ul>

                    <h4 className="text-sm font-semibold mb-4">⚙️ Use Cases</h4>
                    <ul className="list-disc list-inside mb-6 space-y-2">
                        <li>
                            Preventing Abuse: Limit the number of requests an API can handle to prevent abuse or malicious activities.
                        </li>
                        <li>
                            Cost Management: Avoid overages on third-party API calls by rate limiting outgoing requests to those services.
                        </li>
                    </ul>

                    <h4 className="text-sm font-semibold mb-4">⚠️ Important</h4>
                    <ul className="list-disc list-inside mb-6 space-y-2">
                        <li>
                            Current Limitation: Only supports Token Bucket Rate Limiting, which may not suit all scenarios.
                        </li>
                        <li>
                            Under Development: This is a hobby project and still in progress. Not recommended for production use—yet! Stay tuned for v1.0.
                        </li>
                    </ul>
                </div>

            </div>
        </div>
    )
}

export default About