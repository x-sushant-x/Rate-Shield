import { Rule } from '../api/rules'
import modifyRule from '../assets/modify_rule.png'


const jsonData = `[
  {
    "endpoint": "/api/v1/get-user",
    "method": "GET",
    "strategy": "Token Bucket",
    "modifyRules": true
  },
  {
    "endpoint": "/api/v1/generate-otp",
    "method": "POST",
    "strategy": "Sliding Window",
    "modifyRules": true
  }
]
`
const data: Rule[] = JSON.parse(jsonData)

export default function RulesTable() {
    return (
        <div className="px-8 py-8">
            <table className="table-auto w-full text-left">
                <thead>
                    <tr>
                        <th className="text-left" style={{ width: "50%" }}>Endpoint</th>  {/* Left aligned */}
                        <th className="text-center" style={{ width: "15%" }}>Method</th>  {/* Center aligned */}
                        <th className="text-center" style={{ width: "15%" }}>Strategy</th>  {/* Center aligned */}
                        <th className="text-center" style={{ width: "20%" }}>Modify Rules</th>  {/* Center aligned */}
                    </tr>
                </thead>
                <tbody>
                    {data.map((item, index) => (
                        <>
                            <tr key={index}>
                                <td style={{ width: "50%" }} className='pt-6'>{item.endpoint}</td>
                                <td className="text-center pt-6" style={{ width: "15%" }}>{item.method}</td>  {/* Center aligned */}
                                <td className="text-center pt-6" style={{ width: "15%" }}>{item.strategy}</td>  {/* Center aligned */}
                                <td className="text-center pt-6" style={{ width: "20%" }}>
                                    <center>
                                        <img src={modifyRule} />
                                    </center>
                                </td>
                            </tr>

                            <tr>
                                <td colSpan={4}>
                                    <hr className="border-gray-300 mt-4" />
                                </td>
                            </tr>
                        </>
                    ))}
                </tbody>
            </table>
        </div>
    )
}