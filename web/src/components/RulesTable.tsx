import modifyRule from '../assets/modify_rule.png'
import { getAllRules, rule } from '../api/rules';
import { useEffect, useState } from 'react';

export default function RulesTable() {
    const [data, setData] = useState<rule[]>();

    const fetchRules = async () => {
        try {
            const rules = await getAllRules();
            setData(rules);
        } catch (error) {
            console.error("Failed to fetch rules:", error);
        }
    };

    useEffect(() => {
        fetchRules()
    }, [])



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
                    {data === undefined ? "No Rules Created" : data.map((item, index) => (
                        <>
                            <tr key={index}>
                                <td style={{ width: "50%" }} className='pt-6'>{item.endpoint}</td>
                                <td className="text-center pt-6" style={{ width: "15%" }}>{item.http_method}</td>  {/* Center aligned */}
                                <td className="text-center pt-6" style={{ width: "15%" }}>{item.type}</td>  {/* Center aligned */}
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