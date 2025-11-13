import logo from "../assets/logo.svg"

export function AuthPage() {
    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
            <div className="bg-white rounded-2xl border-gray-200 border w-1/3 h-fit">
                {/* Header */}
                <div className="bg-gray-800 p-4 rounded-t-2xl flex justify-between items-center">
                    <img src={logo} alt="Logo" className="h-8" />
                    <p className="text-gray-100 cursor-pointer hover:underline">Create Account</p>
                </div>

                {/* Auth Form */}
                <div className="p-6">
                    {/* Email Field */}
                    <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                        Email Address
                    </label>
                    <input
                        type="text"
                        id="email"
                        className="border border-gray-300 w-2/3 rounded-md p-3 text-gray-800 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all placeholder:text-sm"
                        placeholder="sushant.dhiman9812@gmail.com"
                    />

                    {/* Password Field */}
                    <label htmlFor="password" className="block text-sm font-medium text-gray-700 mt-4 mb-1">
                        Password
                    </label>
                    <input
                        type="password"
                        id="password"
                        className="border border-gray-300 w-2/3 rounded-md p-3 text-gray-800 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
                        placeholder="••••••••"
                    />

                    <br />
                    <br />

                    {/* Button */}
                    <div className="px-8 py-3 rounded-md text-gray-100 bg-gray-800 w-fit cursor-pointer hover:bg-gray-900 transition-all">
                        Go To Dashboard
                    </div>
                </div>
            </div>
        </div>
    )
}
